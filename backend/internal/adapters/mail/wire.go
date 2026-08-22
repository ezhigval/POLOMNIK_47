package mail

import (
	"strings"

	smtpmail "polomnik/internal/adapters/mail/smtp"
	"polomnik/internal/adapters/mail/noop"
	"polomnik/internal/config"
	"polomnik/internal/ports"
)

func New(cfg config.Config) ports.Mailer {
	adapter := strings.ToLower(strings.TrimSpace(cfg.MailAdapter))
	switch adapter {
	case "smtp":
		mailer := smtpmail.New(cfg)
		if mailer.Configured() {
			return mailer
		}
		return noop.New()
	default:
		return noop.New()
	}
}
