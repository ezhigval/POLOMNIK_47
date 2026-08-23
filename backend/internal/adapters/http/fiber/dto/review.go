package dto

import (
	"time"

	"palomnik/internal/domain"
)

type ReviewResponse struct {
	ID               string  `json:"id"`
	TourID           string  `json:"tour_id"`
	ClientName       string  `json:"client_name"`
	Rating           int     `json:"rating"`
	Text             string  `json:"text"`
	CompanyReply     string  `json:"company_reply"`
	CompanyRepliedAt *string `json:"company_replied_at,omitempty"`
	CreatedAt        string  `json:"created_at"`
}

type ManagementReviewResponse struct {
	ReviewResponse
	IsApproved        bool `json:"is_approved"`
	AllowDistribution bool `json:"allow_distribution"`
}

type CreateReviewRequest struct {
	TourID              string `json:"tour_id"`
	ClientName          string `json:"client_name"`
	Rating              int    `json:"rating"`
	Text                string `json:"text"`
	IsApproved          bool   `json:"is_approved"`
	AllowDistribution   bool   `json:"allow_distribution"`
	ConsentPersonalData bool   `json:"consent_personal_data"`
	Website             string `json:"website"`
}

type SetCompanyReplyRequest struct {
	CompanyReply string `json:"company_reply"`
}

func ToReviewResponse(review domain.Review) ReviewResponse {
	return ReviewResponse{
		ID:               review.ID.String(),
		TourID:           review.TourID.String(),
		ClientName:       review.ClientName,
		Rating:           review.Rating,
		Text:             review.Text,
		CompanyReply:     review.CompanyReply,
		CompanyRepliedAt: formatOptionalTime(review.CompanyRepliedAt),
		CreatedAt:        review.CreatedAt.UTC().Format(time.RFC3339),
	}
}

func formatOptionalTime(value *time.Time) *string {
	if value == nil {
		return nil
	}
	formatted := value.UTC().Format(time.RFC3339)
	return &formatted
}

func ToManagementReviewResponse(review domain.Review) ManagementReviewResponse {
	return ManagementReviewResponse{
		ReviewResponse:    ToReviewResponse(review),
		IsApproved:        review.IsApproved,
		AllowDistribution: review.AllowDistribution,
	}
}
