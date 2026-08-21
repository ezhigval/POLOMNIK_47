package domain

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewOutboxEvent(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC)
	entityID := uuid.New()

	event, err := NewOutboxEvent(NewOutboxEventInput{
		ID:         uuid.New(),
		EventType:  OutboxEventCRMPushBooking,
		EntityType: EntityTypeBooking,
		EntityID:   entityID,
		Payload:    json.RawMessage(`{"integration_status":"pending"}`),
		Status:     OutboxStatusPending,
		Now:        now,
	})
	if err != nil {
		t.Fatalf("create outbox event: %v", err)
	}

	if event.EventType != OutboxEventCRMPushBooking {
		t.Fatalf("unexpected event type: %s", event.EventType)
	}
	if event.EntityID != entityID {
		t.Fatalf("unexpected entity id: %s", event.EntityID)
	}
	if event.Status != OutboxStatusPending {
		t.Fatalf("unexpected status: %s", event.Status)
	}
	if !event.CreatedAt.Equal(now) {
		t.Fatalf("unexpected created_at: %s", event.CreatedAt)
	}
}

func TestNewOutboxEventValidation(t *testing.T) {
	t.Parallel()

	_, err := NewOutboxEvent(NewOutboxEventInput{
		EventType:  OutboxEventCRMPushBooking,
		EntityType: EntityTypeBooking,
		EntityID:   uuid.New(),
		Payload:    json.RawMessage(`{}`),
		Status:     OutboxStatusPending,
		Now:        time.Now(),
	})
	if err == nil {
		t.Fatal("expected error for missing id")
	}
}
