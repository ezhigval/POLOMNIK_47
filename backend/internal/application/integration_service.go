package application

import (
	"context"

	"palomnik/internal/ports"
)

type IntegrationService struct {
	refs   ports.IntegrationReferenceRepository
	outbox ports.OutboxRepository
}

func NewIntegrationService(refs ports.IntegrationReferenceRepository, outbox ports.OutboxRepository) *IntegrationService {
	return &IntegrationService{refs: refs, outbox: outbox}
}

func (s *IntegrationService) ListReferences(
	ctx context.Context,
	filters ports.IntegrationReferenceFilters,
	pagination ports.Pagination,
) (ports.IntegrationReferenceList, error) {
	return s.refs.ListReferences(ctx, filters, pagination)
}

func (s *IntegrationService) ListOutboxEvents(
	ctx context.Context,
	filters ports.OutboxFilters,
	pagination ports.Pagination,
) (ports.OutboxList, error) {
	return s.outbox.ListEvents(ctx, filters, pagination)
}

func (s *IntegrationService) OutboxSummary(ctx context.Context) (ports.OutboxSummary, error) {
	return s.outbox.Summarize(ctx)
}
