package dto

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"

	"palomnik/internal/domain"
)

func TestToTourResponseOmitsPriceAndDatesForRegular(t *testing.T) {
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
		t.Fatalf("expected null dates, got start=%#v end=%#v", payload["date_start"], payload["date_end"])
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
