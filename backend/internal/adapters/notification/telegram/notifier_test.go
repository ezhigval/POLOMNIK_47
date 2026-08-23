package telegram

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"

	"palomnik/internal/config"
	"palomnik/internal/domain"
	"palomnik/internal/ports"
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

	msg := FormatNewBookingMessage(booking, tour, "http://localhost:3000/management/bookings")
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

	msg := FormatStatusChangedMessage(booking, tour, domain.BookingStatusNew, "")
	if !strings.Contains(msg, "Новая → На связи") {
		t.Fatalf("unexpected message:\n%s", msg)
	}
}

func TestClientSendMessageRequiresToken(t *testing.T) {
	client := NewClient(config.Config{})
	err := client.SendMessage(t.Context(), "123", "hi")
	if !errors.Is(err, ports.ErrNotificationNotConfigured) {
		t.Fatalf("expected ErrNotificationNotConfigured, got %v", err)
	}
}

func TestClientSendMessagePostsChatID(t *testing.T) {
	var got string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		got = r.Form.Get("chat_id")
		_, _ = io.WriteString(w, `{"ok":true}`)
	}))
	defer server.Close()

	client := NewClient(config.Config{
		TelegramBotToken: "token",
		TelegramAPIBase:  server.URL,
	})
	if err := client.SendMessage(t.Context(), "111", "hi"); err != nil {
		t.Fatalf("send: %v", err)
	}
	if got != "111" {
		t.Fatalf("chat_id: %q", got)
	}
}
