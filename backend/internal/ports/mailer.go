package ports

import (
	"context"
	"errors"
)

var ErrMailerNotConfigured = errors.New("mailer not configured")

type MailMessage struct {
	To      []string
	Subject string
	Text    string
	HTML    string
}

type Mailer interface {
	Configured() bool
	Send(ctx context.Context, message MailMessage) error
}
