package application

import (
	"context"
	"errors"
	"testing"

	"palomnik/internal/adapters/integration/noop"
	"palomnik/internal/adapters/repository/memory"
	"palomnik/internal/domain"
	"palomnik/internal/ports"
)

func TestTourServiceListPublicToursReturnsOnlyActive(t *testing.T) {
	ctx := context.Background()
	store := memory.NewStore()
	service := NewTourService(store, nil, noop.NewCRMAdapter())

	activeTour := testTour()
	inactiveTour := testTour(func(input *domain.NewTourInput) {
		input.ID = testUUID("22222222-2222-2222-2222-222222222222")
		input.Slug = "inactive-tour"
		input.IsActive = false
	})

	if _, err := store.CreateTour(ctx, activeTour); err != nil {
		t.Fatalf("create active tour: %v", err)
	}
	if _, err := store.CreateTour(ctx, inactiveTour); err != nil {
		t.Fatalf("create inactive tour: %v", err)
	}

	list, err := service.ListPublicTours(ctx, ports.TourFilters{}, ports.Pagination{Page: 1, Limit: 20})
	if err != nil {
		t.Fatalf("list public tours: %v", err)
	}
	if len(list.Items) != 1 {
		t.Fatalf("expected 1 active tour, got %d", len(list.Items))
	}
	if list.Items[0].ID != activeTour.ID {
		t.Fatalf("expected active tour id %s, got %s", activeTour.ID, list.Items[0].ID)
	}
}

func TestTourServiceGetPublicTourHidesInactive(t *testing.T) {
	ctx := context.Background()
	store := memory.NewStore()
	service := NewTourService(store, nil, noop.NewCRMAdapter())

	tour := testTour(func(input *domain.NewTourInput) {
		input.IsActive = false
	})
	if _, err := store.CreateTour(ctx, tour); err != nil {
		t.Fatalf("create tour: %v", err)
	}

	_, err := service.GetPublicTour(ctx, tour.ID)
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected not found for inactive tour, got %v", err)
	}
}

func TestTourServiceListPopularTours(t *testing.T) {
	ctx := context.Background()
	store := memory.NewStore()
	service := NewTourService(store, nil, noop.NewCRMAdapter())

	hotTour := testTour(func(input *domain.NewTourInput) {
		input.IsHot = true
	})
	regularTour := testTour(func(input *domain.NewTourInput) {
		input.ID = testUUID("33333333-3333-3333-3333-333333333333")
		input.Slug = "regular-tour"
		input.IsHot = false
	})

	if _, err := store.CreateTour(ctx, hotTour); err != nil {
		t.Fatalf("create hot tour: %v", err)
	}
	if _, err := store.CreateTour(ctx, regularTour); err != nil {
		t.Fatalf("create regular tour: %v", err)
	}

	items, err := service.ListPopularTours(ctx, 10)
	if err != nil {
		t.Fatalf("list popular tours: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 hot tour, got %d", len(items))
	}
}

func TestTourServiceCreateAndUpdateTour(t *testing.T) {
	ctx := context.Background()
	store := memory.NewStore()
	service := NewTourService(store, nil, noop.NewCRMAdapter())

	created, err := service.CreateTour(ctx, CreateTourInput{
		Slug:       "new-tour",
		Title:      "New Tour",
		Price:      5000,
		Currency:   "RUB",
		DateStart:  testDate(2026, 9, 1),
		DateEnd:    testDate(2026, 9, 3),
		SlotsTotal: 8,
		SlotsLeft:  8,
		Location:   "SPB",
		IsActive:   true,
	})
	if err != nil {
		t.Fatalf("create tour: %v", err)
	}

	updated, err := service.UpdateTour(ctx, created.ID, UpdateTourInput{
		Slug:       created.Slug,
		Title:      "Updated Tour",
		Price:      created.Price,
		Currency:   created.Currency,
		DateStart:  created.DateStart,
		DateEnd:    created.DateEnd,
		SlotsTotal: created.SlotsTotal,
		SlotsLeft:  created.SlotsLeft,
		Location:   created.Location,
		IsActive:   created.IsActive,
	})
	if err != nil {
		t.Fatalf("update tour: %v", err)
	}
	if updated.Title != "Updated Tour" {
		t.Fatalf("expected updated title, got %q", updated.Title)
	}
}
