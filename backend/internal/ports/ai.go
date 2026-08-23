package ports

import (
	"context"
	"errors"
)

var ErrAINotConfigured = errors.New("ai adapter not configured")

type AICompletionRequest struct {
	System string
	User   string
}

// AIPort is the v3 language-model port. Live YandexGPT is stage 4/7.
// v4 (calls, autonomous seller) must not be implemented here.
type AIPort interface {
	Configured() bool
	Complete(ctx context.Context, req AICompletionRequest) (string, error)
}
