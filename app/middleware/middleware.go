package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/mobiquai/go_final_project/app/appsettings"
)

func responseWithErrorUnauthorized(w http.ResponseWriter, errorText string, err error) {
	errorResponse := ErrorResponse{fmt.Errorf("%s: %w", errorText, err).Error()}

	errorData, _ := json.Marshal(errorResponse)
	w.WriteHeader(http.StatusUnauthorized)

	_, err = w.Write(errorData)
	if err != nil {
		http.Error(w, fmt.Errorf("error: %w", err).Error(), http.StatusUnauthorized)
	}

}

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pass := appsettings.EnvPassword // смотрим наличие пароля в переменной окружения

		if len(pass) > 0 {
			var jwtTokenString string // JWT-токен из куки

			// получаем куку
			cookie, err := r.Cookie("token")
			if err == nil {
				jwtTokenString = cookie.Value
			}

			jwtToken := jwt.New(jwt.SigningMethodHS256)
			token, err := jwtToken.SignedString([]byte(pass))
			if err != nil {
				responseWithErrorUnauthorized(w, "не удалось провести аутентификацию", fmt.Errorf("failed to sign jwt: %s", err))
				return
			}

			if jwtTokenString != token {
				responseWithErrorUnauthorized(w, "требуется аутентификация", errors.New("authentication required"))
				return
			}

		}

		next(w, r)
	}
}
