package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

const MaxPinnedNews = 3

type NewsArticle struct {
	ID          uuid.UUID
	Slug        string
	Title       string
	Excerpt     string
	Body        string
	ImageURL    string
	PublishedAt time.Time
	IsPublished bool
	IsPinned    bool
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
	IsPinned    bool
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
		IsPinned:    input.IsPinned,
		SortOrder:   input.SortOrder,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// NewsListLess reports whether a should appear before b in public and admin lists:
// pinned articles first (lower sort_order is higher), then published_at descending.
func NewsListLess(a, b NewsArticle) bool {
	if a.IsPinned != b.IsPinned {
		return a.IsPinned
	}
	if a.IsPinned && a.SortOrder != b.SortOrder {
		return a.SortOrder < b.SortOrder
	}
	if !a.PublishedAt.Equal(b.PublishedAt) {
		return a.PublishedAt.After(b.PublishedAt)
	}
	if a.SortOrder != b.SortOrder {
		return a.SortOrder < b.SortOrder
	}
	return a.CreatedAt.After(b.CreatedAt)
}
