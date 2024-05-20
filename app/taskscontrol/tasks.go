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

	"github.com/mobiquai/go_final_project/app/database"
	"github.com/mobiquai/go_final_project/app/service"
)

func responseWithError(w http.ResponseWriter, errorText string, err error) {
	errorResponse := ErrorResponse{fmt.Errorf("%s: %w", errorText, err).Error()}

	errorData, _ := json.Marshal(errorResponse)
	w.WriteHeader(http.StatusBadRequest)

	_, err = w.Write(errorData)
	if err != nil {
		http.Error(w, fmt.Errorf("error: %w", err).Error(), http.StatusBadRequest)
	}

}

func AddTask(w http.ResponseWriter, r *http.Request) {
	var taskData database.Task
	var buffer bytes.Buffer

	if _, err := buffer.ReadFrom(r.Body); err != nil {
		responseWithError(w, "ошибка получения тела запроса", err)
		return
	}

	if err := json.Unmarshal(buffer.Bytes(), &taskData); err != nil {
		responseWithError(w, "json encoding error", err)
		return
	}

	if len(taskData.Date) == 0 {
		taskData.Date = time.Now().Format(dateLayout)

	} else {
		date, err := time.Parse(dateLayout, taskData.Date)
		if err != nil {
			responseWithError(w, "неправильный формат даты", err)
			return
		}

		if date.Before(time.Now()) {
			taskData.Date = time.Now().Format(dateLayout)
		}

	}

	if len(taskData.Title) == 0 {
		responseWithError(w, "неправильное значение title", errors.New("title пуст"))
		return
	}

	if len(taskData.Repeat) > 0 {
		if _, err := service.NextDate(time.Now(), taskData.Date, taskData.Repeat); err != nil {
			responseWithError(w, "неправильный формат repeat", errors.New("не существует такого формата"))
			return
		}
	}

	taskId, err := database.AddTask(taskData)
	if err != nil {
		responseWithError(w, "ошибка создания новой задачи", err)
		return
	}

	taskIdData, err := json.Marshal(TaskIdResponse{Id: taskId})
	if err != nil {
		responseWithError(w, "json decoding error", err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)

	_, err = w.Write(taskIdData)
	if err != nil {
		responseWithError(w, "ошибка записи задачи по id", err)
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
				responseWithError(w, "не удалось найти задачи", err)
				return
			}
		} else {
			tasks, err = database.SearchTasksByDate(date.Format(dateLayout)) // поиск по дате, если параметр соответствует формату
			if err != nil {
				responseWithError(w, "не удалось найти задачи по дате", err)
				return
			}
		}

	} else {
		err := errors.New("")
		tasks, err = database.TasksRead()
		if err != nil {
			responseWithError(w, "ошибка в получении задач", err)
			return
		}
	}

	tasksData, err := json.Marshal(database.Tasks{Tasks: tasks})
	if err != nil {
		responseWithError(w, "json decoding error", err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(tasksData)
	if err != nil {
		responseWithError(w, "ошибка записи задач", err)
	}

	log.Printf("Прочитано %d задач", len(tasks))

}

func ReadTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	task, err := database.ReadTask(id)
	if err != nil {
		responseWithError(w, "не удалось получить задачу", err)
		return
	}

	tasksData, err := json.Marshal(task)
	if err != nil {
		responseWithError(w, "json decoding error", err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(tasksData)
	if err != nil {
		responseWithError(w, "ошибка записи задачи", err)
	}

	log.Printf("Прочитана задача с id=%s", id)

}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	var task database.Task
	var buffer bytes.Buffer

	if _, err := buffer.ReadFrom(r.Body); err != nil {
		responseWithError(w, "ошибка получения тела запроса", err)
		return
	}

	if err := json.Unmarshal(buffer.Bytes(), &task); err != nil {
		responseWithError(w, "json encoding error", err)
		return
	}

	if len(task.Id) == 0 {
		responseWithError(w, "неправильный id", errors.New("id пустой"))
		return
	}

	if _, err := strconv.Atoi(task.Id); err != nil {
		responseWithError(w, "неправильный id", err)
		return
	}

	if _, err := time.Parse(dateLayout, task.Date); err != nil {
		responseWithError(w, "направильный date", err)
		return
	}

	if len(task.Title) == 0 {
		responseWithError(w, "неправильный title", errors.New("title пустой"))
		return
	}

	if len(task.Repeat) > 0 {
		if _, err := service.NextDate(time.Now(), task.Date, task.Repeat); err != nil {
			responseWithError(w, "неправильный формат repeat", errors.New("не существует такого формата"))
			return
		}
	}

	_, err := database.UpdateTask(task)
	if err != nil {
		responseWithError(w, "неправильный title", errors.New("не удалось обновить задачу"))
		return
	}

	taskIdData, err := json.Marshal(task)
	if err != nil {
		responseWithError(w, "json decoding error", err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(taskIdData)
	if err != nil {
		responseWithError(w, "ошибка обновления задачи", err)
		return
	}

	log.Printf("Задача обновлена с id=%s", task.Id)

}

func TaskDone(w http.ResponseWriter, r *http.Request) {
	task, err := database.ReadTask(r.URL.Query().Get("id"))
	if err != nil {
		responseWithError(w, "не удалось получить задачу", err)
		return
	}

	if len(task.Repeat) == 0 {
		err = database.DeleteTask(task.Id)
		if err != nil {
			responseWithError(w, "ну удалось удалить задачу", err)
			return
		}
	} else {
		task.Date, err = service.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			responseWithError(w, "не удалось получить следующую дату", err)
			return
		}

		_, err = database.UpdateTask(task)
		if err != nil {
			responseWithError(w, "не удалось обновить задачу", err)
			return
		}
	}

	tasksData, err := json.Marshal(struct{}{})
	if err != nil {
		responseWithError(w, "json decoding error", err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(tasksData)
	if err != nil {
		responseWithError(w, "ошибка записи задачи", err)
	}

	log.Printf("Выполнена задача с id=%s", task.Id)

}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	err := database.DeleteTask(id)
	if err != nil {
		responseWithError(w, "не удалось удалить задачу", err)
		return
	}

	tasksData, err := json.Marshal(struct{}{})
	if err != nil {
		responseWithError(w, "json decoding error", err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(tasksData)
	if err != nil {
		responseWithError(w, "ошибка записи задачи", err)
		return
	}

	log.Printf("Удалена задача с id=%s", id)

}
