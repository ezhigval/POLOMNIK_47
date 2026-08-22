package application

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
	"polomnik/internal/adapters/integration/noop"
	"polomnik/internal/adapters/integration/recording"
	notificationnoop "polomnik/internal/adapters/notification/noop"
	"polomnik/internal/adapters/repository/memory"
	"polomnik/internal/config"
	"polomnik/internal/domain"
	"polomnik/internal/ports"
)

func TestOutboxWorkerMarksProcessedWhenSynced(t *testing.T) {
	t.Parallel()

	server := newBitrixSuccessServer(t)
	defer server.Close()

	ctx := context.Background()
	store := memory.NewStore()
	crm := recording.NewCRM(bitrix.NewCRMAdapter(config.Config{
		BitrixWebhookURL:       server.URL,
		IntegrationHTTPTimeout: 2 * time.Second,
	}), store, store, recording.CRMOptions{
		ExternalSystem: domain.ExternalSystemBitrix24,
	})

	booking := mustTestBooking(t, store)
	mustPendingOutboxEvent(t, store, domain.OutboxEventCRMPushBooking, domain.EntityTypeBooking, booking.ID)

	worker := NewOutboxWorker(store, store, store, store, crm, noop.NewAccountingAdapter(), notificationnoop.New(), 3)
	if _, err := worker.ProcessBatch(ctx, 10); err != nil {
		t.Fatalf("process batch: %v", err)
	}

	list, err := store.ListEvents(ctx, ports.OutboxFilters{Status: string(domain.OutboxStatusProcessed)}, ports.Pagination{Page: 1, Limit: 10})
	if err != nil {
		t.Fatalf("list processed: %v", err)
	}
	if len(list.Items) != 1 {
		t.Fatalf("expected processed event, got %d", len(list.Items))
	}
}

func TestOutboxWorkerMarksProcessedWhenNotConfigured(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	store := memory.NewStore()
	crm := recording.NewCRM(bitrix.NewCRMAdapter(config.Config{}), store, store, recording.CRMOptions{
		ExternalSystem: domain.ExternalSystemBitrix24,
	})

	booking := mustTestBooking(t, store)
	mustPendingOutboxEvent(t, store, domain.OutboxEventCRMPushBooking, domain.EntityTypeBooking, booking.ID)

	worker := NewOutboxWorker(store, store, store, store, crm, noop.NewAccountingAdapter(), notificationnoop.New(), 3)
	if _, err := worker.ProcessBatch(ctx, 10); err != nil {
		t.Fatalf("process batch: %v", err)
	}

	list, err := store.ListEvents(ctx, ports.OutboxFilters{Status: string(domain.OutboxStatusProcessed)}, ports.Pagination{Page: 1, Limit: 10})
	if err != nil {
		t.Fatalf("list processed: %v", err)
	}
	if len(list.Items) != 1 {
		t.Fatalf("expected processed event, got %d", len(list.Items))
	}
}

func TestOutboxWorkerMarksFailedAfterMaxAttempts(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
	}))
	defer server.Close()

	ctx := context.Background()
	store := memory.NewStore()
	crm := recording.NewCRM(bitrix.NewCRMAdapter(config.Config{
		BitrixWebhookURL:       server.URL,
		IntegrationHTTPTimeout: 2 * time.Second,
	}), store, store, recording.CRMOptions{
		ExternalSystem: domain.ExternalSystemBitrix24,
	})

	booking := mustTestBooking(t, store)
	mustPendingOutboxEvent(t, store, domain.OutboxEventCRMPushBooking, domain.EntityTypeBooking, booking.ID)

	worker := NewOutboxWorker(store, store, store, store, crm, noop.NewAccountingAdapter(), notificationnoop.New(), 2)
	for i := 0; i < 2; i++ {
		if _, err := worker.ProcessBatch(ctx, 10); err != nil {
			t.Fatalf("process batch attempt %d: %v", i+1, err)
		}
	}

	list, err := store.ListEvents(ctx, ports.OutboxFilters{Status: string(domain.OutboxStatusFailed)}, ports.Pagination{Page: 1, Limit: 10})
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if len(list.Items) != 1 {
		t.Fatalf("expected failed event, got %d", len(list.Items))
	}
}

func newBitrixSuccessServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/crm.duplicate.findbycomm"):
			_ = json.NewEncoder(w).Encode(map[string]any{"result": map[string][]any{"CONTACT": {}}})
		case strings.HasSuffix(r.URL.Path, "/crm.contact.add"):
			_ = json.NewEncoder(w).Encode(map[string]any{"result": 1})
		case strings.HasSuffix(r.URL.Path, "/crm.deal.list"):
			_ = json.NewEncoder(w).Encode(map[string]any{"result": []any{}})
		case strings.HasSuffix(r.URL.Path, "/crm.deal.add"):
			_ = json.NewEncoder(w).Encode(map[string]any{"result": 42})
		case strings.HasSuffix(r.URL.Path, "/crm.deal.update"):
			_ = json.NewEncoder(w).Encode(map[string]any{"result": true})
		case strings.HasSuffix(r.URL.Path, "/crm.product.list"):
			_ = json.NewEncoder(w).Encode(map[string]any{"result": []any{}})
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
}

func mustTestBooking(t *testing.T, store *memory.Store) domain.Booking {
	t.Helper()

	ctx := context.Background()
	tour, err := domain.NewTour(domain.NewTourInput{
		ID:         uuid.New(),
		Slug:       "outbox-worker-test",
		Title:      "Outbox Worker Test",
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
	if _, err := store.CreateTour(ctx, tour); err != nil {
		t.Fatalf("persist tour: %v", err)
	}

	booking, err := domain.NewBooking(domain.NewBookingInput{
		ID:          uuid.New(),
		Tour:        tour,
		Name:        "Worker Test",
		Phone:       "+79990000000",
		PeopleCount: 1,
		Source:      "test",
	})
	if err != nil {
		t.Fatalf("create booking: %v", err)
	}
	created, err := store.CreateBooking(ctx, booking)
	if err != nil {
		t.Fatalf("persist booking: %v", err)
	}
	return created
}

func mustPendingOutboxEvent(
	t *testing.T,
	store *memory.Store,
	eventType string,
	entityType string,
	entityID uuid.UUID,
) domain.OutboxEvent {
	t.Helper()

	ctx := context.Background()
	event, err := domain.NewOutboxEvent(domain.NewOutboxEventInput{
		ID:         uuid.New(),
		EventType:  eventType,
		EntityType: entityType,
		EntityID:   entityID,
		Status:     domain.OutboxStatusPending,
		Now:        time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("create outbox event: %v", err)
	}
	created, err := store.Enqueue(ctx, event)
	if err != nil {
		t.Fatalf("enqueue outbox event: %v", err)
	}
	return created
}
