package telegram

import (
	"errors"
	"strings"
	"testing"

	"github.com/google/uuid"

	"polomnik/internal/config"
	"polomnik/internal/domain"
	"polomnik/internal/ports"
)

func TestFormatNewBookingMessage(t *testing.T) {
	booking := domain.Booking{
		ID:          uuid.New(),
		Name:        "Иван",
		Phone:       "+79990000000",
		Email:       "a@b.c",
		PeopleCount: 2,
		TotalPrice:  30000,
		Comment:     "Вопрос",
		Status:      domain.BookingStatusNew,
		Overbooked:  true,
	}
	tour := domain.Tour{Title: "Оптина пустынь"}

	msg := formatNewBookingMessage(booking, tour, "http://localhost:3000/management/bookings")
	for _, want := range []string{"Иван", "+79990000000", "Оптина", "30000", "Overbooking", "management"} {
		if !strings.Contains(msg, want) {
			t.Fatalf("expected %q in message, got:\n%s", want, msg)
		}
	}
}

func TestFormatStatusChangedMessage(t *testing.T) {
	booking := domain.Booking{
		ID:     uuid.New(),
		Name:   "Иван",
		Phone:  "+79990000000",
		Status: domain.BookingStatusContacted,
	}
	tour := domain.Tour{Title: "Дивеево"}

	msg := formatStatusChangedMessage(booking, tour, domain.BookingStatusNew, "")
	if !strings.Contains(msg, "Новая → На связи") {
		t.Fatalf("unexpected message:\n%s", msg)
	}
}

func TestNotifierSkipsWhenNotConfigured(t *testing.T) {
	notifier := New(config.Config{})
	err := notifier.NotifyNewBooking(t.Context(), domain.Booking{}, domain.Tour{})
	if !errors.Is(err, ports.ErrNotificationNotConfigured) {
		t.Fatalf("expected ErrNotificationNotConfigured, got %v", err)
	}
}

func TestClientSendMessageRequiresToken(t *testing.T) {
	client := NewClient(config.Config{})
	err := client.SendMessage(t.Context(), "123", "hi")
	if !errors.Is(err, ports.ErrNotificationNotConfigured) {
		t.Fatalf("expected ErrNotificationNotConfigured, got %v", err)
	}
}
