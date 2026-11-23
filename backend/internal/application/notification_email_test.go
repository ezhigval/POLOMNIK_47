package application

import (
	"context"
	"errors"
	"net/textproto"
	"strings"
	"testing"
	"time"

	"palomnik/internal/domain"
	"palomnik/internal/ports"
)

type emailRouteRepo struct{ routing domain.NotificationRouting }

func (r *emailRouteRepo) GetNotificationRouting(context.Context) (domain.NotificationRouting, error) {
	return r.routing, nil
}
func (r *emailRouteRepo) UpsertNotificationRouting(_ context.Context, routing domain.NotificationRouting) (domain.NotificationRouting, error) {
	r.routing = routing
	return routing, nil
}

type emailRecorder struct {
	configured bool
	sent       []ports.MailMessage
	err        error
}
func (m *emailRecorder) Configured() bool { return m.configured }
func (m *emailRecorder) Send(_ context.Context, message ports.MailMessage) error {
	m.sent = append(m.sent, message)
	return m.err
}

func TestNotificationDispatcherRoutesEmailByEvent(t *testing.T) {
	routing, err := domain.NewNotificationRouting(map[domain.NotificationEventKind][]domain.NotificationRecipient{
		domain.NotificationEventBookingCreated: {
			{Channel: domain.NotificationChannelEmail, Address: "manager@example.ru"},
			{Channel: domain.NotificationChannelEmail, Address: "backup@example.ru"},
		},
	}, time.Time{})
	if err != nil { t.Fatal(err) }

	mail := &emailRecorder{configured: true}
	dispatcher := NewNotificationDispatcher(&emailRouteRepo{routing: routing}, nil, nil, nil, nil, "", "", mail)
	if err := dispatcher.dispatch(context.Background(), domain.NotificationEventBookingCreated, "<b>Новая заявка</b>"); err != nil {
		t.Fatal(err)
	}
	if len(mail.sent) != 2 {
		t.Fatalf("expected two emails, got %d", len(mail.sent))
	}
	for i, address := range []string{"manager@example.ru", "backup@example.ru"} {
		msg := mail.sent[i]
		if len(msg.To) != 1 || msg.To[0] != address || msg.Subject != "Новая заявка" || !strings.Contains(msg.HTML, "Новая заявка") {
			t.Fatalf("unexpected email: %+v", msg)
		}
	}
	if err := dispatcher.dispatch(context.Background(), domain.NotificationEventSupportMessage, "hello"); err == nil {
		t.Fatal("expected no-recipient fallback error without configured Telegram fallback")
	}
	if len(mail.sent) != 2 { t.Fatal("other event must not receive emails") }
}

func TestNotificationDispatcherEmailFailureIsReported(t *testing.T) {
	routing, err := domain.NewNotificationRouting(map[domain.NotificationEventKind][]domain.NotificationRecipient{
		domain.NotificationEventSupportMessage: {{Channel: domain.NotificationChannelEmail, Address: "ops@example.ru"}},
	}, time.Time{})
	if err != nil { t.Fatal(err) }
	expected := errors.New("SMTP connection failed")
	mail := &emailRecorder{configured: true, err: expected}
	dispatcher := NewNotificationDispatcher(&emailRouteRepo{routing: routing}, nil, nil, nil, nil, "", "", mail)
	err = dispatcher.dispatch(context.Background(), domain.NotificationEventSupportMessage, "<b>Support</b>")
	if !errors.Is(err, expected) { t.Fatalf("expected SMTP failure, got %v", err) }
}

func TestNotificationDispatcherUnconfiguredEmailIsNotMarkedSent(t *testing.T) {
	routing, err := domain.NewNotificationRouting(map[domain.NotificationEventKind][]domain.NotificationRecipient{
		domain.NotificationEventBookingCreated: {{Channel: domain.NotificationChannelEmail, Address: "ops@example.ru"}},
	}, time.Time{})
	if err != nil { t.Fatal(err) }
	mail := &emailRecorder{configured: false}
	dispatcher := NewNotificationDispatcher(&emailRouteRepo{routing: routing}, nil, nil, nil, nil, "", "", mail)
	if err := dispatcher.dispatch(context.Background(), domain.NotificationEventBookingCreated, "test"); err != nil {
		t.Fatalf("existing contract skips unconfigured recipients: %v", err)
	}
	if len(mail.sent) != 0 { t.Fatal("unconfigured mailer must not be called") }
}

type maxSenderRecorder struct{ delivered int }
func (m *maxSenderRecorder) Configured() bool { return true }
func (m *maxSenderRecorder) SendMessage(context.Context, string, string) error {
	m.delivered++
	return nil
}

func TestPartialSMTPFailureDoesNotRetrySuccessfulMessenger(t *testing.T) {
	routing, err := domain.NewNotificationRouting(map[domain.NotificationEventKind][]domain.NotificationRecipient{
		domain.NotificationEventBookingCreated: {
			{Channel: domain.NotificationChannelMax, Address: "+79991112233"},
			{Channel: domain.NotificationChannelEmail, Address: "blocked@example.ru"},
		},
	}, time.Time{})
	if err != nil { t.Fatal(err) }
	messenger := &maxSenderRecorder{}
	mailer := &emailRecorder{configured: true, err: errors.New("SMTP unavailable")}
	dispatcher := NewNotificationDispatcher(&emailRouteRepo{routing: routing}, nil, nil, nil, messenger, "", "", mailer)
	if err := dispatcher.dispatch(context.Background(), domain.NotificationEventBookingCreated, "booking"); err != nil {
		t.Fatalf("partial delivery should not retry all channels: %v", err)
	}
	if messenger.delivered != 1 || len(mailer.sent) != 1 {
		t.Fatalf("expected one delivery attempt per channel, messenger=%d email=%d", messenger.delivered, len(mailer.sent))
	}
}

func TestPermanentSMTPBounceCannotBeRetried(t *testing.T) {
	routing, err := domain.NewNotificationRouting(map[domain.NotificationEventKind][]domain.NotificationRecipient{
		domain.NotificationEventBookingCreated: {
			{Channel: domain.NotificationChannelEmail, Address: "blocked@example.ru"},
		},
	}, time.Time{})
	if err != nil { t.Fatal(err) }
	mailer := &emailRecorder{configured: true, err: &textproto.Error{Code: 550, Msg: "bounced check filter"}}
	dispatcher := NewNotificationDispatcher(&emailRouteRepo{routing: routing}, nil, nil, nil, nil, "", "", mailer)
	err = dispatcher.dispatch(context.Background(), domain.NotificationEventBookingCreated, "booking")
	if !errors.Is(err, ports.ErrNotificationPermanent) {
		t.Fatalf("550 must be permanent: %v", err)
	}
}
