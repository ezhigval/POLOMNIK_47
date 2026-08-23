package telegram

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"palomnik/internal/config"
	"palomnik/internal/ports"
)

func TestSenderNotConfiguredWithoutToken(t *testing.T) {
	s := New(config.Config{})
	if s.Configured() {
		t.Fatal("expected unconfigured sender")
	}
	if err := s.Send(context.Background(), "telegram", "1", "hi"); err != ports.ErrMessengerNotConfigured {
		t.Fatalf("got %v", err)
	}
}

func TestSenderPostsSendMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/botTOKEN/sendMessage" {
			t.Fatalf("path %s", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		if string(body) == "" || r.Header.Get("Content-Type") == "" {
			t.Fatalf("empty request")
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	s := New(config.Config{
		TelegramBotToken: "TOKEN",
		TelegramAPIBase:  server.URL,
	})
	if !s.Configured() {
		t.Fatal("expected configured")
	}
	if err := s.Send(context.Background(), "telegram", "123", "hello"); err != nil {
		t.Fatalf("send: %v", err)
	}
}

func TestSenderRejectsOtherChannel(t *testing.T) {
	s := New(config.Config{TelegramBotToken: "TOKEN", TelegramAPIBase: "http://127.0.0.1"})
	if err := s.Send(context.Background(), "whatsapp", "123", "hello"); err != ports.ErrMessengerNotConfigured {
		t.Fatalf("got %v", err)
	}
}
