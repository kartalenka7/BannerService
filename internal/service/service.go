package service

import (
	"avito/internal/model"
	"context"
	"fmt"
)

type Storer interface {
	CreateBanner(ctx context.Context, banner model.BannerCreate) (int, error)
	GetBanners(ctx context.Context, bannersFilters model.BannersFilter) ([]model.BannerCreate, error)
	UpdateBanner(ctx context.Context, banner model.BannerUpdateRequest) (model.BannerCreate, error)
	DeleteBanner(ctx context.Context, bannerId int) error
}

type Cache interface {
	CreateBannerCache(bannerId int, banner model.BannerCreate)
	GetBannerCache(ctx context.Context,
		bannersFilters model.BannersFilter) ([]model.BannerCreate, error)
	DeleteBanner(bannerId int)
}

type Service struct {
	storage Storer
	cache   Cache
}

func NewService(ctx context.Context, storage Storer, cache Cache) *Service {
	return &Service{
		storage: storage,
		cache:   cache,
	}
}

func (s Service) CreateBanner(ctx context.Context, banner model.BannerCreate) (int, error) {
	bannerId, err := s.storage.CreateBanner(ctx, banner)
	if err != nil {
		return 0, err
	}
	banner.BannerId = bannerId
	go s.cache.CreateBannerCache(bannerId, banner)
	return bannerId, nil
}

func (s Service) GetBanners(ctx context.Context,
	bannersFilters model.BannersFilter) ([]model.BannerCreate, error) {

	banners, _ := s.cache.GetBannerCache(ctx, bannersFilters)

	var err error
	if len(banners) == 0 {
		banners, err = s.storage.GetBanners(ctx, bannersFilters)
		if err != nil {
			return nil, err
		}
	}

	fmt.Println(bannersFilters.Limit, bannersFilters.Offset)
	if bannersFilters.Offset > 0 && bannersFilters.Offset < len(banners) {
		banners = banners[bannersFilters.Offset:]
	}
	if bannersFilters.Limit > 0 && bannersFilters.Limit < len(banners) {
		banners = banners[:bannersFilters.Limit]
	}
	return banners, nil
}

func (s Service) GetUserBanner(ctx context.Context, bannersFilters model.BannersFilter) (
	[]model.BannerCreate, error) {

	var banners []model.BannerCreate
	var err error

	if !bannersFilters.UseLastRevision {
		banners, err = s.cache.GetBannerCache(ctx, bannersFilters)
	} else {
		banners, err = s.storage.GetBanners(ctx, bannersFilters)
	}

	if err != nil {
		return nil, err
	}

	// пользователи не могут получать выключенные баннеры
	var userBanner []model.BannerCreate
	for _, b := range banners {
		if b.IsActive {
			userBanner = append(userBanner, b)
		}
	}
	return userBanner, nil
}

func (s Service) UpdateBanner(ctx context.Context, banner model.BannerUpdateRequest) (model.BannerCreate, error) {
	bannerUpd, err := s.storage.UpdateBanner(ctx, banner)
	if err != nil {
		return bannerUpd, err
	}

	go s.cache.CreateBannerCache(banner.BannerId, bannerUpd)
	return bannerUpd, nil
}
func (s Service) DeleteBanner(ctx context.Context, bannerId int) error {
	go s.cache.DeleteBanner(bannerId)
	return s.storage.DeleteBanner(ctx, bannerId)
}
