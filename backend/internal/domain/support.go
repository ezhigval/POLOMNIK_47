package domain

import (
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

type SupportSenderType string

const (
	SupportSenderUser  SupportSenderType = "user"
	SupportSenderStaff SupportSenderType = "staff"
)

type SupportThread struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Subject   string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type SupportMessage struct {
	ID         uuid.UUID
	ThreadID   uuid.UUID
	SenderType SupportSenderType
	Body       string
	CreatedAt  time.Time
}

func NewSupportMessage(input SupportMessage) (SupportMessage, error) {
	body := strings.TrimSpace(input.Body)
	if input.ID == uuid.Nil || input.ThreadID == uuid.Nil {
		return SupportMessage{}, ErrInvalidID
	}
	if body == "" || utf8.RuneCountInString(body) > 4000 {
		return SupportMessage{}, ErrInvalidReviewText
	}
	if input.SenderType != SupportSenderUser && input.SenderType != SupportSenderStaff {
		return SupportMessage{}, ErrInvalidBookingStatus
	}

	now := input.CreatedAt
	if now.IsZero() {
		now = time.Now().UTC()
	}

	return SupportMessage{
		ID:         input.ID,
		ThreadID:   input.ThreadID,
		SenderType: input.SenderType,
		Body:       body,
		CreatedAt:  now,
	}, nil
}

type Favorite struct {
	UserID    uuid.UUID
	TourID    uuid.UUID
	CreatedAt time.Time
}
