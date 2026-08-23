package telegramchannel

import (
	"context"
	"errors"
	"strings"

	telegramnotify "palomnik/internal/adapters/notification/telegram"
	"palomnik/internal/adapters/publisher/postformat"
	"palomnik/internal/config"
	"palomnik/internal/ports"
)

// Publisher posts to a Telegram channel via the same Bot API / Worker as notifications.
// The bot must already be a channel admin; this adapter does not call setWebhook.
type Publisher struct {
	token     string
	channelID string
	client    telegramnotify.Client
}

func New(cfg config.Config) Publisher {
	return Publisher{
		token:     strings.TrimSpace(cfg.TelegramBotToken),
		channelID: strings.TrimSpace(cfg.TelegramChannelID),
		client:    telegramnotify.NewClient(cfg),
	}
}

var _ ports.PublisherPort = Publisher{}

func (p Publisher) Configured() bool {
	return p.token != "" && p.channelID != ""
}

func (p Publisher) Publish(ctx context.Context, channel string, content ports.PublishContent) error {
	if !p.Configured() {
		return ports.ErrPublisherNotConfigured
	}
	if ch := strings.ToLower(strings.TrimSpace(channel)); ch != "" && ch != "telegram_channel" {
		return ports.ErrPublisherNotConfigured
	}

	text := postformat.EscapeHTML(postformat.Post(content))
	if strings.TrimSpace(text) == "" {
		return ports.ErrPublisherNotConfigured
	}

	err := p.client.SendMessage(ctx, p.channelID, text)
	if errors.Is(err, ports.ErrNotificationNotConfigured) {
		return ports.ErrPublisherNotConfigured
	}
	return err
}
