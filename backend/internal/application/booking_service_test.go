package application

import (
	"context"
	"errors"
	"testing"

	"polomnik/internal/adapters/integration/noop"
	notificationnoop "polomnik/internal/adapters/notification/noop"
	"polomnik/internal/adapters/repository/memory"
	"polomnik/internal/domain"
	"polomnik/internal/ports"
)

func TestBookingServiceCreateBookingWithNoopCRMIntegration(t *testing.T) {
	ctx := context.Background()
	store := memory.NewStore()
	service := NewBookingService(store, store, noop.NewCRMAdapter(), noop.NewAccountingAdapter(), notificationnoop.New(), store)

	tour := testTour()
	if _, err := store.CreateTour(ctx, tour); err != nil {
		t.Fatalf("create tour: %v", err)
	}

	result, err := service.CreateBooking(ctx, CreateBookingInput{
		TourID:      tour.ID,
		Name:        "Иван Иванов",
		Phone:       "+79999999999",
		PeopleCount: 1,
	})
	if err != nil {
		t.Fatalf("create booking: %v", err)
	}
	if result.IntegrationStatus != ports.IntegrationStatusNotConfigured {
		t.Fatalf("expected not_configured, got %s", result.IntegrationStatus)
	}
}

func newBookingService(store *memory.Store) *BookingService {
	return NewBookingService(store, store, noop.NewCRMAdapter(), noop.NewAccountingAdapter(), notificationnoop.New(), store)
}

func TestBookingServiceCreateBookingReservesSlots(t *testing.T) {
	ctx := context.Background()
	store := memory.NewStore()
	service := newBookingService(store)

	tour := testTour()
	if _, err := store.CreateTour(ctx, tour); err != nil {
		t.Fatalf("create tour: %v", err)
	}

	result, err := service.CreateBooking(ctx, CreateBookingInput{
		TourID:      tour.ID,
		Name:        "Иван Иванов",
		Phone:       "+79999999999",
		PeopleCount: 2,
	})
	if err != nil {
		t.Fatalf("create booking: %v", err)
	}
	if result.Booking.TotalPrice != 20000 {
		t.Fatalf("expected total price 20000, got %d", result.Booking.TotalPrice)
	}
	if result.IntegrationStatus != ports.IntegrationStatusNotConfigured {
		t.Fatalf("expected not_configured integration status, got %s", result.IntegrationStatus)
	}

	storedTour, err := store.GetTour(ctx, tour.ID)
	if err != nil {
		t.Fatalf("get tour: %v", err)
	}
	if storedTour.SlotsLeft != 8 {
		t.Fatalf("expected 8 slots left, got %d", storedTour.SlotsLeft)
	}
}

func TestBookingServiceCreateBookingRejectsInactiveTour(t *testing.T) {
	ctx := context.Background()
	store := memory.NewStore()
	service := newBookingService(store)

	tour := testTour(func(input *domain.NewTourInput) {
		input.IsActive = false
	})
	if _, err := store.CreateTour(ctx, tour); err != nil {
		t.Fatalf("create tour: %v", err)
	}

	_, err := service.CreateBooking(ctx, CreateBookingInput{
		TourID:      tour.ID,
		Name:        "Иван Иванов",
		Phone:       "+79999999999",
		PeopleCount: 1,
	})
	if !errors.Is(err, ErrTourInactive) {
		t.Fatalf("expected tour inactive, got %v", err)
	}
}

func TestBookingServiceCreateBookingRejectsInsufficientSlots(t *testing.T) {
	ctx := context.Background()
	store := memory.NewStore()
	service := newBookingService(store)

	tour := testTour(func(input *domain.NewTourInput) {
		input.SlotsLeft = 1
		input.OverbookingEnabled = false
	})
	if _, err := store.CreateTour(ctx, tour); err != nil {
		t.Fatalf("create tour: %v", err)
	}

	_, err := service.CreateBooking(ctx, CreateBookingInput{
		TourID:      tour.ID,
		Name:        "Иван Иванов",
		Phone:       "+79999999999",
		PeopleCount: 2,
	})
	if !errors.Is(err, domain.ErrInsufficientSlots) {
		t.Fatalf("expected insufficient slots, got %v", err)
	}
}

func TestBookingServiceCancelBookingReleasesSlots(t *testing.T) {
	ctx := context.Background()
	store := memory.NewStore()
	service := newBookingService(store)

	tour := testTour()
	if _, err := store.CreateTour(ctx, tour); err != nil {
		t.Fatalf("create tour: %v", err)
	}

	result, err := service.CreateBooking(ctx, CreateBookingInput{
		TourID:      tour.ID,
		Name:        "Иван Иванов",
		Phone:       "+79999999999",
		PeopleCount: 3,
	})
	if err != nil {
		t.Fatalf("create booking: %v", err)
	}

	if _, err := service.UpdateBookingStatus(ctx, result.Booking.ID, domain.BookingStatusCancelled); err != nil {
		t.Fatalf("cancel booking: %v", err)
	}

	storedTour, err := store.GetTour(ctx, tour.ID)
	if err != nil {
		t.Fatalf("get tour: %v", err)
	}
	if storedTour.SlotsLeft != 10 {
		t.Fatalf("expected slots released back to 10, got %d", storedTour.SlotsLeft)
	}
}

func TestBookingServiceUpdateBookingStatusFollowsDomainRules(t *testing.T) {
	ctx := context.Background()
	store := memory.NewStore()
	service := newBookingService(store)

	tour := testTour()
	if _, err := store.CreateTour(ctx, tour); err != nil {
		t.Fatalf("create tour: %v", err)
	}

	result, err := service.CreateBooking(ctx, CreateBookingInput{
		TourID:      tour.ID,
		Name:        "Иван Иванов",
		Phone:       "+79999999999",
		PeopleCount: 1,
	})
	if err != nil {
		t.Fatalf("create booking: %v", err)
	}

	_, err = service.UpdateBookingStatus(ctx, result.Booking.ID, domain.BookingStatusConfirmed)
	if !errors.Is(err, domain.ErrInvalidStatusTransition) {
		t.Fatalf("expected invalid transition, got %v", err)
	}
}
