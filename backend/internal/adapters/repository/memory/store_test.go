package memory

import (
	"context"
	"testing"

	"palomnik/internal/domain"
	"palomnik/internal/ports"
)

func TestStoreImplementsRepositoryPorts(t *testing.T) {
	store := NewStore()

	var _ ports.TourRepository = store
	var _ ports.BookingRepository = store
	var _ ports.ReviewRepository = store
	var _ ports.IntegrationReferenceRepository = store
	var _ ports.OutboxRepository = store
}

func TestStoreTourSlotReservation(t *testing.T) {
	ctx := context.Background()
	store := NewStore()
	tour, err := domain.NewTour(domain.NewTourInput{
		ID:         testUUID("11111111-1111-1111-1111-111111111111"),
		Slug:       "test-tour",
		Title:      "Test Tour",
		Price:      100,
		Currency:   "RUB",
		DateStart:  testDate(2026, 6, 1),
		DateEnd:    testDate(2026, 6, 2),
		SlotsTotal: 5,
		SlotsLeft:  5,
		IsActive:   true,
	})
	if err != nil {
		t.Fatalf("create tour: %v", err)
	}

	if _, err := store.CreateTour(ctx, tour); err != nil {
		t.Fatalf("store tour: %v", err)
	}
	if err := store.ReserveSlots(ctx, tour.ID, 2); err != nil {
		t.Fatalf("reserve slots: %v", err)
	}

	stored, err := store.GetTour(ctx, tour.ID)
	if err != nil {
		t.Fatalf("get tour: %v", err)
	}
	if stored.SlotsLeft != 3 {
		t.Fatalf("expected 3 slots left, got %d", stored.SlotsLeft)
	}
}
