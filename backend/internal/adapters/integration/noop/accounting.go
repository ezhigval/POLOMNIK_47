package noop

import (
	"context"

	"polomnik/internal/adapters/integration/stub"
	"polomnik/internal/domain"
	"polomnik/internal/ports"
)

type AccountingAdapter struct{}

func NewAccountingAdapter() AccountingAdapter {
	return AccountingAdapter{}
}

var _ ports.AccountingPort = AccountingAdapter{}

func (AccountingAdapter) ExportBooking(context.Context, domain.Booking) (ports.IntegrationResult, error) {
	return stub.NotConfigured("accounting"), nil
}

func (AccountingAdapter) SyncCounterparty(context.Context, domain.Booking) (ports.IntegrationResult, error) {
	return stub.NotConfigured("accounting"), nil
}
