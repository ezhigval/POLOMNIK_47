package notification

import (
	"polomnik/internal/adapters/notification/noop"
	"polomnik/internal/adapters/notification/recording"
	"polomnik/internal/adapters/notification/telegram"
	"polomnik/internal/config"
	"polomnik/internal/ports"
)

func Inner(cfg config.Config) ports.NotificationPort {
	switch cfg.NotificationAdapter {
	case "telegram":
		return telegram.New(cfg)
	default:
		return noop.New()
	}
}

func New(cfg config.Config, outbox ports.OutboxRepository) ports.NotificationPort {
	if !Enabled(cfg) {
		return noop.New()
	}

	inner := telegram.New(cfg)
	return recording.New(inner, outbox)
}

func Enabled(cfg config.Config) bool {
	return cfg.NotificationAdapter == "telegram"
}
