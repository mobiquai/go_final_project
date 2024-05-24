package handler

import (
	"log"
	"net/http"

	"github.com/mobiquai/go_final_project/app/appsettings"
)

func GetFrontEnd() http.Handler {
	log.Printf("Загружены файлы фронтенда, расположенные по пути: %s\n", appsettings.WebDir)

	return http.FileServer(http.Dir(appsettings.WebDir))
}
