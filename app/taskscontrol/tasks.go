package taskscontrol

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/mobiquai/go_final_project/app/appsettings"
	"github.com/mobiquai/go_final_project/app/database"
	"github.com/mobiquai/go_final_project/app/service"
)

func responseWithError(w http.ResponseWriter, errorText string, err error, statusCode int) {
	errorResponse := ErrorResponse{fmt.Errorf("%s: %w", errorText, err).Error()}

	errorData, _ := json.Marshal(errorResponse)
	w.WriteHeader(statusCode)

	_, err = w.Write(errorData)
	if err != nil {
		http.Error(w, fmt.Errorf("error: %w", err).Error(), statusCode)
	}

}

func AddTask(w http.ResponseWriter, r *http.Request) {
	var taskData database.Task
	var buffer bytes.Buffer

	if _, err := buffer.ReadFrom(r.Body); err != nil {
		responseWithError(w, "ошибка получения тела запроса", err, http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(buffer.Bytes(), &taskData); err != nil {
		responseWithError(w, "json encoding error", err, http.StatusBadRequest)
		return
	}

	if len(taskData.Date) == 0 {
		taskData.Date = time.Now().Format(appsettings.DateLayout)

	} else {
		date, err := time.Parse(appsettings.DateLayout, taskData.Date)
		if err != nil {
			responseWithError(w, "неправильный формат даты", err, http.StatusBadRequest)
			return
		}

		if date.Before(time.Now()) {
			taskData.Date = time.Now().Format(appsettings.DateLayout)
		}

	}

	if len(taskData.Title) == 0 {
		responseWithError(w, "неправильное значение title", errors.New("title is empty"), http.StatusBadRequest)
		return
	}

	if len(taskData.Repeat) > 0 {
		if _, err := service.NextDate(time.Now(), taskData.Date, taskData.Repeat); err != nil {
			responseWithError(w, "неправильный формат repeat", errors.New("there is no such format"), http.StatusBadRequest)
			return
		}
	}

	taskId, err := database.AddTask(taskData)
	if err != nil {
		responseWithError(w, "ошибка создания новой задачи", err, http.StatusInternalServerError)
		return
	}

	taskIdData, err := json.Marshal(TaskIdResponse{Id: taskId})
	if err != nil {
		responseWithError(w, "json decoding error", err, http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)

	_, err = w.Write(taskIdData)
	if err != nil {
		responseWithError(w, "ошибка записи задачи по id", err, http.StatusBadRequest)
	}

	log.Printf("Успешно добавлена новая задача с id=%d\n", taskId)

}

func ReadTasks(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search") // получаем параметр запроса "search"

	var tasks []database.Task

	if len(search) > 0 { // параметр "search" задан
		date, err := time.Parse("02.01.2006", search)
		if err != nil {
			tasks, err = database.SearchTasks(search) // поиск по строке
			if err != nil {
				responseWithError(w, "не удалось найти задачи", err, http.StatusInternalServerError)
				return
			}
		} else {
			tasks, err = database.SearchTasksByDate(date.Format(appsettings.DateLayout)) // поиск по дате, если параметр соответствует формату
			if err != nil {
				responseWithError(w, "не удалось найти задачи по дате", err, http.StatusInternalServerError)
				return
			}
		}

	} else {
		err := errors.New("")
		tasks, err = database.TasksRead()
		if err != nil {
			responseWithError(w, "ошибка получения списка задач", err, http.StatusInternalServerError)
			return
		}
	}

	tasksData, err := json.Marshal(database.Tasks{Tasks: tasks})
	if err != nil {
		responseWithError(w, "json decoding error", err, http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(tasksData)
	if err != nil {
		responseWithError(w, "ошибка записи задач", err, http.StatusBadRequest)
	}

	log.Printf("Прочитано %d задач", len(tasks))

}

func ReadTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	task, err := database.ReadTask(id)
	if err != nil {
		responseWithError(w, "не удалось получить задачу", err, http.StatusInternalServerError)
		return
	}

	tasksData, err := json.Marshal(task)
	if err != nil {
		responseWithError(w, "json decoding error", err, http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(tasksData)
	if err != nil {
		responseWithError(w, "ошибка записи задачи", err, http.StatusBadRequest)
	}

	log.Printf("Прочитана задача с id=%s", id)

}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	var task database.Task
	var buffer bytes.Buffer

	if _, err := buffer.ReadFrom(r.Body); err != nil {
		responseWithError(w, "ошибка получения тела запроса", err, http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(buffer.Bytes(), &task); err != nil {
		responseWithError(w, "json encoding error", err, http.StatusBadRequest)
		return
	}

	if len(task.Id) == 0 {
		responseWithError(w, "неправильный id", errors.New("id is empty"), http.StatusBadRequest)
		return
	}

	if _, err := strconv.Atoi(task.Id); err != nil {
		responseWithError(w, "неправильный id", err, http.StatusBadRequest)
		return
	}

	if _, err := time.Parse(appsettings.DateLayout, task.Date); err != nil {
		responseWithError(w, "направильный date", err, http.StatusBadRequest)
		return
	}

	if len(task.Title) == 0 {
		responseWithError(w, "неправильный title", errors.New("title is empty"), http.StatusBadRequest)
		return
	}

	if len(task.Repeat) > 0 {
		if _, err := service.NextDate(time.Now(), task.Date, task.Repeat); err != nil {
			responseWithError(w, "неправильный формат repeat", errors.New("there is no such format"), http.StatusBadRequest)
			return
		}
	}

	_, err := database.UpdateTask(task)
	if err != nil {
		responseWithError(w, "не удалось обновить задачу", errors.New("failed to update the task"), http.StatusInternalServerError)
		return
	}

	taskIdData, err := json.Marshal(task)
	if err != nil {
		responseWithError(w, "json decoding error", err, http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(taskIdData)
	if err != nil {
		responseWithError(w, "ошибка обновления задачи", err, http.StatusBadRequest)
		return
	}

	log.Printf("Задача обновлена с id=%s", task.Id)

}

func TaskDone(w http.ResponseWriter, r *http.Request) {
	task, err := database.ReadTask(r.URL.Query().Get("id"))
	if err != nil {
		responseWithError(w, "не удалось получить задачу", err, http.StatusInternalServerError)
		return
	}

	if len(task.Repeat) == 0 {
		err = database.DeleteTask(task.Id)
		if err != nil {
			responseWithError(w, "ну удалось удалить задачу", err, http.StatusInternalServerError)
			return
		}
	} else {
		task.Date, err = service.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			responseWithError(w, "не удалось получить следующую дату", err, http.StatusBadRequest)
			return
		}

		_, err = database.UpdateTask(task)
		if err != nil {
			responseWithError(w, "не удалось обновить задачу", err, http.StatusInternalServerError)
			return
		}
	}

	tasksData, err := json.Marshal(struct{}{})
	if err != nil {
		responseWithError(w, "json decoding error", err, http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(tasksData)
	if err != nil {
		responseWithError(w, "ошибка записи задачи", err, http.StatusBadRequest)
	}

	log.Printf("Выполнена задача с id=%s", task.Id)

}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	err := database.DeleteTask(id)
	if err != nil {
		responseWithError(w, "не удалось удалить задачу", err, http.StatusInternalServerError)
		return
	}

	tasksData, err := json.Marshal(struct{}{})
	if err != nil {
		responseWithError(w, "json decoding error", err, http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(tasksData)
	if err != nil {
		responseWithError(w, "ошибка записи задачи", err, http.StatusBadRequest)
		return
	}

	log.Printf("Удалена задача с id=%s", id)

}
