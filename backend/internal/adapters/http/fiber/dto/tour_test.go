package dto

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"

	"palomnik/internal/domain"
)

func TestToTourResponseOmitsPriceWhenZero(t *testing.T) {
	tour, err := domain.NewTour(domain.NewTourInput{
		ID:         uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		Slug:       "regular",
		Title:      "Regular",
		Price:      0,
		Currency:   "RUB",
		SlotsTotal: 10,
		SlotsLeft:  10,
		IsActive:   true,
		IsRegular:  true,
	})
	if err != nil {
		t.Fatalf("create tour: %v", err)
	}

	raw, err := json.Marshal(ToTourResponse(tour))
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(raw, &payload); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if payload["is_regular"] != true {
		t.Fatalf("expected is_regular true, got %#v", payload["is_regular"])
	}
	if payload["price"] != nil {
		t.Fatalf("expected null price, got %#v", payload["price"])
	}
	if payload["date_start"] != nil || payload["date_end"] != nil {
		t.Fatalf("expected null dates for regular, got start=%#v end=%#v", payload["date_start"], payload["date_end"])
	}
}

func TestToTourResponseIncludesPriceForRegularWithPrice(t *testing.T) {
	tour, err := domain.NewTour(domain.NewTourInput{
		ID:         uuid.MustParse("44444444-4444-4444-4444-444444444444"),
		Slug:       "regular-priced",
		Title:      "Regular priced",
		Price:      5000,
		Currency:   "RUB",
		SlotsTotal: 10,
		SlotsLeft:  10,
		IsActive:   true,
		IsRegular:  true,
	})
	if err != nil {
		t.Fatalf("create tour: %v", err)
	}

	resp := ToTourResponse(tour)
	if resp.Price == nil || *resp.Price != 5000 {
		t.Fatalf("expected price 5000 for regular tour, got %#v", resp.Price)
	}
	if resp.DateStart != nil || resp.DateEnd != nil {
		t.Fatalf("expected null dates for regular, got start=%#v end=%#v", resp.DateStart, resp.DateEnd)
	}
}

func TestToTourResponseIncludesOverbookingFlag(t *testing.T) {
	tour, err := domain.NewTour(domain.NewTourInput{
		ID:                 uuid.MustParse("33333333-3333-3333-3333-333333333333"),
		Slug:               "overbook",
		Title:              "Overbook",
		Price:              10000,
		Currency:           "RUB",
		DateStart:          time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC),
		DateEnd:            time.Date(2026, 9, 3, 0, 0, 0, 0, time.UTC),
		SlotsTotal:         10,
		SlotsLeft:          0,
		IsActive:           true,
		OverbookingEnabled: true,
	})
	if err != nil {
		t.Fatalf("create tour: %v", err)
	}

	resp := ToTourResponse(tour)
	if !resp.OverbookingEnabled {
		t.Fatal("expected overbooking_enabled true on public tour response")
	}
}

func TestToPublicTourResponseBurningDiscount(t *testing.T) {
	today := time.Date(2026, 8, 26, 0, 0, 0, 0, time.UTC)
	tour, err := domain.NewTour(domain.NewTourInput{
		ID:                 uuid.MustParse("55555555-5555-5555-5555-555555555555"),
		Slug:               "burning",
		Title:              "Burning",
		Price:              10000,
		Currency:           "RUB",
		DateStart:          today,
		DateEnd:            time.Date(2026, 8, 28, 0, 0, 0, 0, time.UTC),
		SlotsTotal:         10,
		SlotsLeft:          10,
		IsActive:           true,
		HotDiscountPercent: 10,
	})
	if err != nil {
		t.Fatalf("create tour: %v", err)
	}

	resp := ToPublicTourResponse(tour, domain.TourCatalogContext{Today: today})
	if !resp.IsBurning {
		t.Fatal("expected is_burning true")
	}
	if resp.Price == nil || *resp.Price != 9000 {
		t.Fatalf("expected discounted price 9000, got %#v", resp.Price)
	}
	if resp.OriginalPrice == nil || *resp.OriginalPrice != 10000 {
		t.Fatalf("expected original price 10000, got %#v", resp.OriginalPrice)
	}
}

func TestToTourResponseKeepsPriceAndDatesForDated(t *testing.T) {
	tour, err := domain.NewTour(domain.NewTourInput{
		ID:         uuid.MustParse("22222222-2222-2222-2222-222222222222"),
		Slug:       "dated",
		Title:      "Dated",
		Price:      15000,
		Currency:   "RUB",
		DateStart:  time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC),
		DateEnd:    time.Date(2026, 9, 3, 0, 0, 0, 0, time.UTC),
		SlotsTotal: 10,
		SlotsLeft:  10,
		IsActive:   true,
	})
	if err != nil {
		t.Fatalf("create tour: %v", err)
	}

	resp := ToTourResponse(tour)
	if resp.IsRegular {
		t.Fatal("expected dated tour not regular")
	}
	if resp.Price == nil || *resp.Price != 15000 {
		t.Fatalf("expected price 15000, got %#v", resp.Price)
	}
	if resp.DateStart == nil || *resp.DateStart != "2026-09-01" {
		t.Fatalf("expected date_start 2026-09-01, got %#v", resp.DateStart)
	}
	if resp.DateEnd == nil || *resp.DateEnd != "2026-09-03" {
		t.Fatalf("expected date_end 2026-09-03, got %#v", resp.DateEnd)
	}
}
