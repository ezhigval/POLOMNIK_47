package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Review struct {
	ID         uuid.UUID
	TourID     uuid.UUID
	ClientName string
	Rating     int
	Text       string
	IsApproved bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type NewReviewInput struct {
	ID         uuid.UUID
	TourID     uuid.UUID
	ClientName string
	Rating     int
	Text       string
	IsApproved bool
	Now        time.Time
}

func NewReview(input NewReviewInput) (Review, error) {
	if input.ID == uuid.Nil || input.TourID == uuid.Nil {
		return Review{}, ErrInvalidID
	}
	if strings.TrimSpace(input.ClientName) == "" {
		return Review{}, ErrInvalidClientName
	}
	if input.Rating < 1 || input.Rating > 5 {
		return Review{}, ErrInvalidRating
	}
	if strings.TrimSpace(input.Text) == "" {
		return Review{}, ErrInvalidReviewText
	}

	now := input.Now
	if now.IsZero() {
		now = time.Now().UTC()
	}

	return Review{
		ID:         input.ID,
		TourID:     input.TourID,
		ClientName: strings.TrimSpace(input.ClientName),
		Rating:     input.Rating,
		Text:       strings.TrimSpace(input.Text),
		IsApproved: input.IsApproved,
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}

func (r *Review) Approve() {
	r.IsApproved = true
	r.UpdatedAt = time.Now().UTC()
}

func (r *Review) Reject() {
	r.IsApproved = false
	r.UpdatedAt = time.Now().UTC()
}
