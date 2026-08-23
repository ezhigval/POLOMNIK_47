package dto

import (
	"time"

	"palomnik/internal/domain"
)

type CreateUserPhotoRequest struct {
	URL                 string `json:"url"`
	Caption             string `json:"caption"`
	AllowDistribution   bool   `json:"allow_distribution"`
	ConsentPersonalData bool   `json:"consent_personal_data"`
	Website             string `json:"website"`
}

type UserPhotoResponse struct {
	ID                string    `json:"id"`
	URL               string    `json:"url"`
	Caption           string    `json:"caption"`
	AllowDistribution bool      `json:"allow_distribution"`
	CreatedAt         time.Time `json:"created_at"`
}

func ToUserPhotoResponse(photo domain.UserPhoto) UserPhotoResponse {
	return UserPhotoResponse{
		ID:                photo.ID.String(),
		URL:               photo.URL,
		Caption:           photo.Caption,
		AllowDistribution: photo.AllowDistribution,
		CreatedAt:         photo.CreatedAt,
	}
}

func ToUserPhotoResponses(items []domain.UserPhoto) []UserPhotoResponse {
	out := make([]UserPhotoResponse, 0, len(items))
	for _, item := range items {
		out = append(out, ToUserPhotoResponse(item))
	}
	return out
}
