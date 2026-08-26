package application

import (
	"context"
	"errors"
	"testing"

	"palomnik/internal/adapters/integration/noop"
	notificationnoop "palomnik/internal/adapters/notification/noop"
	"palomnik/internal/adapters/repository/memory"
	"palomnik/internal/domain"
	"palomnik/internal/ports"
)

func TestBookingServiceCreateBookingWithNoopCRMIntegration(t *testing.T) {
	ctx := context.Background()
	store := memory.NewStore()
	service := NewBookingService(store, store, noop.NewCRMAdapter(), noop.NewAccountingAdapter(), notificationnoop.New(), store, 0)

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
	return NewBookingService(store, store, noop.NewCRMAdapter(), noop.NewAccountingAdapter(), notificationnoop.New(), store, 0)
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

func TestBookingServiceCreateBookingAcceptsRegularTour(t *testing.T) {
	ctx := context.Background()
	store := memory.NewStore()
	service := newBookingService(store)

	tour := testTour(func(input *domain.NewTourInput) {
		input.IsRegular = true
		input.Price = 0
		input.DateStart = testDate(2020, 1, 1)
		input.DateEnd = testDate(2020, 1, 2)
	})
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
	if result.Booking.Status != domain.BookingStatusNew {
		t.Fatalf("expected NEW, got %s", result.Booking.Status)
	}
	if result.Booking.TotalPrice != 0 {
		t.Fatalf("expected total 0, got %d", result.Booking.TotalPrice)
	}
}

func TestBookingServiceCreateBookingRejectsExpiredDatedTour(t *testing.T) {
	ctx := context.Background()
	store := memory.NewStore()
	service := newBookingService(store)

	tour := testTour(func(input *domain.NewTourInput) {
		input.DateStart = testDate(2020, 1, 1)
		input.DateEnd = testDate(2020, 1, 5)
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
	if !errors.Is(err, ErrTourExpired) {
		t.Fatalf("expected tour expired, got %v", err)
	}
}

func TestBookingServiceListAllBookingsFiltersByStatus(t *testing.T) {
	ctx := context.Background()
	store := memory.NewStore()
	service := newBookingService(store)

	tour := testTour()
	if _, err := store.CreateTour(ctx, tour); err != nil {
		t.Fatalf("create tour: %v", err)
	}

	first, err := service.CreateBooking(ctx, CreateBookingInput{
		TourID:      tour.ID,
		Name:        "Андрей",
		Phone:       "+79990000001",
		PeopleCount: 1,
	})
	if err != nil {
		t.Fatalf("create first: %v", err)
	}
	if _, err := service.CreateBooking(ctx, CreateBookingInput{
		TourID:      tour.ID,
		Name:        "Борис",
		Phone:       "+79990000002",
		PeopleCount: 1,
	}); err != nil {
		t.Fatalf("create second: %v", err)
	}
	if _, err := service.UpdateBookingStatus(ctx, first.Booking.ID, domain.BookingStatusContacted); err != nil {
		t.Fatalf("contact first: %v", err)
	}

	status := domain.BookingStatusNew
	got, err := service.ListAllBookings(ctx, ports.BookingFilters{Status: &status})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(got) != 1 || got[0].Name != "Борис" {
		t.Fatalf("expected only NEW booking Борис, got %+v", got)
	}
}
