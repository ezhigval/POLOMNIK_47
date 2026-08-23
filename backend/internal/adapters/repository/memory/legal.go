package memory

import (
	"context"
	"sort"
	"strings"

	"github.com/google/uuid"

	"palomnik/internal/domain"
	"palomnik/internal/ports"
)

func (s *Store) HasAnyDocuments(_ context.Context) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.legalDocuments) > 0, nil
}

func (s *Store) ListDocuments(_ context.Context, filters ports.LegalDocumentFilters) ([]domain.LegalDocument, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	items := make([]domain.LegalDocument, 0)
	for _, doc := range s.legalDocuments {
		if filters.Type != "" && doc.Type != filters.Type {
			continue
		}
		if filters.ActiveOnly && !filters.IncludeAll && !doc.IsActive {
			continue
		}
		items = append(items, doc)
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].Type != items[j].Type {
			return items[i].Type < items[j].Type
		}
		return items[i].PublishedAt.After(items[j].PublishedAt)
	})
	return items, nil
}

func (s *Store) GetDocument(_ context.Context, id uuid.UUID) (domain.LegalDocument, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	doc, ok := s.legalDocuments[id]
	if !ok {
		return domain.LegalDocument{}, domain.ErrNotFound
	}
	return doc, nil
}

func (s *Store) GetActiveByType(_ context.Context, docType domain.LegalDocumentType) (domain.LegalDocument, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, doc := range s.legalDocuments {
		if doc.Type == docType && doc.IsActive {
			return doc, nil
		}
	}
	return domain.LegalDocument{}, domain.ErrNotFound
}

func (s *Store) GetByTypeAndVersion(_ context.Context, docType domain.LegalDocumentType, version string) (domain.LegalDocument, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, doc := range s.legalDocuments {
		if doc.Type == docType && doc.Version == version {
			return doc, nil
		}
	}
	return domain.LegalDocument{}, domain.ErrNotFound
}

func (s *Store) CreateDocument(_ context.Context, doc domain.LegalDocument) (domain.LegalDocument, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, existing := range s.legalDocuments {
		if existing.Type == doc.Type && existing.Version == doc.Version {
			return domain.LegalDocument{}, domain.ErrDuplicateSlug
		}
	}
	s.legalDocuments[doc.ID] = doc
	return doc, nil
}

func (s *Store) DeactivateType(_ context.Context, docType domain.LegalDocumentType) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for id, doc := range s.legalDocuments {
		if doc.Type == docType && doc.IsActive {
			doc.IsActive = false
			s.legalDocuments[id] = doc
		}
	}
	return nil
}

func (s *Store) CreateConsent(_ context.Context, consent domain.Consent) (domain.Consent, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.consents[consent.ID] = consent
	return consent, nil
}

func (s *Store) ListConsents(_ context.Context, filters ports.ConsentFilters, pagination ports.Pagination) (ports.ConsentList, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	items := make([]domain.Consent, 0)
	for _, c := range s.consents {
		if filters.UserID != nil && (c.UserID == nil || *c.UserID != *filters.UserID) {
			continue
		}
		if filters.RequestID != nil && (c.RequestID == nil || *c.RequestID != *filters.RequestID) {
			continue
		}
		if filters.ConsentType != "" && c.ConsentType != filters.ConsentType {
			continue
		}
		items = append(items, c)
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].AcceptedAt.After(items[j].AcceptedAt)
	})
	pageItems, meta := page(items, pagination)
	return ports.ConsentList{Items: pageItems, Meta: meta}, nil
}

func legalDocKey(docType domain.LegalDocumentType, version string) string {
	return string(docType) + ":" + strings.TrimSpace(version)
}
