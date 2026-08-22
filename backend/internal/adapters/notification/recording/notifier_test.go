package recording

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"polomnik/internal/adapters/notification/noop"
	"polomnik/internal/adapters/repository/memory"
	"polomnik/internal/domain"
	"polomnik/internal/ports"
)

type failingNotifier struct{}

func (failingNotifier) NotifyNewBooking(context.Context, domain.Booking, domain.Tour) error {
	return errors.New("telegram down")
}

func (failingNotifier) NotifyBookingStatusChanged(context.Context, domain.Booking, domain.Tour, domain.BookingStatus) error {
	return errors.New("telegram down")
}

func (failingNotifier) NotifySupportMessage(context.Context, domain.SupportNotification) error {
	return errors.New("telegram down")
}

func TestRecordingNotifierSkipsNotConfigured(t *testing.T) {
	ctx := context.Background()
	store := memory.NewStore()
	notifier := New(noop.New(), store)

	err := notifier.NotifyNewBooking(ctx, domain.Booking{ID: uuid.New()}, domain.Tour{})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	events, listErr := store.ListEvents(ctx, ports.OutboxFilters{}, ports.Pagination{Page: 1, Limit: 10})
	if listErr != nil {
		t.Fatalf("list outbox: %v", listErr)
	}
	if len(events.Items) != 0 {
		t.Fatalf("expected no outbox events, got %d", len(events.Items))
	}
}

func TestRecordingNotifierEnqueuesOnFailure(t *testing.T) {
	ctx := context.Background()
	store := memory.NewStore()
	bookingID := uuid.New()
	notifier := New(failingNotifier{}, store)

	err := notifier.NotifyNewBooking(ctx, domain.Booking{ID: bookingID}, domain.Tour{Title: "Test"})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	events, listErr := store.ListEvents(ctx, ports.OutboxFilters{}, ports.Pagination{Page: 1, Limit: 10})
	if listErr != nil {
		t.Fatalf("list outbox: %v", listErr)
	}
	if len(events.Items) != 1 {
		t.Fatalf("expected 1 outbox event, got %d", len(events.Items))
	}
	if events.Items[0].EventType != domain.OutboxEventNotificationBookingCreated {
		t.Fatalf("unexpected event type %s", events.Items[0].EventType)
	}
}
