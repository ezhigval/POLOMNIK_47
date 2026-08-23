package ports

import (
	"context"

	"github.com/google/uuid"

	"palomnik/internal/domain"
)

type IntegrationStatus string

const (
	IntegrationStatusNotConfigured IntegrationStatus = "not_configured"
	IntegrationStatusSynced        IntegrationStatus = "synced"
	IntegrationStatusPending       IntegrationStatus = "pending"
	IntegrationStatusFailed        IntegrationStatus = "failed"
)

type IntegrationResult struct {
	Status     IntegrationStatus
	ExternalID string
	Message    string
}

type CRMPort interface {
	SyncTour(ctx context.Context, tour domain.Tour) (IntegrationResult, error)
	PullTours(ctx context.Context) ([]domain.Tour, error)
	PushBooking(ctx context.Context, booking domain.Booking) (IntegrationResult, error)
	UpdateBookingStatus(ctx context.Context, booking domain.Booking) (IntegrationResult, error)
	PushReview(ctx context.Context, review domain.Review) (IntegrationResult, error)
	PullReviews(ctx context.Context) ([]domain.Review, error)
}

// CRMInboundPort resolves external CRM events back to local bookings.
type CRMInboundPort interface {
	ResolveBookingStatusFromDeal(ctx context.Context, dealID string) (uuid.UUID, domain.BookingStatus, error)
}
