package memory

import (
	"context"
	"sort"

	"github.com/google/uuid"

	"palomnik/internal/domain"
)

func (s *Store) ListUserPhotos(_ context.Context, userID uuid.UUID) ([]domain.UserPhoto, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	items := make([]domain.UserPhoto, 0)
	for _, photo := range s.userPhotos {
		if photo.UserID == userID {
			items = append(items, photo)
		}
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].CreatedAt.Equal(items[j].CreatedAt) {
			return items[i].ID.String() < items[j].ID.String()
		}
		return items[i].CreatedAt.After(items[j].CreatedAt)
	})
	return items, nil
}

func (s *Store) CreateUserPhoto(_ context.Context, photo domain.UserPhoto) (domain.UserPhoto, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.userPhotos[photo.ID] = photo
	return photo, nil
}

func (s *Store) DeleteUserPhoto(_ context.Context, userID, id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	existing, ok := s.userPhotos[id]
	if !ok || existing.UserID != userID {
		return domain.ErrNotFound
	}
	delete(s.userPhotos, id)
	return nil
}
