package application

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"palomnik/internal/domain"
	"palomnik/internal/ports"
)

const (
	supportDraftSystem = "Ты помощник менеджера паломнической службы. Напиши черновик ответа сотруднику. Не выдумывай цены, даты, слоты и богословские утверждения. Если в переписке нет факта — напиши, что менеджер уточнит. Это черновик: клиенту его отправлять нельзя. Эскалация человеку обязательна."
	recommendSystem    = "Выбери только id туров из переданного списка. Ответ — UUID через запятую, без текста. Не добавляй туры вне списка. Не пиши цены, описания и богословие."
	metricsSystem      = "Суммируй только переданные числа. Не выдумывай визиты, конверсии и выручку. Не давай богословских советов."
	watchdogSystem     = "Перечисли только переданные факты. Не предлагай рестарт прода и не меняй инфраструктуру."
)

type SupportDraft struct {
	Configured bool
	Escalate   bool
	Draft      string
	Note       string
}

type MetricsDigest struct {
	Configured         bool
	BookingsByStatus   map[string]int
	ActiveTours        int
	OpenSupportThreads int
	OutboxPending      int
	OutboxFailed       int
	Summary            string
}

type AIWatchdogReport struct {
	WatchdogReport
	Configured bool
	Summary    string
}

type AIFeaturesService struct {
	ai           ports.AIPort
	tours        *TourService
	bookings     *BookingService
	support      *SupportService
	integrations *IntegrationService
	watchdog     *Watchdog
}

func NewAIFeaturesService(
	ai ports.AIPort,
	tours *TourService,
	bookings *BookingService,
	support *SupportService,
	integrations *IntegrationService,
	watchdog *Watchdog,
) *AIFeaturesService {
	return &AIFeaturesService{
		ai:           ai,
		tours:        tours,
		bookings:     bookings,
		support:      support,
		integrations: integrations,
		watchdog:     watchdog,
	}
}

func (s *AIFeaturesService) Configured() bool {
	return s != nil && s.ai != nil && s.ai.Configured()
}

func (s *AIFeaturesService) SupportDraft(ctx context.Context, threadID uuid.UUID) (SupportDraft, error) {
	out := SupportDraft{
		Configured: s.Configured(),
		Escalate:   true,
		Note:       "Черновик для менеджера. Клиенту не отправляется. Ответить должен человек. Цены и богословие не выдумывать.",
	}
	if s == nil || s.support == nil {
		return out, domain.ErrNotFound
	}
	_, messages, err := s.support.GetThreadByID(ctx, threadID)
	if err != nil {
		return SupportDraft{}, err
	}
	if !out.Configured {
		return out, nil
	}
	var b strings.Builder
	for _, message := range messages {
		fmt.Fprintf(&b, "%s: %s\n", message.SenderType, message.Body)
	}
	text, err := s.ai.Complete(ctx, ports.AICompletionRequest{
		System: supportDraftSystem,
		User:   strings.TrimSpace(b.String()),
	})
	if err != nil {
		return out, nil
	}
	out.Draft = strings.TrimSpace(text)
	return out, nil
}

func (s *AIFeaturesService) RecommendPublishedTours(ctx context.Context, seedID uuid.UUID, limit int) ([]domain.Tour, error) {
	if s == nil || s.tours == nil {
		return nil, nil
	}
	if limit < 1 {
		limit = 3
	}
	if limit > 8 {
		limit = 8
	}
	if _, err := s.tours.GetPublicTour(ctx, seedID); err != nil {
		return nil, err
	}
	active := true
	list, err := s.tours.ListPublicTours(ctx, ports.TourFilters{IsActive: &active}, ports.Pagination{Page: 1, Limit: 50})
	if err != nil {
		return nil, err
	}
	published := make([]domain.Tour, 0, len(list.Items))
	byID := make(map[uuid.UUID]domain.Tour, len(list.Items))
	for _, tour := range list.Items {
		if !tour.IsActive || tour.ID == seedID {
			continue
		}
		published = append(published, tour)
		byID[tour.ID] = tour
	}
	if len(published) == 0 {
		return nil, nil
	}
	if !s.Configured() {
		return published[:min(limit, len(published))], nil
	}
	var catalog strings.Builder
	for _, tour := range published {
		fmt.Fprintf(
			&catalog,
			"%s\t%s\t%s\t%d\t%s\n",
			tour.ID,
			tour.Title,
			tour.DateStart.UTC().Format("2006-01-02"),
			tour.Price,
			tour.Location,
		)
	}
	text, err := s.ai.Complete(ctx, ports.AICompletionRequest{
		System: recommendSystem,
		User:   catalog.String(),
	})
	if err != nil || strings.TrimSpace(text) == "" {
		return published[:min(limit, len(published))], nil
	}
	picked := make([]domain.Tour, 0, limit)
	seen := map[uuid.UUID]bool{}
	for _, part := range strings.FieldsFunc(text, func(r rune) bool {
		return r == ',' || r == ' ' || r == '\n' || r == ';'
	}) {
		id, parseErr := uuid.Parse(strings.TrimSpace(part))
		if parseErr != nil || id == seedID || seen[id] {
			continue
		}
		tour, ok := byID[id]
		if !ok || !tour.IsActive {
			continue
		}
		seen[id] = true
		picked = append(picked, tour)
		if len(picked) >= limit {
			break
		}
	}
	if len(picked) == 0 {
		return published[:min(limit, len(published))], nil
	}
	return picked, nil
}

