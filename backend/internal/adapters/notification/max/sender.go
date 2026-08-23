package maxnotify

import (
	"context"
	"strings"

	"palomnik/internal/config"
	"palomnik/internal/ports"
)

// Sender is a Max messenger stub. It never invents an API: without credentials
// it is not configured; with credentials it still no-ops send until a real adapter exists.
type Sender struct {
	token string
}

func New(cfg config.Config) Sender {
	return Sender{token: strings.TrimSpace(cfg.MaxBotToken)}
}

var _ ports.MaxMessageSender = Sender{}

func (s Sender) Configured() bool {
	return s.token != ""
}

func (s Sender) SendMessage(context.Context, string, string) error {
	// Credentials alone are not enough — Max HTTP API is not wired yet.
	return ports.ErrNotificationNotConfigured
}
