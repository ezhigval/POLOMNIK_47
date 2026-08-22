package application

import "errors"

var (
	ErrTourInactive = errors.New("tour inactive")
	ErrTourExpired  = errors.New("tour expired")
	ErrUnauthorized = errors.New("unauthorized")
)
