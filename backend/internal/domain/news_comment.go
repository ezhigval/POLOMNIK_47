package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type NewsComment struct {
	ID        uuid.UUID
	NewsID    uuid.UUID
	UserID    uuid.UUID
	Author    string
	Body      string
	CreatedAt time.Time
}

type NewNewsCommentInput struct {
	ID     uuid.UUID
	NewsID uuid.UUID
	UserID uuid.UUID
	Body   string
	Now    time.Time
}

func NewNewsComment(input NewNewsCommentInput) (NewsComment, error) {
	if input.ID == uuid.Nil {
		return NewsComment{}, ErrInvalidID
	}
	if input.NewsID == uuid.Nil || input.UserID == uuid.Nil {
		return NewsComment{}, ErrInvalidID
	}

	body, err := ValidateNewsCommentBody(input.Body)
	if err != nil {
		return NewsComment{}, err
	}

	now := input.Now
	if now.IsZero() {
		now = time.Now().UTC()
	}

	return NewsComment{
		ID:        input.ID,
		NewsID:    input.NewsID,
		UserID:    input.UserID,
		Body:      body,
		CreatedAt: now,
	}, nil
}

func DisplayNewsCommentAuthor(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return "Паломник"
	}
	return name
}
