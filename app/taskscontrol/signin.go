package taskscontrol

import (
	"bytes"
	"encoding/json"
	"fmt"

	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/mobiquai/go_final_project/app/appsettings"
)

func Sign(w http.ResponseWriter, r *http.Request) {
	var signData SignIn
	var buffer bytes.Buffer

	if _, err := buffer.ReadFrom(r.Body); err != nil {
		responseWithError(w, "ошибка получения тела запроса", err, http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(buffer.Bytes(), &signData); err != nil {
		responseWithError(w, "json encoding error", err, http.StatusBadRequest)
		return
	}

	envPassword := appsettings.EnvPassword // получаем значение переменной окружения

	if signData.Password == envPassword {
		jwtToken := jwt.New(jwt.SigningMethodHS256)
		token, err := jwtToken.SignedString([]byte(envPassword))
		if err != nil {
			http.Error(w, fmt.Errorf("failed to sign jwt: %s", err).Error(), http.StatusUnauthorized)
		}

		taskIdData, err := json.Marshal(AuthToken{Token: token})
		if err != nil {
			responseWithError(w, "json decoding error", err, http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)

		_, err = w.Write(taskIdData)
		if err != nil {
			http.Error(w, fmt.Errorf("error: %w", err).Error(), http.StatusUnauthorized)
		}

		return

	}

	errorResponse := ErrorResponse{Error: "пароль неверный"}
	errorData, _ := json.Marshal(errorResponse)

	w.WriteHeader(http.StatusUnauthorized)
	_, err := w.Write(errorData)
	if err != nil {
		http.Error(w, fmt.Errorf("error: %w", err).Error(), http.StatusUnauthorized)
	}

}
