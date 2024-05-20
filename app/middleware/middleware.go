package middleware

import (
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
)

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// смотрим наличие пароля в переменной окружения
		pass := os.Getenv("TODO_PASSWORD")

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
				http.Error(w, fmt.Errorf("failed to sign jwt: %s", err).Error(), http.StatusUnauthorized)
			}

			if jwtTokenString != token {
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}

		}

		next(w, r)
	}
}
