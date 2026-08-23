package ports

import (
	"context"
	"errors"
)

var (
	ErrCaptchaNotConfigured = errors.New("captcha adapter not configured")
	ErrCaptchaFailed        = errors.New("captcha verification failed")
)

// CaptchaPort verifies an anti-bot token. Without keys Configured() is false
// and the HTTP layer keeps honeypot + rate limit only.
type CaptchaPort interface {
	Configured() bool
	Provider() string
	ClientKey() string
	Verify(ctx context.Context, token string, ip string) error
}
