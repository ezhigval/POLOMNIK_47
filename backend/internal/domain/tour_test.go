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

func TestNewTourRegularAllowsMissingDates(t *testing.T) {
	tour, err := NewTour(validTourInput(func(input *NewTourInput) {
		input.IsRegular = true
		input.Price = 0
		input.DateStart = time.Time{}
		input.DateEnd = time.Time{}
	}))
	if err != nil {
		t.Fatalf("create regular tour: %v", err)
	}
	if !tour.IsRegular {
		t.Fatal("expected is_regular")
	}
	if !tour.DateStart.IsZero() || !tour.DateEnd.IsZero() {
		t.Fatalf("expected empty dates, got %v – %v", tour.DateStart, tour.DateEnd)
	}
}

func TestNewTourRegularClearsProvidedDates(t *testing.T) {
	tour, err := NewTour(validTourInput(func(input *NewTourInput) {
		input.IsRegular = true
	}))
	if err != nil {
		t.Fatalf("create regular tour: %v", err)
	}
	if !tour.DateStart.IsZero() || !tour.DateEnd.IsZero() {
		t.Fatal("regular tour must not keep a schedule")
	}
	if tour.HasPublicSchedule() {
		t.Fatal("regular tour must not have a public schedule")
	}
}

func TestNewTourStillRequiresDatesWhenNotRegular(t *testing.T) {
	_, err := NewTour(validTourInput(func(input *NewTourInput) {
		input.DateStart = time.Time{}
		input.DateEnd = time.Time{}
	}))
	if !errors.Is(err, ErrInvalidDateRange) {
		t.Fatalf("expected invalid date range, got %v", err)
	}
}

func TestTourBookingTotalZeroWhenRegular(t *testing.T) {
	tour, err := NewTour(validTourInput(func(input *NewTourInput) {
		input.IsRegular = true
		input.Price = 3500
	}))
	if err != nil {
		t.Fatalf("create regular tour: %v", err)
	}
	if got := tour.BookingTotal(2); got != 0 {
		t.Fatalf("expected booking total 0 for regular tour, got %d", got)
	}
}

func TestCompareToursByCatalogPutsRegularAfterDated(t *testing.T) {
	datedLater, err := NewTour(validTourInput(func(input *NewTourInput) {
		input.ID = uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
		input.Slug = "dated-later"
		input.DateStart = time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)
		input.DateEnd = time.Date(2026, 10, 3, 0, 0, 0, 0, time.UTC)
	}))
	if err != nil {
		t.Fatalf("dated later: %v", err)
	}
	datedSooner, err := NewTour(validTourInput(func(input *NewTourInput) {
		input.ID = uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")
		input.Slug = "dated-sooner"
		input.DateStart = time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
		input.DateEnd = time.Date(2026, 9, 3, 0, 0, 0, 0, time.UTC)
	}))
	if err != nil {
		t.Fatalf("dated sooner: %v", err)
	}
	regular, err := NewTour(validTourInput(func(input *NewTourInput) {
		input.ID = uuid.MustParse("00000000-0000-0000-0000-000000000001")
		input.Slug = "regular"
		input.IsRegular = true
	}))
	if err != nil {
		t.Fatalf("regular: %v", err)
	}

	if CompareToursByCatalog(datedSooner, datedLater) >= 0 {
		t.Fatal("expected sooner dated tour before later dated tour")
	}
	if CompareToursByCatalog(datedLater, regular) >= 0 {
		t.Fatal("expected dated tour before regular tour")
	}
	if CompareToursByCatalog(regular, datedSooner) <= 0 {
		t.Fatal("expected regular tour after dated tour")
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
