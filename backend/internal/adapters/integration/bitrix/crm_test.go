package bitrix

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"

	"polomnik/internal/config"
	"polomnik/internal/domain"
	"polomnik/internal/ports"
)

func TestCRMAdapterReturnsNotConfiguredWithoutWebhook(t *testing.T) {
	t.Parallel()

	adapter := NewCRMAdapter(config.Config{})
	result, err := adapter.PushBooking(t.Context(), mustBooking(t))
	if err != nil {
		t.Fatalf("push booking: %v", err)
	}
	if result.Status != ports.IntegrationStatusNotConfigured {
		t.Fatalf("expected not_configured, got %s", result.Status)
	}
}

func TestCRMAdapterPushBookingCreatesDeal(t *testing.T) {
	t.Parallel()

	booking := mustBooking(t)
	server := newBitrixTestServer(t, bitrixHandlers{
		onDuplicate: func(w http.ResponseWriter, _ *http.Request) {
			writeBitrixResult(w, map[string][]any{"CONTACT": {}})
		},
		onContactAdd: func(w http.ResponseWriter, _ *http.Request) {
			writeBitrixResult(w, 501)
		},
		onDealList: func(w http.ResponseWriter, _ *http.Request) {
			writeBitrixResult(w, []any{})
		},
		onDealAdd: func(w http.ResponseWriter, _ *http.Request) {
			writeBitrixResult(w, 9001)
		},
		onProductList: func(w http.ResponseWriter, _ *http.Request) {
			writeBitrixResult(w, []any{})
		},
	})
	defer server.Close()

	adapter := NewCRMAdapter(config.Config{
		BitrixWebhookURL:       server.URL,
		IntegrationHTTPTimeout: 2 * time.Second,
	})
	result, err := adapter.PushBooking(t.Context(), booking)
	if err != nil {
		t.Fatalf("push booking: %v", err)
	}
	if result.Status != ports.IntegrationStatusSynced {
		t.Fatalf("expected synced, got %s", result.Status)
	}
	if result.ExternalID != "9001" {
		t.Fatalf("expected deal id 9001, got %q", result.ExternalID)
	}
}

func TestCRMAdapterUpdateBookingStatusUpdatesExistingDeal(t *testing.T) {
	t.Parallel()

	booking := mustBooking(t)
	booking.Status = domain.BookingStatusContacted

	server := newBitrixTestServer(t, bitrixHandlers{
		onDealList: func(w http.ResponseWriter, _ *http.Request) {
			writeBitrixResult(w, []map[string]any{{"ID": "777"}})
		},
		onDealUpdate: func(w http.ResponseWriter, _ *http.Request) {
			writeBitrixResult(w, true)
		},
	})
	defer server.Close()

	adapter := NewCRMAdapter(config.Config{
		BitrixWebhookURL:       server.URL,
		BitrixStageContacted:   "UC_CONTACTED",
		IntegrationHTTPTimeout: 2 * time.Second,
	})
	result, err := adapter.UpdateBookingStatus(t.Context(), booking)
	if err != nil {
		t.Fatalf("update booking status: %v", err)
	}
	if result.Status != ports.IntegrationStatusSynced {
		t.Fatalf("expected synced, got %s", result.Status)
	}
	if result.ExternalID != "777" {
		t.Fatalf("expected deal id 777, got %q", result.ExternalID)
	}
}

func TestCRMAdapterSyncTourCreatesProduct(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/crm.product.list"):
			writeBitrixResult(w, []any{})
		case strings.HasSuffix(r.URL.Path, "/crm.product.add"):
			writeBitrixResult(w, 501)
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	adapter := NewCRMAdapter(config.Config{
		BitrixWebhookURL:       server.URL,
		IntegrationHTTPTimeout: 2 * time.Second,
	})
	tour, err := domain.NewTour(domain.NewTourInput{
		ID:         uuid.New(),
		Slug:       "sync-tour",
		Title:      "Tour",
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

	result, err := adapter.SyncTour(t.Context(), tour)
	if err != nil {
		t.Fatalf("sync tour: %v", err)
	}
	if result.Status != ports.IntegrationStatusSynced {
		t.Fatalf("expected synced, got %s", result.Status)
	}
	if result.ExternalID != "501" {
		t.Fatalf("expected product id 501, got %q", result.ExternalID)
	}
}

type bitrixHandlers struct {
	onDuplicate       http.HandlerFunc
	onContactAdd      http.HandlerFunc
	onDealList        http.HandlerFunc
	onDealAdd         http.HandlerFunc
	onDealUpdate      http.HandlerFunc
	onProductList     http.HandlerFunc
	onProductRowsSet  http.HandlerFunc
}

func newBitrixTestServer(t *testing.T, handlers bitrixHandlers) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/crm.duplicate.findbycomm"):
			handlers.onDuplicate(w, r)
		case strings.HasSuffix(r.URL.Path, "/crm.contact.add"):
			handlers.onContactAdd(w, r)
		case strings.HasSuffix(r.URL.Path, "/crm.deal.list"):
			handlers.onDealList(w, r)
		case strings.HasSuffix(r.URL.Path, "/crm.deal.add"):
			handlers.onDealAdd(w, r)
		case strings.HasSuffix(r.URL.Path, "/crm.deal.update"):
			handlers.onDealUpdate(w, r)
		case strings.HasSuffix(r.URL.Path, "/crm.product.list"):
			if handlers.onProductList != nil {
				handlers.onProductList(w, r)
			} else {
				writeBitrixResult(w, []any{})
			}
		case strings.HasSuffix(r.URL.Path, "/crm.deal.productrows.set"):
			if handlers.onProductRowsSet != nil {
				handlers.onProductRowsSet(w, r)
			} else {
				writeBitrixResult(w, true)
			}
		default:
			t.Fatalf("unexpected bitrix path: %s", r.URL.Path)
		}
	}))
}

func writeBitrixResult(w http.ResponseWriter, result any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"result": result})
}

func mustBooking(t *testing.T) domain.Booking {
	t.Helper()

	tour, err := domain.NewTour(domain.NewTourInput{
		ID:         uuid.New(),
		Slug:       "bitrix-test",
		Title:      "Bitrix Test",
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
