package service

import (
	"avito/internal/model"
	"context"
)

type Storer interface {
	CreateBanner(ctx context.Context, banner model.BannerCreate) (int, error)
	GetBanners(ctx context.Context, bannersFilters model.BannersFilter) ([]model.BannerCreate, error)
	UpdateBanner(ctx context.Context, banner model.BannerUpdateRequest) (model.BannerCreate, error)
	DeleteBanner()
}

type Service struct {
	storage Storer
}

func NewService(ctx context.Context, storage Storer) *Service {
	return &Service{
		storage: storage,
	}
}

func (s Service) CreateBanner(ctx context.Context, banner model.BannerCreate) (int, error) {
	return s.storage.CreateBanner(ctx, banner)
}
func (s Service) GetBanners(ctx context.Context,
	bannersFilters model.BannersFilter) ([]model.BannerCreate, error) {

	banners, err := s.storage.GetBanners(ctx, bannersFilters)
	if err != nil {
		return nil, err
	}
	if bannersFilters.Offset > 0 && bannersFilters.Offset < len(banners) {
		banners = banners[bannersFilters.Offset+1:]
	}
	if bannersFilters.Limit > 0 && bannersFilters.Limit < len(banners) {
		banners = banners[:bannersFilters.Limit]
	}
	return banners, nil
}

func (s Service) GetUserBanner(ctx context.Context, bannersFilters model.BannersFilter) (
	[]model.BannerCreate, error) {
	banner, err := s.storage.GetBanners(ctx, bannersFilters)
	if err != nil {
		return nil, err
	}

	// пользователи не могут получать выключенные баннеры
	var userBanner []model.BannerCreate
	for _, b := range banner {
		if b.IsActive {
			userBanner = append(userBanner, b)
		}
	}
	return userBanner, nil
}

func (s Service) UpdateBanner(ctx context.Context, banner model.BannerUpdateRequest) (model.BannerCreate, error) {
	return s.storage.UpdateBanner(ctx, banner)
}
func (s Service) DeleteBanner() {

}
