package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

func NewRouter(service ServiceInterface, log *logrus.Logger) chi.Router {
	router := chi.NewRouter()

	server := &Server{
		service: service,
		log:     log,
	}

	router.Route("/banner", func(r chi.Router) {
		r.Use(server.checkAdminAuth)
		r.Get("/", server.handlerGetBanners)
		r.Post("/", server.handlerCreateBanner)
		r.Patch("/{id}", server.handlerUpdateBanner)
		r.Delete("/{id}", server.handlerDeleteBanner)
	})

	router.Group(func(r chi.Router) {
		r.Use(server.checkUserAuth)
		r.Get("/user_banner", server.handlerGetUserBanner)
	})

	router.Post("/auth", server.handlerAuthentification)

	return router
}
