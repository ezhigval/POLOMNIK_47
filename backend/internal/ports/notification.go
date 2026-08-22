package ports

import (
	"context"
	"errors"

	"polomnik/internal/domain"
)

var ErrNotificationNotConfigured = errors.New("notification adapter not configured")

type TelegramBotSender interface {
	SendMessage(ctx context.Context, chatID string, text string) error
}

type NotificationPort interface {
	NotifyNewBooking(ctx context.Context, booking domain.Booking, tour domain.Tour) error
	NotifyBookingStatusChanged(ctx context.Context, booking domain.Booking, tour domain.Tour, previousStatus domain.BookingStatus) error
	NotifySupportMessage(ctx context.Context, note domain.SupportNotification) error
}
