package noop

import (
	"context"

	"palomnik/internal/ports"
)

type Verifier struct{}

func New() *Verifier {
	return &Verifier{}
}

func (v *Verifier) Available() bool {
	return false
}

func (v *Verifier) Start(_ context.Context, _ string) (ports.PhoneChallenge, error) {
	return ports.PhoneChallenge{}, ports.ErrPhoneVerifierNotConfigured
}

func (v *Verifier) Status(_ context.Context, _ string) (ports.PhoneCheckStatus, error) {
	return "", ports.ErrPhoneVerifierNotConfigured
}
