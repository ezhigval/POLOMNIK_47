package dto

import (
	"time"

	"palomnik/internal/domain"
)

type OAuthLoginRequest struct {
	Provider            string `json:"provider"`
	Subject             string `json:"subject"`
	Email               string `json:"email"`
	Name                string `json:"name"`
	Phone               string `json:"phone"`
	ConsentPersonalData bool   `json:"consent_personal_data"`
	ConsentTerms        bool   `json:"consent_terms"`
}

type SendSupportMessageRequest struct {
	Body                string `json:"body"`
	Website             string `json:"website"`
	ConsentPersonalData bool   `json:"consent_personal_data"`
}

type SupportMessageResponse struct {
	ID         string    `json:"id"`
	SenderType string    `json:"sender_type"`
	Body       string    `json:"body"`
	CreatedAt  time.Time `json:"created_at"`
}

type SupportThreadResponse struct {
	ID        string                   `json:"id"`
	UserID    string                   `json:"user_id,omitempty"`
	Subject   string                   `json:"subject"`
	Status    string                   `json:"status"`
	Messages  []SupportMessageResponse `json:"messages,omitempty"`
	CreatedAt time.Time                `json:"created_at,omitempty"`
	UpdatedAt time.Time                `json:"updated_at"`
}

func ToSupportThreadResponse(thread domain.SupportThread, messages []domain.SupportMessage) SupportThreadResponse {
	items := make([]SupportMessageResponse, 0, len(messages))
	for _, message := range messages {
		items = append(items, SupportMessageResponse{
			ID:         message.ID.String(),
			SenderType: string(message.SenderType),
			Body:       message.Body,
			CreatedAt:  message.CreatedAt,
		})
	}
	return SupportThreadResponse{
		ID:        thread.ID.String(),
		Subject:   thread.Subject,
		Status:    thread.Status,
		Messages:  items,
		UpdatedAt: thread.UpdatedAt,
	}
}

func ToManagementSupportThreadSummary(thread domain.SupportThread) SupportThreadResponse {
	return SupportThreadResponse{
		ID:        thread.ID.String(),
		UserID:    thread.UserID.String(),
		Subject:   thread.Subject,
		Status:    thread.Status,
		CreatedAt: thread.CreatedAt,
		UpdatedAt: thread.UpdatedAt,
	}
}

func ToManagementSupportThreadResponse(thread domain.SupportThread, messages []domain.SupportMessage) SupportThreadResponse {
	resp := ToSupportThreadResponse(thread, messages)
	resp.UserID = thread.UserID.String()
	resp.CreatedAt = thread.CreatedAt
	return resp
}
