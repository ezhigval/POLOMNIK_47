package memory

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"polomnik/internal/domain"
	"polomnik/internal/ports"
)

func (s *Store) UpsertReference(_ context.Context, ref domain.IntegrationReference) (domain.IntegrationReference, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := integrationReferenceKey(ref.LocalEntityType, ref.LocalEntityID, ref.ExternalSystem, ref.ExternalEntityType)
	if existing, ok := s.integrationRefs[key]; ok {
		updated := ref
		updated.ID = existing.ID
		updated.CreatedAt = existing.CreatedAt
		s.integrationRefs[key] = cloneIntegrationReference(updated)
		return cloneIntegrationReference(updated), nil
	}

	s.integrationRefs[key] = cloneIntegrationReference(ref)
	return cloneIntegrationReference(ref), nil
}

func (s *Store) GetReference(_ context.Context, query ports.IntegrationReferenceQuery) (domain.IntegrationReference, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	key := integrationReferenceKey(
		query.LocalEntityType,
		query.LocalEntityID,
		query.ExternalSystem,
		query.ExternalEntityType,
	)
	ref, ok := s.integrationRefs[key]
	if !ok {
		return domain.IntegrationReference{}, domain.ErrNotFound
	}
	return cloneIntegrationReference(ref), nil
}

func (s *Store) ListReferences(_ context.Context, filters ports.IntegrationReferenceFilters, pagination ports.Pagination) (ports.IntegrationReferenceList, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	items := make([]domain.IntegrationReference, 0, len(s.integrationRefs))
	for _, ref := range s.integrationRefs {
		if !matchesIntegrationReferenceFilters(ref, filters) {
			continue
		}
		items = append(items, cloneIntegrationReference(ref))
	}

	sortIntegrationReferences(items)
	pageItems, meta := page(items, pagination)
	return ports.IntegrationReferenceList{Items: pageItems, Meta: meta}, nil
}

func matchesIntegrationReferenceFilters(ref domain.IntegrationReference, filters ports.IntegrationReferenceFilters) bool {
	if filters.ExternalSystem != "" && ref.ExternalSystem != filters.ExternalSystem {
		return false
	}
	if filters.LocalEntityType != "" && ref.LocalEntityType != filters.LocalEntityType {
		return false
	}
	if filters.SyncStatus != "" && string(ref.SyncStatus) != filters.SyncStatus {
		return false
	}
	return true
}

func sortIntegrationReferences(items []domain.IntegrationReference) {
	for i := 0; i < len(items); i++ {
		for j := i + 1; j < len(items); j++ {
			if items[j].UpdatedAt.After(items[i].UpdatedAt) {
				items[i], items[j] = items[j], items[i]
			}
		}
	}
}

func integrationReferenceKey(
	localEntityType string,
	localEntityID uuid.UUID,
	externalSystem string,
	externalEntityType string,
) string {
	return fmt.Sprintf(
		"%s:%s:%s:%s",
		localEntityType,
		localEntityID,
		externalSystem,
		externalEntityType,
	)
}

func cloneIntegrationReference(ref domain.IntegrationReference) domain.IntegrationReference {
	cloned := ref
	if ref.LastSyncAt != nil {
		syncAt := *ref.LastSyncAt
		cloned.LastSyncAt = &syncAt
	}
	return cloned
}
