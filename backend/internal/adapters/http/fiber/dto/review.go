package dto

import (
	"time"

	"polomnik/internal/domain"
)

type ReviewResponse struct {
	ID         string `json:"id"`
	TourID     string `json:"tour_id"`
	ClientName string `json:"client_name"`
	Rating     int    `json:"rating"`
	Text       string `json:"text"`
	CreatedAt  string `json:"created_at"`
}

type ManagementReviewResponse struct {
	ReviewResponse
	IsApproved bool `json:"is_approved"`
}

type CreateReviewRequest struct {
	TourID     string `json:"tour_id"`
	ClientName string `json:"client_name"`
	Rating     int    `json:"rating"`
	Text       string `json:"text"`
	IsApproved bool   `json:"is_approved"`
}

func ToReviewResponse(review domain.Review) ReviewResponse {
	return ReviewResponse{
		ID:         review.ID.String(),
		TourID:     review.TourID.String(),
		ClientName: review.ClientName,
		Rating:     review.Rating,
		Text:       review.Text,
		CreatedAt:  review.CreatedAt.UTC().Format(time.RFC3339),
	}
}

func ToManagementReviewResponse(review domain.Review) ManagementReviewResponse {
	return ManagementReviewResponse{
		ReviewResponse: ToReviewResponse(review),
		IsApproved:     review.IsApproved,
	}
}
