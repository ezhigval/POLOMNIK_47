package maxfeed

import (
	"context"
	"errors"
	"strings"
	"time"

	maxmsg "palomnik/internal/adapters/messenger/max"
	"palomnik/internal/adapters/publisher/postformat"
	"palomnik/internal/ports"
)

// Publisher posts to a Max channel via official POST /messages?chat_id=.
type Publisher struct {
	chatID string
	sender maxmsg.Sender
}

func New(token, apiBase, chatID string, timeout time.Duration) Publisher {
	return Publisher{
		chatID: strings.TrimSpace(chatID),
		sender: maxmsg.New(token, apiBase, timeout),
	}
}

var _ ports.PublisherPort = Publisher{}

func (p Publisher) Configured() bool {
	return p.chatID != "" && p.sender.Configured()
}

func (p Publisher) Publish(ctx context.Context, channel string, content ports.PublishContent) error {
	if !p.Configured() {
		return ports.ErrPublisherNotConfigured
	}
	if ch := strings.ToLower(strings.TrimSpace(channel)); ch != "" && ch != "max_feed" {
		return ports.ErrPublisherNotConfigured
	}

	text := postformat.Post(content)
	if strings.TrimSpace(text) == "" {
		return ports.ErrPublisherNotConfigured
	}

	err := p.sender.Send(ctx, "max", "chat:"+p.chatID, text)
	if errors.Is(err, ports.ErrMessengerNotConfigured) {
		return ports.ErrPublisherNotConfigured
	}
	return err
}
