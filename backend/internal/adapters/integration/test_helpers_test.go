package integration

import (
	"testing"
	"time"

	"github.com/google/uuid"

	"polomnik/internal/domain"
)

func mustBooking(t *testing.T) domain.Booking {
	t.Helper()

	tour, err := domain.NewTour(domain.NewTourInput{
		ID:         uuid.New(),
		Slug:       "test-tour",
		Title:      "Test Tour",
		Price:      1000,
		Currency:   "RUB",
		DateStart:  time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC),
		DateEnd:    time.Date(2026, 6, 2, 0, 0, 0, 0, time.UTC),
		SlotsTotal: 10,
		SlotsLeft:  10,
		IsActive:   true,
	})
	if err != nil {
		t.Fatalf("create tour: %v", err)
	}

	booking, err := domain.NewBooking(domain.NewBookingInput{
		ID:          uuid.New(),
		Tour:        tour,
		Name:        "Test User",
		Phone:       "+79990000000",
		PeopleCount: 1,
		Source:      "test",
	})
	if err != nil {
		t.Fatalf("create booking: %v", err)
	}

	return booking
}
