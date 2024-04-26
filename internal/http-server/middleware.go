package server

import (
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func (s Server) checkAdminAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.log.Info("Проверяем токен админа")

		token := r.Header.Get("token")
		if token == `` {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		jwtToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("HTTP_SERVER_PASSWORD")), nil
		})

		if err != nil {
			s.log.Error(err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if !jwtToken.Valid {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		claims, ok := jwtToken.Claims.(jwt.MapClaims)
		if !ok || !claims.VerifyExpiresAt(time.Now().Unix(), true) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s Server) checkUserAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.log.Info("Проверяем токен пользователя")

		token := r.Header.Get("token")
		if token == `` {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		jwtToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
			return []byte(""), nil
		})
		if err != nil {
			s.log.Error(err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims, ok := jwtToken.Claims.(jwt.MapClaims)
		if !ok || !claims.VerifyExpiresAt(time.Now().Unix(), true) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
