package publisher

import (
	"context"
	"strings"

	"palomnik/internal/ports"
)

// Mux routes Publish by channel name. Failure of one channel is the caller's
// concern (stage 6 SMM does not roll back the others).
type Mux struct {
	channels map[string]ports.PublisherPort
}

func newMux(channels map[string]ports.PublisherPort) Mux {
	return Mux{channels: channels}
}

var _ ports.PublisherPort = Mux{}

func (m Mux) Configured() bool {
	for _, port := range m.channels {
		if port != nil && port.Configured() {
			return true
		}
	}
	return false
}

func (m Mux) Publish(ctx context.Context, channel string, content ports.PublishContent) error {
	name := strings.ToLower(strings.TrimSpace(channel))
	port := m.channels[name]
	if port == nil || !port.Configured() {
		return ports.ErrPublisherNotConfigured
	}
	return port.Publish(ctx, name, content)
}
