package tests

import (
	cacheredis "avito/internal/cacheredis"
	"avito/internal/logger"
	"avito/internal/service"
	"avito/internal/storage"
	"context"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
)

type BannerSuite struct {
	suite.Suite
	storage *storage.Storage
	cache   *cacheredis.Redis
	service *service.Service
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(BannerSuite))
}

func (suite *BannerSuite) SetupSuite() {
	// ожидание готовности контейнера с БД
	time.Sleep(2 * time.Second)
	ctx := context.Background()
	godotenv.Load()

	log := logger.InitLogger(os.Getenv("LOGGER_LEVEL"))

	if err := suite.isDBAvailable(os.Getenv("STORAGE_PATH")); err != nil {
		suite.T().Skipf("skip db tests: database is not available: %v", err)
	}

	suite.storage, _ = storage.NewStorage(ctx, os.Getenv("STORAGE_PATH"), log)
	clientRedis := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_DSN"),
	})
	suite.cache = cacheredis.NewRedis(clientRedis, log)
	suite.service = service.NewService(ctx, suite.storage, suite.cache)
}

func (suite *BannerSuite) TearDownSuite() {
	suite.cache.Close()
	suite.storage.Close()
}

func (suite *BannerSuite) isDBAvailable(dsn string) error {
	ctx := context.Background()
	pgxPool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return err
	}
	defer pgxPool.Close()
	return pgxPool.Ping(ctx)
}
