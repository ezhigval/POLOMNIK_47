package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"polomnik/internal/domain"
	"polomnik/internal/ports"
)

func (s *Store) Enqueue(ctx context.Context, event domain.OutboxEvent) (domain.OutboxEvent, error) {
	row := s.db.QueryRowContext(ctx, `
INSERT INTO outbox_events (
    id, event_type, entity_type, entity_id, payload, status, attempts, last_error, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
ON CONFLICT (event_type, entity_id) WHERE (status = 'pending')
DO UPDATE SET
    payload = EXCLUDED.payload,
    last_error = EXCLUDED.last_error,
    updated_at = EXCLUDED.updated_at
RETURNING id, event_type, entity_type, entity_id, payload, status, attempts, last_error, created_at, updated_at
`, event.ID, event.EventType, event.EntityType, event.EntityID, event.Payload,
		event.Status, event.Attempts, event.LastError, event.CreatedAt, event.UpdatedAt)
	return scanOutboxEvent(row)
}

func (s *Store) ListEvents(ctx context.Context, filters ports.OutboxFilters, pagination ports.Pagination) (ports.OutboxList, error) {
	pagination = ports.NormalizePagination(pagination.Page, pagination.Limit)
	args := outboxFilterArgs(filters)

	total, err := count(ctx, s.db, "SELECT COUNT(*) FROM outbox_events WHERE "+outboxWhereClause, args...)
	if err != nil {
		return ports.OutboxList{}, err
	}

	rows, err := s.db.QueryContext(ctx, `
SELECT id, event_type, entity_type, entity_id, payload, status, attempts, last_error, created_at, updated_at
FROM outbox_events
WHERE `+outboxWhereClause+`
ORDER BY created_at DESC, id ASC
LIMIT $4 OFFSET $5
`, append(args, pagination.Limit, offset(pagination))...)
	if err != nil {
		return ports.OutboxList{}, fmt.Errorf("list outbox events: %w", err)
	}
	defer rows.Close()

	items, err := scanOutboxEvents(rows)
	if err != nil {
		return ports.OutboxList{}, err
	}

	return ports.OutboxList{Items: items, Meta: pageMeta(pagination, total)}, nil
}

func (s *Store) MarkProcessed(ctx context.Context, id uuid.UUID) error {
	result, err := s.db.ExecContext(ctx, `
UPDATE outbox_events
SET status = 'processed', updated_at = NOW()
WHERE id = $1
`, id)
	if err != nil {
		return fmt.Errorf("mark outbox processed: %w", err)
	}
	return requireAffected(result)
}

func (s *Store) ListPending(ctx context.Context, limit int) ([]domain.OutboxEvent, error) {
	if limit <= 0 {
		limit = 10
	}

	rows, err := s.db.QueryContext(ctx, `
SELECT id, event_type, entity_type, entity_id, payload, status, attempts, last_error, created_at, updated_at
FROM outbox_events
WHERE status = 'pending'
ORDER BY created_at ASC, id ASC
LIMIT $1
`, limit)
	if err != nil {
		return nil, fmt.Errorf("list pending outbox events: %w", err)
	}
	defer rows.Close()

	return scanOutboxEvents(rows)
}

func (s *Store) Summarize(ctx context.Context) (ports.OutboxSummary, error) {
	var summary ports.OutboxSummary
	var oldestPending sql.NullTime

	err := s.db.QueryRowContext(ctx, `
SELECT
    COUNT(*) FILTER (WHERE status = 'pending')::int,
    COUNT(*) FILTER (WHERE status = 'failed')::int,
    COUNT(*) FILTER (WHERE status = 'processed')::int,
    MIN(created_at) FILTER (WHERE status = 'pending')
FROM outbox_events
`).Scan(&summary.Pending, &summary.Failed, &summary.Processed, &oldestPending)
	if err != nil {
		return ports.OutboxSummary{}, fmt.Errorf("summarize outbox: %w", err)
	}
	if oldestPending.Valid {
		t := oldestPending.Time.UTC()
		summary.OldestPendingAt = &t
	}

	var latestFailedAt sql.NullTime
	var latestFailedError sql.NullString
	err = s.db.QueryRowContext(ctx, `
SELECT last_error, updated_at
FROM outbox_events
WHERE status = 'failed'
ORDER BY updated_at DESC, id DESC
LIMIT 1
`).Scan(&latestFailedError, &latestFailedAt)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return ports.OutboxSummary{}, fmt.Errorf("summarize outbox failed sample: %w", err)
	}
	if latestFailedAt.Valid {
		t := latestFailedAt.Time.UTC()
		summary.LatestFailedAt = &t
	}
	if latestFailedError.Valid {
		summary.LatestFailedError = latestFailedError.String
	}

	return summary, nil
}

func (s *Store) RecordAttempt(ctx context.Context, id uuid.UUID, lastError string, maxAttempts int) error {
	if maxAttempts <= 0 {
		maxAttempts = 1
	}

	result, err := s.db.ExecContext(ctx, `
UPDATE outbox_events
SET
    attempts = attempts + 1,
    last_error = $2,
    status = CASE WHEN attempts + 1 >= $3 THEN 'failed' ELSE 'pending' END,
    updated_at = NOW()
WHERE id = $1 AND status = 'pending'
`, id, lastError, maxAttempts)
	if err != nil {
		return fmt.Errorf("record outbox attempt: %w", err)
	}
	return requireAffected(result)
}

const outboxWhereClause = `
($1::text = '' OR status = $1) AND
($2::text = '' OR entity_type = $2) AND
($3::text = '' OR event_type = $3)
`

func outboxFilterArgs(filters ports.OutboxFilters) []any {
	return []any{filters.Status, filters.EntityType, filters.EventType}
}

func scanOutboxEvent(row scanner) (domain.OutboxEvent, error) {
	var event domain.OutboxEvent
	var payload []byte
	err := row.Scan(
		&event.ID,
		&event.EventType,
		&event.EntityType,
		&event.EntityID,
		&payload,
		&event.Status,
		&event.Attempts,
		&event.LastError,
		&event.CreatedAt,
		&event.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.OutboxEvent{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.OutboxEvent{}, fmt.Errorf("scan outbox event: %w", err)
	}
	if len(payload) == 0 {
		event.Payload = json.RawMessage(`{}`)
	} else {
		event.Payload = json.RawMessage(payload)
	}
	return event, nil
}

func scanOutboxEvents(rows *sql.Rows) ([]domain.OutboxEvent, error) {
	var items []domain.OutboxEvent
	for rows.Next() {
		event, err := scanOutboxEvent(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, event)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate outbox events: %w", err)
	}
	return items, nil
}
