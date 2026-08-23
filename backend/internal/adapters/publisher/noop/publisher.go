package noop

import (
	"context"

	"palomnik/internal/ports"
)

type Publisher struct{}

func New() Publisher { return Publisher{} }

var _ ports.PublisherPort = Publisher{}

func (Publisher) Configured() bool { return false }

func (Publisher) Publish(context.Context, string, ports.PublishContent) error {
	return ports.ErrPublisherNotConfigured
}
