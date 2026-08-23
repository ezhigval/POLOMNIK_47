package ports

import (
	"context"
	"time"

	"github.com/google/uuid"

	"palomnik/internal/domain"
)

type OutboxFilters struct {
	Status     string
	EntityType string
	EventType  string
}

type OutboxList struct {
	Items []domain.OutboxEvent
	Meta  PageMeta
}

type OutboxSummary struct {
	Pending           int
	Failed            int
	Processed         int
	OldestPendingAt   *time.Time
	LatestFailedAt    *time.Time
	LatestFailedError string
}

type OutboxRepository interface {
	Enqueue(ctx context.Context, event domain.OutboxEvent) (domain.OutboxEvent, error)
	ListEvents(ctx context.Context, filters OutboxFilters, pagination Pagination) (OutboxList, error)
	ListPending(ctx context.Context, limit int) ([]domain.OutboxEvent, error)
	Summarize(ctx context.Context) (OutboxSummary, error)
	MarkProcessed(ctx context.Context, id uuid.UUID) error
	RecordAttempt(ctx context.Context, id uuid.UUID, lastError string, maxAttempts int) error
}
