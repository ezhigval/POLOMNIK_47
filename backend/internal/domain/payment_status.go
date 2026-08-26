package domain

type PaymentStatus string

const (
	PaymentStatusUnpaid          PaymentStatus = "UNPAID"
	PaymentStatusAwaitingPayment PaymentStatus = "AWAITING_PAYMENT"
	PaymentStatusPaid            PaymentStatus = "PAID"
	PaymentStatusNotRequired     PaymentStatus = "NOT_REQUIRED"
)

func (s PaymentStatus) Valid() bool {
	switch s {
	case PaymentStatusUnpaid,
		PaymentStatusAwaitingPayment,
		PaymentStatusPaid,
		PaymentStatusNotRequired:
		return true
	default:
		return false
	}
}

func DefaultPaymentStatus(totalPrice int) PaymentStatus {
	if totalPrice <= 0 {
		return PaymentStatusNotRequired
	}
	return PaymentStatusUnpaid
}
