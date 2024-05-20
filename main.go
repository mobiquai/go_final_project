package main

import (
	"log"

	_ "github.com/mattn/go-sqlite3"

	"github.com/go-chi/chi/v5"

	"github.com/mobiquai/go_final_project/app/database"
	"github.com/mobiquai/go_final_project/app/handler"
	"github.com/mobiquai/go_final_project/app/middleware"
	"github.com/mobiquai/go_final_project/app/server"
	"github.com/mobiquai/go_final_project/app/taskscontrol"
)

func main() {
	database.InitiateDb()

	r := chi.NewRouter()

	r.Mount("/", handler.GetFrontEnd()) // Файл-сервер

	r.Get("/api/nextdate", taskscontrol.NextDate) // Служебная функция nexdate

	r.Post("/api/task", middleware.Auth(taskscontrol.AddTask))       // POST-запрос по добавлению новой задачи
	r.Get("/api/tasks", middleware.Auth(taskscontrol.ReadTasks))     // GET-запрос на получение существующих задач
	r.Get("/api/task", middleware.Auth(taskscontrol.ReadTask))       // GET-запрос на получение данных задачи по ее id
	r.Put("/api/task", middleware.Auth(taskscontrol.UpdateTask))     // PUT-запрос на редактирование данных задачи по ее id
	r.Post("/api/task/done", middleware.Auth(taskscontrol.TaskDone)) // POST-запрос по установке задачи выполненной по ее id
	r.Delete("/api/task", middleware.Auth(taskscontrol.DeleteTask))  // DELETE-запрос на удаление задачи по ее id

	r.Post("/api/signin", taskscontrol.Sign) // Проверка введеного пароля

	server := new(server.Server)
	if err := server.Start(r); err != nil {
		log.Fatalf("Невозможно запустить сервер: %v", err)
		return
	}

	log.Println("Сервер остановлен")

}
