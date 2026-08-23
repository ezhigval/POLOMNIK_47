package noop

import (
	"context"

	"palomnik/internal/domain"
	"palomnik/internal/ports"
)

type Gateway struct{}

func New() Gateway { return Gateway{} }

var _ ports.PaymentPort = Gateway{}

func (Gateway) Configured() bool { return false }

func (Gateway) CreatePayment(context.Context, domain.Booking) (ports.PaymentSession, error) {
	return ports.PaymentSession{}, ports.ErrPaymentNotConfigured
}
