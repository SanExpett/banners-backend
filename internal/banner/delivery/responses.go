package delivery

import "github.com/SanExpett/banners-backend/pkg/models"

const (
	ResponseSuccessfulDeleteBanner = "Баннер успешно удален"
	ResponseSuccessfulUpdateBanner = "Баннер успешно обновлен"
)

type BannerResponse struct {
	Status int             `json:"status"`
	Body   *models.Content `json:"body"`
}

func NewBannerResponse(status int, body *models.Content) *BannerResponse {
	return &BannerResponse{
		Status: status,
		Body:   body,
	}
}

type BannerListResponse struct {
	Status int              `json:"status"`
	Body   []*models.Banner `json:"body"`
}

func NewBannerListResponse(status int, body []*models.Banner) *BannerListResponse {
	return &BannerListResponse{
		Status: status,
		Body:   body,
	}
}
