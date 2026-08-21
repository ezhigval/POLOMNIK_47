package application

import (
	"context"
	"errors"
	"testing"

	"polomnik/internal/adapters/integration/noop"
	"polomnik/internal/adapters/repository/memory"
	"polomnik/internal/domain"
	"polomnik/internal/ports"
)

func TestReviewServiceListPublicReviewsReturnsOnlyApproved(t *testing.T) {
	ctx := context.Background()
	store := memory.NewStore()
	service := NewReviewService(store, store, noop.NewCRMAdapter())

	tour := testTour()
	if _, err := store.CreateTour(ctx, tour); err != nil {
		t.Fatalf("create tour: %v", err)
	}

	approved, err := domain.NewReview(domain.NewReviewInput{
		ID:         testUUID("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"),
		TourID:     tour.ID,
		ClientName: "Anna",
		Rating:     5,
		Text:       "Great tour",
		IsApproved: true,
	})
	if err != nil {
		t.Fatalf("create approved review domain: %v", err)
	}
	pending, err := domain.NewReview(domain.NewReviewInput{
		ID:         testUUID("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"),
		TourID:     tour.ID,
		ClientName: "Bob",
		Rating:     4,
		Text:       "Good tour",
		IsApproved: false,
	})
	if err != nil {
		t.Fatalf("create pending review domain: %v", err)
	}

	if _, err := store.CreateReview(ctx, approved); err != nil {
		t.Fatalf("store approved review: %v", err)
	}
	if _, err := store.CreateReview(ctx, pending); err != nil {
		t.Fatalf("store pending review: %v", err)
	}

	list, err := service.ListPublicReviews(ctx, ports.ReviewFilters{}, ports.Pagination{Page: 1, Limit: 20})
	if err != nil {
		t.Fatalf("list public reviews: %v", err)
	}
	if len(list.Items) != 1 {
		t.Fatalf("expected 1 approved review, got %d", len(list.Items))
	}
}

func TestReviewServiceCreateReviewRequiresExistingTour(t *testing.T) {
	ctx := context.Background()
	store := memory.NewStore()
	service := NewReviewService(store, store, noop.NewCRMAdapter())

	_, err := service.CreateReview(ctx, CreateReviewInput{
		TourID:     testUUID("99999999-9999-9999-9999-999999999999"),
		ClientName: "Anna",
		Rating:     5,
		Text:       "Great tour",
	})
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestReviewServiceApproveReview(t *testing.T) {
	ctx := context.Background()
	store := memory.NewStore()
	service := NewReviewService(store, store, noop.NewCRMAdapter())

	tour := testTour()
	if _, err := store.CreateTour(ctx, tour); err != nil {
		t.Fatalf("create tour: %v", err)
	}

	created, err := service.CreateReview(ctx, CreateReviewInput{
		TourID:     tour.ID,
		ClientName: "Anna",
		Rating:     5,
		Text:       "Great tour",
	})
	if err != nil {
		t.Fatalf("create review: %v", err)
	}

	approved, err := service.ApproveReview(ctx, created.ID)
	if err != nil {
		t.Fatalf("approve review: %v", err)
	}
	if !approved.IsApproved {
		t.Fatal("expected review to be approved")
	}
}
