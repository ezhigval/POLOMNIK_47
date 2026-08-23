package memory

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"

	"palomnik/internal/domain"
	"palomnik/internal/ports"
)

func TestStoreEnqueueOutboxEvent(t *testing.T) {
	ctx := context.Background()
	store := NewStore()

	event, err := domain.NewOutboxEvent(domain.NewOutboxEventInput{
		ID:         uuid.New(),
		EventType:  domain.OutboxEventCRMPushBooking,
		EntityType: domain.EntityTypeBooking,
		EntityID:   uuid.New(),
		Payload:    json.RawMessage(`{"integration_status":"pending"}`),
		Status:     domain.OutboxStatusPending,
		Now:        time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("create outbox event: %v", err)
	}

	if _, err := store.Enqueue(ctx, event); err != nil {
		t.Fatalf("enqueue outbox event: %v", err)
	}

	dup, err := domain.NewOutboxEvent(domain.NewOutboxEventInput{
		ID:         uuid.New(),
		EventType:  domain.OutboxEventCRMPushBooking,
		EntityType: domain.EntityTypeBooking,
		EntityID:   event.EntityID,
		Payload:    json.RawMessage(`{"integration_status":"failed"}`),
		Status:     domain.OutboxStatusPending,
		Now:        time.Date(2026, 6, 10, 13, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("create duplicate outbox event: %v", err)
	}
	if _, err := store.Enqueue(ctx, dup); err != nil {
		t.Fatalf("enqueue duplicate outbox event: %v", err)
	}

	list, err := store.ListEvents(ctx, ports.OutboxFilters{}, ports.Pagination{Page: 1, Limit: 10})
	if err != nil {
		t.Fatalf("list outbox events: %v", err)
	}
	if len(list.Items) != 1 {
		t.Fatalf("expected 1 outbox event, got %d", len(list.Items))
	}
}

func TestStoreSummarizeOutbox(t *testing.T) {
	ctx := context.Background()
	store := NewStore()
	now := time.Date(2026, 7, 15, 12, 0, 0, 0, time.UTC)

	pendingOlder, err := domain.NewOutboxEvent(domain.NewOutboxEventInput{
		ID:         uuid.New(),
		EventType:  domain.OutboxEventCRMPushBooking,
		EntityType: domain.EntityTypeBooking,
		EntityID:   uuid.New(),
		Status:     domain.OutboxStatusPending,
		Now:        now.Add(-time.Hour),
	})
	if err != nil {
		t.Fatalf("create pending: %v", err)
	}
	pendingNewer, err := domain.NewOutboxEvent(domain.NewOutboxEventInput{
		ID:         uuid.New(),
		EventType:  domain.OutboxEventNotificationBookingCreated,
		EntityType: domain.EntityTypeBooking,
		EntityID:   uuid.New(),
		Status:     domain.OutboxStatusPending,
		Now:        now,
	})
	if err != nil {
		t.Fatalf("create pending newer: %v", err)
	}
	failed, err := domain.NewOutboxEvent(domain.NewOutboxEventInput{
		ID:         uuid.New(),
		EventType:  domain.OutboxEventAccountingExport,
		EntityType: domain.EntityTypeBooking,
		EntityID:   uuid.New(),
		Status:     domain.OutboxStatusFailed,
		LastError:  "upstream timeout",
		Now:        now.Add(-30 * time.Minute),
	})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	failed.UpdatedAt = now.Add(-10 * time.Minute)

	processed, err := domain.NewOutboxEvent(domain.NewOutboxEventInput{
		ID:         uuid.New(),
		EventType:  domain.OutboxEventCRMSyncTour,
		EntityType: domain.EntityTypeTour,
		EntityID:   uuid.New(),
		Status:     domain.OutboxStatusProcessed,
		Now:        now.Add(-2 * time.Hour),
	})
	if err != nil {
		t.Fatalf("create processed: %v", err)
	}

	for _, event := range []domain.OutboxEvent{pendingOlder, pendingNewer, failed, processed} {
		if _, err := store.Enqueue(ctx, event); err != nil {
			t.Fatalf("enqueue: %v", err)
		}
	}

	summary, err := store.Summarize(ctx)
	if err != nil {
		t.Fatalf("summarize: %v", err)
	}
	if summary.Pending != 2 || summary.Failed != 1 || summary.Processed != 1 {
		t.Fatalf("unexpected summary counts: %+v", summary)
	}
	if summary.OldestPendingAt == nil || !summary.OldestPendingAt.Equal(pendingOlder.CreatedAt) {
		t.Fatalf("unexpected oldest pending: %v", summary.OldestPendingAt)
	}
	if summary.LatestFailedError != "upstream timeout" {
		t.Fatalf("unexpected latest failed error: %q", summary.LatestFailedError)
	}
}
