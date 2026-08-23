package application

import (
	"context"

	"palomnik/internal/domain"
	"palomnik/internal/ports"
)

func resolveIntegrationStatus(result ports.IntegrationResult, err error) ports.IntegrationStatus {
	if err != nil {
		return ports.IntegrationStatusFailed
	}
	return result.Status
}

func shouldExportBookingToAccounting(status domain.BookingStatus) bool {
	return status == domain.BookingStatusConfirmed || status == domain.BookingStatusCompleted
}

func runAccountingSync(accounting ports.AccountingPort, ctx context.Context, booking domain.Booking) {
	_, _ = accounting.SyncCounterparty(ctx, booking)
	if shouldExportBookingToAccounting(booking.Status) {
		_, _ = accounting.ExportBooking(ctx, booking)
	}
}
