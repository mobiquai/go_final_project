package server

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/mobiquai/go_final_project/tests"
)

const host = "localhost"

type Server struct {
	httpServer *http.Server
	Handler    http.Handler
}

func (serv *Server) Start(router http.Handler) error {
	serv.httpServer = &http.Server{
		Addr:    getAddr(),
		Handler: router,
	}

	log.Printf("Сервер запущен по пути: %s", serv.httpServer.Addr)

	return serv.httpServer.ListenAndServe()
}

func getAddr() string {
	port := tests.Port
	envPort := os.Getenv("TODO_PORT") // получаем значение переменной окружения

	if len(envPort) > 0 {
		if eport, err := strconv.ParseInt(envPort, 10, 32); err == nil {
			port = int(eport)
			log.Printf("Получен номер порта из переменной окружения TODO_PORT: %d", port)
		}
	} else {
		log.Printf("Получен номер порта из файла settings.go: %d", port)
	}

	return host + ":" + strconv.Itoa(port)
}
