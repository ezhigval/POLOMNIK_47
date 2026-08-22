package application

import "errors"

var (
	ErrTourInactive                   = errors.New("tour inactive")
	ErrTourExpired                    = errors.New("tour expired")
	ErrUnauthorized                   = errors.New("unauthorized")
	ErrPhoneVerificationUnavailable   = errors.New("phone verification unavailable")
	ErrPhoneVerificationRequired      = errors.New("phone verification required")
	ErrPhoneVerificationNotConfirmed  = errors.New("phone verification not confirmed")
	ErrPhoneUserNotFound              = errors.New("phone user not found")
)
