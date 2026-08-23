package ports

import (
	"context"

	"github.com/google/uuid"

	"palomnik/internal/domain"
)

type IntegrationReferenceQuery struct {
	LocalEntityType    string
	LocalEntityID      uuid.UUID
	ExternalSystem     string
	ExternalEntityType string
}

type IntegrationReferenceFilters struct {
	ExternalSystem  string
	LocalEntityType string
	SyncStatus      string
}

type IntegrationReferenceList struct {
	Items []domain.IntegrationReference
	Meta  PageMeta
}

type IntegrationReferenceRepository interface {
	UpsertReference(ctx context.Context, ref domain.IntegrationReference) (domain.IntegrationReference, error)
	GetReference(ctx context.Context, query IntegrationReferenceQuery) (domain.IntegrationReference, error)
	ListReferences(ctx context.Context, filters IntegrationReferenceFilters, pagination Pagination) (IntegrationReferenceList, error)
}
