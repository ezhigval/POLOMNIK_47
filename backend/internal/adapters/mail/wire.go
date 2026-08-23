package mail

import (
	"strings"

	smtpmail "palomnik/internal/adapters/mail/smtp"
	"palomnik/internal/adapters/mail/noop"
	"palomnik/internal/config"
	"palomnik/internal/ports"
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
