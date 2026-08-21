package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	EntityTypeBooking = "booking"
	EntityTypeTour    = "tour"
	EntityTypeReview  = "review"

	ExternalSystemBitrix24 = "bitrix24"
	ExternalSystemOneC     = "onec"

	ExternalEntityTypeDeal    = "deal"
	ExternalEntityTypeOrder   = "order"
	ExternalEntityTypeTour    = "tour"
	ExternalEntityTypeProduct = "product"
	ExternalEntityTypeReview  = "review"
)

type SyncStatus string

const (
	SyncStatusNotConfigured SyncStatus = "not_configured"
	SyncStatusSynced        SyncStatus = "synced"
	SyncStatusPending       SyncStatus = "pending"
	SyncStatusFailed        SyncStatus = "failed"
)

type IntegrationReference struct {
	ID                 uuid.UUID
	LocalEntityType    string
	LocalEntityID      uuid.UUID
	ExternalSystem     string
	ExternalEntityType string
	ExternalEntityID   string
	SyncStatus         SyncStatus
	LastSyncAt         *time.Time
	LastError          string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type NewIntegrationReferenceInput struct {
	ID                 uuid.UUID
	LocalEntityType    string
	LocalEntityID      uuid.UUID
	ExternalSystem     string
	ExternalEntityType string
	ExternalEntityID   string
	SyncStatus         SyncStatus
	LastError          string
	Now                time.Time
}

func NewIntegrationReference(input NewIntegrationReferenceInput) (IntegrationReference, error) {
	if input.ID == uuid.Nil {
		return IntegrationReference{}, ErrInvalidID
	}
	if input.LocalEntityID == uuid.Nil {
		return IntegrationReference{}, ErrInvalidID
	}
	if strings.TrimSpace(input.LocalEntityType) == "" {
		return IntegrationReference{}, ErrInvalidEntityType
	}
	if strings.TrimSpace(input.ExternalSystem) == "" {
		return IntegrationReference{}, ErrInvalidExternalSystem
	}
	if strings.TrimSpace(input.ExternalEntityType) == "" {
		return IntegrationReference{}, ErrInvalidExternalEntityType
	}
	if strings.TrimSpace(input.ExternalEntityID) == "" {
		return IntegrationReference{}, ErrInvalidExternalEntityID
	}
	if !isValidSyncStatus(input.SyncStatus) {
		return IntegrationReference{}, ErrInvalidSyncStatus
	}

	now := input.Now
	if now.IsZero() {
		now = time.Now().UTC()
	}

	return IntegrationReference{
		ID:                 input.ID,
		LocalEntityType:    strings.TrimSpace(input.LocalEntityType),
		LocalEntityID:      input.LocalEntityID,
		ExternalSystem:     strings.TrimSpace(input.ExternalSystem),
		ExternalEntityType: strings.TrimSpace(input.ExternalEntityType),
		ExternalEntityID:   strings.TrimSpace(input.ExternalEntityID),
		SyncStatus:         input.SyncStatus,
		LastSyncAt:         &now,
		LastError:          strings.TrimSpace(input.LastError),
		CreatedAt:          now,
		UpdatedAt:          now,
	}, nil
}

func (ref IntegrationReference) WithSyncUpdate(
	status SyncStatus,
	externalEntityID string,
	lastError string,
	now time.Time,
) (IntegrationReference, error) {
	if !isValidSyncStatus(status) {
		return IntegrationReference{}, ErrInvalidSyncStatus
	}
	if strings.TrimSpace(externalEntityID) == "" {
		return IntegrationReference{}, ErrInvalidExternalEntityID
	}
	if now.IsZero() {
		now = time.Now().UTC()
	}

	updated := ref
	updated.SyncStatus = status
	updated.ExternalEntityID = strings.TrimSpace(externalEntityID)
	updated.LastError = strings.TrimSpace(lastError)
	updated.LastSyncAt = &now
	updated.UpdatedAt = now
	return updated, nil
}

func isValidSyncStatus(status SyncStatus) bool {
	switch status {
	case SyncStatusNotConfigured, SyncStatusSynced, SyncStatusPending, SyncStatusFailed:
		return true
	default:
		return false
	}
}

func PlaceholderExternalID(localEntityID uuid.UUID) string {
	return "local:" + localEntityID.String()
}
