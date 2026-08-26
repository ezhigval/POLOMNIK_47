package application

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"palomnik/internal/adapters/integration/noop"
	"palomnik/internal/adapters/repository/memory"
	"palomnik/internal/domain"
)

type recordingTourNotifier struct {
	hidden []domain.Tour
}

func (n *recordingTourNotifier) NotifyNewBooking(context.Context, domain.Booking, domain.Tour) error {
	return nil
}

func (n *recordingTourNotifier) NotifyBookingStatusChanged(context.Context, domain.Booking, domain.Tour, domain.BookingStatus) error {
	return nil
}

func (n *recordingTourNotifier) NotifySupportMessage(context.Context, domain.SupportNotification) error {
	return nil
}

func (n *recordingTourNotifier) NotifyTourHidden(_ context.Context, tour domain.Tour) error {
	n.hidden = append(n.hidden, tour)
	return nil
}

func TestTourLifecycleServiceHidesExpiredTour(t *testing.T) {
	store := memory.NewStore()
	tours := NewTourService(store, nil, noop.NewCRMAdapter())
	ctx := context.Background()

	yesterday := time.Date(2026, 8, 25, 0, 0, 0, 0, time.UTC)
	today := time.Date(2026, 8, 26, 0, 0, 0, 0, time.UTC)
	expired, err := domain.NewTour(domain.NewTourInput{
		ID:         uuid.New(),
		Slug:       "expired",
		Title:      "Expired tour",
		Price:      5000,
		Currency:   "RUB",
		DateStart:  yesterday.Add(-2 * 24 * time.Hour),
		DateEnd:    yesterday,
		SlotsTotal: 10,
		SlotsLeft:  5,
		IsActive:   true,
	})
	if err != nil {
		t.Fatalf("create expired tour: %v", err)
	}
	if _, err := store.CreateTour(ctx, expired); err != nil {
		t.Fatalf("persist expired tour: %v", err)
	}

	activeRegular, err := domain.NewTour(domain.NewTourInput{
		ID:         uuid.New(),
		Slug:       "regular",
		Title:      "Regular tour",
		Price:      0,
		Currency:   "RUB",
		SlotsTotal: 10,
		SlotsLeft:  5,
		IsActive:   true,
		IsRegular:  true,
	})
	if err != nil {
		t.Fatalf("create regular tour: %v", err)
	}
	if _, err := store.CreateTour(ctx, activeRegular); err != nil {
		t.Fatalf("persist regular tour: %v", err)
	}

	notifier := &recordingTourNotifier{}
	lifecycle := NewTourLifecycleService(tours, notifier)
	hidden, err := lifecycle.HideExpiredActiveTours(ctx, today)
	if err != nil {
		t.Fatalf("hide expired tours: %v", err)
	}
	if hidden != 1 {
		t.Fatalf("expected 1 hidden tour, got %d", hidden)
	}
	if len(notifier.hidden) != 1 || notifier.hidden[0].ID != expired.ID {
		t.Fatalf("expected notification for expired tour, got %#v", notifier.hidden)
	}

	got, err := store.GetTour(ctx, expired.ID)
	if err != nil {
		t.Fatalf("get expired tour: %v", err)
	}
	if got.IsActive {
		t.Fatal("expected expired tour to be inactive")
	}

	gotRegular, err := store.GetTour(ctx, activeRegular.ID)
	if err != nil {
		t.Fatalf("get regular tour: %v", err)
	}
	if !gotRegular.IsActive {
		t.Fatal("regular tour must stay active")
	}
}
