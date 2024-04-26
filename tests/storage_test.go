package tests

import (
	"avito/internal/logger"
	"avito/internal/model"
	"avito/internal/storage"
	"context"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func BenchmarkCreateBanner(b *testing.B) {

	log := logger.InitLogger(os.Getenv("LOGGER_LEVEL"))
	ctx := context.Background()
	godotenv.Load()

	s, err := storage.NewStorage(ctx, os.Getenv("STORAGE_PATH"), log)
	if err != nil {
		log.Error(err.Error())
		return
	}

	banner := model.BannerCreate{
		FeatureId: model.RandomTestInt(),
		Tags:      []int{model.RandomTestInt(), model.RandomTestInt()},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		id, _ := s.CreateBanner(ctx, banner)
		b.StopTimer()
		s.DeleteBanner(ctx, id)
		b.StartTimer()
	}
}
