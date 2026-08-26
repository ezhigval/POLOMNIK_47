package application

import (
	"context"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"

	"palomnik/internal/adapters/repository/memory"
	"palomnik/internal/domain"
)

type recordingSupportNotifier struct {
	mu    sync.Mutex
	notes []domain.SupportNotification
}

func (n *recordingSupportNotifier) NotifyNewBooking(context.Context, domain.Booking, domain.Tour) error {
	return nil
}

func (n *recordingSupportNotifier) NotifyBookingStatusChanged(context.Context, domain.Booking, domain.Tour, domain.BookingStatus) error {
	return nil
}

func (n *recordingSupportNotifier) NotifySupportMessage(_ context.Context, note domain.SupportNotification) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.notes = append(n.notes, note)
	return nil
}

func (n *recordingSupportNotifier) NotifyTourHidden(context.Context, domain.Tour) error {
	return nil
}

func TestSendUserMessageTriggersNotifier(t *testing.T) {
	store := memory.NewStore()
	notifier := &recordingSupportNotifier{}
	svc := NewSupportService(store, notifier)

	userID := uuid.New()
	messages, err := svc.SendUserMessage(context.Background(), userID, "Нужна помощь с заявкой")
	if err != nil {
		t.Fatalf("send: %v", err)
	}
	if len(messages) == 0 {
		t.Fatal("expected saved messages")
	}

	notifier.mu.Lock()
	defer notifier.mu.Unlock()
	if len(notifier.notes) != 1 {
		t.Fatalf("expected one notify call, got %d", len(notifier.notes))
	}
	if notifier.notes[0].Body != "Нужна помощь с заявкой" {
		t.Fatalf("body: %q", notifier.notes[0].Body)
	}
	if notifier.notes[0].UserID != userID {
		t.Fatalf("user: %s", notifier.notes[0].UserID)
	}
}

func TestSendStaffMessageDoesNotNotify(t *testing.T) {
	store := memory.NewStore()
	notifier := &recordingSupportNotifier{}
	svc := NewSupportService(store, notifier)

	userID := uuid.New()
	if _, err := svc.SendUserMessage(context.Background(), userID, "Вопрос от паломника"); err != nil {
		t.Fatalf("user send: %v", err)
	}
	notifier.mu.Lock()
	notifier.notes = nil
	notifier.mu.Unlock()

	threads, err := svc.ListThreads(context.Background())
	if err != nil || len(threads) != 1 {
		t.Fatalf("list threads: %v len=%d", err, len(threads))
	}

	thread, messages, err := svc.SendStaffMessage(context.Background(), threads[0].ID, "Ответ менеджера")
	if err != nil {
		t.Fatalf("staff send: %v", err)
	}
	if thread.ID != threads[0].ID {
		t.Fatalf("thread id mismatch")
	}
	foundStaff := false
	for _, msg := range messages {
		if msg.SenderType == domain.SupportSenderStaff && msg.Body == "Ответ менеджера" {
			foundStaff = true
		}
	}
	if !foundStaff {
		t.Fatal("expected staff reply in messages")
	}

	notifier.mu.Lock()
	defer notifier.mu.Unlock()
	if len(notifier.notes) != 0 {
		t.Fatalf("staff reply must not notify managers, got %d", len(notifier.notes))
	}
}

func TestListThreadsOrderedByUpdatedAt(t *testing.T) {
	store := memory.NewStore()
	svc := NewSupportService(store, nil)

	older := uuid.New()
	newer := uuid.New()
	if _, err := svc.SendUserMessage(context.Background(), older, "Старое"); err != nil {
		t.Fatalf("older: %v", err)
	}
	time.Sleep(5 * time.Millisecond)
	if _, err := svc.SendUserMessage(context.Background(), newer, "Новое"); err != nil {
		t.Fatalf("newer: %v", err)
	}

	threads, err := svc.ListThreads(context.Background())
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(threads) != 2 {
		t.Fatalf("want 2 threads, got %d", len(threads))
	}
	if threads[0].UserID != newer {
		t.Fatalf("expected newest first, got %s", threads[0].UserID)
	}
}

func TestManagementSupportThreadURL(t *testing.T) {
	id := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	got := ManagementSupportThreadURL("https://tikhvin-palomnik.ru/management/bookings", id)
	want := "https://tikhvin-palomnik.ru/management/support/11111111-1111-1111-1111-111111111111"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestFormatSupportMessageDeepLink(t *testing.T) {
	id := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	text := formatSupportMessage(domain.SupportNotification{
		ThreadID: id,
		Body:     "Помогите",
	}, "https://tikhvin-palomnik.ru/management/bookings")
	if !strings.Contains(text, "/management/support/22222222-2222-2222-2222-222222222222") {
		t.Fatalf("missing deep link: %s", text)
	}
	if strings.Contains(text, "нет экрана") {
		t.Fatalf("old wording still present: %s", text)
	}
}
