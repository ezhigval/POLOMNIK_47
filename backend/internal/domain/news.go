package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type NewsArticle struct {
	ID          uuid.UUID
	Slug        string
	Title       string
	Excerpt     string
	Body        string
	ImageURL    string
	PublishedAt time.Time
	IsPublished bool
	SortOrder   int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type NewNewsArticleInput struct {
	ID          uuid.UUID
	Slug        string
	Title       string
	Excerpt     string
	Body        string
	ImageURL    string
	PublishedAt time.Time
	IsPublished bool
	SortOrder   int
	Now         time.Time
}

func NewNewsArticle(input NewNewsArticleInput) (NewsArticle, error) {
	if input.ID == uuid.Nil {
		return NewsArticle{}, ErrInvalidID
	}

	slug := normalizeSlug(input.Slug)
	title := strings.TrimSpace(input.Title)
	excerpt := strings.TrimSpace(input.Excerpt)
	body := strings.TrimSpace(input.Body)
	imageURL := strings.TrimSpace(input.ImageURL)

	if slug == "" {
		return NewsArticle{}, ErrInvalidSlug
	}
	if title == "" {
		return NewsArticle{}, ErrInvalidTitle
	}
	if excerpt == "" {
		return NewsArticle{}, ErrInvalidExcerpt
	}
	if body == "" {
		return NewsArticle{}, ErrInvalidArticleBody
	}
	if input.PublishedAt.IsZero() {
		return NewsArticle{}, ErrInvalidPublishedAt
	}

	now := input.Now
	if now.IsZero() {
		now = time.Now().UTC()
	}

	return NewsArticle{
		ID:          input.ID,
		Slug:        slug,
		Title:       title,
		Excerpt:     excerpt,
		Body:        body,
		ImageURL:    imageURL,
		PublishedAt: input.PublishedAt.UTC().Truncate(24 * time.Hour),
		IsPublished: input.IsPublished,
		SortOrder:   input.SortOrder,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}
