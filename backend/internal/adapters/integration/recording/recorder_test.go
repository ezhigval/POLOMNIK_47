package recording

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"

	"polomnik/internal/adapters/integration/bitrix"
	"polomnik/internal/adapters/repository/memory"
	"polomnik/internal/config"
	"polomnik/internal/domain"
	"polomnik/internal/ports"
)

func TestRecordingCRMPersistsBookingSyncReference(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	store := memory.NewStore()
	adapter := NewCRM(bitrix.NewCRMAdapter(config.Config{}), store, store, CRMOptions{
		ExternalSystem: domain.ExternalSystemBitrix24,
	})

	booking := mustBooking(t)
	result, err := adapter.PushBooking(ctx, booking)
	if err != nil {
		t.Fatalf("push booking: %v", err)
	}
	if result.Status != ports.IntegrationStatusNotConfigured {
		t.Fatalf("expected not_configured, got %s", result.Status)
	}

	ref, err := store.GetReference(ctx, ports.IntegrationReferenceQuery{
		LocalEntityType:    domain.EntityTypeBooking,
		LocalEntityID:      booking.ID,
		ExternalSystem:     domain.ExternalSystemBitrix24,
		ExternalEntityType: domain.ExternalEntityTypeDeal,
	})
	if err != nil {
		t.Fatalf("get reference: %v", err)
	}
	if ref.SyncStatus != domain.SyncStatusNotConfigured {
		t.Fatalf("expected stored not_configured, got %s", ref.SyncStatus)
	}
}

func TestRecordingCRMDoesNotEnqueueOutboxOnSynced(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/crm.duplicate.findbycomm"):
			_ = json.NewEncoder(w).Encode(map[string]any{"result": map[string][]any{"CONTACT": {}}})
		case strings.HasSuffix(r.URL.Path, "/crm.contact.add"):
			_ = json.NewEncoder(w).Encode(map[string]any{"result": 1})
		case strings.HasSuffix(r.URL.Path, "/crm.deal.list"):
			_ = json.NewEncoder(w).Encode(map[string]any{"result": []any{}})
		case strings.HasSuffix(r.URL.Path, "/crm.deal.add"):
			_ = json.NewEncoder(w).Encode(map[string]any{"result": 42})
		case strings.HasSuffix(r.URL.Path, "/crm.product.list"):
			_ = json.NewEncoder(w).Encode(map[string]any{"result": []any{}})
		case strings.HasSuffix(r.URL.Path, "/crm.deal.productrows.set"):
			_ = json.NewEncoder(w).Encode(map[string]any{"result": true})
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	ctx := context.Background()
	store := memory.NewStore()
	adapter := NewCRM(bitrix.NewCRMAdapter(config.Config{
		BitrixWebhookURL:       server.URL,
		IntegrationHTTPTimeout: 2 * time.Second,
	}), store, store, CRMOptions{
		ExternalSystem: domain.ExternalSystemBitrix24,
	})

	booking := mustBooking(t)
	if _, err := adapter.PushBooking(ctx, booking); err != nil {
		t.Fatalf("push booking: %v", err)
	}

	list, err := store.ListEvents(ctx, ports.OutboxFilters{}, ports.Pagination{Page: 1, Limit: 10})
	if err != nil {
		t.Fatalf("list outbox: %v", err)
	}
	if len(list.Items) != 0 {
		t.Fatalf("expected no outbox events on synced, got %d", len(list.Items))
	}
}

func TestRecordingCRMEnqueuesOutboxOnFailed(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"INTERNAL","error_description":"boom"}`))
	}))
	defer server.Close()

	ctx := context.Background()
	store := memory.NewStore()
	adapter := NewCRM(bitrix.NewCRMAdapter(config.Config{
		BitrixWebhookURL:       server.URL,
		IntegrationHTTPTimeout: 2 * time.Second,
	}), store, store, CRMOptions{
		ExternalSystem: domain.ExternalSystemBitrix24,
	})

	booking := mustBooking(t)
	if _, err := adapter.PushBooking(ctx, booking); err == nil {
		t.Fatal("expected push booking error")
	}

	list, err := store.ListEvents(ctx, ports.OutboxFilters{
		EventType: domain.OutboxEventCRMPushBooking,
	}, ports.Pagination{Page: 1, Limit: 10})
	if err != nil {
		t.Fatalf("list outbox: %v", err)
	}
	if len(list.Items) != 1 {
		t.Fatalf("expected 1 outbox event, got %d", len(list.Items))
	}
	if list.Items[0].Status != domain.OutboxStatusPending {
		t.Fatalf("expected pending outbox status, got %s", list.Items[0].Status)
	}
}

func mustBooking(t *testing.T) domain.Booking {
	t.Helper()

	tour, err := domain.NewTour(domain.NewTourInput{
		ID:         uuid.New(),
		Slug:       "recording-test",
		Title:      "Recording Test",
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
