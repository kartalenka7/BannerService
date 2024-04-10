package server

import (
	"avito/internal/config"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

func NewRouter(service ServiceInterface, log *logrus.Logger, cfg config.Config) chi.Router {
	router := chi.NewRouter()

	server := &Server{
		service: service,
		log:     log,
		config:  cfg,
	}

	router.Route("/banner", func(r chi.Router) {
		r.Use(server.checkUserAuth)
		r.Get("/", server.handlerGetBanners)
		r.Post("/", server.handlerCreateBanner)
		r.Patch("/{id}", server.handlerUpdateBanner)
		r.Delete("/{id}", server.handlerDeleteBanner)
	})

	router.Group(func(r chi.Router) {
		r.Post("/auth", server.handlerAuthentification)
		r.Get("/user_banner", server.handlerGetUserBanner)
	})
	return router
}
