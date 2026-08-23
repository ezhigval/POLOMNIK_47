package application

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"palomnik/internal/domain"
	"palomnik/internal/ports"
)

type OutboxWorker struct {
	outbox        ports.OutboxRepository
	tours         ports.TourRepository
	bookings      ports.BookingRepository
	reviews       ports.ReviewRepository
	crm           ports.CRMPort
	accounting    ports.AccountingPort
	notifications ports.NotificationPort
	maxAttempts   int
}

func NewOutboxWorker(
	outbox ports.OutboxRepository,
	tours ports.TourRepository,
	bookings ports.BookingRepository,
	reviews ports.ReviewRepository,
	crm ports.CRMPort,
	accounting ports.AccountingPort,
	notifications ports.NotificationPort,
	maxAttempts int,
) *OutboxWorker {
	if maxAttempts <= 0 {
		maxAttempts = 5
	}
	return &OutboxWorker{
		outbox:        outbox,
		tours:         tours,
		bookings:      bookings,
		reviews:       reviews,
		crm:           crm,
		accounting:    accounting,
		notifications: notifications,
		maxAttempts:   maxAttempts,
	}
}

func (w *OutboxWorker) ProcessBatch(ctx context.Context, limit int) (int, error) {
	events, err := w.outbox.ListPending(ctx, limit)
	if err != nil {
		return 0, err
	}

	processed := 0
	for _, event := range events {
		if err := w.ProcessOne(ctx, event); err != nil {
			return processed, err
		}
		processed++
	}
	return processed, nil
}

func (w *OutboxWorker) ProcessOne(ctx context.Context, event domain.OutboxEvent) error {
	if isNotificationEvent(event.EventType) {
		callErr := w.dispatchNotification(ctx, event)
		if callErr == nil {
			return w.outbox.MarkProcessed(ctx, event.ID)
		}
		return w.outbox.RecordAttempt(ctx, event.ID, callErr.Error(), w.maxAttempts)
	}

	result, callErr := w.dispatch(ctx, event)
	if outboxAttemptComplete(result, callErr) {
		return w.outbox.MarkProcessed(ctx, event.ID)
	}

	lastError := integrationErrorMessage(result, callErr)
	return w.outbox.RecordAttempt(ctx, event.ID, lastError, w.maxAttempts)
}

func isNotificationEvent(eventType string) bool {
	switch eventType {
	case domain.OutboxEventNotificationBookingCreated, domain.OutboxEventNotificationBookingStatus, domain.OutboxEventNotificationSupport:
		return true
	default:
		return false
	}
}

func (w *OutboxWorker) dispatchNotification(ctx context.Context, event domain.OutboxEvent) error {
	if event.EventType == domain.OutboxEventNotificationSupport {
		return w.dispatchSupportNotification(ctx, event)
	}

	booking, err := w.bookings.GetBooking(ctx, event.EntityID)
	if err != nil {
		return err
	}

	tour, tourErr := w.tours.GetTour(ctx, booking.TourID)
	if tourErr != nil {
		tour = domain.Tour{ID: booking.TourID}
	}

	switch event.EventType {
	case domain.OutboxEventNotificationBookingCreated:
		return w.notifications.NotifyNewBooking(ctx, booking, tour)
	case domain.OutboxEventNotificationBookingStatus:
		previousStatus := domain.BookingStatusNew
		if len(event.Payload) > 0 {
			var payload struct {
				PreviousStatus string `json:"previous_status"`
			}
			if json.Unmarshal(event.Payload, &payload) == nil && payload.PreviousStatus != "" {
				previousStatus = domain.BookingStatus(payload.PreviousStatus)
			}
		}
		return w.notifications.NotifyBookingStatusChanged(ctx, booking, tour, previousStatus)
	default:
		return fmt.Errorf("unsupported notification event type: %s", event.EventType)
	}
}

func (w *OutboxWorker) dispatchSupportNotification(ctx context.Context, event domain.OutboxEvent) error {
	var payload struct {
		ThreadID  string `json:"thread_id"`
		MessageID string `json:"message_id"`
		UserID    string `json:"user_id"`
		Body      string `json:"body"`
	}
	if len(event.Payload) > 0 {
		_ = json.Unmarshal(event.Payload, &payload)
	}

	note := domain.SupportNotification{Body: payload.Body}
	if parsed, err := uuid.Parse(payload.ThreadID); err == nil {
		note.ThreadID = parsed
	}
	if parsed, err := uuid.Parse(payload.MessageID); err == nil {
		note.MessageID = parsed
	} else {
		note.MessageID = event.EntityID
	}
	if parsed, err := uuid.Parse(payload.UserID); err == nil {
		note.UserID = parsed
	}
	return w.notifications.NotifySupportMessage(ctx, note)
}

func (w *OutboxWorker) dispatch(ctx context.Context, event domain.OutboxEvent) (ports.IntegrationResult, error) {
	switch event.EventType {
	case domain.OutboxEventCRMPushBooking:
		booking, err := w.bookings.GetBooking(ctx, event.EntityID)
		if err != nil {
			return ports.IntegrationResult{}, err
		}
		return w.crm.PushBooking(ctx, booking)

	case domain.OutboxEventCRMUpdateBookingStatus:
		booking, err := w.bookings.GetBooking(ctx, event.EntityID)
		if err != nil {
			return ports.IntegrationResult{}, err
		}
		return w.crm.UpdateBookingStatus(ctx, booking)

	case domain.OutboxEventCRMSyncTour:
		tour, err := w.tours.GetTour(ctx, event.EntityID)
		if err != nil {
			return ports.IntegrationResult{}, err
		}
		return w.crm.SyncTour(ctx, tour)

	case domain.OutboxEventCRMPushReview:
		review, err := w.reviews.GetReview(ctx, event.EntityID)
		if err != nil {
			return ports.IntegrationResult{}, err
		}
		return w.crm.PushReview(ctx, review)

	case domain.OutboxEventAccountingExport:
		booking, err := w.bookings.GetBooking(ctx, event.EntityID)
		if err != nil {
			return ports.IntegrationResult{}, err
		}
		return w.accounting.ExportBooking(ctx, booking)

	case domain.OutboxEventAccountingCounterparty:
		booking, err := w.bookings.GetBooking(ctx, event.EntityID)
		if err != nil {
			return ports.IntegrationResult{}, err
		}
		return w.accounting.SyncCounterparty(ctx, booking)

	default:
		return ports.IntegrationResult{}, fmt.Errorf("unsupported outbox event type: %s", event.EventType)
	}
}

func outboxAttemptComplete(result ports.IntegrationResult, callErr error) bool {
	if callErr != nil {
		return false
	}
	return result.Status == ports.IntegrationStatusSynced ||
		result.Status == ports.IntegrationStatusNotConfigured
}

func integrationErrorMessage(result ports.IntegrationResult, callErr error) string {
	if callErr != nil {
		return callErr.Error()
	}
	if result.Message != "" {
		return result.Message
	}
	return string(result.Status)
}
