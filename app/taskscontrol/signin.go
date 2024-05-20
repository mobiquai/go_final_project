package taskscontrol

import (
	"bytes"
	"encoding/json"
	"fmt"

	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
)

func Sign(w http.ResponseWriter, r *http.Request) {
	var signData SignIn
	var buffer bytes.Buffer

	if _, err := buffer.ReadFrom(r.Body); err != nil {
		responseWithError(w, "ошибка получения тела запроса", err)
		return
	}

	if err := json.Unmarshal(buffer.Bytes(), &signData); err != nil {
		responseWithError(w, "json encoding error", err)
		return
	}

	envPassword := os.Getenv("TODO_PASSWORD")

	if signData.Password == envPassword {
		jwtToken := jwt.New(jwt.SigningMethodHS256)
		token, err := jwtToken.SignedString([]byte(envPassword))
		if err != nil {
			http.Error(w, fmt.Errorf("failed to sign jwt: %s", err).Error(), http.StatusUnauthorized)
		}

		taskIdData, err := json.Marshal(AuthToken{Token: token})
		if err != nil {
			responseWithError(w, "json decoding error", err)
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)

		_, err = w.Write(taskIdData)
		if err != nil {
			http.Error(w, fmt.Errorf("error: %w", err).Error(), http.StatusUnauthorized)
		}

	} else {
		errorResponse := ErrorResponse{Error: "пароль неверный"}
		errorData, _ := json.Marshal(errorResponse)

		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write(errorData)
		if err != nil {
			http.Error(w, fmt.Errorf("error: %w", err).Error(), http.StatusUnauthorized)
		}

	}
}
