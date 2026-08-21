package noop

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"polomnik/internal/domain"
)

type CRMInboundAdapter struct{}

func NewCRMInboundAdapter() CRMInboundAdapter {
	return CRMInboundAdapter{}
}

func (CRMInboundAdapter) ResolveBookingStatusFromDeal(context.Context, string) (uuid.UUID, domain.BookingStatus, error) {
	return uuid.Nil, "", fmt.Errorf("crm inbound is not configured")
}
