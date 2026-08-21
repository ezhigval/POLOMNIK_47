package noop

import (
	"context"

	"polomnik/internal/adapters/integration/stub"
	"polomnik/internal/domain"
	"polomnik/internal/ports"
)

type CRMAdapter struct{}

func NewCRMAdapter() CRMAdapter {
	return CRMAdapter{}
}

var _ ports.CRMPort = CRMAdapter{}

func (CRMAdapter) SyncTour(context.Context, domain.Tour) (ports.IntegrationResult, error) {
	return stub.NotConfigured("crm"), nil
}

func (CRMAdapter) PullTours(context.Context) ([]domain.Tour, error) {
	return nil, nil
}

func (CRMAdapter) PushBooking(context.Context, domain.Booking) (ports.IntegrationResult, error) {
	return stub.NotConfigured("crm"), nil
}

func (CRMAdapter) UpdateBookingStatus(context.Context, domain.Booking) (ports.IntegrationResult, error) {
	return stub.NotConfigured("crm"), nil
}

func (CRMAdapter) PushReview(context.Context, domain.Review) (ports.IntegrationResult, error) {
	return stub.NotConfigured("crm"), nil
}

func (CRMAdapter) PullReviews(context.Context) ([]domain.Review, error) {
	return nil, nil
}
