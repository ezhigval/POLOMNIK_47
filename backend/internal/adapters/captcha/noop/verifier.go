package noop

import (
	"context"

	"palomnik/internal/ports"
)

type Verifier struct{}

func New() Verifier {
	return Verifier{}
}

var _ ports.CaptchaPort = Verifier{}

func (Verifier) Configured() bool { return false }

func (Verifier) Provider() string { return "noop" }

func (Verifier) ClientKey() string { return "" }

func (Verifier) Verify(context.Context, string, string) error {
	return ports.ErrCaptchaNotConfigured
}
