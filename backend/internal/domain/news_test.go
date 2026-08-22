package domain

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewNewsArticleRequiresTitleAndBody(t *testing.T) {
	_, err := NewNewsArticle(validNewsInput(func(input *NewNewsArticleInput) {
		input.Title = "  "
		input.Body = ""
	}))
	if !errors.Is(err, ErrInvalidTitle) {
		t.Fatalf("expected invalid title, got %v", err)
	}
}

func TestNewNewsArticleNormalizesSlug(t *testing.T) {
	article, err := NewNewsArticle(validNewsInput(func(input *NewNewsArticleInput) {
		input.Slug = " Tikhvin Path "
	}))
	if err != nil {
		t.Fatalf("create article: %v", err)
	}
	if article.Slug != "tikhvin-path" {
		t.Fatalf("expected normalized slug, got %q", article.Slug)
	}
}

func validNewsInput(mutators ...func(*NewNewsArticleInput)) NewNewsArticleInput {
	input := NewNewsArticleInput{
		ID:          uuid.New(),
		Slug:        "tikhvin-path",
		Title:       "Тихвинский путь",
		Excerpt:     "Анонс статьи",
		Body:        "Полный текст статьи.",
		PublishedAt: time.Date(2026, 8, 21, 0, 0, 0, 0, time.UTC),
		IsPublished: true,
	}
	for _, mutate := range mutators {
		mutate(&input)
	}
	return input
}
