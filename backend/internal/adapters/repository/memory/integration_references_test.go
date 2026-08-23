package memory

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"palomnik/internal/domain"
	"palomnik/internal/ports"
)

func TestStoreUpsertIntegrationReference(t *testing.T) {
	ctx := context.Background()
	store := NewStore()
	localID := uuid.New()

	first, err := domain.NewIntegrationReference(domain.NewIntegrationReferenceInput{
		ID:                 uuid.New(),
		LocalEntityType:    domain.EntityTypeBooking,
		LocalEntityID:      localID,
		ExternalSystem:     domain.ExternalSystemBitrix24,
		ExternalEntityType: domain.ExternalEntityTypeDeal,
		ExternalEntityID:   domain.PlaceholderExternalID(localID),
		SyncStatus:         domain.SyncStatusPending,
		LastError:          "waiting",
		Now:                time.Date(2026, 6, 10, 10, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("create reference: %v", err)
	}

	if _, err := store.UpsertReference(ctx, first); err != nil {
		t.Fatalf("upsert reference: %v", err)
	}

	second, err := first.WithSyncUpdate(
		domain.SyncStatusSynced,
		"deal-42",
		"",
		time.Date(2026, 6, 10, 11, 0, 0, 0, time.UTC),
	)
	if err != nil {
		t.Fatalf("update reference: %v", err)
	}

	if _, err := store.UpsertReference(ctx, second); err != nil {
		t.Fatalf("upsert updated reference: %v", err)
	}

	stored, err := store.GetReference(ctx, ports.IntegrationReferenceQuery{
		LocalEntityType:    domain.EntityTypeBooking,
		LocalEntityID:      localID,
		ExternalSystem:     domain.ExternalSystemBitrix24,
		ExternalEntityType: domain.ExternalEntityTypeDeal,
	})
	if err != nil {
		t.Fatalf("get reference: %v", err)
	}
	if stored.SyncStatus != domain.SyncStatusSynced {
		t.Fatalf("expected synced status, got %s", stored.SyncStatus)
	}
	if stored.ExternalEntityID != "deal-42" {
		t.Fatalf("expected external id deal-42, got %q", stored.ExternalEntityID)
	}
}
