package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewIntegrationReferenceUsesPlaceholderExternalIDPattern(t *testing.T) {
	localID := uuid.New()
	ref, err := NewIntegrationReference(NewIntegrationReferenceInput{
		ID:                 uuid.New(),
		LocalEntityType:    EntityTypeBooking,
		LocalEntityID:      localID,
		ExternalSystem:     ExternalSystemBitrix24,
		ExternalEntityType: ExternalEntityTypeDeal,
		ExternalEntityID:   PlaceholderExternalID(localID),
		SyncStatus:         SyncStatusPending,
		Now:                time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("create integration reference: %v", err)
	}
	if ref.SyncStatus != SyncStatusPending {
		t.Fatalf("expected pending status, got %s", ref.SyncStatus)
	}
}

func TestNewIntegrationReferenceRejectsEmptyExternalEntityID(t *testing.T) {
	_, err := NewIntegrationReference(NewIntegrationReferenceInput{
		ID:                 uuid.New(),
		LocalEntityType:    EntityTypeBooking,
		LocalEntityID:      uuid.New(),
		ExternalSystem:     ExternalSystemOneC,
		ExternalEntityType: ExternalEntityTypeOrder,
		SyncStatus:         SyncStatusNotConfigured,
	})
	if err != ErrInvalidExternalEntityID {
		t.Fatalf("expected invalid external entity id, got %v", err)
	}
}