func (s *AIFeaturesService) MetricsDigest(ctx context.Context) (MetricsDigest, error) {
	out := MetricsDigest{
		Configured:       s.Configured(),
		BookingsByStatus: map[string]int{},
	}
	if s == nil {
		return out, nil
	}
	if s.bookings != nil {
		bookings, err := s.bookings.ListAllBookings(ctx, ports.BookingFilters{})
		if err != nil {
			return MetricsDigest{}, err
		}
		for _, booking := range bookings {
			out.BookingsByStatus[string(booking.Status)]++
		}
	}
	if s.tours != nil {
		active := true
		list, err := s.tours.ListPublicTours(ctx, ports.TourFilters{IsActive: &active}, ports.Pagination{Page: 1, Limit: 1})
		if err != nil {
			return MetricsDigest{}, err
		}
		out.ActiveTours = list.Meta.Total
	}
	if s.support != nil {
		threads, err := s.support.ListThreads(ctx)
		if err != nil {
			return MetricsDigest{}, err
		}
		for _, thread := range threads {
			if thread.Status == "open" {
				out.OpenSupportThreads++
			}
		}
	}
	if s.integrations != nil {
		summary, err := s.integrations.OutboxSummary(ctx)
		if err != nil {
			return MetricsDigest{}, err
		}
		out.OutboxPending = summary.Pending
		out.OutboxFailed = summary.Failed
	}
	if !out.Configured {
		return out, nil
	}
	var b strings.Builder
	fmt.Fprintf(&b, "active_tours=%d\n", out.ActiveTours)
	fmt.Fprintf(&b, "open_support_threads=%d\n", out.OpenSupportThreads)
	fmt.Fprintf(&b, "outbox_pending=%d outbox_failed=%d\n", out.OutboxPending, out.OutboxFailed)
	for _, status := range []domain.BookingStatus{
		domain.BookingStatusNew,
		domain.BookingStatusContacted,
		domain.BookingStatusConfirmed,
		domain.BookingStatusCompleted,
		domain.BookingStatusCancelled,
	} {
		fmt.Fprintf(&b, "bookings_%s=%d\n", status, out.BookingsByStatus[string(status)])
	}
	text, err := s.ai.Complete(ctx, ports.AICompletionRequest{System: metricsSystem, User: b.String()})
	if err != nil {
		return out, nil
	}
	out.Summary = strings.TrimSpace(text)
	return out, nil
}

func (s *AIFeaturesService) WatchdogReport(ctx context.Context) AIWatchdogReport {
	out := AIWatchdogReport{Configured: s.Configured()}
	if s != nil && s.watchdog != nil {
		out.WatchdogReport = s.watchdog.Snapshot(ctx)
	} else {
		out.WatchdogReport = WatchdogReport{At: time.Now().UTC(), RestartAttempted: false}
	}
	out.RestartAttempted = false
	if !out.Configured {
		return out
	}
	user := fmt.Sprintf(
		"database=%s disk_percent=%d outbox_pending=%d outbox_failed=%d status_5xx=%d backup_overdue=%t restart_attempted=false",
		out.Database, out.DiskPercent, out.OutboxPending, out.OutboxFailed, out.Status5xx, out.BackupOverdue,
	)
	text, err := s.ai.Complete(ctx, ports.AICompletionRequest{System: watchdogSystem, User: user})
	if err != nil {
		return out
	}
	out.Summary = strings.TrimSpace(text)
	return out
}
