package telegram

import (
	"context"
	"fmt"
	"strings"

	"polomnik/internal/config"
	"polomnik/internal/domain"
	"polomnik/internal/ports"
)

type Notifier struct {
	client   Client
	chatID   string
	mgmtBase string
}

func New(cfg config.Config) Notifier {
	return Notifier{
		client:   NewClient(cfg),
		chatID:   cfg.TelegramChatID,
		mgmtBase: strings.TrimRight(cfg.ManagementBaseURL, "/"),
	}
}

var _ ports.NotificationPort = Notifier{}

func (n Notifier) NotifyNewBooking(ctx context.Context, booking domain.Booking, tour domain.Tour) error {
	if !n.configured() {
		return ports.ErrNotificationNotConfigured
	}
	return n.client.SendMessage(ctx, n.chatID, formatNewBookingMessage(booking, tour, n.mgmtBase))
}

func (n Notifier) NotifyBookingStatusChanged(
	ctx context.Context,
	booking domain.Booking,
	tour domain.Tour,
	previousStatus domain.BookingStatus,
) error {
	if !n.configured() {
		return ports.ErrNotificationNotConfigured
	}
	if booking.Status == previousStatus {
		return nil
	}
	return n.client.SendMessage(ctx, n.chatID, formatStatusChangedMessage(booking, tour, previousStatus, n.mgmtBase))
}

func (n Notifier) configured() bool {
	return n.client.Configured() && n.chatID != ""
}

func formatNewBookingMessage(booking domain.Booking, tour domain.Tour, mgmtBase string) string {
	var b strings.Builder
	b.WriteString("<b>Новая заявка</b>\n")
	appendBookingDetails(&b, booking, tour)
	if mgmtBase != "" {
		b.WriteString(fmt.Sprintf("\n<a href=\"%s\">Открыть заявки</a>", mgmtBase))
	}
	return b.String()
}

func formatStatusChangedMessage(
	booking domain.Booking,
	tour domain.Tour,
	previousStatus domain.BookingStatus,
	mgmtBase string,
) string {
	var b strings.Builder
	b.WriteString("<b>Статус заявки изменён</b>\n")
	b.WriteString(fmt.Sprintf("%s → %s\n\n", formatStatus(previousStatus), formatStatus(booking.Status)))
	appendBookingDetails(&b, booking, tour)
	if mgmtBase != "" {
		b.WriteString(fmt.Sprintf("\n<a href=\"%s\">Открыть заявки</a>", mgmtBase))
	}
	return b.String()
}

func appendBookingDetails(b *strings.Builder, booking domain.Booking, tour domain.Tour) {
	b.WriteString(fmt.Sprintf("Тур: %s\n", escapeHTML(tour.Title)))
	b.WriteString(fmt.Sprintf("Клиент: %s\n", escapeHTML(booking.Name)))
	b.WriteString(fmt.Sprintf("Телефон: %s\n", escapeHTML(booking.Phone)))
	if booking.Email != "" {
		b.WriteString(fmt.Sprintf("Email: %s\n", escapeHTML(booking.Email)))
	}
	b.WriteString(fmt.Sprintf("Участников: %d\n", booking.PeopleCount))
	b.WriteString(fmt.Sprintf("Сумма: %s ₽\n", formatAmount(booking.TotalPrice)))
	b.WriteString(fmt.Sprintf("Статус: %s\n", formatStatus(booking.Status)))
	if booking.Comment != "" {
		b.WriteString(fmt.Sprintf("Комментарий: %s\n", escapeHTML(booking.Comment)))
	}
	if booking.Overbooked {
		b.WriteString("⚠ Overbooking\n")
	}
	b.WriteString(fmt.Sprintf("ID: <code>%s</code>", booking.ID.String()))
}

func formatStatus(status domain.BookingStatus) string {
	switch status {
	case domain.BookingStatusNew:
		return "Новая"
	case domain.BookingStatusContacted:
		return "На связи"
	case domain.BookingStatusConfirmed:
		return "Подтверждена"
	case domain.BookingStatusCompleted:
		return "Завершена"
	case domain.BookingStatusCancelled:
		return "Отменена"
	default:
		return string(status)
	}
}

func escapeHTML(value string) string {
	replacer := strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;")
	return replacer.Replace(value)
}

func formatAmount(amount int) string {
	return fmt.Sprintf("%d", amount)
}
