package notification

import (
	maxnotify "polomnik/internal/adapters/notification/max"
	"polomnik/internal/adapters/notification/noop"
	"polomnik/internal/adapters/notification/recording"
	"polomnik/internal/adapters/notification/telegram"
	"polomnik/internal/application"
	"polomnik/internal/config"
	"polomnik/internal/ports"
)

type Deps struct {
	Routing ports.NotificationRoutingRepository
	Legacy  ports.TelegramRecipientsRepository
	Chats   ports.TelegramChatMapRepository
}

func Dispatcher(cfg config.Config, deps Deps) *application.NotificationDispatcher {
	var telegramSender ports.TelegramBotSender
	client := telegram.NewClient(cfg)
	if client.Configured() {
		telegramSender = client
	}
	return application.NewNotificationDispatcher(
		deps.Routing,
		deps.Legacy,
		deps.Chats,
		telegramSender,
		maxnotify.New(cfg),
		cfg.TelegramChatID,
		cfg.ManagementBaseURL,
	)
}

func Inner(cfg config.Config, deps Deps) ports.NotificationPort {
	if !Enabled(cfg) {
		return noop.New()
	}
	return Dispatcher(cfg, deps)
}

func New(cfg config.Config, outbox ports.OutboxRepository, deps Deps) ports.NotificationPort {
	if !Enabled(cfg) {
		return noop.New()
	}
	return recording.New(Dispatcher(cfg, deps), outbox)
}

func Enabled(cfg config.Config) bool {
	return cfg.NotificationAdapter == "telegram"
}
