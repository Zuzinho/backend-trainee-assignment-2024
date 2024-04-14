package handler

type postBannerResponse struct {
	BannerID int64 `json:"banner_id"`
}

type postUserResponse struct {
	Token string `json:"token"`
}
