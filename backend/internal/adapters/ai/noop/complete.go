package noop

import (
	"context"

	"polomnik/internal/ports"
)

type Client struct{}

func New() Client { return Client{} }

var _ ports.AIPort = Client{}

func (Client) Configured() bool { return false }

func (Client) Complete(context.Context, ports.AICompletionRequest) (string, error) {
	return "", ports.ErrAINotConfigured
}
