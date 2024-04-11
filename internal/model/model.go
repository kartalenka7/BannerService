package model

type BannersFilter struct {
	FeatureId       int
	TagId           int
	Limit           int
	Offset          int
	UseLastRevision bool
}

type UserBannerRequest struct {
	FeatureId       int `validate:"required"`
	TagId           int `validate:"required"`
	UseLastRevision bool
}
type BannerContent struct {
	Title string `json:"title"`
	Text  string `json:"text"`
	URL   string `json:"url"`
}
type BannerUpdateRequest struct {
	BannerId  int
	Tags      []int         `json:"tag_ids"`
	FeatureId int           `json:"feature_id"`
	Content   BannerContent `json:"content"`
	IsActive  interface{}   `json:"is_active"`
}

type BannerCreate struct {
	BannerId  int
	Tags      []int         `json:"tag_ids" validate:"required"`
	FeatureId int           `json:"feature_id" validate:"required"`
	Content   BannerContent `json:"content" validate:"required"`
	IsActive  bool          `json:"is_active"`
}

type BannerCreatedResp struct {
	BannerId int `json:"banner_id"`
}

type AdminAuth struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Response struct {
	Error string `json:"error"`
}

func (AdminAuth) Valid() {

}
