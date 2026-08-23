package application

import (
	"context"
	"time"

	"github.com/google/uuid"

	"palomnik/internal/domain"
	"palomnik/internal/ports"
)

type ConsentService struct {
	consents  ports.ConsentRepository
	documents ports.LegalDocumentRepository
}

func NewConsentService(consents ports.ConsentRepository, documents ports.LegalDocumentRepository) *ConsentService {
	return &ConsentService{consents: consents, documents: documents}
}

type RecordConsentInput struct {
	ConsentType domain.ConsentType
	UserID      *uuid.UUID
	RequestID   *uuid.UUID
	IP          string
	UserAgent   string
}

func documentTypeForConsent(consentType domain.ConsentType) domain.LegalDocumentType {
	switch consentType {
	case domain.ConsentTypeMarketing:
		return domain.LegalDocTypeMarketing
	case domain.ConsentTypeDistribution:
		return domain.LegalDocTypeDistribution
	case domain.ConsentTypeCookieAll, domain.ConsentTypeCookieEssential, domain.ConsentTypeCookieReject:
		return domain.LegalDocTypeCookie
	default:
		return domain.LegalDocTypePersonalData
	}
}

func (s *ConsentService) RecordConsent(ctx context.Context, input RecordConsentInput) (domain.Consent, error) {
	docType := documentTypeForConsent(input.ConsentType)
	doc, err := s.documents.GetActiveByType(ctx, docType)
	if err != nil {
		return domain.Consent{}, err
	}

	now := time.Now().UTC()
	consent, err := domain.NewConsent(domain.NewConsentInput{
		ID:              uuid.New(),
		UserID:          input.UserID,
		RequestID:       input.RequestID,
		ConsentType:     input.ConsentType,
		DocumentID:      doc.ID,
		DocumentVersion: doc.Version,
		AcceptedAt:      now,
		IP:              input.IP,
		UserAgent:       input.UserAgent,
	})
	if err != nil {
		return domain.Consent{}, err
	}
	return s.consents.CreateConsent(ctx, consent)
}

func (s *ConsentService) ListUserConsents(ctx context.Context, userID uuid.UUID, pagination ports.Pagination) (ports.ConsentList, error) {
	return s.consents.ListConsents(ctx, ports.ConsentFilters{UserID: &userID}, pagination)
}

func (s *ConsentService) ListRequestConsents(ctx context.Context, requestID uuid.UUID) ([]domain.Consent, error) {
	list, err := s.consents.ListConsents(ctx, ports.ConsentFilters{RequestID: &requestID}, ports.Pagination{Page: 1, Limit: 100})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (s *ConsentService) ListManagementConsents(ctx context.Context, filters ports.ConsentFilters, pagination ports.Pagination) (ports.ConsentList, error) {
	return s.consents.ListConsents(ctx, filters, pagination)
}

func (s *ConsentService) GetActiveDocumentForConsent(ctx context.Context, consentType domain.ConsentType) (domain.LegalDocument, error) {
	return s.documents.GetActiveByType(ctx, documentTypeForConsent(consentType))
}
