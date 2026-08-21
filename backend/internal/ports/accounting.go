package ports

import (
	"context"

	"polomnik/internal/domain"
)

type AccountingPort interface {
	ExportBooking(ctx context.Context, booking domain.Booking) (IntegrationResult, error)
	SyncCounterparty(ctx context.Context, booking domain.Booking) (IntegrationResult, error)
}
