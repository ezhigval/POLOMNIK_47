package memory

import (
	"context"
	"sort"
	"strings"

	"github.com/google/uuid"

	"palomnik/internal/domain"
	"palomnik/internal/ports"
)

func (s *Store) ListNews(_ context.Context, filters ports.NewsFilters, pagination ports.Pagination) (ports.NewsList, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	items := make([]domain.NewsArticle, 0, len(s.news))
	for _, article := range s.news {
		if filters.PublishedOnly && !article.IsPublished {
			continue
		}
		items = append(items, article)
	}

	sort.Slice(items, func(i, j int) bool {
		if !items[i].PublishedAt.Equal(items[j].PublishedAt) {
			return items[i].PublishedAt.After(items[j].PublishedAt)
		}
		if items[i].SortOrder != items[j].SortOrder {
			return items[i].SortOrder < items[j].SortOrder
		}
		return items[i].CreatedAt.After(items[j].CreatedAt)
	})

	pageItems, meta := page(items, pagination)
	return ports.NewsList{Items: pageItems, Meta: meta}, nil
}

func (s *Store) GetNewsBySlug(_ context.Context, slug string) (domain.NewsArticle, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	slug = strings.TrimSpace(slug)
	for _, article := range s.news {
		if article.Slug == slug {
			return article, nil
		}
	}
	return domain.NewsArticle{}, domain.ErrNotFound
}

func (s *Store) GetNews(_ context.Context, id uuid.UUID) (domain.NewsArticle, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	article, ok := s.news[id]
	if !ok {
		return domain.NewsArticle{}, domain.ErrNotFound
	}
	return article, nil
}

func (s *Store) CreateNews(_ context.Context, article domain.NewsArticle) (domain.NewsArticle, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.newsHasSlug(article.Slug, uuid.Nil) {
		return domain.NewsArticle{}, domain.ErrDuplicateSlug
	}
	s.news[article.ID] = article
	return article, nil
}

func (s *Store) UpdateNews(_ context.Context, article domain.NewsArticle) (domain.NewsArticle, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.news[article.ID]; !ok {
		return domain.NewsArticle{}, domain.ErrNotFound
	}
	if s.newsHasSlug(article.Slug, article.ID) {
		return domain.NewsArticle{}, domain.ErrDuplicateSlug
	}
	s.news[article.ID] = article
	return article, nil
}

func (s *Store) DeleteNews(_ context.Context, id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.news[id]; !ok {
		return domain.ErrNotFound
	}
	delete(s.news, id)
	return nil
}

func (s *Store) newsHasSlug(slug string, except uuid.UUID) bool {
	for _, article := range s.news {
		if article.Slug == slug && article.ID != except {
			return true
		}
	}
	return false
}
