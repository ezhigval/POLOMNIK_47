package application

import (
	"context"
	"sync"
	"testing"

	"github.com/google/uuid"

	"polomnik/internal/adapters/repository/memory"
	"polomnik/internal/domain"
)

type recordingSupportNotifier struct {
	mu    sync.Mutex
	notes []domain.SupportNotification
}

func (n *recordingSupportNotifier) NotifyNewBooking(context.Context, domain.Booking, domain.Tour) error {
	return nil
}

func (n *recordingSupportNotifier) NotifyBookingStatusChanged(context.Context, domain.Booking, domain.Tour, domain.BookingStatus) error {
	return nil
}

func (n *recordingSupportNotifier) NotifySupportMessage(_ context.Context, note domain.SupportNotification) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.notes = append(n.notes, note)
	return nil
}

func TestSendUserMessageTriggersNotifier(t *testing.T) {
	store := memory.NewStore()
	notifier := &recordingSupportNotifier{}
	svc := NewSupportService(store, notifier)

	userID := uuid.New()
	messages, err := svc.SendUserMessage(context.Background(), userID, "Нужна помощь с заявкой")
	if err != nil {
		t.Fatalf("send: %v", err)
	}
	if len(messages) == 0 {
		t.Fatal("expected saved messages")
	}

	notifier.mu.Lock()
	defer notifier.mu.Unlock()
	if len(notifier.notes) != 1 {
		t.Fatalf("expected one notify call, got %d", len(notifier.notes))
	}
	if notifier.notes[0].Body != "Нужна помощь с заявкой" {
		t.Fatalf("body: %q", notifier.notes[0].Body)
	}
	if notifier.notes[0].UserID != userID {
		t.Fatalf("user: %s", notifier.notes[0].UserID)
	}
}
