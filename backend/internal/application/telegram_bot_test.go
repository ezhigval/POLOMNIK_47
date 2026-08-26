package application

import (
	"context"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"

	"palomnik/internal/adapters/integration/noop"
	notificationnoop "palomnik/internal/adapters/notification/noop"
	"palomnik/internal/adapters/repository/memory"
	"palomnik/internal/domain"
	"palomnik/internal/ports"
)

type recordingBotSender struct {
	mu       sync.Mutex
	messages []string
}

func (s *recordingBotSender) SendMessage(_ context.Context, _, text string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.messages = append(s.messages, text)
	return nil
}

func (s *recordingBotSender) last() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.messages) == 0 {
		return ""
	}
	return s.messages[len(s.messages)-1]
}

type recordingMessenger struct {
	configured bool
	mu         sync.Mutex
	sent       [][3]string
}

func (m *recordingMessenger) Configured() bool { return m.configured }

func (m *recordingMessenger) Send(_ context.Context, channel, address, text string) error {
	if !m.configured {
		return ports.ErrMessengerNotConfigured
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sent = append(m.sent, [3]string{channel, address, text})
	return nil
}

func botTestUser(t *testing.T, store *memory.Store, telegramID string) domain.User {
	t.Helper()
	user, err := domain.NewOAuthUser(domain.OAuthUserInput{
		ID:       uuid.New(),
		Provider: "telegram",
		Subject:  telegramID,
		Name:     "Менеджер",
		Now:      time.Now().UTC(),
	})
	if err != nil {
		t.Fatalf("user: %v", err)
	}
	if _, err := store.CreateUser(context.Background(), user); err != nil {
		t.Fatalf("create user: %v", err)
	}
	identity, err := domain.NewUserIdentity(user.ID, "telegram", telegramID, time.Now().UTC())
	if err != nil {
		t.Fatalf("identity: %v", err)
	}
	if err := store.AddIdentity(context.Background(), identity); err != nil {
		t.Fatalf("add identity: %v", err)
	}
	return user
}

func botTestRole(t *testing.T, store *memory.Store, name string, userID uuid.UUID, perms []domain.Permission) {
	t.Helper()
	role, err := domain.NewAdminRole(uuid.New(), name, "password1", perms, time.Now().UTC())
	if err != nil {
		t.Fatalf("role: %v", err)
	}
	if _, err := store.CreateAdminRole(context.Background(), role); err != nil {
		t.Fatalf("create role: %v", err)
	}
	if err := store.AssignUserToRole(context.Background(), domain.AdminRoleAssignment{
		RoleID:    role.ID,
		UserID:    userID,
		CreatedAt: time.Now().UTC(),
	}); err != nil {
		t.Fatalf("assign: %v", err)
	}
}

func newBotService(t *testing.T, store *memory.Store, sender ports.TelegramBotSender, messenger ports.MessengerPort) *TelegramService {
	t.Helper()
	settings := NewNotificationSettingsService(store, store, store, true, false)
	support := NewSupportService(store, nil).WithClientMessenger(messenger, store)
	tours := NewTourService(store, nil, noop.NewCRMAdapter())
	bookings := NewBookingService(store, store, noop.NewCRMAdapter(), noop.NewAccountingAdapter(), notificationnoop.New(), store, 0)
	adminRoles := NewAdminRoleService(store, store, "admin-token", "jwt-secret")
	return NewTelegramService(settings, store, sender, "").WithBot(TelegramBotDeps{
		Support:    support,
		Bookings:   bookings,
		Tours:      tours,
		AdminRoles: adminRoles,
		Users:      store,
	})
}

func TestBotReplyFromSupportRecipient(t *testing.T) {
	store := memory.NewStore()
	sender := &recordingBotSender{}
	svc := newBotService(t, store, sender, &recordingMessenger{})

	if _, err := svc.UpdateRecipients(context.Background(), "", "ezhigval"); err != nil {
		t.Fatalf("recipients: %v", err)
	}

	pilgrim := uuid.New()
	support := NewSupportService(store, nil)
	if _, err := support.SendUserMessage(context.Background(), pilgrim, "Нужна помощь"); err != nil {
		t.Fatalf("user message: %v", err)
	}
	threads, err := support.ListThreads(context.Background())
	if err != nil || len(threads) != 1 {
		t.Fatalf("threads: %v len=%d", err, len(threads))
	}

	if err := svc.HandleInboundUpdate(context.Background(), TelegramInboundUpdate{
		Username: "EzhigVal",
		UserID:   "1001",
		ChatID:   "1001",
		Text:     "/reply " + threads[0].ID.String() + " Отвечаем в боте",
	}); err != nil {
		t.Fatalf("inbound: %v", err)
	}
	if sender.last() != "Ответ записан в диалог." {
		t.Fatalf("reply ack: %q", sender.last())
	}

	_, messages, err := support.GetThreadByID(context.Background(), threads[0].ID)
	if err != nil {
		t.Fatalf("get thread: %v", err)
	}
	found := false
	for _, msg := range messages {
		if msg.SenderType == domain.SupportSenderStaff && msg.Body == "Отвечаем в боте" {
			found = true
		}
	}
	if !found {
		t.Fatal("expected staff message in thread")
	}
}

func TestBotReplyToNotificationExtractsThreadID(t *testing.T) {
	store := memory.NewStore()
	sender := &recordingBotSender{}
	svc := newBotService(t, store, sender, &recordingMessenger{})
	if _, err := svc.UpdateRecipients(context.Background(), "", "ezhigval"); err != nil {
		t.Fatalf("recipients: %v", err)
	}

	support := NewSupportService(store, nil)
	if _, err := support.SendUserMessage(context.Background(), uuid.New(), "Вопрос"); err != nil {
		t.Fatalf("user message: %v", err)
	}
	threads, err := support.ListThreads(context.Background())
	if err != nil || len(threads) != 1 {
		t.Fatalf("threads: %v", err)
	}

	if err := svc.HandleInboundUpdate(context.Background(), TelegramInboundUpdate{
		Username:    "ezhigval",
		ChatID:      "1001",
		Text:        "Сейчас ответим",
		ReplyToText: "Сообщение в поддержку\nДиалог: " + threads[0].ID.String(),
	}); err != nil {
		t.Fatalf("inbound: %v", err)
	}
	if sender.last() != "Ответ записан в диалог." {
		t.Fatalf("ack: %q", sender.last())
	}
}

func TestBotReplyForbiddenForStranger(t *testing.T) {
	store := memory.NewStore()
	sender := &recordingBotSender{}
	svc := newBotService(t, store, sender, &recordingMessenger{})
	if _, err := svc.UpdateRecipients(context.Background(), "", "ezhigval"); err != nil {
		t.Fatalf("recipients: %v", err)
	}
	support := NewSupportService(store, nil)
	if _, err := support.SendUserMessage(context.Background(), uuid.New(), "Вопрос"); err != nil {
		t.Fatalf("user message: %v", err)
	}
	threads, _ := support.ListThreads(context.Background())

	if err := svc.HandleInboundUpdate(context.Background(), TelegramInboundUpdate{
		Username: "other_user",
		ChatID:   "2002",
		Text:     "/reply " + threads[0].ID.String() + " нет",
	}); err != nil {
		t.Fatalf("inbound: %v", err)
	}
	if sender.last() != botForbiddenText {
		t.Fatalf("got %q", sender.last())
	}
}

func TestBotTourCommandsRequireManageTours(t *testing.T) {
	store := memory.NewStore()
	sender := &recordingBotSender{}
	svc := newBotService(t, store, sender, &recordingMessenger{})

	tour := testTour()
	if _, err := store.CreateTour(context.Background(), tour); err != nil {
		t.Fatalf("create tour: %v", err)
	}

	if err := svc.HandleInboundUpdate(context.Background(), TelegramInboundUpdate{
		Username: "nobody1",
		UserID:   "1",
		ChatID:   "1",
		Text:     "/tour test-tour off",
	}); err != nil {
		t.Fatalf("inbound: %v", err)
	}
	if sender.last() != botForbiddenText {
		t.Fatalf("forbidden: %q", sender.last())
	}

	user := botTestUser(t, store, "7001")
	botTestRole(t, store, "smm_bot", user.ID, []domain.Permission{domain.PermManageTours})

	if err := svc.HandleInboundUpdate(context.Background(), TelegramInboundUpdate{
		Username: "manager",
		UserID:   "7001",
		ChatID:   "7001",
		Text:     "/tour test-tour price 15000",
	}); err != nil {
		t.Fatalf("price: %v", err)
	}
	if !strings.Contains(sender.last(), "15000") {
		t.Fatalf("price reply: %q", sender.last())
	}

	if err := svc.HandleInboundUpdate(context.Background(), TelegramInboundUpdate{
		Username: "manager",
		UserID:   "7001",
		ChatID:   "7001",
		Text:     "/tour test-tour slots 3 8",
	}); err != nil {
		t.Fatalf("slots: %v", err)
	}
	saved, err := store.GetTour(context.Background(), tour.ID)
	if err != nil {
		t.Fatalf("get tour: %v", err)
	}
	if saved.Price != 15000 || saved.SlotsLeft != 3 || saved.SlotsTotal != 8 {
		t.Fatalf("tour ops: %+v", saved)
	}

	if err := svc.HandleInboundUpdate(context.Background(), TelegramInboundUpdate{
		Username: "manager",
		UserID:   "7001",
		ChatID:   "7001",
		Text:     "/tour test-tour off",
	}); err != nil {
		t.Fatalf("off: %v", err)
	}
	saved, _ = store.GetTour(context.Background(), tour.ID)
	if saved.IsActive {
		t.Fatal("expected tour off")
	}
}

func TestBotBookingsMaskPhone(t *testing.T) {
	store := memory.NewStore()
	sender := &recordingBotSender{}
	svc := newBotService(t, store, sender, &recordingMessenger{})

	tour := testTour()
	if _, err := store.CreateTour(context.Background(), tour); err != nil {
		t.Fatalf("create tour: %v", err)
	}
	bookings := NewBookingService(store, store, noop.NewCRMAdapter(), noop.NewAccountingAdapter(), notificationnoop.New(), store, 0)
	created, err := bookings.CreateBooking(context.Background(), CreateBookingInput{
		TourID:      tour.ID,
		Name:        "Иван Иванов",
		Phone:       "+79991234567",
		PeopleCount: 1,
	})
	if err != nil {
		t.Fatalf("booking: %v", err)
	}

	user := botTestUser(t, store, "8001")
	botTestRole(t, store, "book_bot", user.ID, []domain.Permission{domain.PermManageBookings})

	if err := svc.HandleInboundUpdate(context.Background(), TelegramInboundUpdate{
		Username: "booker",
		UserID:   "8001",
		ChatID:   "8001",
		Text:     "/booking " + created.Booking.ID.String(),
	}); err != nil {
		t.Fatalf("inbound: %v", err)
	}
	if strings.Contains(sender.last(), "79991234567") || strings.Contains(sender.last(), "+79991234567") {
		t.Fatalf("phone must be masked: %q", sender.last())
	}
	if !strings.Contains(sender.last(), domain.MaskPhone("+79991234567")) {
		t.Fatalf("expected masked phone in %q", sender.last())
	}
}

func TestStaffReplyFansOutToBoundIdentity(t *testing.T) {
	store := memory.NewStore()
	messenger := &recordingMessenger{configured: true}
	support := NewSupportService(store, nil).WithClientMessenger(messenger, store)

	user := botTestUser(t, store, "555")
	if _, err := support.SendUserMessage(context.Background(), user.ID, "Вопрос с сайта"); err != nil {
		t.Fatalf("user: %v", err)
	}
	threads, err := support.ListThreads(context.Background())
	if err != nil || len(threads) != 1 {
		t.Fatalf("threads: %v", err)
	}
	if _, _, err := support.SendStaffMessage(context.Background(), threads[0].ID, "Ответ менеджера"); err != nil {
		t.Fatalf("staff: %v", err)
	}
	messenger.mu.Lock()
	defer messenger.mu.Unlock()
	if len(messenger.sent) != 1 {
		t.Fatalf("fan-out: %+v", messenger.sent)
	}
	if messenger.sent[0][0] != "telegram" || messenger.sent[0][1] != "555" || messenger.sent[0][2] != "Ответ менеджера" {
		t.Fatalf("sent: %+v", messenger.sent[0])
	}
}

func TestStaffReplyNoFanOutWhenMessengerNoop(t *testing.T) {
	store := memory.NewStore()
	messenger := &recordingMessenger{configured: false}
	support := NewSupportService(store, nil).WithClientMessenger(messenger, store)
	user := botTestUser(t, store, "556")
	if _, err := support.SendUserMessage(context.Background(), user.ID, "Вопрос"); err != nil {
		t.Fatalf("user: %v", err)
	}
	threads, _ := support.ListThreads(context.Background())
	if _, _, err := support.SendStaffMessage(context.Background(), threads[0].ID, "Ответ"); err != nil {
		t.Fatalf("staff: %v", err)
	}
	if len(messenger.sent) != 0 {
		t.Fatalf("noop must not send: %+v", messenger.sent)
	}
}
