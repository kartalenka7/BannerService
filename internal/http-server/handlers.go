package server

import (
	"avito/internal/model"
	"context"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type ServiceInterface interface {
	CreateBanner(ctx context.Context, banner model.BannerCreate) (int, error)
	GetBanners(ctx context.Context, bannersFilters model.BannersFilter) ([]model.BannerCreate, error)
	GetUserBanner(ctx context.Context, bannersFilters model.BannersFilter) ([]model.BannerCreate, error)
	UpdateBanner(ctx context.Context, banner model.BannerUpdateRequest) (model.BannerCreate, error)
	DeleteBanner(ctx context.Context, bannerId int) error
}

type Server struct {
	service ServiceInterface
	log     *logrus.Logger
}

func ErrorResponse(err error, rw http.ResponseWriter, r *http.Request, status int) {
	response := model.Response{
		Error: err.Error(),
	}
	render.Status(r, status)
	render.JSON(rw, r, response)

}

func (server Server) handlerAuthentification(rw http.ResponseWriter, r *http.Request) {
	var admin model.AdminAuth

	server.log.Info("Аутентификация")
	if err := render.DecodeJSON(r.Body, &admin); err != nil {
		ErrorResponse(err, rw, r, http.StatusBadRequest)
		return
	}

	var secret []byte

	password := os.Getenv("HTTP_SERVER_PASSWORD")
	if admin.Login == os.Getenv("HTTP_SERVER_LOGIN") &&
		admin.Password == password {
		secret = []byte(password)
	}
	token, err := model.GetToken(secret)
	if err != nil {
		ErrorResponse(err, rw, r, http.StatusInternalServerError)
		return
	}
	rw.Header().Add("Authorization", token)
	rw.WriteHeader(http.StatusOK)
}

func (server Server) handlerGetUserBanner(rw http.ResponseWriter, r *http.Request) {
	server.log.Info("Получаем баннер пользователя")
	bannerFilters, err := model.ParseQuery(r)
	if err != nil {
		ErrorResponse(err, rw, r, http.StatusBadRequest)
		return
	}

	bannerRequest := model.UserBannerRequest{
		FeatureId:       bannerFilters.FeatureId,
		TagId:           bannerFilters.TagId,
		UseLastRevision: bannerFilters.UseLastRevision,
	}

	if err := validator.New().Struct(bannerRequest); err != nil {
		validateErr := err.(validator.ValidationErrors)
		render.Status(r, http.StatusBadRequest)
		render.JSON(rw, r, model.ValidationError(validateErr))
		return
	}

	bannersResponse, err := server.service.GetUserBanner(r.Context(), bannerFilters)
	if err != nil {
		ErrorResponse(err, rw, r, http.StatusInternalServerError)
		return
	}
	if len(bannersResponse) == 0 {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	render.JSON(rw, r, bannersResponse)
}

func (server Server) handlerGetBanners(rw http.ResponseWriter, r *http.Request) {
	server.log.Info("Получаем все баннеры")
	bannerFilters, err := model.ParseQuery(r)
	if err != nil {
		ErrorResponse(err, rw, r, http.StatusBadRequest)
		return
	}

	bannersResponse, err := server.service.GetBanners(r.Context(), bannerFilters)
	if err != nil {
		ErrorResponse(err, rw, r, http.StatusInternalServerError)
		return
	}
	if len(bannersResponse) == 0 {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	render.JSON(rw, r, bannersResponse)
}

func (server Server) handlerCreateBanner(rw http.ResponseWriter, r *http.Request) {
	server.log.Info("Создаем баннер")
	var banner model.BannerCreate
	var err error
	if err = render.DecodeJSON(r.Body, &banner); err != nil {
		ErrorResponse(err, rw, r, http.StatusBadRequest)
		return
	}

	if err := validator.New().Struct(banner); err != nil {
		validateErr := err.(validator.ValidationErrors)
		render.Status(r, http.StatusBadRequest)
		render.JSON(rw, r, model.ValidationError(validateErr))
		return
	}

	var bannerResp model.BannerCreatedResp
	bannerResp.BannerId, err = server.service.CreateBanner(r.Context(), banner)
	if err != nil {
		ErrorResponse(err, rw, r, http.StatusInternalServerError)
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(rw, r, bannerResp)
}

func (server Server) handlerUpdateBanner(rw http.ResponseWriter, r *http.Request) {
	server.log.Info("Обновление баннера")

	// TODO добавить валидации на запросы
	var bannerUpd model.BannerUpdateRequest
	var err error
	if err = render.DecodeJSON(r.Body, &bannerUpd); err != nil {
		ErrorResponse(err, rw, r, http.StatusBadRequest)
		return
	}

	bannerUpd.BannerId, err = strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		ErrorResponse(err, rw, r, http.StatusBadRequest)
		return
	}

	var banner model.BannerCreate
	if banner, err = server.service.UpdateBanner(r.Context(), bannerUpd); err != nil {
		ErrorResponse(err, rw, r, http.StatusInternalServerError)
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(rw, r, banner)
}

func (server Server) handlerDeleteBanner(rw http.ResponseWriter, r *http.Request) {

	bannerId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		ErrorResponse(err, rw, r, http.StatusBadRequest)
		return
	}

	err = server.service.DeleteBanner(r.Context(), bannerId)
	if err != nil {
		ErrorResponse(err, rw, r, http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
