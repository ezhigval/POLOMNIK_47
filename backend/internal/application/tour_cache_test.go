package application

import (
	"context"
	"sync"
	"testing"
	"time"

	"polomnik/internal/adapters/integration/noop"
	"polomnik/internal/adapters/repository/memory"
	"polomnik/internal/domain"
	"polomnik/internal/ports"
)

type memoryCache struct {
	mu    sync.RWMutex
	items map[string][]byte
}

func newMemoryCache() *memoryCache {
	return &memoryCache{items: make(map[string][]byte)}
}

func (c *memoryCache) Get(_ context.Context, key string) ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, ok := c.items[key]
	if !ok {
		return nil, ports.ErrCacheMiss
	}
	return value, nil
}

func (c *memoryCache) Set(_ context.Context, key string, value []byte, _ time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = value
	return nil
}

func (c *memoryCache) Delete(_ context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
	return nil
}

func TestTourServiceCacheHitAndInvalidation(t *testing.T) {
	ctx := context.Background()
	store := memory.NewStore()
	cache := newMemoryCache()
	service := NewTourService(store, cache, noop.NewCRMAdapter())

	tour := testTour()
	if _, err := store.CreateTour(ctx, tour); err != nil {
		t.Fatalf("create tour: %v", err)
	}

	first, err := service.ListPublicToursCached(ctx, ports.TourFilters{}, ports.Pagination{Page: 1, Limit: 20})
	if err != nil {
		t.Fatalf("first list: %v", err)
	}
	if len(first.Items) != 1 {
		t.Fatalf("expected 1 tour, got %d", len(first.Items))
	}

	if _, err := service.CreateTour(ctx, CreateTourInput{
		Slug:       "another-tour",
		Title:      "Another Tour",
		Price:      1000,
		Currency:   "RUB",
		DateStart:  testDate(2026, 9, 1),
		DateEnd:    testDate(2026, 9, 2),
		SlotsTotal: 5,
		SlotsLeft:  5,
		IsActive:   true,
	}); err != nil {
		t.Fatalf("create another tour: %v", err)
	}

	second, err := service.ListPublicToursCached(ctx, ports.TourFilters{}, ports.Pagination{Page: 1, Limit: 20})
	if err != nil {
		t.Fatalf("second list: %v", err)
	}
	if len(second.Items) != 2 {
		t.Fatalf("expected cache invalidation to return 2 tours, got %d", len(second.Items))
	}
}

func TestTourServiceWorksWithoutCache(t *testing.T) {
	ctx := context.Background()
	store := memory.NewStore()
	service := NewTourService(store, nil, noop.NewCRMAdapter())

	tour := testTour(func(input *domain.NewTourInput) {
		input.IsHot = true
	})
	if _, err := store.CreateTour(ctx, tour); err != nil {
		t.Fatalf("create tour: %v", err)
	}

	items, err := service.ListPopularToursCached(ctx, 10)
	if err != nil {
		t.Fatalf("list popular: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 popular tour, got %d", len(items))
	}
}
