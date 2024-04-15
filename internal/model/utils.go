package model

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
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

	last := q.Get("use_last_revision")
	if last != `` {
		l, err := strconv.ParseBool(last)
		if err != nil {
			return bannerFilters, err
		}
		bannerFilters.UseLastRevision = l
	}

	return bannerFilters, nil
}

func GetToken(secret []byte) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Duration(time.Hour * 24)).Unix()

	tokenString, err := token.SignedString(secret)
	return tokenString, err
}

func ValidationError(errs validator.ValidationErrors) Response {
	var errMsgs []string
	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is requires", err.Field()))
		case "url":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid URL", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid", err.Field()))
		}
	}
	return Response{
		Error: strings.Join(errMsgs, ", "),
	}
}
