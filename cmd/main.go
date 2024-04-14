package main

import (
	cache "avito/internal/cacheredis"
	"avito/internal/config"
	server "avito/internal/http-server"
	"avito/internal/logger"
	"avito/internal/service"
	"avito/internal/storage"
	"context"
	"fmt"
	"net/http"

	"github.com/redis/go-redis/v9"
)

func main() {
	log := logger.InitLogger()

	cfg, err := config.InitConfig()
	if err != nil {
		log.Error(err.Error())
		return
	}

	fmt.Println(cfg)

	ctx := context.Background()
	storage, err := storage.NewStorage(ctx, cfg.StoragePath, log)
	if err != nil {
		log.Error(err.Error())
		return
	}
	defer storage.Close()

	clientRedis := redis.NewClient(&redis.Options{
		Addr: cfg.RedisDSN,
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
	r := server.NewRouter(service, log, cfg)

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
