package telegram

import (
	"context"
	"errors"
	"strings"

	telegramnotify "palomnik/internal/adapters/notification/telegram"
	"palomnik/internal/config"
	"palomnik/internal/ports"
)

// Sender implements MessengerPort over the existing Telegram Bot HTTP client.
// Production must keep TELEGRAM_API_BASE on the outbound Worker, same as notifications.
type Sender struct {
	token  string
	client telegramnotify.Client
}

func New(cfg config.Config) Sender {
	return Sender{
		token:  strings.TrimSpace(cfg.TelegramBotToken),
		client: telegramnotify.NewClient(cfg),
	}
}

var _ ports.MessengerPort = Sender{}

func (s Sender) Configured() bool {
	return s.token != ""
}

func (s Sender) Send(ctx context.Context, channel, address, text string) error {
	if !s.Configured() {
		return ports.ErrMessengerNotConfigured
	}
	if ch := strings.ToLower(strings.TrimSpace(channel)); ch != "" && ch != "telegram" {
		return ports.ErrMessengerNotConfigured
	}
	err := s.client.SendMessage(ctx, strings.TrimSpace(address), text)
	if errors.Is(err, ports.ErrNotificationNotConfigured) {
		return ports.ErrMessengerNotConfigured
	}
	return err
}
