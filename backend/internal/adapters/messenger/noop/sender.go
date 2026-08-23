package noop

import (
	"context"

	"polomnik/internal/ports"
)

type Sender struct{}

func New() Sender { return Sender{} }

var _ ports.MessengerPort = Sender{}

func (Sender) Configured() bool { return false }

func (Sender) Send(context.Context, string, string, string) error {
	return ports.ErrMessengerNotConfigured
}
