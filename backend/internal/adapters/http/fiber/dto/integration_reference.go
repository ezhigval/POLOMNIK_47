package dto

import (
	"time"

	"palomnik/internal/domain"
)

type IntegrationReferenceResponse struct {
	ID                 string  `json:"id"`
	LocalEntityType    string  `json:"local_entity_type"`
	LocalEntityID      string  `json:"local_entity_id"`
	ExternalSystem     string  `json:"external_system"`
	ExternalEntityType string  `json:"external_entity_type"`
	ExternalEntityID   string  `json:"external_entity_id"`
	SyncStatus         string  `json:"sync_status"`
	LastSyncAt         *string `json:"last_sync_at"`
	LastError          string  `json:"last_error"`
	CreatedAt          string  `json:"created_at"`
	UpdatedAt          string  `json:"updated_at"`
}

func ToIntegrationReferenceResponse(ref domain.IntegrationReference) IntegrationReferenceResponse {
	response := IntegrationReferenceResponse{
		ID:                 ref.ID.String(),
		LocalEntityType:    ref.LocalEntityType,
		LocalEntityID:      ref.LocalEntityID.String(),
		ExternalSystem:     ref.ExternalSystem,
		ExternalEntityType: ref.ExternalEntityType,
		ExternalEntityID:   ref.ExternalEntityID,
		SyncStatus:         string(ref.SyncStatus),
		LastError:          ref.LastError,
		CreatedAt:          ref.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:          ref.UpdatedAt.UTC().Format(time.RFC3339),
	}
	if ref.LastSyncAt != nil {
		formatted := ref.LastSyncAt.UTC().Format(time.RFC3339)
		response.LastSyncAt = &formatted
	}
	return response
}
