package cache

import (
	"avito/internal/model"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type Redis struct {
	client *redis.Client
	log    *logrus.Logger
}

func NewRedis(client *redis.Client, log *logrus.Logger) *Redis {
	return &Redis{client: client,
		log: log}
}

func (r *Redis) CreateBannerCache(bannerId int, banner model.BannerCreate) {
	r.log.WithFields(logrus.Fields{"bannerId": bannerId}).Info("Запись баннера в кэш")
	ctx := context.Background()
	bannerJSON, err := json.Marshal(banner)
	if err != nil {
		r.log.Error(err.Error())
		return
	}

	for _, tag := range banner.Tags {
		err = r.client.HSet(ctx, fmt.Sprintf("bytag_%d", tag),
			fmt.Sprint(banner.FeatureId), fmt.Sprint(bannerId)).Err()
		if err != nil {
			r.log.Error(err.Error())
		}

		err = r.client.HSet(ctx, fmt.Sprintf("byfeature_%d", banner.FeatureId),
			fmt.Sprint(tag), fmt.Sprint(bannerId)).Err()
		if err != nil {
			r.log.Error(err.Error())
		}
	}

	err = r.client.HSet(ctx, "mybanner", fmt.Sprint(bannerId), bannerJSON).Err()
	if err != nil {
		r.log.Error(err.Error())
	}

}

func (r *Redis) GetBannerCache(ctx context.Context,
	bannersFilters model.BannersFilter) ([]model.BannerCreate, error) {
	var bannerId string
	var err error
	var banners []model.BannerCreate

	var idList []string

	r.log.Info("Получаем баннеры из кэша")

	if bannersFilters.FeatureId != 0 && bannersFilters.TagId != 0 {
		bannerId, err = r.client.HGet(ctx, fmt.Sprintf("byfeature_%d", bannersFilters.FeatureId),
			fmt.Sprint(bannersFilters.TagId)).Result()
		if err != nil {
			r.log.Error(err.Error())
			return nil, err
		}

		idList = append(idList, bannerId)
	} else if bannersFilters.FeatureId != 0 {

		bannerIds, err := r.client.HGetAll(ctx,
			fmt.Sprintf("byfeature_%d", bannersFilters.FeatureId)).Result()
		if err != nil {
			r.log.Error(err.Error())
			return nil, err
		}

		for _, id := range bannerIds {
			idList = append(idList, id)
		}

	} else if bannersFilters.TagId != 0 {

		bannerIds, err := r.client.HGetAll(ctx, fmt.Sprintf("bytag_%d", bannersFilters.TagId)).Result()
		if err != nil {
			r.log.Error(err.Error())
			return nil, err
		}

		for _, id := range bannerIds {
			idList = append(idList, id)
		}
	} else {
		bannersJSON, err := r.client.HGetAll(ctx, "mybanner").Result()

		if err != nil {
			r.log.Error(err.Error())
			return nil, err
		}
		for _, bannerJSON := range bannersJSON {
			var bannerUnmarshalled model.BannerCreate
			if err = json.Unmarshal([]byte(bannerJSON), &bannerUnmarshalled); err != nil {
				r.log.Error(err.Error())
				return nil, err
			}
			banners = append(banners, bannerUnmarshalled)
		}
		if len(banners) == 0 {
			//TODO вынести ошибку в model
			return nil, errors.New("Не найдено")
		}
		return banners, nil
	}

	for _, b := range idList {
		var bannerUnmarshalled model.BannerCreate
		bannerJSON, err := r.client.HGet(ctx, "mybanner", b).Result()
		if err != nil {
			r.log.Error(err.Error())
			return nil, err
		}
		if err = json.Unmarshal([]byte(bannerJSON), &bannerUnmarshalled); err != nil {
			r.log.Error(err.Error())
			return nil, err
		}
		bannerUnmarshalled.BannerId, _ = strconv.Atoi(b)
		banners = append(banners, bannerUnmarshalled)
	}

	if len(banners) == 0 {
		//TODO вынести ошибку в model
		return nil, errors.New("Не найдено")
	}

	return banners, nil
}

func (r *Redis) DeleteBanner(bannerId int) {
	r.log.Info("Удалить баннер из кэша")
	result := r.client.HDel(context.Background(), "mybanner", fmt.Sprintf("%d", bannerId))
	if err := result.Err(); err != nil {
		r.log.Error(err.Error())
	}
}

func (r *Redis) Close() {
	r.client.Close()
}
