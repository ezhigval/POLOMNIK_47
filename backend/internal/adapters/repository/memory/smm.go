package memory

import (
	"context"
	"sort"
	"time"

	"github.com/google/uuid"

	"palomnik/internal/domain"
	"palomnik/internal/ports"
)

func (s *Store) CreateSMMPost(_ context.Context, post domain.SMMPost) (domain.SMMPost, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.smmPosts == nil {
		s.smmPosts = make(map[uuid.UUID]domain.SMMPost)
	}
	s.smmPosts[post.ID] = cloneSMMPost(post)
	return cloneSMMPost(post), nil
}

func (s *Store) GetSMMPost(_ context.Context, id uuid.UUID) (domain.SMMPost, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	post, ok := s.smmPosts[id]
	if !ok {
		return domain.SMMPost{}, domain.ErrNotFound
	}
	return cloneSMMPost(post), nil
}

func (s *Store) ListSMMPosts(_ context.Context, pagination ports.Pagination) (ports.SMMPostList, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	items := make([]domain.SMMPost, 0, len(s.smmPosts))
	for _, post := range s.smmPosts {
		items = append(items, cloneSMMPost(post))
	}
	sort.Slice(items, func(i, j int) bool {
		if !items[i].PublishAt.Equal(items[j].PublishAt) {
			return items[i].PublishAt.After(items[j].PublishAt)
		}
		return items[i].CreatedAt.After(items[j].CreatedAt)
	})
	pageItems, meta := page(items, pagination)
	return ports.SMMPostList{Items: pageItems, Meta: meta}, nil
}

func (s *Store) ListDueSMMPosts(_ context.Context, now time.Time) ([]domain.SMMPost, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var due []domain.SMMPost
	for _, post := range s.smmPosts {
		if post.PublishedAt == nil && !post.PublishAt.After(now) {
			due = append(due, cloneSMMPost(post))
		}
	}
	sort.Slice(due, func(i, j int) bool { return due[i].PublishAt.Before(due[j].PublishAt) })
	return due, nil
}

func (s *Store) SaveSMMPost(_ context.Context, post domain.SMMPost) (domain.SMMPost, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.smmPosts[post.ID]; !ok {
		return domain.SMMPost{}, domain.ErrNotFound
	}
	s.smmPosts[post.ID] = cloneSMMPost(post)
	return cloneSMMPost(post), nil
}

func (s *Store) DeleteSMMPost(_ context.Context, id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.smmPosts[id]; !ok {
		return domain.ErrNotFound
	}
	delete(s.smmPosts, id)
	return nil
}

func cloneSMMPost(post domain.SMMPost) domain.SMMPost {
	post.Channels = append([]string(nil), post.Channels...)
	post.Results = append([]domain.SMMChannelResult(nil), post.Results...)
	if post.PublishedAt != nil {
		ts := *post.PublishedAt
		post.PublishedAt = &ts
	}
	return post
}
