package memory

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"

	"palomnik/internal/domain"
	"palomnik/internal/ports"
)

type newsLikeKey struct {
	newsID    uuid.UUID
	visitorID string
}

func (s *Store) CountNewsLikes(_ context.Context, newsID uuid.UUID) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	count := 0
	for key := range s.newsLikes {
		if key.newsID == newsID {
			count++
		}
	}
	return count, nil
}

func (s *Store) HasNewsLike(_ context.Context, newsID uuid.UUID, visitorID string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.newsLikes[newsLikeKey{newsID: newsID, visitorID: strings.TrimSpace(visitorID)}]
	return ok, nil
}

func (s *Store) AddNewsLike(_ context.Context, newsID uuid.UUID, visitorID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.newsLikes[newsLikeKey{newsID: newsID, visitorID: strings.TrimSpace(visitorID)}] = time.Now().UTC()
	return nil
}

func (s *Store) RemoveNewsLike(_ context.Context, newsID uuid.UUID, visitorID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.newsLikes, newsLikeKey{newsID: newsID, visitorID: strings.TrimSpace(visitorID)})
	return nil
}

func (s *Store) CountNewsComments(_ context.Context, newsID uuid.UUID) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	count := 0
	for _, comment := range s.newsComments {
		if comment.NewsID == newsID {
			count++
		}
	}
	return count, nil
}

func (s *Store) ListNewsComments(_ context.Context, newsID uuid.UUID, pagination ports.Pagination) ([]domain.NewsComment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	pagination = ports.NormalizePagination(pagination.Page, pagination.Limit)
	items := make([]domain.NewsComment, 0)
	for _, comment := range s.newsComments {
		if comment.NewsID == newsID {
			items = append(items, comment)
		}
	}

	sort.SliceStable(items, func(i, j int) bool {
		return items[i].CreatedAt.Before(items[j].CreatedAt)
	})

	start := (pagination.Page - 1) * pagination.Limit
	if start >= len(items) {
		return []domain.NewsComment{}, nil
	}
	end := start + pagination.Limit
	if end > len(items) {
		end = len(items)
	}
	return items[start:end], nil
}

func (s *Store) AddNewsComment(_ context.Context, comment domain.NewsComment) (domain.NewsComment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.newsComments[comment.ID]; ok {
		return domain.NewsComment{}, fmt.Errorf("duplicate comment id")
	}
	s.newsComments[comment.ID] = comment
	return comment, nil
}

var _ ports.NewsEngagementRepository = (*Store)(nil)
