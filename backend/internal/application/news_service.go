package application

import (
	"context"
	"time"

	"github.com/google/uuid"

	"palomnik/internal/domain"
	"palomnik/internal/ports"
)

type NewsService struct {
	news  ports.NewsRepository
	cache ports.CachePort
}

func NewNewsService(news ports.NewsRepository, cache ports.CachePort) *NewsService {
	return &NewsService{news: news, cache: cache}
}

type NewsArticleInput struct {
	Slug        string
	Title       string
	Excerpt     string
	Body        string
	ImageURL    string
	PublishedAt time.Time
	IsPublished bool
	SortOrder   int
}

func (s *NewsService) ListPublicNews(ctx context.Context, pagination ports.Pagination) (ports.NewsList, error) {
	return s.news.ListNews(ctx, ports.NewsFilters{PublishedOnly: true}, pagination)
}

func (s *NewsService) ListPublicNewsCached(ctx context.Context, pagination ports.Pagination) (ports.NewsList, error) {
	if s.cache == nil {
		return s.ListPublicNews(ctx, pagination)
	}

	key, err := s.cacheKey(ctx, "list:public", newsListCacheID(pagination))
	if err != nil {
		return s.ListPublicNews(ctx, pagination)
	}

	var cached ports.NewsList
	if ok, _ := s.readCache(ctx, key, &cached); ok {
		return cached, nil
	}

	list, err := s.ListPublicNews(ctx, pagination)
	if err != nil {
		return ports.NewsList{}, err
	}
	_ = s.writeCache(ctx, key, list)
	return list, nil
}

func (s *NewsService) ListNews(ctx context.Context, pagination ports.Pagination) (ports.NewsList, error) {
	return s.news.ListNews(ctx, ports.NewsFilters{}, pagination)
}

func (s *NewsService) GetNews(ctx context.Context, id uuid.UUID) (domain.NewsArticle, error) {
	return s.news.GetNews(ctx, id)
}

func (s *NewsService) GetPublishedNewsBySlug(ctx context.Context, slug string) (domain.NewsArticle, error) {
	article, err := s.news.GetNewsBySlug(ctx, slug)
	if err != nil {
		return domain.NewsArticle{}, err
	}
	if !article.IsPublished {
		return domain.NewsArticle{}, domain.ErrNotFound
	}
	return article, nil
}

func (s *NewsService) CreateNews(ctx context.Context, input NewsArticleInput) (domain.NewsArticle, error) {
	article, err := domain.NewNewsArticle(domain.NewNewsArticleInput{
		ID:          uuid.New(),
		Slug:        input.Slug,
		Title:       input.Title,
		Excerpt:     input.Excerpt,
		Body:        input.Body,
		ImageURL:    input.ImageURL,
		PublishedAt: input.PublishedAt,
		IsPublished: input.IsPublished,
		SortOrder:   input.SortOrder,
	})
	if err != nil {
		return domain.NewsArticle{}, err
	}
	created, err := s.news.CreateNews(ctx, article)
	if err != nil {
		return domain.NewsArticle{}, err
	}
	s.invalidateNewsCache(ctx)
	return created, nil
}

func (s *NewsService) UpdateNews(ctx context.Context, id uuid.UUID, input NewsArticleInput) (domain.NewsArticle, error) {
	existing, err := s.news.GetNews(ctx, id)
	if err != nil {
		return domain.NewsArticle{}, err
	}

	updated, err := domain.NewNewsArticle(domain.NewNewsArticleInput{
		ID:          existing.ID,
		Slug:        input.Slug,
		Title:       input.Title,
		Excerpt:     input.Excerpt,
		Body:        input.Body,
		ImageURL:    input.ImageURL,
		PublishedAt: input.PublishedAt,
		IsPublished: input.IsPublished,
		SortOrder:   input.SortOrder,
		Now:         existing.CreatedAt,
	})
	if err != nil {
		return domain.NewsArticle{}, err
	}
	updated.CreatedAt = existing.CreatedAt
	updated.UpdatedAt = time.Now().UTC()
	saved, err := s.news.UpdateNews(ctx, updated)
	if err != nil {
		return domain.NewsArticle{}, err
	}
	s.invalidateNewsCache(ctx)
	return saved, nil
}

func (s *NewsService) DeleteNews(ctx context.Context, id uuid.UUID) error {
	if err := s.news.DeleteNews(ctx, id); err != nil {
		return err
	}
	s.invalidateNewsCache(ctx)
	return nil
}
