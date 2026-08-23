package smtpmail

import (
	"context"
	"fmt"
	"net/smtp"
	"strings"
	"time"

	"palomnik/internal/config"
	"palomnik/internal/ports"
)

type Mailer struct {
	host     string
	port     string
	username string
	password string
	from     string
	timeout  time.Duration
}

func New(cfg config.Config) *Mailer {
	host := strings.TrimSpace(cfg.SMTPHost)
	from := strings.TrimSpace(cfg.SMTPFrom)
	if host == "" || from == "" {
		return &Mailer{}
	}
	port := strings.TrimSpace(cfg.SMTPPort)
	if port == "" {
		port = "587"
	}
	timeout := cfg.SMTPTimeout
	if timeout <= 0 {
		timeout = 15 * time.Second
	}
	return &Mailer{
		host:     host,
		port:     port,
		username: strings.TrimSpace(cfg.SMTPUsername),
		password: cfg.SMTPPassword,
		from:     from,
		timeout:  timeout,
	}
}

func (m *Mailer) Configured() bool {
	return m != nil && m.host != "" && m.from != ""
}

func (m *Mailer) Send(ctx context.Context, message ports.MailMessage) error {
	if !m.Configured() {
		return ports.ErrMailerNotConfigured
	}
	if len(message.To) == 0 {
		return fmt.Errorf("mail: empty recipients")
	}

	addr := m.host + ":" + m.port
	body := buildMIME(m.from, message)

	auth := smtp.PlainAuth("", m.username, m.password, m.host)
	done := make(chan error, 1)
	go func() {
		done <- smtp.SendMail(addr, auth, m.from, message.To, []byte(body))
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return err
	case <-time.After(m.timeout):
		return fmt.Errorf("mail: smtp timeout")
	}
}

func buildMIME(from string, message ports.MailMessage) string {
	var b strings.Builder
	b.WriteString("From: ")
	b.WriteString(from)
	b.WriteString("\r\nTo: ")
	b.WriteString(strings.Join(message.To, ", "))
	b.WriteString("\r\nSubject: ")
	b.WriteString(message.Subject)
	b.WriteString("\r\nMIME-Version: 1.0\r\n")
	if strings.TrimSpace(message.HTML) != "" {
		b.WriteString("Content-Type: text/html; charset=UTF-8\r\n\r\n")
		b.WriteString(message.HTML)
	} else {
		b.WriteString("Content-Type: text/plain; charset=UTF-8\r\n\r\n")
		b.WriteString(message.Text)
	}
	return b.String()
}
