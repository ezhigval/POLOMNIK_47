package ports

import (
	"context"

	"github.com/google/uuid"

	"palomnik/internal/domain"
)

type LegalDocumentFilters struct {
	Type         domain.LegalDocumentType
	ActiveOnly   bool
	IncludeAll   bool
}

type LegalDocumentRepository interface {
	ListDocuments(ctx context.Context, filters LegalDocumentFilters) ([]domain.LegalDocument, error)
	GetDocument(ctx context.Context, id uuid.UUID) (domain.LegalDocument, error)
	GetActiveByType(ctx context.Context, docType domain.LegalDocumentType) (domain.LegalDocument, error)
	GetByTypeAndVersion(ctx context.Context, docType domain.LegalDocumentType, version string) (domain.LegalDocument, error)
	CreateDocument(ctx context.Context, doc domain.LegalDocument) (domain.LegalDocument, error)
	DeactivateType(ctx context.Context, docType domain.LegalDocumentType) error
	HasAnyDocuments(ctx context.Context) (bool, error)
}

type ConsentFilters struct {
	UserID      *uuid.UUID
	RequestID   *uuid.UUID
	ConsentType domain.ConsentType
}

type ConsentRepository interface {
	CreateConsent(ctx context.Context, consent domain.Consent) (domain.Consent, error)
	ListConsents(ctx context.Context, filters ConsentFilters, pagination Pagination) (ConsentList, error)
}

type ConsentList struct {
	Items []domain.Consent
	Meta  PageMeta
}
