package config

import (
	"time"

	"github.com/caarlos0/env"
)

type Config struct {
	ServerAddr    string        `env:"RUN_ADDR" envDefault:"localhost:3000"`
	ServerTimeout time.Duration `env:"SERVER_TIMEOUT" envDefault:"100ms"`
	StoragePath   string        `env:"STORAGE_PATH" envDefault:"postgres://habruser:habr@localhost:5432/habrdb?sslmode=disable"`
	// StoragePath string `env:"STORAGE_PATH" envDefault:"postgres://avitouser:avitopass@localhost:5432/banner_db?sslmode=disable"`
	AdminLogin    string `env:"HTTP_SERVER_LOGIN" envDefault:"admin"`
	AdminPassword string `env:"HTTP_SERVER_PASSWORD" envDefault:"admin"`
	RedisDSN      string `env:"REDIS_DSN" envDefault:"localhost:6379"`
	RedisPassword string `env:"REDIS_PASSWORD"`
	RedisDB       int
}

func InitConfig() (Config, error) {

	var cfg Config

	err := env.Parse(&cfg)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}
