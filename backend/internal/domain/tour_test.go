package domain

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewTourValidatesSlots(t *testing.T) {
	_, err := NewTour(validTourInput(func(input *NewTourInput) {
		input.SlotsTotal = 10
		input.SlotsLeft = 11
	}))

	if !errors.Is(err, ErrInvalidSlots) {
		t.Fatalf("expected invalid slots, got %v", err)
	}
}

func TestTourReserveAndReleaseSlots(t *testing.T) {
	tour, err := NewTour(validTourInput(func(input *NewTourInput) {
		input.SlotsTotal = 10
		input.SlotsLeft = 5
	}))
	if err != nil {
		t.Fatalf("create tour: %v", err)
	}

	if err := tour.ReserveSlots(3); err != nil {
		t.Fatalf("reserve slots: %v", err)
	}
	if tour.SlotsLeft != 2 {
		t.Fatalf("expected 2 slots left, got %d", tour.SlotsLeft)
	}

	if err := tour.ReleaseSlots(20); err != nil {
		t.Fatalf("release slots: %v", err)
	}
	if tour.SlotsLeft != 10 {
		t.Fatalf("expected slots left capped at total 10, got %d", tour.SlotsLeft)
	}
}

func TestTourOverbookingAllowsReservationWithoutNegativeSlots(t *testing.T) {
	tour, err := NewTour(validTourInput(func(input *NewTourInput) {
		input.SlotsTotal = 3
		input.SlotsLeft = 1
		input.OverbookingEnabled = true
	}))
	if err != nil {
		t.Fatalf("create tour: %v", err)
	}

	if !tour.BookingWouldBeOverbooked(2) {
		t.Fatal("expected booking to be overbooked")
	}
	if err := tour.ReserveSlots(2); err != nil {
		t.Fatalf("reserve overbooked slots: %v", err)
	}
	if tour.SlotsLeft != 1 {
		t.Fatalf("expected overbooking reservation not to make slots negative, got %d", tour.SlotsLeft)
	}
}

func TestNewTourEmptyImagesNotNil(t *testing.T) {
	tour, err := NewTour(validTourInput())
	if err != nil {
		t.Fatalf("create tour: %v", err)
	}
	if tour.Images == nil {
		t.Fatal("expected non-nil empty images slice")
	}
	if len(tour.Images) != 0 {
		t.Fatalf("expected empty images, got %v", tour.Images)
	}
}

func validTourInput(mutators ...func(*NewTourInput)) NewTourInput {
	input := NewTourInput{
		ID:         uuid.New(),
		Slug:       "optina-pustyn",
		Title:      "Оптина пустынь",
		Price:      15000,
		Currency:   "RUB",
		DateStart:  time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC),
		DateEnd:    time.Date(2026, 6, 5, 0, 0, 0, 0, time.UTC),
		SlotsTotal: 20,
		SlotsLeft:  10,
		IsActive:   true,
		Now:        time.Date(2026, 5, 20, 0, 0, 0, 0, time.UTC),
	}
	for _, mutate := range mutators {
		mutate(&input)
	}
	return input
}
