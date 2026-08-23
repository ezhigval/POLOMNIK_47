package publisher

import (
	"context"
	"testing"

	"palomnik/internal/adapters/repository/memory"
	"palomnik/internal/application"
	"palomnik/internal/config"
	"palomnik/internal/ports"
)

func TestNewNoopByDefault(t *testing.T) {
	port := New(config.Config{}, nil)
	if port.Configured() {
		t.Fatal("expected noop")
	}
}

func TestNewLiveFallsBackWithoutKeys(t *testing.T) {
	port := New(config.Config{PublisherAdapter: "live"}, nil)
	if port.Configured() {
		t.Fatal("live without keys should be noop")
	}
}

func TestNewLiveSiteNewsOnly(t *testing.T) {
	news := application.NewNewsService(memory.NewStore(), nil)
	port := New(config.Config{PublisherAdapter: "live"}, news)
	if !port.Configured() {
		t.Fatal("expected live configured via site_news")
	}
	if err := port.Publish(context.Background(), "vk_wall", ports.PublishContent{Title: "A"}); err != ports.ErrPublisherNotConfigured {
		t.Fatalf("unconfigured channel: %v", err)
	}
	if err := port.Publish(context.Background(), "site_news", ports.PublishContent{Title: "A", Body: "B"}); err != nil {
		t.Fatalf("site_news: %v", err)
	}
}

func TestNewSingleAdapterWithoutKeys(t *testing.T) {
	if port := New(config.Config{PublisherAdapter: "telegram_channel"}, nil); port.Configured() {
		t.Fatal("telegram_channel without keys should be noop")
	}
	if port := New(config.Config{PublisherAdapter: "vk_wall"}, nil); port.Configured() {
		t.Fatal("vk_wall without keys should be noop")
	}
	if port := New(config.Config{PublisherAdapter: "max_feed"}, nil); port.Configured() {
		t.Fatal("max_feed without keys should be noop")
	}
}

func TestPublisherUsesTelegram(t *testing.T) {
	cfg := config.Config{PublisherAdapter: "live", TelegramBotToken: "t", TelegramChannelID: "@ch"}
	if !PublisherUsesTelegram(cfg) {
		t.Fatal("expected telegram publisher live")
	}
	cfg.PublisherAdapter = "vk_wall"
	if PublisherUsesTelegram(cfg) {
		t.Fatal("vk_wall should not require telegram worker")
	}
}
