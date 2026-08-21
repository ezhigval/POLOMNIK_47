package domain

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	OutboxEventCRMPushBooking         = "crm.push_booking"
	OutboxEventCRMUpdateBookingStatus = "crm.update_booking_status"
	OutboxEventCRMSyncTour            = "crm.sync_tour"
	OutboxEventAccountingExport       = "accounting.export_booking"
	OutboxEventAccountingCounterparty = "accounting.sync_counterparty"
	OutboxEventCRMPushReview          = "crm.push_review"
	OutboxEventNotificationBookingCreated = "notification.booking.created"
	OutboxEventNotificationBookingStatus  = "notification.booking.status_changed"
)

type OutboxStatus string

const (
	OutboxStatusPending   OutboxStatus = "pending"
	OutboxStatusProcessed OutboxStatus = "processed"
	OutboxStatusFailed    OutboxStatus = "failed"
)

type OutboxEvent struct {
	ID         uuid.UUID
	EventType  string
	EntityType string
	EntityID   uuid.UUID
	Payload    json.RawMessage
	Status     OutboxStatus
	Attempts   int
	LastError  string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type NewOutboxEventInput struct {
	ID         uuid.UUID
	EventType  string
	EntityType string
	EntityID   uuid.UUID
	Payload    json.RawMessage
	Status     OutboxStatus
	LastError  string
	Now        time.Time
}

func NewOutboxEvent(input NewOutboxEventInput) (OutboxEvent, error) {
	if input.ID == uuid.Nil || input.EntityID == uuid.Nil {
		return OutboxEvent{}, ErrInvalidID
	}
	if strings.TrimSpace(input.EventType) == "" {
		return OutboxEvent{}, ErrInvalidEventType
	}
	if strings.TrimSpace(input.EntityType) == "" {
		return OutboxEvent{}, ErrInvalidEntityType
	}
	if !isValidOutboxStatus(input.Status) {
		return OutboxEvent{}, ErrInvalidOutboxStatus
	}

	payload := input.Payload
	if len(payload) == 0 {
		payload = json.RawMessage(`{}`)
	}

	now := input.Now
	if now.IsZero() {
		now = time.Now().UTC()
	}

	return OutboxEvent{
		ID:         input.ID,
		EventType:  strings.TrimSpace(input.EventType),
		EntityType: strings.TrimSpace(input.EntityType),
		EntityID:   input.EntityID,
		Payload:    payload,
		Status:     input.Status,
		Attempts:   0,
		LastError:  strings.TrimSpace(input.LastError),
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}

func isValidOutboxStatus(status OutboxStatus) bool {
	switch status {
	case OutboxStatusPending, OutboxStatusProcessed, OutboxStatusFailed:
		return true
	default:
		return false
	}
}
