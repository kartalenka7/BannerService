package tests

import (
	"avito/internal/model"
	"avito/internal/service"
	mock "avito/mocks/storage"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_GetBanner(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
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
		storage        *mock.MockStorer
		args           args
		expectedErrors errors
	}{
		{
			name:    "Success",
			storage: mock.NewMockStorer(ctrl),
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
			tt.storage.EXPECT().CreateBanner(ctx,
				tt.args.bannerCreated).Return(tt.args.bannerId, tt.expectedErrors.createError)
			tt.storage.EXPECT().GetBanners(ctx,
				tt.args.bannersFilter).Return([]model.BannerCreate{tt.args.bannerCreated}, tt.expectedErrors.getError)

			s := service.NewService(ctx, tt.storage)
			id, err := s.CreateBanner(ctx, tt.args.bannerCreated)
			assert.Equal(t, id, tt.args.bannerId)
			assert.Equal(t, err, tt.expectedErrors.createError)
			if err == nil {
				banners, err := s.GetBanners(ctx, tt.args.bannersFilter)
				assert.Equal(t, err, tt.expectedErrors.getError)
				assert.Equal(t, banners, []model.BannerCreate{tt.args.bannerCreated})
			}
		})
	}
}
