package application

import (
	"context"
	"testing"
	"time"

	"palomnik/internal/adapters/repository/memory"
	"palomnik/internal/ports"
)

func TestNewsServiceCacheHitAndInvalidation(t *testing.T) {
	ctx := context.Background()
	store := memory.NewStore()
	cache := newMemoryCache()
	service := NewNewsService(store, cache)

	if _, err := service.CreateNews(ctx, NewsArticleInput{
		Slug:        "first",
		Title:       "First",
		Excerpt:     "Excerpt",
		Body:        "Body text",
		PublishedAt: time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC),
		IsPublished: true,
	}); err != nil {
		t.Fatalf("create: %v", err)
	}

	first, err := service.ListPublicNewsCached(ctx, ports.Pagination{Page: 1, Limit: 20})
	if err != nil {
		t.Fatalf("first list: %v", err)
	}
	if len(first.Items) != 1 {
		t.Fatalf("expected 1 article, got %d", len(first.Items))
	}

	if _, err := service.CreateNews(ctx, NewsArticleInput{
		Slug:        "second",
		Title:       "Second",
		Excerpt:     "Excerpt",
		Body:        "Body text",
		PublishedAt: time.Date(2026, 10, 2, 0, 0, 0, 0, time.UTC),
		IsPublished: true,
	}); err != nil {
		t.Fatalf("create second: %v", err)
	}

	second, err := service.ListPublicNewsCached(ctx, ports.Pagination{Page: 1, Limit: 20})
	if err != nil {
		t.Fatalf("second list: %v", err)
	}
	if len(second.Items) != 2 {
		t.Fatalf("expected cache invalidation, got %d", len(second.Items))
	}
}
