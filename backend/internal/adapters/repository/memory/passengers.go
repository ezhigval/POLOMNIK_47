package memory

import (
	"context"
	"sort"
	"time"

	"github.com/google/uuid"

	"palomnik/internal/domain"
)

func (s *Store) ListPassengers(_ context.Context, userID uuid.UUID) ([]domain.Passenger, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	items := make([]domain.Passenger, 0)
	for _, passenger := range s.passengers {
		if passenger.UserID == userID {
			items = append(items, passenger)
		}
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].CreatedAt.Equal(items[j].CreatedAt) {
			return items[i].ID.String() < items[j].ID.String()
		}
		return items[i].CreatedAt.Before(items[j].CreatedAt)
	})
	return items, nil
}

func (s *Store) GetPassenger(_ context.Context, userID, id uuid.UUID) (domain.Passenger, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	passenger, ok := s.passengers[id]
	if !ok || passenger.UserID != userID {
		return domain.Passenger{}, domain.ErrNotFound
	}
	return passenger, nil
}

func (s *Store) CreatePassenger(_ context.Context, passenger domain.Passenger) (domain.Passenger, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.passengers[passenger.ID] = passenger
	return passenger, nil
}

func (s *Store) UpdatePassenger(_ context.Context, passenger domain.Passenger) (domain.Passenger, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing, ok := s.passengers[passenger.ID]
	if !ok || existing.UserID != passenger.UserID {
		return domain.Passenger{}, domain.ErrNotFound
	}
	s.passengers[passenger.ID] = passenger
	return passenger, nil
}

func (s *Store) DeletePassenger(_ context.Context, userID, id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing, ok := s.passengers[id]
	if !ok || existing.UserID != userID {
		return domain.ErrNotFound
	}
	delete(s.passengers, id)
	return nil
}

func mergePassengersLocked(s *Store, targetID, sourceID uuid.UUID) {
	for id, passenger := range s.passengers {
		if passenger.UserID == sourceID {
			passenger.UserID = targetID
			passenger.UpdatedAt = time.Now().UTC()
			s.passengers[id] = passenger
		}
	}
}
