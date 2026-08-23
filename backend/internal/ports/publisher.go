package ports

import (
	"context"
	"errors"
)

var ErrPublisherNotConfigured = errors.New("publisher adapter not configured")

type PublishContent struct {
	Title string
	Body  string
	URL   string
}

// PublisherPort is the v3 SMM channel. Stage 4 adapters publish to site/news/VK/Max.
type PublisherPort interface {
	Configured() bool
	Publish(ctx context.Context, channel string, content PublishContent) error
}
