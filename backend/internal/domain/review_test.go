package domain

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewReviewValidatesRating(t *testing.T) {
	_, err := NewReview(validReviewInput(func(input *NewReviewInput) {
		input.Rating = 6
	}))

	if !errors.Is(err, ErrInvalidRating) {
		t.Fatalf("expected invalid rating, got %v", err)
	}
}

func TestReviewApproveAndReject(t *testing.T) {
	review, err := NewReview(validReviewInput(func(input *NewReviewInput) {
		input.IsApproved = false
	}))
	if err != nil {
		t.Fatalf("create review: %v", err)
	}

	review.Approve()
	if !review.IsApproved {
		t.Fatal("expected review approved")
	}

	review.Reject()
	if review.IsApproved {
		t.Fatal("expected review rejected")
	}
}

func TestReviewSetCompanyReply(t *testing.T) {
	review, err := NewReview(validReviewInput())
	if err != nil {
		t.Fatalf("create review: %v", err)
	}

	review.SetCompanyReply("  Благодарим за отзыв.  ")
	if review.CompanyReply != "Благодарим за отзыв." {
		t.Fatalf("unexpected reply %q", review.CompanyReply)
	}
	if review.CompanyRepliedAt == nil {
		t.Fatal("expected replied_at to be set")
	}

	review.SetCompanyReply("   ")
	if review.CompanyReply != "" {
		t.Fatalf("expected empty reply, got %q", review.CompanyReply)
	}
	if review.CompanyRepliedAt != nil {
		t.Fatal("expected replied_at to be cleared")
	}
}

func validReviewInput(mutators ...func(*NewReviewInput)) NewReviewInput {
	input := NewReviewInput{
		ID:         uuid.New(),
		TourID:     uuid.New(),
		ClientName: "Мария",
		Rating:     5,
		Text:       "Очень хорошая поездка",
		Now:        time.Date(2026, 5, 20, 0, 0, 0, 0, time.UTC),
	}
	for _, mutate := range mutators {
		mutate(&input)
	}
	return input
}
