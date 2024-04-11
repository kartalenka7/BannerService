package main

import (
	"avito/internal/config"
	server "avito/internal/http-server"
	"avito/internal/logger"
	"avito/internal/service"
	"avito/internal/storage"
	"context"
	"fmt"
	"net/http"
)

func main() {
	// TODO: init logger
	log := logger.InitLogger()

	// TODO: init config
	cfg, err := config.InitConfig()
	if err != nil {
		log.Error(err.Error())
		return
	}

	fmt.Println(cfg)

	// TODO: init storage
	ctx := context.Background()
	storage, err := storage.NewStorage(ctx, cfg.StoragePath, log)
	if err != nil {
		log.Error(err.Error())
		return
	}
	defer storage.Close()

	service := service.NewService(ctx, storage)

	// TODO: init router
	log.Info("Инициализируем роутер")
	r := server.NewRouter(service, log, cfg)

	// TODO: run server
	srv := &http.Server{
		Addr:         cfg.ServerAddr,
		ReadTimeout:  cfg.ServerTimeout,
		WriteTimeout: cfg.ServerTimeout,
		Handler:      r,
	}
	log.Info("Запускаем сервер")
	if err := srv.ListenAndServe(); err != nil {
		log.Error(err.Error())
	}
}
