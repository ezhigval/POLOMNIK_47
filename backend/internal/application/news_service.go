package application

import (
	"context"
	"time"

	"github.com/google/uuid"

	"polomnik/internal/domain"
	"polomnik/internal/ports"
)

type NewsService struct {
	news ports.NewsRepository
}

func NewNewsService(news ports.NewsRepository) *NewsService {
	return &NewsService{news: news}
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

func (s *NewsService) ListNews(ctx context.Context, pagination ports.Pagination) (ports.NewsList, error) {
	return s.news.ListNews(ctx, ports.NewsFilters{}, pagination)
}

func (s *NewsService) GetNews(ctx context.Context, id uuid.UUID) (domain.NewsArticle, error) {
	return s.news.GetNews(ctx, id)
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
	return s.news.CreateNews(ctx, article)
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
	return s.news.UpdateNews(ctx, updated)
}

func (s *NewsService) DeleteNews(ctx context.Context, id uuid.UUID) error {
	return s.news.DeleteNews(ctx, id)
}
