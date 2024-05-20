package handler

import (
	"log"
	"net/http"
)

const webDir = "./web"

func GetFrontEnd() http.Handler {
	log.Printf("Загружены файлы фронтенда, расположенные по пути: %s\n", webDir)

	return http.FileServer(http.Dir(webDir))
}
