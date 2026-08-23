package application

import (
	"context"
	"strings"
	"sync"
	"testing"
	"time"

	"palomnik/internal/adapters/repository/memory"
	"palomnik/internal/ports"
)

type recordingMailer struct {
	mu       sync.Mutex
	messages []ports.MailMessage
	ok       bool
}

func (m *recordingMailer) Configured() bool { return m.ok }

func (m *recordingMailer) Send(_ context.Context, message ports.MailMessage) error {
	if !m.ok {
		return ports.ErrMailerNotConfigured
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.messages = append(m.messages, message)
	return nil
}

func TestRequestPasswordResetUnavailableWithoutMailer(t *testing.T) {
	store := memory.NewStore()
	svc := NewAuthService(store, store, nil, &recordingMailer{ok: false}, SocialAuthConfig{}, "secret", time.Hour, "http://localhost:3000")
	err := svc.RequestPasswordReset(context.Background(), "user@example.com")
	if err != ErrPasswordResetUnavailable {
		t.Fatalf("want unavailable, got %v", err)
	}
}

func TestPasswordResetFlow(t *testing.T) {
	store := memory.NewStore()
	mailer := &recordingMailer{ok: true}
	svc := NewAuthService(store, store, nil, mailer, SocialAuthConfig{}, "secret", time.Hour, "https://example.test")

	_, err := svc.Register(context.Background(), RegisterInput{
		Email:    "pilgrim@example.com",
		Phone:    "+79991234567",
		Name:     "Иван",
		Password: "oldpassword",
	})
	if err != nil {
		t.Fatalf("register: %v", err)
	}

	if err := svc.RequestPasswordReset(context.Background(), "pilgrim@example.com"); err != nil {
		t.Fatalf("request reset: %v", err)
	}
	mailer.mu.Lock()
	if len(mailer.messages) < 1 {
		mailer.mu.Unlock()
		t.Fatal("expected reset email")
	}
	body := mailer.messages[len(mailer.messages)-1].Text
	mailer.mu.Unlock()
	idx := strings.Index(body, "token=")
	if idx < 0 {
		t.Fatalf("no token in mail: %s", body)
	}
	token := strings.TrimSpace(body[idx+len("token="):])
	if cut := strings.IndexAny(token, " \n\r"); cut >= 0 {
		token = token[:cut]
	}

	if err := svc.ResetPassword(context.Background(), token, "newpassword1"); err != nil {
		t.Fatalf("reset: %v", err)
	}
	if _, err := svc.Login(context.Background(), LoginInput{Login: "pilgrim@example.com", Password: "newpassword1"}); err != nil {
		t.Fatalf("login with new password: %v", err)
	}
	if err := svc.ResetPassword(context.Background(), token, "anotherpass"); err != ErrInvalidPasswordResetToken {
		t.Fatalf("reused token should fail, got %v", err)
	}
}

func TestRequestPasswordResetUnknownEmailSilent(t *testing.T) {
	store := memory.NewStore()
	mailer := &recordingMailer{ok: true}
	svc := NewAuthService(store, store, nil, mailer, SocialAuthConfig{}, "secret", time.Hour, "http://localhost:3000")
	if err := svc.RequestPasswordReset(context.Background(), "nobody@example.com"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	mailer.mu.Lock()
	defer mailer.mu.Unlock()
	if len(mailer.messages) != 0 {
		t.Fatalf("should not email unknown addresses")
	}
}
