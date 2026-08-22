package telegram

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"

	"polomnik/internal/adapters/repository/memory"
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
	notifier := New(config.Config{}, nil, nil)
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

func TestNotifierResolvesUsernamesAndSkipsUnknown(t *testing.T) {
	var mu sync.Mutex
	var chatIDs []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("parse form: %v", err)
		}
		mu.Lock()
		chatIDs = append(chatIDs, r.Form.Get("chat_id"))
		mu.Unlock()
		_, _ = io.WriteString(w, `{"ok":true}`)
	}))
	defer server.Close()

	store := memory.NewStore()
	settings, err := domain.NewTelegramRecipients("ezhigval, unknown_user", "", time.Time{})
	if err != nil {
		t.Fatalf("settings: %v", err)
	}
	if _, err := store.UpsertTelegramRecipients(t.Context(), settings); err != nil {
		t.Fatalf("upsert settings: %v", err)
	}
	binding, err := domain.NewTelegramChatBinding("ezhigval", "111", time.Time{})
	if err != nil {
		t.Fatalf("binding: %v", err)
	}
	if err := store.UpsertTelegramChatBinding(t.Context(), binding); err != nil {
		t.Fatalf("upsert binding: %v", err)
	}

	notifier := New(config.Config{
		TelegramBotToken: "token",
		TelegramAPIBase:  server.URL,
		TelegramChatID:   "999",
	}, store, store)

	if err := notifier.NotifyNewBooking(t.Context(), domain.Booking{Name: "Иван"}, domain.Tour{Title: "Тур"}); err != nil {
		t.Fatalf("notify: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()
	if len(chatIDs) != 1 || chatIDs[0] != "111" {
		t.Fatalf("expected only resolved chat 111, got %v", chatIDs)
	}
}

func TestNotifierUsesFallbackChatWhenSettingsEmpty(t *testing.T) {
	var got string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		got = r.Form.Get("chat_id")
		_, _ = io.WriteString(w, `{"ok":true}`)
	}))
	defer server.Close()

	notifier := New(config.Config{
		TelegramBotToken: "token",
		TelegramAPIBase:  server.URL,
		TelegramChatID:   "999",
	}, memory.NewStore(), memory.NewStore())

	if err := notifier.NotifySupportMessage(t.Context(), domain.SupportNotification{Body: "Помогите"}); err != nil {
		t.Fatalf("notify: %v", err)
	}
	if got != "999" {
		t.Fatalf("fallback chat: %q", got)
	}
}

func TestNotifierSkipsUnknownWithoutFallback(t *testing.T) {
	called := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	store := memory.NewStore()
	settings, err := domain.NewTelegramRecipients("unknown_user", "", time.Time{})
	if err != nil {
		t.Fatalf("settings: %v", err)
	}
	if _, err := store.UpsertTelegramRecipients(t.Context(), settings); err != nil {
		t.Fatalf("upsert: %v", err)
	}

	notifier := New(config.Config{
		TelegramBotToken: "token",
		TelegramAPIBase:  server.URL,
		TelegramChatID:   "999",
	}, store, store)

	if err := notifier.NotifyNewBooking(t.Context(), domain.Booking{}, domain.Tour{}); err != nil {
		t.Fatalf("notify: %v", err)
	}
	if called {
		t.Fatal("should skip unknown username instead of sending")
	}
}
