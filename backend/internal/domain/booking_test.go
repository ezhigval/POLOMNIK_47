package domain

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewBookingCalculatesTotalPrice(t *testing.T) {
	tour, err := NewTour(validTourInput(func(input *NewTourInput) {
		input.Price = 12000
		input.SlotsLeft = 5
	}))
	if err != nil {
		t.Fatalf("create tour: %v", err)
	}

	booking, err := NewBooking(validBookingInput(tour, func(input *NewBookingInput) {
		input.PeopleCount = 3
	}))
	if err != nil {
		t.Fatalf("create booking: %v", err)
	}

	if booking.TotalPrice != 36000 {
		t.Fatalf("expected total price 36000, got %d", booking.TotalPrice)
	}
	if booking.Status != BookingStatusNew {
		t.Fatalf("expected status NEW, got %s", booking.Status)
	}
}

func TestNewBookingRegularTourWithPriceCalculatesTotal(t *testing.T) {
	tour, err := NewTour(validTourInput(func(input *NewTourInput) {
		input.IsRegular = true
		input.Price = 3500
		input.SlotsLeft = 5
	}))
	if err != nil {
		t.Fatalf("create tour: %v", err)
	}

	booking, err := NewBooking(validBookingInput(tour, func(input *NewBookingInput) {
		input.PeopleCount = 2
	}))
	if err != nil {
		t.Fatalf("create booking: %v", err)
	}
	if booking.TotalPrice != 7000 {
		t.Fatalf("expected total 7000 for regular tour with price, got %d", booking.TotalPrice)
	}
	if booking.PaymentStatus != PaymentStatusUnpaid {
		t.Fatalf("expected UNPAID payment status, got %s", booking.PaymentStatus)
	}
	if booking.Status != BookingStatusNew {
		t.Fatalf("expected status NEW, got %s", booking.Status)
	}
}

func TestNewBookingZeroTotalHasNotRequiredPaymentStatus(t *testing.T) {
	tour, err := NewTour(validTourInput(func(input *NewTourInput) {
		input.IsRegular = true
		input.Price = 0
		input.SlotsLeft = 5
	}))
	if err != nil {
		t.Fatalf("create tour: %v", err)
	}

	booking, err := NewBooking(validBookingInput(tour, func(input *NewBookingInput) {
		input.PeopleCount = 2
	}))
	if err != nil {
		t.Fatalf("create booking: %v", err)
	}
	if booking.TotalPrice != 0 {
		t.Fatalf("expected total 0, got %d", booking.TotalPrice)
	}
	if booking.PaymentStatus != PaymentStatusNotRequired {
		t.Fatalf("expected NOT_REQUIRED payment status, got %s", booking.PaymentStatus)
	}
}

func TestNewBookingRejectsInsufficientSlots(t *testing.T) {
	tour, err := NewTour(validTourInput(func(input *NewTourInput) {
		input.SlotsLeft = 1
		input.OverbookingEnabled = false
	}))
	if err != nil {
		t.Fatalf("create tour: %v", err)
	}

	_, err = NewBooking(validBookingInput(tour, func(input *NewBookingInput) {
		input.PeopleCount = 2
	}))

	if !errors.Is(err, ErrInsufficientSlots) {
		t.Fatalf("expected insufficient slots, got %v", err)
	}
}

func TestNewBookingMarksOverbooked(t *testing.T) {
	tour, err := NewTour(validTourInput(func(input *NewTourInput) {
		input.SlotsLeft = 1
		input.OverbookingEnabled = true
	}))
	if err != nil {
		t.Fatalf("create tour: %v", err)
	}

	booking, err := NewBooking(validBookingInput(tour, func(input *NewBookingInput) {
		input.PeopleCount = 2
	}))
	if err != nil {
		t.Fatalf("create booking: %v", err)
	}

	if !booking.Overbooked {
		t.Fatal("expected booking to be marked overbooked")
	}
}

func TestBookingStatusTransitions(t *testing.T) {
	booking := Booking{Status: BookingStatusNew}

	for _, next := range []BookingStatus{
		BookingStatusContacted,
		BookingStatusConfirmed,
		BookingStatusCompleted,
	} {
		if err := booking.ChangeStatus(next); err != nil {
			t.Fatalf("change status to %s: %v", next, err)
		}
	}

	if !booking.IsTerminal() {
		t.Fatal("expected completed booking to be terminal")
	}
	if err := booking.ChangeStatus(BookingStatusCancelled); !errors.Is(err, ErrInvalidStatusTransition) {
		t.Fatalf("expected invalid transition from completed, got %v", err)
	}
}

func TestBookingCanCancelFromActiveStatuses(t *testing.T) {
	for _, status := range []BookingStatus{
		BookingStatusNew,
		BookingStatusContacted,
		BookingStatusConfirmed,
	} {
		booking := Booking{Status: status}
		if err := booking.ChangeStatus(BookingStatusCancelled); err != nil {
			t.Fatalf("expected cancellation from %s to be allowed: %v", status, err)
		}
	}
}

func TestBookingRejectsSkippedTransition(t *testing.T) {
	booking := Booking{Status: BookingStatusNew}

	err := booking.ChangeStatus(BookingStatusConfirmed)

	if !errors.Is(err, ErrInvalidStatusTransition) {
		t.Fatalf("expected invalid transition, got %v", err)
	}
}

func validBookingInput(tour Tour, mutators ...func(*NewBookingInput)) NewBookingInput {
	input := NewBookingInput{
		ID:          uuid.New(),
		Tour:        tour,
		Name:        "Иван Иванов",
		Phone:       "+79999999999",
		Email:       "mail@test.com",
		PeopleCount: 2,
		Source:      "web",
		Now:         time.Date(2026, 5, 20, 0, 0, 0, 0, time.UTC),
	}
	for _, mutate := range mutators {
		mutate(&input)
	}
	return input
}
