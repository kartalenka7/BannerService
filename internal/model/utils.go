package model

import (
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v4"
)

func ParseQuery(r *http.Request) (BannersFilter, error) {
	var bannerFilters BannersFilter
	var err error

	q := r.URL.Query()

	featureId := q.Get("feature_id")
	if featureId != `` {
		bannerFilters.FeatureId, err = strconv.Atoi(featureId)
		if err != nil {
			return bannerFilters, err
		}
	}

	tagId := q.Get("tag_id")
	if tagId != `` {
		bannerFilters.TagId, err = strconv.Atoi(tagId)
		if err != nil {
			return bannerFilters, err
		}
	}

	limit := q.Get("limit")
	if limit != `` {
		bannerFilters.Limit, err = strconv.Atoi(limit)
		if err != nil {
			return bannerFilters, err
		}
	}

	offset := q.Get("offset")
	if offset != `` {
		bannerFilters.Offset, err = strconv.Atoi(offset)
		if err != nil {
			return bannerFilters, err
		}
	}

	return bannerFilters, nil
}

func GetToken(secret []byte) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	tokenString, err := token.SignedString(secret)
	return tokenString, err
}
