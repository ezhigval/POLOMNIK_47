package ports

import (
	"context"
	"errors"
)

var ErrMessengerNotConfigured = errors.New("messenger adapter not configured")

// MessengerPort is the v3 chat/bot channel. Stage 4 adapters (telegram/max/whatsapp)
// implement it; until then Configured() is false and the site stays up.
type MessengerPort interface {
	Configured() bool
	Send(ctx context.Context, channel string, address string, text string) error
}
