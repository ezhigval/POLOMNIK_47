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

func TestNewsListLessPutsPinnedFirst(t *testing.T) {
	olderPinned := mustNews(t, func(input *NewNewsArticleInput) {
		input.Slug = "pinned"
		input.Title = "Закреплена"
		input.IsPinned = true
		input.SortOrder = 2
		input.PublishedAt = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	})
	newer := mustNews(t, func(input *NewNewsArticleInput) {
		input.Slug = "fresh"
		input.Title = "Свежая"
		input.PublishedAt = time.Date(2026, 8, 21, 0, 0, 0, 0, time.UTC)
	})
	mainPinned := mustNews(t, func(input *NewNewsArticleInput) {
		input.Slug = "main"
		input.Title = "Главная"
		input.IsPinned = true
		input.SortOrder = 1
		input.PublishedAt = time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)
	})
	if !NewsListLess(olderPinned, newer) {
		t.Fatal("pinned should precede later unpinned news")
	}
	if !NewsListLess(mainPinned, olderPinned) {
		t.Fatal("lower sort_order should be the main pinned article")
	}
	if NewsListLess(newer, mainPinned) {
		t.Fatal("unpinned should not precede pinned")
	}
}

func mustNews(t *testing.T, mutate func(*NewNewsArticleInput)) NewsArticle {
	t.Helper()
	article, err := NewNewsArticle(validNewsInput(mutate))
	if err != nil {
		t.Fatalf("news: %v", err)
	}
	return article
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
