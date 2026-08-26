package ports

import (
	"context"
	"errors"

	"palomnik/internal/domain"
)

var ErrNotificationNotConfigured = errors.New("notification adapter not configured")

type TelegramBotSender interface {
	SendMessage(ctx context.Context, chatID string, text string) error
}

// MaxMessageSender delivers a prepared message to a Max address (phone).
// Implementations must be noop until credentials are configured.
type MaxMessageSender interface {
	Configured() bool
	SendMessage(ctx context.Context, address string, text string) error
}

type NotificationPort interface {
	NotifyNewBooking(ctx context.Context, booking domain.Booking, tour domain.Tour) error
	NotifyBookingStatusChanged(ctx context.Context, booking domain.Booking, tour domain.Tour, previousStatus domain.BookingStatus) error
	NotifySupportMessage(ctx context.Context, note domain.SupportNotification) error
	NotifyTourHidden(ctx context.Context, tour domain.Tour) error
}
