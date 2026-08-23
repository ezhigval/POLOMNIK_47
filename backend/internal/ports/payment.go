package ports

import (
	"context"
	"errors"

	"palomnik/internal/domain"
)

var ErrPaymentNotConfigured = errors.New("payment adapter not configured")

type PaymentSession struct {
	ID          string
	RedirectURL string
}

// PaymentPort is the v3 checkout port. Stage 8 adapters (sber / yookassa)
// choose PAYMENT_ADAPTER. Amount is always booking.TotalPrice. No refunds.
type PaymentPort interface {
	Configured() bool
	CreatePayment(ctx context.Context, booking domain.Booking) (PaymentSession, error)
}
