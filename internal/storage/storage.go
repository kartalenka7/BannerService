package storage

import (
	"avito/internal/model"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

var (
	// DROP TABLE IF EXISTS tags;
	// DROP TABLE IF EXISTS banners;
	createBannerTable = `
	    DROP TABLE IF EXISTS tags;
	    DROP TABLE IF EXISTS banners;
		CREATE TABLE IF NOT EXISTS banners(
		banner_id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
		tags_group_id INTEGER UNIQUE GENERATED ALWAYS AS IDENTITY,
		feature_id INTEGER,
		title VARCHAR(128),
		text VARCHAR(512),
		url VARCHAR(256),
		is_active BOOLEAN NOT NULL DEFAULT FALSE
	);`

	createTagsTable = `
		CREATE TABLE IF NOT EXISTS tags(
		tag_id INTEGER ,
		tags_group_id INTEGER,
		banner_id INTEGER,
		CONSTRAINT fk_tags_group_id FOREIGN KEY (tags_group_id) 
	    REFERENCES banners(tags_group_id),
		PRIMARY KEY (tag_id, tags_group_id)
	);`

	createBanner = `INSERT INTO banners(feature_id, title, text, url, is_active) 
					VALUES ($1, $2, $3, $4, $5)
					RETURNING banner_id, tags_group_id`

	createTagGroup = `INSERT INTO tags(tag_id, tags_group_id, banner_id)
					  VALUES ($1, $2, $3)`

	GetBanners = `SELECT DISTINCT b.banner_id, b.feature_id, b.tags_group_id, b.title, b.text, b.url, b.is_active
				  FROM banners b
				  JOIN tags t ON b.tags_group_id = t.tags_group_id
				  WHERE ($1 = 0 OR b.feature_id = $1)
				  	AND ($2 = 0 OR t.tag_id = $2)`
	GetTags = `SELECT tag_id
			   FROM tags
			   WHERE tags_group_id = $1`

	updateBanner = `
		UPDATE banners
		SET feature_id = COALESCE(NULLIF($1, 0), feature_id),
			title = CASE WHEN $2 <> '' THEN $2 ELSE title END,
			text = CASE WHEN $3 <> '' THEN $3 ELSE text END,
			url = CASE WHEN $4 <> '' THEN $4 ELSE url END,
			is_active = $5
		WHERE banner_id = $6
		RETURNING tags_group_id, feature_id, title, text, url, is_active;`

	updateBannerBool = `
		UPDATE banners
		SET feature_id = COALESCE(NULLIF($1, 0), feature_id),
			title = CASE WHEN $2 <> '' THEN $2 ELSE title END,
			text = CASE WHEN $3 <> '' THEN $3 ELSE text END,
			url = CASE WHEN $4 <> '' THEN $4 ELSE url END
		WHERE banner_id = $5
		RETURNING tags_group_id, feature_id, title, text, url, is_active;`
)

type Storage struct {
	pgxPool *pgxpool.Pool
	log     *logrus.Logger
}

func NewStorage(ctx context.Context, connString string,
	log *logrus.Logger) (*Storage, error) {
	//TODO: add logging
	log.Info("Запускаем инициализацию хранилища")
	pgxPool, err := pgxpool.New(ctx, connString)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	if _, err = pgxPool.Exec(ctx, createBannerTable); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	if _, err = pgxPool.Exec(ctx, createTagsTable); err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return &Storage{
		pgxPool: pgxPool,
		log:     log,
	}, nil
}

func (s *Storage) CreateBanner(ctx context.Context, banner model.BannerCreate) (int, error) {
	var bannerId int
	var tagsGroupId int

	row := s.pgxPool.QueryRow(ctx, createBanner, banner.FeatureId,
		banner.Content.Title, banner.Content.Text, banner.Content.URL,
		banner.IsActive)
	err := row.Scan(&bannerId, &tagsGroupId)
	if err != nil {
		s.log.Error(err.Error())
		return 0, err
	}

	for _, tag := range banner.Tags {
		_, err := s.pgxPool.Exec(ctx, createTagGroup, tag, tagsGroupId, bannerId)
		if err != nil {
			s.log.Error(err.Error())
			return 0, err
		}
	}

	return bannerId, nil
}

func (s *Storage) GetBanners(ctx context.Context, bannersFilters model.BannersFilter) (
	[]model.BannerCreate, error) {

	rows, err := s.pgxPool.Query(ctx, GetBanners, bannersFilters.FeatureId,
		bannersFilters.TagId)
	defer rows.Close()
	if err != nil {
		s.log.Error(err.Error())
		return nil, err
	}
	var banners []model.BannerCreate
	for rows.Next() {
		var banner model.BannerCreate
		var tagsGroupId int
		err = rows.Scan(&banner.BannerId, &banner.FeatureId, &tagsGroupId,
			&banner.Content.Title, &banner.Content.Text, &banner.Content.URL,
			&banner.IsActive)
		if err != nil {
			s.log.Error(err.Error())
			return nil, err
		}

		tagsRows, err := s.pgxPool.Query(ctx, GetTags, &tagsGroupId)
		defer tagsRows.Close()
		if err != nil {
			s.log.Error(err.Error())
			return nil, err
		}
		for tagsRows.Next() {
			var tagId int
			err = tagsRows.Scan(&tagId)
			if err != nil {
				s.log.Error(err.Error())
				return nil, err
			}
			banner.Tags = append(banner.Tags, tagId)
		}
		banners = append(banners, banner)
	}
	if bannersFilters.Offset > 0 && bannersFilters.Offset < len(banners) {
		banners = banners[bannersFilters.Offset+1:]
	}
	if bannersFilters.Limit > 0 && bannersFilters.Limit < len(banners) {
		banners = banners[:bannersFilters.Limit]
	}
	return banners, nil
}

func (s *Storage) UpdateBanner(ctx context.Context, banner model.BannerUpdateRequest) (
	model.BannerCreate, error) {

	var row pgx.Row

	if banner.IsActive == "true" || banner.IsActive == "false" {
		IsActive, ok := banner.IsActive.(bool)
		if ok {
			row = s.pgxPool.QueryRow(ctx, updateBanner, banner.FeatureId,
				banner.Content.Title, banner.Content.Text, banner.Content.URL, IsActive,
				banner.BannerId)
		}
	} else {
		row = s.pgxPool.QueryRow(ctx, updateBannerBool, banner.FeatureId,
			banner.Content.Title, banner.Content.Text, banner.Content.URL, banner.BannerId)
	}

	var updatedBanner model.BannerCreate
	var tagsGroupId int

	err := row.Scan(&tagsGroupId, &updatedBanner.FeatureId, &updatedBanner.Content.Title,
		&updatedBanner.Content.Text, &updatedBanner.Content.URL, &updatedBanner.IsActive)
	if err != nil {
		s.log.Error(err.Error())
	}
	return updatedBanner, err
}

func (s *Storage) DeleteBanner() {

}
