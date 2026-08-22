package application

import (
	"context"
	"testing"

	"polomnik/internal/adapters/repository/memory"
	"polomnik/internal/domain"
)

func TestTelegramSettingsValidation(t *testing.T) {
	svc := NewTelegramService(memory.NewStore(), memory.NewStore(), nil, "")
	if _, err := svc.UpdateRecipients(context.Background(), "ab", ""); err != domain.ErrInvalidTelegramUsername {
		t.Fatalf("expected invalid username, got %v", err)
	}

	view, err := svc.UpdateRecipients(context.Background(), "@EzhigVal", "other_user")
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if len(view.Recipients.BookingUsernames) != 1 || view.Recipients.BookingUsernames[0] != "ezhigval" {
		t.Fatalf("booking: %v", view.Recipients.BookingUsernames)
	}
	if len(view.Statuses) != 2 {
		t.Fatalf("statuses: %d", len(view.Statuses))
	}
	for _, item := range view.Statuses {
		if item.ChatBound {
			t.Fatalf("expected waiting /start for %s", item.Username)
		}
	}
}

func TestTelegramInboundUpdatePersistsChatMap(t *testing.T) {
	store := memory.NewStore()
	svc := NewTelegramService(store, store, nil, "")

	if _, err := svc.UpdateRecipients(context.Background(), "ezhigval", ""); err != nil {
		t.Fatalf("update: %v", err)
	}
	if err := svc.HandleInboundUpdate(context.Background(), TelegramInboundUpdate{
		Username: "EzhigVal",
		ChatID:   "111",
		Text:     "/start",
	}); err != nil {
		t.Fatalf("inbound: %v", err)
	}

	view, err := svc.Settings(context.Background())
	if err != nil {
		t.Fatalf("settings: %v", err)
	}
	if len(view.Statuses) != 1 || !view.Statuses[0].ChatBound {
		t.Fatalf("expected bound status, got %+v", view.Statuses)
	}

	if err := svc.HandleInboundUpdate(context.Background(), TelegramInboundUpdate{
		Username: "ezhigval",
		ChatID:   "222",
		Text:     "ещё раз",
	}); err != nil {
		t.Fatalf("update again: %v", err)
	}
	bindings, err := store.ListTelegramChatBindings(context.Background(), []string{"ezhigval"})
	if err != nil {
		t.Fatalf("bindings: %v", err)
	}
	if bindings["ezhigval"] != "222" {
		t.Fatalf("chat map should update, got %v", bindings)
	}
}

func TestTelegramInboundUpdateIgnoresInvalidUsername(t *testing.T) {
	store := memory.NewStore()
	svc := NewTelegramService(store, store, nil, "")
	if err := svc.HandleInboundUpdate(context.Background(), TelegramInboundUpdate{
		Username: "ab",
		ChatID:   "111",
		Text:     "/start",
	}); err != nil {
		t.Fatalf("invalid username must not fail webhook: %v", err)
	}
	bindings, err := store.ListTelegramChatBindings(context.Background(), []string{"ab"})
	if err != nil {
		t.Fatalf("bindings: %v", err)
	}
	if len(bindings) != 0 {
		t.Fatalf("expected empty map, got %v", bindings)
	}
}
