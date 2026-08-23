package telegramchannel

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"palomnik/internal/config"
	"palomnik/internal/ports"
)

func TestPublisherNotConfiguredWithoutKeys(t *testing.T) {
	p := New(config.Config{TelegramBotToken: "token"})
	if p.Configured() {
		t.Fatal("expected unconfigured without channel id")
	}
	if err := p.Publish(context.Background(), "telegram_channel", ports.PublishContent{Title: "A"}); err != ports.ErrPublisherNotConfigured {
		t.Fatalf("got %v", err)
	}
}

func TestPublisherPostsSendMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/botTOKEN/sendMessage" {
			t.Fatalf("path %s", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		values, err := url.ParseQuery(string(body))
		if err != nil {
			t.Fatalf("form: %v", err)
		}
		if values.Get("chat_id") != "@tikhvin_path" {
			t.Fatalf("chat_id %q", values.Get("chat_id"))
		}
		if !strings.Contains(values.Get("text"), "Заголовок") {
			t.Fatalf("text %q", values.Get("text"))
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	p := New(config.Config{
		TelegramBotToken:  "TOKEN",
		TelegramAPIBase:   server.URL,
		TelegramChannelID: "@tikhvin_path",
	})
	if !p.Configured() {
		t.Fatal("expected configured")
	}
	if err := p.Publish(context.Background(), "telegram_channel", ports.PublishContent{Title: "Заголовок", Body: "Текст"}); err != nil {
		t.Fatalf("publish: %v", err)
	}
}

func TestPublisherRejectsOtherChannel(t *testing.T) {
	p := New(config.Config{TelegramBotToken: "TOKEN", TelegramChannelID: "-1001", TelegramAPIBase: "http://127.0.0.1"})
	if err := p.Publish(context.Background(), "vk_wall", ports.PublishContent{Title: "A"}); err != ports.ErrPublisherNotConfigured {
		t.Fatalf("got %v", err)
	}
}
