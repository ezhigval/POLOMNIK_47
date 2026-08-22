package ports

import (
	"context"
	"errors"
)

var (
	ErrPhoneVerifierNotConfigured = errors.New("phone verifier not configured")
	ErrPhoneChallengeFailed       = errors.New("phone challenge failed")
	ErrPhoneChallengePending      = errors.New("phone challenge pending")
	ErrPhoneChallengeExpired      = errors.New("phone challenge expired")
)

// PhoneVerifier checks that a user controls a phone number (channel-agnostic).
// Current live adapter: sms.ru callcheck (user calls a provided number).
type PhoneVerifier interface {
	Available() bool
	Start(ctx context.Context, phone string) (PhoneChallenge, error)
	Status(ctx context.Context, checkID string) (PhoneCheckStatus, error)
}

type PhoneChallenge struct {
	CheckID         string
	CallPhone       string
	CallPhonePretty string
	ExpiresIn       int // seconds; sms.ru allows ~5 minutes
}

type PhoneCheckStatus string

const (
	PhoneCheckPending   PhoneCheckStatus = "pending"
	PhoneCheckConfirmed PhoneCheckStatus = "confirmed"
	PhoneCheckExpired   PhoneCheckStatus = "expired"
)
