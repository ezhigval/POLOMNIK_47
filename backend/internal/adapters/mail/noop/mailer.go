package noop

import (
	"context"

	"polomnik/internal/ports"
)

type Mailer struct{}

func New() Mailer { return Mailer{} }

func (Mailer) Configured() bool { return false }

func (Mailer) Send(context.Context, ports.MailMessage) error {
	return ports.ErrMailerNotConfigured
}
