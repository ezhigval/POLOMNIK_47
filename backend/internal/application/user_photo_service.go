package application

import (
	"context"
	"time"

	"github.com/google/uuid"

	"palomnik/internal/domain"
	"palomnik/internal/ports"
)

type UserPhotoService struct {
	photos ports.UserPhotoRepository
}

func NewUserPhotoService(photos ports.UserPhotoRepository) *UserPhotoService {
	return &UserPhotoService{photos: photos}
}

type UserPhotoInput struct {
	URL               string
	Caption           string
	AllowDistribution bool
}

func (s *UserPhotoService) List(ctx context.Context, userID uuid.UUID) ([]domain.UserPhoto, error) {
	if s == nil || s.photos == nil {
		return nil, domain.ErrNotFound
	}
	if userID == uuid.Nil {
		return nil, domain.ErrInvalidID
	}
	return s.photos.ListUserPhotos(ctx, userID)
}

func (s *UserPhotoService) Create(ctx context.Context, userID uuid.UUID, input UserPhotoInput) (domain.UserPhoto, error) {
	if s == nil || s.photos == nil {
		return domain.UserPhoto{}, domain.ErrNotFound
	}
	photo, err := domain.NewUserPhoto(domain.NewUserPhotoInput{
		ID:                uuid.New(),
		UserID:            userID,
		URL:               input.URL,
		Caption:           input.Caption,
		AllowDistribution: input.AllowDistribution,
		Now:               time.Now().UTC(),
	})
	if err != nil {
		return domain.UserPhoto{}, err
	}
	return s.photos.CreateUserPhoto(ctx, photo)
}

func (s *UserPhotoService) Delete(ctx context.Context, userID, id uuid.UUID) error {
	if s == nil || s.photos == nil {
		return domain.ErrNotFound
	}
	return s.photos.DeleteUserPhoto(ctx, userID, id)
}
