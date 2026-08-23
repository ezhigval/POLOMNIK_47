package onec

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"

	"palomnik/internal/config"
	"palomnik/internal/domain"
	"palomnik/internal/ports"
)

func TestAccountingAdapterReturnsNotConfiguredWithoutBaseURL(t *testing.T) {
	t.Parallel()

	adapter := NewAccountingAdapter(config.Config{})
	result, err := adapter.ExportBooking(t.Context(), mustBooking(t))
	if err != nil {
		t.Fatalf("export booking: %v", err)
	}
	if result.Status != ports.IntegrationStatusNotConfigured {
		t.Fatalf("expected not_configured, got %s", result.Status)
	}
}

func TestAccountingAdapterExportBooking(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/hs/palomnik/booking" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var payload bookingPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload.Phone == "" {
			t.Fatal("expected phone in payload")
		}
		_ = json.NewEncoder(w).Encode(exchangeResponse{
			ExternalID: "DOC-1001",
			Status:     "ok",
			Message:    "accepted",
		})
	}))
	defer server.Close()

	adapter := NewAccountingAdapter(config.Config{
		OneCBaseURL:            server.URL,
		IntegrationHTTPTimeout: 2 * time.Second,
	})
	booking := mustBooking(t)
	booking.Status = domain.BookingStatusConfirmed
	result, err := adapter.ExportBooking(t.Context(), booking)
	if err != nil {
		t.Fatalf("export booking: %v", err)
	}
	if result.Status != ports.IntegrationStatusSynced {
		t.Fatalf("expected synced, got %s", result.Status)
	}
	if result.ExternalID != "DOC-1001" {
		t.Fatalf("expected external id DOC-1001, got %q", result.ExternalID)
	}
}

func TestAccountingAdapterSyncCounterparty(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/hs/palomnik/counterparty" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(exchangeResponse{
			ExternalID: "CP-42",
			Status:     "ok",
		})
	}))
	defer server.Close()

	adapter := NewAccountingAdapter(config.Config{
		OneCBaseURL:            server.URL,
		IntegrationHTTPTimeout: 2 * time.Second,
	})
	result, err := adapter.SyncCounterparty(t.Context(), mustBooking(t))
	if err != nil {
		t.Fatalf("sync counterparty: %v", err)
	}
	if result.ExternalID != "CP-42" {
		t.Fatalf("expected CP-42, got %q", result.ExternalID)
	}
}

func mustBooking(t *testing.T) domain.Booking {
	t.Helper()

	tour, err := domain.NewTour(domain.NewTourInput{
		ID:         uuid.New(),
		Slug:       "onec-test",
		Title:      "1C Test",
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
