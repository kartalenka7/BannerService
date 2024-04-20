package tests

import (
	cache "avito/internal/cacheredis"
	"avito/internal/logger"
	"avito/internal/model"
	"avito/internal/service"
	"avito/internal/storage"
	"context"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func Test_GetBanner(t *testing.T) {
	ctx := context.Background()
	godotenv.Load()

	type args struct {
		bannerCreated model.BannerCreate
		bannersFilter model.BannersFilter
		bannerId      int
	}
	type errors struct {
		createError error
		getError    error
	}
	tests := []struct {
		name           string
		args           args
		log            *logrus.Logger
		expectedErrors errors
	}{
		{
			name: "Success",
			log:  logger.InitLogger(),
			args: args{
				bannerCreated: model.BannerCreate{
					FeatureId: 1,
					Tags:      []int{1, 2},
					Content: model.BannerContent{
						Title: "Title1",
						URL:   "http://test",
					},
					IsActive: true,
				},
				bannersFilter: model.BannersFilter{
					FeatureId: 1,
				},
				bannerId: 1,
			},
			expectedErrors: errors{
				createError: nil,
				getError:    nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			storage, err := storage.NewStorage(ctx, os.Getenv("STORAGE_PATH"), tt.log)
			assert.NoError(t, err)

			clientRedis := redis.NewClient(&redis.Options{
				Addr: os.Getenv("REDIS_DSN"),
			})
			cache := cache.NewRedis(clientRedis, tt.log)
			defer cache.Close()

			status := clientRedis.Ping(ctx)
			assert.NoError(t, status.Err())

			s := service.NewService(ctx, storage, cache)

			id, err := s.CreateBanner(ctx, tt.args.bannerCreated)
			assert.NotNil(t, id)
			assert.Equal(t, err, tt.expectedErrors.createError)
			if err == nil {
				tt.args.bannerCreated.BannerId = id
				banners, err := s.GetBanners(ctx, tt.args.bannersFilter)
				assert.Equal(t, err, tt.expectedErrors.getError)
				assert.Equal(t, banners, []model.BannerCreate{tt.args.bannerCreated})
			}
		})
	}
}
