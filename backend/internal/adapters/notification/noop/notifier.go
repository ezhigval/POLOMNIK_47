package noop

import (
	"context"

	"polomnik/internal/domain"
	"polomnik/internal/ports"
)

type Notifier struct{}

func New() Notifier {
	return Notifier{}
}

var _ ports.NotificationPort = Notifier{}

func (Notifier) NotifyNewBooking(context.Context, domain.Booking, domain.Tour) error {
	return ports.ErrNotificationNotConfigured
}

func (Notifier) NotifyBookingStatusChanged(context.Context, domain.Booking, domain.Tour, domain.BookingStatus) error {
	return ports.ErrNotificationNotConfigured
}
