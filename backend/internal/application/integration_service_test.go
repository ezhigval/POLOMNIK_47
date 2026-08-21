package application

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"polomnik/internal/adapters/repository/memory"
	"polomnik/internal/domain"
	"polomnik/internal/ports"
)

func TestIntegrationServiceListReferences(t *testing.T) {
	ctx := context.Background()
	store := memory.NewStore()
	service := NewIntegrationService(store, store)

	localID := uuid.New()
	ref, err := domain.NewIntegrationReference(domain.NewIntegrationReferenceInput{
		ID:                 uuid.New(),
		LocalEntityType:    domain.EntityTypeBooking,
		LocalEntityID:      localID,
		ExternalSystem:     domain.ExternalSystemBitrix24,
		ExternalEntityType: domain.ExternalEntityTypeDeal,
		ExternalEntityID:   domain.PlaceholderExternalID(localID),
		SyncStatus:         domain.SyncStatusPending,
		Now:                time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("create reference: %v", err)
	}
	if _, err := store.UpsertReference(ctx, ref); err != nil {
		t.Fatalf("upsert reference: %v", err)
	}

	list, err := service.ListReferences(ctx, ports.IntegrationReferenceFilters{
		ExternalSystem: domain.ExternalSystemBitrix24,
	}, ports.Pagination{Page: 1, Limit: 20})
	if err != nil {
		t.Fatalf("list references: %v", err)
	}
	if len(list.Items) != 1 {
		t.Fatalf("expected 1 reference, got %d", len(list.Items))
	}
}
