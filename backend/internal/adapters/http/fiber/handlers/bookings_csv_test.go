package handlers

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"

	"palomnik/internal/domain"
)

func TestManagementBookingsCSVIncludesExistingFields(t *testing.T) {
	id := uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	tourID := uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")
	created := time.Date(2026, 8, 23, 12, 0, 0, 0, time.UTC)
	body, err := managementBookingsCSV([]domain.Booking{{
		ID:          id,
		TourID:      tourID,
		Name:        "Иван",
		Phone:       "+7999",
		Email:       "a@b.c",
		PeopleCount: 2,
		Status:        domain.BookingStatusNew,
		PaymentStatus: domain.PaymentStatusUnpaid,
		TotalPrice:    15000,
		Comment:     "нужен трансфер",
		Source:      "web",
		CreatedAt:   created,
	}})
	if err != nil {
		t.Fatal(err)
	}
	text := string(body)
	if !strings.Contains(text, "id,created_at,status,payment_status,name,phone,email,tour_id") {
		t.Fatalf("header: %q", text)
	}
	if !strings.Contains(text, "Иван") || !strings.Contains(text, "NEW") || !strings.Contains(text, "15000") {
		t.Fatalf("row: %q", text)
	}
}
