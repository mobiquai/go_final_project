package taskscontrol

import (
	"fmt"
	"net/http"
	"time"

	"github.com/mobiquai/go_final_project/app/service"
)

func NextDate(w http.ResponseWriter, r *http.Request) {
	now, err := time.Parse(dateLayout, r.FormValue("now"))
	if err != nil {
		http.Error(w, "'now' ошибка формата даты", http.StatusBadRequest)
		return
	}

	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	nextDate, err := service.NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)

	_, err = w.Write([]byte(nextDate))
	if err != nil {
		http.Error(w, fmt.Errorf("ошибка записи в data: %w", err).Error(), http.StatusBadRequest)
	}

}
