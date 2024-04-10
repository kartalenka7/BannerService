package server

import (
	"net/http"

	"github.com/golang-jwt/jwt/v4"
)

func (s Server) checkUserAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.log.Info("Проверяем токен админа")

		token := r.Header.Get("token")
		if token == `` {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		jwtToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
			// Возвращаем секрет для проверки подписи
			return []byte(s.config.AdminPassword), nil
		})

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if !jwtToken.Valid {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
