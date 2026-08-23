package messenger

import (
	"testing"

	"palomnik/internal/config"
)

func TestNewFallsBackToNoopWithoutKeys(t *testing.T) {
	if port := New(config.Config{MessengerAdapter: "telegram"}); port.Configured() {
		t.Fatal("telegram without token should be noop")
	}
	if port := New(config.Config{MessengerAdapter: "max"}); port.Configured() {
		t.Fatal("max without token should be noop")
	}
	if port := New(config.Config{MessengerAdapter: "whatsapp"}); port.Configured() {
		t.Fatal("whatsapp without keys should be noop")
	}
}

func TestNewNoopByDefault(t *testing.T) {
	port := New(config.Config{})
	if port.Configured() {
		t.Fatal("expected noop")
	}
}

func TestNewTelegramConfiguredWithToken(t *testing.T) {
	port := New(config.Config{MessengerAdapter: "telegram", TelegramBotToken: "token"})
	if !port.Configured() {
		t.Fatal("expected telegram configured")
	}
}
