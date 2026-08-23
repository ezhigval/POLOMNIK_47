package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type UserPhoto struct {
	ID                uuid.UUID
	UserID            uuid.UUID
	URL               string
	Caption           string
	AllowDistribution bool
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type NewUserPhotoInput struct {
	ID                uuid.UUID
	UserID            uuid.UUID
	URL               string
	Caption           string
	AllowDistribution bool
	Now               time.Time
}

func NewUserPhoto(input NewUserPhotoInput) (UserPhoto, error) {
	if input.ID == uuid.Nil || input.UserID == uuid.Nil {
		return UserPhoto{}, ErrInvalidID
	}
	url := strings.TrimSpace(input.URL)
	if url == "" {
		return UserPhoto{}, ErrInvalidPhotoURL
	}
	now := input.Now
	if now.IsZero() {
		now = time.Now().UTC()
	}
	return UserPhoto{
		ID:                input.ID,
		UserID:            input.UserID,
		URL:               url,
		Caption:           strings.TrimSpace(input.Caption),
		AllowDistribution: input.AllowDistribution,
		CreatedAt:         now,
		UpdatedAt:         now,
	}, nil
}
