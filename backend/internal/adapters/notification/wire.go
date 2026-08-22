package notification

import (
	"polomnik/internal/adapters/notification/noop"
	"polomnik/internal/adapters/notification/recording"
	"polomnik/internal/adapters/notification/telegram"
	"polomnik/internal/config"
	"polomnik/internal/ports"
)

type Deps struct {
	Recipients ports.TelegramRecipientsRepository
	Chats      ports.TelegramChatMapRepository
}

func Inner(cfg config.Config, deps Deps) ports.NotificationPort {
	switch cfg.NotificationAdapter {
	case "telegram":
		return telegram.New(cfg, deps.Recipients, deps.Chats)
	default:
		return noop.New()
	}
}

func New(cfg config.Config, outbox ports.OutboxRepository, deps Deps) ports.NotificationPort {
	if !Enabled(cfg) {
		return noop.New()
	}

	inner := telegram.New(cfg, deps.Recipients, deps.Chats)
	return recording.New(inner, outbox)
}

func Enabled(cfg config.Config) bool {
	return cfg.NotificationAdapter == "telegram"
}
