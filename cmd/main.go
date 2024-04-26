package main

import (
	cache "avito/internal/cacheredis"
	server "avito/internal/http-server"
	"avito/internal/logger"
	"avito/internal/service"
	"avito/internal/storage"
	"context"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func main() {
	godotenv.Load()

	log := logger.InitLogger(os.Getenv("LOGGER_LEVEL"))

	ctx := context.Background()
	storage, err := storage.NewStorage(ctx, os.Getenv("STORAGE_PATH"), log)
	if err != nil {
		log.Error(err.Error())
		return
	}
	defer storage.Close()

	clientRedis := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_DSN"),
	})

	cache := cache.NewRedis(clientRedis, log)
	defer cache.Close()
	log.Info("Redis запущен")
	status := clientRedis.Ping(ctx)
	if status.Err() != nil {
		log.Error(status.Err().Error())
		return
	}
	log.Info("Redis пингуется")

	service := service.NewService(ctx, storage, cache)

	log.Info("Инициализируем роутер")
	r := server.NewRouter(service, log)

	srv := &http.Server{
		Addr: os.Getenv("RUN_ADDR"),
		// ReadTimeout:  cfg.ServerTimeout,
		// WriteTimeout: cfg.ServerTimeout,
		Handler: r,
	}
	log.Info("Запускаем сервер")
	if err := srv.ListenAndServe(); err != nil {
		log.Error(err.Error())
	}
}
