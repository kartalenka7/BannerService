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
	return s.storage.GetBanners(ctx, bannersFilters)
}
func (s Service) UpdateBanner(ctx context.Context, banner model.BannerUpdateRequest) (model.BannerCreate, error) {
	return s.storage.UpdateBanner(ctx, banner)
}
func (s Service) DeleteBanner() {

}
