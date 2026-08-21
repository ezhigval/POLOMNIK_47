package dto

import (
	"time"

	"polomnik/internal/domain"
)

type OAuthLoginRequest struct {
	Provider string `json:"provider"`
	Subject  string `json:"subject"`
	Email    string `json:"email"`
	Name     string `json:"name"`
}

type SendSupportMessageRequest struct {
	Body string `json:"body"`
}

type SupportMessageResponse struct {
	ID         string    `json:"id"`
	SenderType string    `json:"sender_type"`
	Body       string    `json:"body"`
	CreatedAt  time.Time `json:"created_at"`
}

type SupportThreadResponse struct {
	ID        string                   `json:"id"`
	Subject   string                   `json:"subject"`
	Status    string                   `json:"status"`
	Messages  []SupportMessageResponse `json:"messages"`
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
