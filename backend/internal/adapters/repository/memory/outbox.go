package memory

import (
	"context"
	"encoding/json"
	"sort"

	"github.com/google/uuid"

	"polomnik/internal/domain"
	"polomnik/internal/ports"
)

func (s *Store) Enqueue(_ context.Context, event domain.OutboxEvent) (domain.OutboxEvent, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for id, existing := range s.outboxEvents {
		if existing.Status == domain.OutboxStatusPending &&
			existing.EventType == event.EventType &&
			existing.EntityID == event.EntityID {
			updated := cloneOutboxEvent(event)
			updated.ID = id
			updated.Attempts = existing.Attempts
			updated.CreatedAt = existing.CreatedAt
			s.outboxEvents[id] = updated
			return cloneOutboxEvent(updated), nil
		}
	}

	cloned := cloneOutboxEvent(event)
	s.outboxEvents[cloned.ID] = cloned
	return cloneOutboxEvent(cloned), nil
}

func (s *Store) ListEvents(_ context.Context, filters ports.OutboxFilters, pagination ports.Pagination) (ports.OutboxList, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	items := make([]domain.OutboxEvent, 0, len(s.outboxEvents))
	for _, event := range s.outboxEvents {
		if !matchesOutboxFilters(event, filters) {
			continue
		}
		items = append(items, cloneOutboxEvent(event))
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].CreatedAt.Equal(items[j].CreatedAt) {
			return items[i].ID.String() < items[j].ID.String()
		}
		return items[i].CreatedAt.After(items[j].CreatedAt)
	})

	pageItems, meta := page(items, pagination)
	return ports.OutboxList{Items: pageItems, Meta: meta}, nil
}

func (s *Store) MarkProcessed(_ context.Context, id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	event, ok := s.outboxEvents[id]
	if !ok {
		return domain.ErrNotFound
	}
	event.Status = domain.OutboxStatusProcessed
	s.outboxEvents[id] = event
	return nil
}

func (s *Store) ListPending(_ context.Context, limit int) ([]domain.OutboxEvent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if limit <= 0 {
		limit = 10
	}

	items := make([]domain.OutboxEvent, 0)
	for _, event := range s.outboxEvents {
		if event.Status != domain.OutboxStatusPending {
			continue
		}
		items = append(items, cloneOutboxEvent(event))
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].CreatedAt.Equal(items[j].CreatedAt) {
			return items[i].ID.String() < items[j].ID.String()
		}
		return items[i].CreatedAt.Before(items[j].CreatedAt)
	})

	if len(items) > limit {
		items = items[:limit]
	}
	return items, nil
}

func (s *Store) Summarize(_ context.Context) (ports.OutboxSummary, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var summary ports.OutboxSummary
	for _, event := range s.outboxEvents {
		switch event.Status {
		case domain.OutboxStatusPending:
			summary.Pending++
			if summary.OldestPendingAt == nil || event.CreatedAt.Before(*summary.OldestPendingAt) {
				created := event.CreatedAt
				summary.OldestPendingAt = &created
			}
		case domain.OutboxStatusFailed:
			summary.Failed++
			if summary.LatestFailedAt == nil || event.UpdatedAt.After(*summary.LatestFailedAt) {
				updated := event.UpdatedAt
				summary.LatestFailedAt = &updated
				summary.LatestFailedError = event.LastError
			}
		case domain.OutboxStatusProcessed:
			summary.Processed++
		}
	}
	return summary, nil
}

func (s *Store) RecordAttempt(_ context.Context, id uuid.UUID, lastError string, maxAttempts int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if maxAttempts <= 0 {
		maxAttempts = 1
	}

	event, ok := s.outboxEvents[id]
	if !ok {
		return domain.ErrNotFound
	}
	if event.Status != domain.OutboxStatusPending {
		return nil
	}

	event.Attempts++
	event.LastError = lastError
	if event.Attempts >= maxAttempts {
		event.Status = domain.OutboxStatusFailed
	}
	s.outboxEvents[id] = event
	return nil
}

func matchesOutboxFilters(event domain.OutboxEvent, filters ports.OutboxFilters) bool {
	if filters.Status != "" && string(event.Status) != filters.Status {
		return false
	}
	if filters.EntityType != "" && event.EntityType != filters.EntityType {
		return false
	}
	if filters.EventType != "" && event.EventType != filters.EventType {
		return false
	}
	return true
}

func cloneOutboxEvent(event domain.OutboxEvent) domain.OutboxEvent {
	cloned := event
	if len(event.Payload) > 0 {
		cloned.Payload = append(json.RawMessage(nil), event.Payload...)
	}
	return cloned
}
