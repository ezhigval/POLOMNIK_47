package application

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"

	"palomnik/internal/adapters/integration/noop"
	notificationnoop "palomnik/internal/adapters/notification/noop"
	"palomnik/internal/adapters/repository/memory"
	"palomnik/internal/domain"
	"palomnik/internal/ports"
)

type stubAI struct {
	configured bool
	text       string
	err        error
	last       ports.AICompletionRequest
}

func (s *stubAI) Configured() bool { return s != nil && s.configured }

func (s *stubAI) Complete(_ context.Context, req ports.AICompletionRequest) (string, error) {
	s.last = req
	if s.err != nil {
		return "", s.err
	}
	return s.text, nil
}

func TestSupportDraftNoopEscalatesWithoutWriting(t *testing.T) {
	store := memory.NewStore()
	support := NewSupportService(store, nil)
	userID := uuid.New()
	if _, err := support.SendUserMessage(context.Background(), userID, "Когда ближайший тур?"); err != nil {
		t.Fatalf("message: %v", err)
	}
	thread, err := support.GetOrCreateThread(context.Background(), userID)
	if err != nil {
		t.Fatalf("thread: %v", err)
	}

	svc := NewAIFeaturesService(&stubAI{}, nil, nil, support, nil, nil)
	draft, err := svc.SupportDraft(context.Background(), thread.ID)
	if err != nil {
		t.Fatalf("draft: %v", err)
	}
	if draft.Configured || !draft.Escalate || draft.Draft != "" {
		t.Fatalf("expected noop escalate, got %+v", draft)
	}

	_, messages, err := support.GetThreadByID(context.Background(), thread.ID)
	if err != nil {
		t.Fatalf("messages: %v", err)
	}
	for _, message := range messages {
		if strings.Contains(message.Body, "черновик модели") {
			t.Fatalf("draft leaked into thread: %+v", message)
		}
	}
}

func TestSupportDraftConfiguredDoesNotPersist(t *testing.T) {
	store := memory.NewStore()
	support := NewSupportService(store, nil)
	userID := uuid.New()
	if _, err := support.SendUserMessage(context.Background(), userID, "Есть места?"); err != nil {
		t.Fatalf("message: %v", err)
	}
	thread, err := support.GetOrCreateThread(context.Background(), userID)
	if err != nil {
		t.Fatalf("thread: %v", err)
	}

	ai := &stubAI{configured: true, text: "Уточните даты у менеджера."}
	svc := NewAIFeaturesService(ai, nil, nil, support, nil, nil)
	draft, err := svc.SupportDraft(context.Background(), thread.ID)
	if err != nil {
		t.Fatalf("draft: %v", err)
	}
	if !draft.Configured || !draft.Escalate || draft.Draft != "Уточните даты у менеджера." {
		t.Fatalf("unexpected draft %+v", draft)
	}
	if !strings.Contains(ai.last.System, "Не выдумывай цены") {
		t.Fatalf("system prompt missing price rule: %q", ai.last.System)
	}

	_, messages, err := support.GetThreadByID(context.Background(), thread.ID)
	if err != nil {
		t.Fatalf("messages: %v", err)
	}
	for _, message := range messages {
		if message.Body == draft.Draft && message.SenderType == domain.SupportSenderStaff {
			t.Fatalf("AI draft was stored as staff message")
		}
	}
}

