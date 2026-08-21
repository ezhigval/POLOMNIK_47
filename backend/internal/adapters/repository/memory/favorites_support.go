package memory

import (
	"context"
	"time"

	"github.com/google/uuid"

	"polomnik/internal/domain"
)

func (s *Store) AddFavorite(_ context.Context, userID, tourID uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	key := favoriteKey(userID, tourID)
	if _, ok := s.favorites[key]; ok {
		return nil
	}
	s.favorites[key] = domain.Favorite{UserID: userID, TourID: tourID, CreatedAt: time.Now().UTC()}
	return nil
}

func (s *Store) RemoveFavorite(_ context.Context, userID, tourID uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	key := favoriteKey(userID, tourID)
	if _, ok := s.favorites[key]; !ok {
		return domain.ErrNotFound
	}
	delete(s.favorites, key)
	return nil
}

func (s *Store) ListFavoriteTourIDs(_ context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var ids []uuid.UUID
	for _, favorite := range s.favorites {
		if favorite.UserID == userID {
			ids = append(ids, favorite.TourID)
		}
	}
	return ids, nil
}

func (s *Store) IsFavorite(_ context.Context, userID, tourID uuid.UUID) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.favorites[favoriteKey(userID, tourID)]
	return ok, nil
}

func (s *Store) GetOpenThread(_ context.Context, userID uuid.UUID) (domain.SupportThread, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, thread := range s.supportThreads {
		if thread.UserID == userID && thread.Status == "open" {
			return thread, nil
		}
	}
	return domain.SupportThread{}, domain.ErrNotFound
}

func (s *Store) CreateThread(_ context.Context, thread domain.SupportThread) (domain.SupportThread, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.supportThreads[thread.ID] = thread
	return thread, nil
}

func (s *Store) ListMessages(_ context.Context, threadID uuid.UUID) ([]domain.SupportMessage, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var messages []domain.SupportMessage
	for _, message := range s.supportMessages {
		if message.ThreadID == threadID {
			messages = append(messages, message)
		}
	}
	return messages, nil
}

func (s *Store) AddMessage(_ context.Context, message domain.SupportMessage) (domain.SupportMessage, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.supportMessages[message.ID] = message
	return message, nil
}

func (s *Store) TouchThread(_ context.Context, threadID uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	thread, ok := s.supportThreads[threadID]
	if !ok {
		return domain.ErrNotFound
	}
	thread.UpdatedAt = time.Now().UTC()
	s.supportThreads[threadID] = thread
	return nil
}

func favoriteKey(userID, tourID uuid.UUID) string {
	return userID.String() + ":" + tourID.String()
}
