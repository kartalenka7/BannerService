package model

type BannersFilter struct {
	FeatureId int
	TagId     int
	Limit     int
	Offset    int
}

type BannerUpdateRequest struct {
	BannerId  int
	Tags      []int `json:"tag_ids"`
	FeatureId int   `json:"feature_id"`
	Content   struct {
		Title string `json:"title"`
		Text  string `json:"text"`
		URL   string `json:"url"`
	} `json:"content"`
	IsActive interface{} `json:"is_active"`
}

type BannerCreate struct {
	BannerId  int
	Tags      []int `json:"tag_ids"`
	FeatureId int   `json:"feature_id"`
	Content   struct {
		Title string `json:"title"`
		Text  string `json:"text"`
		URL   string `json:"url"`
	} `json:"content"`
	IsActive bool `json:"is_active"`
}

type BannerCreatedResp struct {
	BannerId int `json:"banner_id"`
}

type AdminAuth struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (AdminAuth) Valid() {

}
