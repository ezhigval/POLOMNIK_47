package dto

import (
	"encoding/json"
	"time"

	"polomnik/internal/domain"
)

type OutboxEventResponse struct {
	ID         string          `json:"id"`
	EventType  string          `json:"event_type"`
	EntityType string          `json:"entity_type"`
	EntityID   string          `json:"entity_id"`
	Payload    json.RawMessage `json:"payload"`
	Status     string          `json:"status"`
	Attempts   int             `json:"attempts"`
	LastError  string          `json:"last_error"`
	CreatedAt  string          `json:"created_at"`
	UpdatedAt  string          `json:"updated_at"`
}

func ToOutboxEventResponse(event domain.OutboxEvent) OutboxEventResponse {
	payload := event.Payload
	if len(payload) == 0 {
		payload = json.RawMessage(`{}`)
	}
	return OutboxEventResponse{
		ID:         event.ID.String(),
		EventType:  event.EventType,
		EntityType: event.EntityType,
		EntityID:   event.EntityID.String(),
		Payload:    payload,
		Status:     string(event.Status),
		Attempts:   event.Attempts,
		LastError:  event.LastError,
		CreatedAt:  event.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:  event.UpdatedAt.UTC().Format(time.RFC3339),
	}
}
