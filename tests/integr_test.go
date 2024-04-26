package tests

import (
	"avito/internal/model"
	"context"
)

func (suite *BannerSuite) TestGetBanner() {
	suite.NotNil(suite.service)

	testCases := []struct {
		name    string
		banner  model.BannerCreate
		wantErr error
	}{
		{
			name: "success",
			banner: model.BannerCreate{
				FeatureId: 1,
				Tags:      []int{1, 2},
				Content: model.BannerContent{
					Title: "Title1",
					URL:   "http://test",
				},
				IsActive: true,
			},
			wantErr: nil,
		},
	}

	ctx := context.Background()

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			id, err := suite.service.CreateBanner(ctx, tc.banner)
			suite.NoError(err)
			suite.NotNil(id)

			tc.banner.BannerId = id

			banners, err := suite.service.GetBanners(ctx, model.BannersFilter{FeatureId: 1})
			suite.Equal(tc.wantErr, err)
			if err == nil {
				suite.Equal([]model.BannerCreate{tc.banner}, banners)
			}
		})
	}
}
