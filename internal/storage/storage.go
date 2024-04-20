package storage

import (
	"avito/internal/model"
	"context"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

var (
	createBannerTable = `
		CREATE TABLE IF NOT EXISTS banners(
		banner_id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
		title VARCHAR(128),
		text VARCHAR(512),
		url VARCHAR(256),
		is_active BOOLEAN NOT NULL DEFAULT FALSE
	);`

	createGroupTable = `
		CREATE TABLE IF NOT EXISTS groupTable(
		tag_id INTEGER,
		feature_id INTEGER,
		banner_id INTEGER,
		PRIMARY KEY (tag_id, feature_id),
		CONSTRAINT fk_banner_id FOREIGN KEY (banner_id) REFERENCES banners(banner_id));`

	updateGroupTable = `
	UPDATE groupTable
		SET tag_id = $1,
		    feature_id = CASE WHEN $3 != '' THEN $3
			                ELSE feature_id END
		WHERE banner_id = $2
		RETURNING tag_id, feature_id`

	updateGroupFeature = `
	UPDATE groupTable
		SET feature_id = $1
		WHERE banner_id = $2
		RETURNING tag_id, feature_id`

	deleteGroup = `DELETE FROM groupTable WHERE banner_id = $1 AND tag_id NOT IN $2;`

	delete = `
	DELETE FROM groupTable WHERE banner_id = $1;
	DELETE FROM banner WHERE banner_id = $1;`

	createBanner = `INSERT INTO banners(title, text, url, is_active) 
					VALUES ($1, $2, $3, $4)
					RETURNING banner_id`

	createGroup = `INSERT INTO groupTable(tag_id, feature_id, banner_id)
					VALUES ($1, $2, $3)`

	GetBanners = `SELECT DISTINCT banner_id, title, text, url, is_active
				  FROM banners 
				  WHERE banner_id = $1`

	GetGroup = `SELECT banner_id, feature_id, tag_id
			   FROM groupTable
			   WHERE ($1 = 0 OR feature_id = $1)
			   AND ($2 = 0 OR tag_id = $2)`

	updateBanner = `
		UPDATE banners
		SET title = CASE WHEN $1 <> '' THEN $1 ELSE title END,
			text = CASE WHEN $2 <> '' THEN $2 ELSE text END,
			url = CASE WHEN $3 <> '' THEN $3 ELSE url END,
			is_active = $4
		WHERE banner_id = $5
		RETURNING banner_id, title, text, url, is_active;`

	updateBannerBool = `
		UPDATE banners
		SET title = CASE WHEN $1 <> '' THEN $1 ELSE title END,
			text = CASE WHEN $2 <> '' THEN $2 ELSE text END,
			url = CASE WHEN $3 <> '' THEN $3 ELSE url END
		WHERE banner_id = $4
		RETURNING  banner_id, title, text, url, is_active;`
)

type Storage struct {
	pgxPool *pgxpool.Pool
	log     *logrus.Logger
}

func NewStorage(ctx context.Context, connString string,
	log *logrus.Logger) (*Storage, error) {

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

	if _, err = pgxPool.Exec(ctx, createGroupTable); err != nil {
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
	s.log.Info("Запись баннера в бд")
	tx, err := s.pgxPool.Begin(ctx)
	if err != nil {
		s.log.Error(err.Error())
		return 0, err
	}

	row := tx.QueryRow(ctx, createBanner, banner.Content.Title,
		banner.Content.Text, banner.Content.URL, banner.IsActive)
	err = row.Scan(&bannerId)
	if err != nil {
		s.log.Error(err.Error())
		tx.Rollback(ctx)
		return 0, err
	}

	for _, tag := range banner.Tags {
		_, err := tx.Exec(ctx, createGroup, tag, banner.FeatureId, bannerId)
		if err != nil {
			s.log.Error(err.Error())
			tx.Rollback(ctx)
			return 0, err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		s.log.Error(err.Error())
		tx.Rollback(ctx)
		return 0, err
	}

	return bannerId, nil
}

func (s *Storage) GetBanners(ctx context.Context, bannersFilters model.BannersFilter) (
	[]model.BannerCreate, error) {

	s.log.Info("Получаем баннеры из бд")

	var group model.Group

	groupRows, err := s.pgxPool.Query(ctx, GetGroup, bannersFilters.FeatureId,
		bannersFilters.TagId)
	defer groupRows.Close()
	if err != nil {
		s.log.Error(err.Error())
		return nil, err
	}

	var banners []model.BannerCreate
	var groups []model.Group
	mapId := make(map[string]int)

	for groupRows.Next() {
		err = groupRows.Scan(&group.BannerId, &group.FeatureId, &group.TagId)
		if err != nil {
			s.log.Error(err.Error())
			return nil, err
		}
		groups = append(groups, group)
		mapId[group.BannerId] = group.FeatureId
	}

	for bannerId, featureId := range mapId {

		row := s.pgxPool.QueryRow(ctx, GetBanners, bannerId)

		var banner model.BannerCreate
		err = row.Scan(&banner.BannerId, &banner.Content.Title, &banner.Content.Text,
			&banner.Content.URL, &banner.IsActive)
		if err != nil {
			s.log.Error(err.Error())
			return nil, err
		}

		for _, g := range groups {
			if g.BannerId == bannerId {
				banner.Tags = append(banner.Tags, g.TagId)
			}
		}

		banner.FeatureId = featureId
		if err != nil {
			s.log.Error(err.Error())
			return nil, err
		}
		banners = append(banners, banner)
	}

	return banners, nil
}
func (s *Storage) UpdateBanner(ctx context.Context, banner model.BannerUpdateRequest) (
	model.BannerCreate, error) {

	var row pgx.Row

	var updatedBanner model.BannerCreate
	for _, tag := range banner.Tags {
		row = s.pgxPool.QueryRow(ctx, updateGroupTable, tag, banner.BannerId, banner.FeatureId)

		var tag string
		var feature string

		err := row.Scan(&tag, &feature)
		if err != nil {
			s.log.Error(err.Error())
			return updatedBanner, err
		}
		tagId, _ := strconv.Atoi(tag)
		updatedBanner.Tags = append(updatedBanner.Tags, tagId)

		featureId, _ := strconv.Atoi(feature)
		updatedBanner.FeatureId = featureId
	}

	updatedBanner.BannerId = banner.BannerId
	if len(updatedBanner.Tags) != 0 {
		_, err := s.pgxPool.Exec(ctx, deleteGroup, updatedBanner.BannerId, updatedBanner.Tags)
		if err != nil {
			s.log.Error(err.Error())
			return updatedBanner, err
		}
	} else {
		_, err := s.pgxPool.Exec(ctx, updateGroupFeature, banner.BannerId, banner.FeatureId)
		if err != nil {
			s.log.Error(err.Error())
			return updatedBanner, err
		}
		updatedBanner.FeatureId = banner.FeatureId
	}

	if banner.IsActive == "true" || banner.IsActive == "false" {
		IsActive, ok := banner.IsActive.(bool)
		if ok {
			row = s.pgxPool.QueryRow(ctx, updateBanner,
				banner.Content.Title, banner.Content.Text, banner.Content.URL, IsActive,
				banner.BannerId)
		}
	} else {
		row = s.pgxPool.QueryRow(ctx, updateBannerBool,
			banner.Content.Title, banner.Content.Text, banner.Content.URL, banner.BannerId)
	}

	err := row.Scan(&updatedBanner.BannerId, &updatedBanner.Content.Title,
		&updatedBanner.Content.Text, &updatedBanner.Content.URL, &updatedBanner.IsActive)
	if err != nil {
		s.log.Error(err.Error())
	}
	return updatedBanner, err
}

func (s *Storage) DeleteBanner(ctx context.Context, bannerId int) error {
	s.log.Info("Удалить баннер из бд")
	_, err := s.pgxPool.Exec(ctx, delete, bannerId)
	if err != nil {
		s.log.Error(err.Error())
	}
	return err
}

func (s *Storage) Close() {
	s.pgxPool.Close()
}
