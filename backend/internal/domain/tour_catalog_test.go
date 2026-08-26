package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestTourIsBurningOnWindow(t *testing.T) {
	start := time.Date(2026, 8, 26, 0, 0, 0, 0, time.UTC)
	tour, err := NewTour(validTourInput(func(input *NewTourInput) {
		input.DateStart = start
		input.DateEnd = time.Date(2026, 8, 28, 0, 0, 0, 0, time.UTC)
	}))
	if err != nil {
		t.Fatalf("create tour: %v", err)
	}

	windowStart := start.AddDate(0, 0, -BurningWindowDays)
	if !tour.IsBurningOn(windowStart) {
		t.Fatal("expected burning on first day of window")
	}
	if !tour.IsBurningOn(start) {
		t.Fatal("expected burning on departure day")
	}
	if !tour.IsBurningOn(time.Date(2026, 8, 26, 23, 59, 0, 0, time.UTC)) {
		t.Fatal("expected burning late on departure day")
	}
	if tour.IsBurningOn(windowStart.AddDate(0, 0, -1)) {
		t.Fatal("expected not burning before window")
	}
	if tour.IsBurningOn(start.AddDate(0, 0, 1)) {
		t.Fatal("expected not burning after departure day even when tour still runs")
	}
}

func TestRegularTourNeverBurning(t *testing.T) {
	tour, err := NewTour(validTourInput(func(input *NewTourInput) {
		input.IsRegular = true
	}))
	if err != nil {
		t.Fatalf("create tour: %v", err)
	}
	if tour.IsBurningOn(time.Now().UTC()) {
		t.Fatal("regular tour must not be burning")
	}
}

func TestApplyPercentDiscount(t *testing.T) {
	if got := ApplyPercentDiscount(10000, 10); got != 9000 {
		t.Fatalf("expected 9000, got %d", got)
	}
	if got := ApplyPercentDiscount(10000, 0); got != 10000 {
		t.Fatalf("expected unchanged price, got %d", got)
	}
}

func TestTourUnitPriceInBurningWithDiscount(t *testing.T) {
	today := time.Date(2026, 8, 26, 0, 0, 0, 0, time.UTC)
	tour, err := NewTour(validTourInput(func(input *NewTourInput) {
		input.Price = 10000
		input.HotDiscountPercent = 15
		input.DateStart = today
		input.DateEnd = time.Date(2026, 8, 28, 0, 0, 0, 0, time.UTC)
	}))
	if err != nil {
		t.Fatalf("create tour: %v", err)
	}
	catalog := TourCatalogContext{Today: today}
	if got := tour.UnitPriceIn(catalog); got != 8500 {
		t.Fatalf("expected discounted unit price 8500, got %d", got)
	}
	if got := tour.BookingTotalIn(2, catalog); got != 17000 {
		t.Fatalf("expected booking total 17000, got %d", got)
	}
}

func TestTourBurningWithoutDiscountKeepsPrice(t *testing.T) {
	today := time.Date(2026, 8, 26, 0, 0, 0, 0, time.UTC)
	tour, err := NewTour(validTourInput(func(input *NewTourInput) {
		input.ID = uuid.New()
		input.Price = 12000
		input.DateStart = today
		input.DateEnd = time.Date(2026, 8, 27, 0, 0, 0, 0, time.UTC)
	}))
	if err != nil {
		t.Fatalf("create tour: %v", err)
	}
	catalog := TourCatalogContext{Today: today}
	if !tour.IsBurningIn(catalog) {
		t.Fatal("expected burning flag without configured discount")
	}
	if got := tour.UnitPriceIn(catalog); got != 12000 {
		t.Fatalf("expected full price when discount percent is 0, got %d", got)
	}
}

func TestScheduleEndedOnSkipsRegular(t *testing.T) {
	tour, err := NewTour(validTourInput(func(input *NewTourInput) {
		input.IsRegular = true
	}))
	if err != nil {
		t.Fatalf("create tour: %v", err)
	}
	if tour.ScheduleEndedOn(time.Date(2099, 1, 1, 0, 0, 0, 0, time.UTC)) {
		t.Fatal("regular tour schedule must not end")
	}
}
