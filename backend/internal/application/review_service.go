package application

import (
	"context"

	"github.com/google/uuid"

	"polomnik/internal/domain"
	"polomnik/internal/ports"
)

type ReviewService struct {
	reviews ports.ReviewRepository
	tours   ports.TourRepository
	crm     ports.CRMPort
}

func NewReviewService(reviews ports.ReviewRepository, tours ports.TourRepository, crm ports.CRMPort) *ReviewService {
	return &ReviewService{reviews: reviews, tours: tours, crm: crm}
}

type CreateReviewInput struct {
	TourID     uuid.UUID
	ClientName string
	Rating     int
	Text       string
	IsApproved bool
}

func (s *ReviewService) ListPublicReviews(ctx context.Context, filters ports.ReviewFilters, pagination ports.Pagination) (ports.ReviewList, error) {
	approved := true
	filters.IsApproved = &approved
	return s.reviews.ListReviews(ctx, filters, pagination)
}

func (s *ReviewService) ListReviews(ctx context.Context, filters ports.ReviewFilters, pagination ports.Pagination) (ports.ReviewList, error) {
	return s.reviews.ListReviews(ctx, filters, pagination)
}

func (s *ReviewService) CreateReview(ctx context.Context, input CreateReviewInput) (domain.Review, error) {
	if _, err := s.tours.GetTour(ctx, input.TourID); err != nil {
		return domain.Review{}, err
	}

	review, err := domain.NewReview(domain.NewReviewInput{
		ID:         uuid.New(),
		TourID:     input.TourID,
		ClientName: input.ClientName,
		Rating:     input.Rating,
		Text:       input.Text,
		IsApproved: input.IsApproved,
	})
	if err != nil {
		return domain.Review{}, err
	}
	return s.reviews.CreateReview(ctx, review)
}

func (s *ReviewService) ApproveReview(ctx context.Context, id uuid.UUID) (domain.Review, error) {
	review, err := s.reviews.ApproveReview(ctx, id)
	if err != nil {
		return domain.Review{}, err
	}
	_, _ = s.crm.PushReview(ctx, review)
	return review, nil
}

func (s *ReviewService) RejectReview(ctx context.Context, id uuid.UUID) (domain.Review, error) {
	return s.reviews.RejectReview(ctx, id)
}

func (s *ReviewService) DeleteReview(ctx context.Context, id uuid.UUID) error {
	return s.reviews.DeleteReview(ctx, id)
}