func TestRecommendPublishedToursFiltersInactiveAndUnknown(t *testing.T) {
	store := memory.NewStore()
	tours := NewTourService(store, nil, noop.NewCRMAdapter())
	seed := testTour()
	if _, err := store.CreateTour(context.Background(), seed); err != nil {
		t.Fatalf("seed: %v", err)
	}

	other := testTour(func(in *domain.NewTourInput) {
		in.ID = testUUID("22222222-2222-2222-2222-222222222222")
		in.Slug = "other-tour"
		in.Title = "Other Tour"
	})
	if _, err := store.CreateTour(context.Background(), other); err != nil {
		t.Fatalf("other: %v", err)
	}

	hidden := testTour(func(in *domain.NewTourInput) {
		in.ID = testUUID("33333333-3333-3333-3333-333333333333")
		in.Slug = "hidden-tour"
		in.Title = "Hidden Tour"
		in.IsActive = false
	})
	if _, err := store.CreateTour(context.Background(), hidden); err != nil {
		t.Fatalf("hidden: %v", err)
	}

	ai := &stubAI{configured: true, text: hidden.ID.String() + ", 99999999-9999-9999-9999-999999999999, " + other.ID.String()}
	svc := NewAIFeaturesService(ai, tours, nil, nil, nil, nil)
	got, err := svc.RecommendPublishedTours(context.Background(), seed.ID, 3)
	if err != nil {
		t.Fatalf("recommend: %v", err)
	}
	if len(got) != 1 || got[0].ID != other.ID {
		t.Fatalf("expected only published other tour, got %+v", got)
	}
}

func TestRecommendPublishedToursNoopReturnsPublishedOnly(t *testing.T) {
	store := memory.NewStore()
	tours := NewTourService(store, nil, noop.NewCRMAdapter())
	seed := testTour()
	if _, err := store.CreateTour(context.Background(), seed); err != nil {
		t.Fatalf("seed: %v", err)
	}
	other := testTour(func(in *domain.NewTourInput) {
		in.ID = testUUID("22222222-2222-2222-2222-222222222222")
		in.Slug = "other-tour"
		in.Title = "Other Tour"
	})
	if _, err := store.CreateTour(context.Background(), other); err != nil {
		t.Fatalf("other: %v", err)
	}

	svc := NewAIFeaturesService(&stubAI{}, tours, nil, nil, nil, nil)
	got, err := svc.RecommendPublishedTours(context.Background(), seed.ID, 3)
	if err != nil {
		t.Fatalf("recommend: %v", err)
	}
	if len(got) != 1 || got[0].ID != other.ID || !got[0].IsActive {
		t.Fatalf("expected published fallback, got %+v", got)
	}
}

func TestMetricsDigestHasNoVisits(t *testing.T) {
	store := memory.NewStore()
	tours := NewTourService(store, nil, noop.NewCRMAdapter())
	bookings := NewBookingService(store, store, noop.NewCRMAdapter(), noop.NewAccountingAdapter(), notificationnoop.New(), store, 0)
	seed := testTour()
	if _, err := store.CreateTour(context.Background(), seed); err != nil {
		t.Fatalf("tour: %v", err)
	}
	if _, err := bookings.CreateBooking(context.Background(), CreateBookingInput{
		TourID: seed.ID, Name: "Иван", Phone: "+79990000001", PeopleCount: 1,
	}); err != nil {
		t.Fatalf("booking: %v", err)
	}

	svc := NewAIFeaturesService(&stubAI{configured: true, text: "1 заявка NEW"}, tours, bookings, NewSupportService(store, nil), NewIntegrationService(store, store), nil)
	digest, err := svc.MetricsDigest(context.Background())
	if err != nil {
		t.Fatalf("digest: %v", err)
	}
	if digest.BookingsByStatus["NEW"] != 1 {
		t.Fatalf("expected NEW=1, got %+v", digest.BookingsByStatus)
	}
	if digest.Summary != "1 заявка NEW" {
		t.Fatalf("summary: %q", digest.Summary)
	}
}

func TestWatchdogNeverRestarts(t *testing.T) {
	fixed := time.Date(2026, 8, 23, 20, 0, 0, 0, time.UTC)
	w := &Watchdog{
		Now:        func() time.Time { return fixed },
		PingDB:     func(context.Context) error { return errors.New("down") },
		Status5xx:  func() uint64 { return 4 },
		BackupPath: "",
	}
	svc := NewAIFeaturesService(&stubAI{configured: true, text: "диск ок"}, nil, nil, nil, nil, w)
	report := svc.WatchdogReport(context.Background())
	if report.RestartAttempted {
		t.Fatal("watchdog must not restart prod")
	}
	if !report.BackupOverdue {
		t.Fatal("expected overdue when backup is missing")
	}
	if report.Database != "error" || report.Status5xx != 4 {
		t.Fatalf("unexpected report %+v", report)
	}
}
