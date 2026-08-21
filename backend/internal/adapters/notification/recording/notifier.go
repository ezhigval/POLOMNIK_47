package recording

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"

	"polomnik/internal/domain"
	"polomnik/internal/ports"
)

type Notifier struct {
	inner  ports.NotificationPort
	outbox ports.OutboxRepository
}

func New(inner ports.NotificationPort, outbox ports.OutboxRepository) Notifier {
	return Notifier{inner: inner, outbox: outbox}
}

var _ ports.NotificationPort = Notifier{}

func (n Notifier) NotifyNewBooking(ctx context.Context, booking domain.Booking, tour domain.Tour) error {
	err := n.inner.NotifyNewBooking(ctx, booking, tour)
	n.record(ctx, domain.OutboxEventNotificationBookingCreated, booking.ID, err)
	return nil
}

func (n Notifier) NotifyBookingStatusChanged(
	ctx context.Context,
	booking domain.Booking,
	tour domain.Tour,
	previousStatus domain.BookingStatus,
) error {
	err := n.inner.NotifyBookingStatusChanged(ctx, booking, tour, previousStatus)
	payload, _ := json.Marshal(map[string]string{
		"previous_status": string(previousStatus),
		"current_status":  string(booking.Status),
	})
	n.recordWithPayload(ctx, domain.OutboxEventNotificationBookingStatus, booking.ID, payload, err)
	return nil
}

func (n Notifier) record(ctx context.Context, eventType string, entityID uuid.UUID, callErr error) {
	n.recordWithPayload(ctx, eventType, entityID, json.RawMessage(`{}`), callErr)
}

func (n Notifier) recordWithPayload(
	ctx context.Context,
	eventType string,
	entityID uuid.UUID,
	payload json.RawMessage,
	callErr error,
) {
	if n.outbox == nil || shouldSkipNotificationError(callErr) {
		return
	}

	lastError := callErr.Error()
	event, err := domain.NewOutboxEvent(domain.NewOutboxEventInput{
		ID:         uuid.New(),
		EventType:  eventType,
		EntityType: domain.EntityTypeBooking,
		EntityID:   entityID,
		Payload:    payload,
		Status:     domain.OutboxStatusPending,
		LastError:  lastError,
		Now:        time.Now().UTC(),
	})
	if err != nil {
		return
	}

	_, _ = n.outbox.Enqueue(ctx, event)
}

func shouldSkipNotificationError(err error) bool {
	return err == nil || errors.Is(err, ports.ErrNotificationNotConfigured)
}
