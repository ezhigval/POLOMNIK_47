package application

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"palomnik/internal/domain"
	"palomnik/internal/legal/content"
	"palomnik/internal/legal/operator"
	"palomnik/internal/ports"
)

type LegalDocumentService struct {
	repo ports.LegalDocumentRepository
	op   operator.Config
}

func NewLegalDocumentService(repo ports.LegalDocumentRepository, op operator.Config) *LegalDocumentService {
	return &LegalDocumentService{repo: repo, op: op}
}

func (s *LegalDocumentService) BootstrapInitialDocuments(ctx context.Context) error {
	if s.repo == nil {
		return nil
	}

	now := time.Now().UTC()
	for _, spec := range content.AllInitial(s.op) {
		_, err := s.repo.GetActiveByType(ctx, spec.Type)
		if err == nil {
			continue
		}
		if !errors.Is(err, domain.ErrNotFound) {
			return err
		}
		doc, err := domain.NewLegalDocument(domain.NewLegalDocumentInput{
			ID:          uuid.New(),
			Type:        spec.Type,
			Version:     spec.Version,
			Title:       spec.Title,
			Content:     spec.Content,
			PublishedAt: now,
			IsActive:    true,
			Now:         now,
		})
		if err != nil {
			return err
		}
		if _, err := s.repo.CreateDocument(ctx, doc); err != nil {
			return err
		}
	}
	return nil
}

func (s *LegalDocumentService) ListActiveDocuments(ctx context.Context) ([]domain.LegalDocument, error) {
	return s.repo.ListDocuments(ctx, ports.LegalDocumentFilters{ActiveOnly: true})
}

func (s *LegalDocumentService) GetActiveByType(ctx context.Context, docType domain.LegalDocumentType) (domain.LegalDocument, error) {
	return s.repo.GetActiveByType(ctx, docType)
}

func (s *LegalDocumentService) GetByTypeAndVersion(ctx context.Context, docType domain.LegalDocumentType, version string) (domain.LegalDocument, error) {
	return s.repo.GetByTypeAndVersion(ctx, docType, version)
}

func (s *LegalDocumentService) ListAllVersions(ctx context.Context, docType domain.LegalDocumentType) ([]domain.LegalDocument, error) {
	return s.repo.ListDocuments(ctx, ports.LegalDocumentFilters{Type: docType, IncludeAll: true})
}

type PublishLegalDocumentInput struct {
	Type    domain.LegalDocumentType
	Version string
	Title   string
	Content string
}

func (s *LegalDocumentService) PublishNewVersion(ctx context.Context, input PublishLegalDocumentInput) (domain.LegalDocument, error) {
	now := time.Now().UTC()
	doc, err := domain.NewLegalDocument(domain.NewLegalDocumentInput{
		ID:          uuid.New(),
		Type:        input.Type,
		Version:     input.Version,
		Title:       input.Title,
		Content:     input.Content,
		PublishedAt: now,
		IsActive:    true,
		Now:         now,
	})
	if err != nil {
		return domain.LegalDocument{}, err
	}

	if err := s.repo.DeactivateType(ctx, input.Type); err != nil {
		return domain.LegalDocument{}, err
	}
	return s.repo.CreateDocument(ctx, doc)
}

func (s *LegalDocumentService) ListDocuments(ctx context.Context, filters ports.LegalDocumentFilters) ([]domain.LegalDocument, error) {
	return s.repo.ListDocuments(ctx, filters)
}

func (s *LegalDocumentService) GetDocument(ctx context.Context, id uuid.UUID) (domain.LegalDocument, error) {
	return s.repo.GetDocument(ctx, id)
}
