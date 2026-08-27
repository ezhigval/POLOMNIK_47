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

func TestTourServiceFindTourBySlug(t *testing.T) {
	ctx := context.Background()
	store := memory.NewStore()
	service := NewTourService(store, nil, noop.NewCRMAdapter())

	tour := testTour(func(input *domain.NewTourInput) {
		input.Slug = "Optina-Pustyn"
	})
	if _, err := store.CreateTour(ctx, tour); err != nil {
		t.Fatalf("create tour: %v", err)
	}

	found, err := service.FindTour(ctx, "optina-pustyn")
	if err != nil {
		t.Fatalf("find tour by slug: %v", err)
	}
	if found.ID != tour.ID {
		t.Fatalf("expected tour %s, got %s", tour.ID, found.ID)
	}
}

func TestTourServiceFindPublicTourHidesInactive(t *testing.T) {
	ctx := context.Background()
	store := memory.NewStore()
	service := NewTourService(store, nil, noop.NewCRMAdapter())

	tour := testTour(func(input *domain.NewTourInput) {
		input.Slug = "inactive-slug"
		input.IsActive = false
	})
	if _, err := store.CreateTour(ctx, tour); err != nil {
		t.Fatalf("create tour: %v", err)
	}

	_, err := service.FindPublicTour(ctx, "inactive-slug")
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected not found for inactive slug, got %v", err)
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

func TestTourServiceCreateRegularTourFillsRemainingSlots(t *testing.T) {
	ctx := context.Background()
	store := memory.NewStore()
	service := NewTourService(store, nil, noop.NewCRMAdapter())

	created, err := service.CreateTour(ctx, CreateTourInput{
		Slug:       "regular-seats",
		Title:      "Regular seats",
		Price:      0,
		Currency:   "RUB",
		SlotsTotal: 50,
		SlotsLeft:  0,
		Location:   "SPB",
		IsActive:   true,
		IsRegular:  true,
	})
	if err != nil {
		t.Fatalf("create tour: %v", err)
	}
	if created.SlotsLeft != 50 {
		t.Fatalf("expected slots_left 50, got %d", created.SlotsLeft)
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

func TestTourServiceListPublicToursPutsRegularAfterDated(t *testing.T) {
	ctx := context.Background()
	store := memory.NewStore()
	service := NewTourService(store, nil, noop.NewCRMAdapter())

	regular := testTour(func(input *domain.NewTourInput) {
		input.ID = testUUID("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
		input.Slug = "regular-listed"
		input.IsRegular = true
		input.DateStart = testDate(2026, 1, 1)
		input.DateEnd = testDate(2026, 1, 2)
	})
	laterDated := testTour(func(input *domain.NewTourInput) {
		input.ID = testUUID("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")
		input.Slug = "dated-later"
		input.DateStart = testDate(2026, 11, 1)
		input.DateEnd = testDate(2026, 11, 3)
	})
	soonerDated := testTour(func(input *domain.NewTourInput) {
		input.ID = testUUID("cccccccc-cccc-cccc-cccc-cccccccccccc")
		input.Slug = "dated-sooner"
		input.DateStart = testDate(2026, 9, 1)
		input.DateEnd = testDate(2026, 9, 3)
	})

	for _, tour := range []domain.Tour{regular, laterDated, soonerDated} {
		if _, err := store.CreateTour(ctx, tour); err != nil {
			t.Fatalf("create tour %s: %v", tour.Slug, err)
		}
	}

	list, err := service.ListPublicTours(ctx, ports.TourFilters{}, ports.Pagination{Page: 1, Limit: 20})
	if err != nil {
		t.Fatalf("list public tours: %v", err)
	}
	if len(list.Items) != 3 {
		t.Fatalf("expected 3 tours, got %d", len(list.Items))
	}
	if list.Items[0].Slug != "dated-sooner" || list.Items[1].Slug != "dated-later" || list.Items[2].Slug != "regular-listed" {
		t.Fatalf("unexpected order: %s, %s, %s", list.Items[0].Slug, list.Items[1].Slug, list.Items[2].Slug)
	}
}
