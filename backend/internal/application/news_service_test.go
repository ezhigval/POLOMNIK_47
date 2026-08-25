package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"palomnik/internal/adapters/repository/memory"
	"palomnik/internal/domain"
	"palomnik/internal/ports"
)

func TestNewsServiceRejectsFourthPin(t *testing.T) {
	store := memory.NewStore()
	svc := NewNewsService(store, nil)
	ctx := context.Background()
	day := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)

	for i := 0; i < domain.MaxPinnedNews; i++ {
		if _, err := svc.CreateNews(ctx, NewsArticleInput{
			Slug:        "pin-" + string(rune('a'+i)),
			Title:       "Закреплённая",
			Excerpt:     "Анонс",
			Body:        "Текст статьи.",
			PublishedAt: day.AddDate(0, 0, i),
			IsPublished: true,
			IsPinned:    true,
			SortOrder:   i + 1,
		}); err != nil {
			t.Fatalf("pin %d: %v", i, err)
		}
	}

	_, err := svc.CreateNews(ctx, NewsArticleInput{
		Slug:        "pin-d",
		Title:       "Четвёртая",
		Excerpt:     "Анонс",
		Body:        "Текст статьи.",
		PublishedAt: day,
		IsPublished: true,
		IsPinned:    true,
		SortOrder:   4,
	})
	if !errors.Is(err, domain.ErrTooManyPinnedNews) {
		t.Fatalf("expected too many pinned, got %v", err)
	}
}

func TestNewsServiceAllowsUpdatingExistingPin(t *testing.T) {
	store := memory.NewStore()
	svc := NewNewsService(store, nil)
	ctx := context.Background()
	day := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)

	var first domain.NewsArticle
	for i := 0; i < domain.MaxPinnedNews; i++ {
		article, err := svc.CreateNews(ctx, NewsArticleInput{
			Slug:        "pin-" + string(rune('a'+i)),
			Title:       "Закреплённая",
			Excerpt:     "Анонс",
			Body:        "Текст статьи.",
			PublishedAt: day.AddDate(0, 0, i),
			IsPublished: true,
			IsPinned:    true,
			SortOrder:   i + 1,
		})
		if err != nil {
			t.Fatalf("pin %d: %v", i, err)
		}
		if i == 0 {
			first = article
		}
	}

	updated, err := svc.UpdateNews(ctx, first.ID, NewsArticleInput{
		Slug:        first.Slug,
		Title:       "Главная закреплённая",
		Excerpt:     first.Excerpt,
		Body:        first.Body,
		PublishedAt: first.PublishedAt,
		IsPublished: true,
		IsPinned:    true,
		SortOrder:   1,
	})
	if err != nil {
		t.Fatalf("update existing pin: %v", err)
	}
	if updated.Title != "Главная закреплённая" || !updated.IsPinned {
		t.Fatalf("unexpected update: %+v", updated)
	}
}

func TestNewsServiceListsPinnedBeforeNewer(t *testing.T) {
	store := memory.NewStore()
	svc := NewNewsService(store, nil)
	ctx := context.Background()

	if _, err := svc.CreateNews(ctx, NewsArticleInput{
		Slug:        "fresh",
		Title:       "Свежая",
		Excerpt:     "Анонс",
		Body:        "Текст статьи.",
		PublishedAt: time.Date(2026, 8, 21, 0, 0, 0, 0, time.UTC),
		IsPublished: true,
	}); err != nil {
		t.Fatalf("fresh: %v", err)
	}
	if _, err := svc.CreateNews(ctx, NewsArticleInput{
		Slug:        "pinned",
		Title:       "Закреплённая",
		Excerpt:     "Анонс",
		Body:        "Текст статьи.",
		PublishedAt: time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC),
		IsPublished: true,
		IsPinned:    true,
		SortOrder:   1,
	}); err != nil {
		t.Fatalf("pinned: %v", err)
	}

	list, err := svc.ListPublicNews(ctx, ports.Pagination{Page: 1, Limit: 20})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(list.Items) != 2 || list.Items[0].Slug != "pinned" || list.Items[1].Slug != "fresh" {
		got := make([]string, 0, len(list.Items))
		for _, item := range list.Items {
			got = append(got, item.Slug)
		}
		t.Fatalf("order: %v", got)
	}
}
