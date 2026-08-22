package telegram

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"polomnik/internal/config"
	"polomnik/internal/domain"
	"polomnik/internal/ports"
)

type Notifier struct {
	client     Client
	recipients ports.TelegramRecipientsRepository
	chats      ports.TelegramChatMapRepository
	fallback   string
	mgmtBase   string
}

func New(cfg config.Config, recipients ports.TelegramRecipientsRepository, chats ports.TelegramChatMapRepository) Notifier {
	return Notifier{
		client:     NewClient(cfg),
		recipients: recipients,
		chats:      chats,
		fallback:   strings.TrimSpace(cfg.TelegramChatID),
		mgmtBase:   strings.TrimRight(cfg.ManagementBaseURL, "/"),
	}
}

var _ ports.NotificationPort = Notifier{}

func (n Notifier) NotifyNewBooking(ctx context.Context, booking domain.Booking, tour domain.Tour) error {
	return n.sendToKind(ctx, domain.TelegramRecipientBooking, formatNewBookingMessage(booking, tour, n.mgmtBase))
}

func (n Notifier) NotifyBookingStatusChanged(
	ctx context.Context,
	booking domain.Booking,
	tour domain.Tour,
	previousStatus domain.BookingStatus,
) error {
	if booking.Status == previousStatus {
		return nil
	}
	return n.sendToKind(ctx, domain.TelegramRecipientBooking, formatStatusChangedMessage(booking, tour, previousStatus, n.mgmtBase))
}

func (n Notifier) NotifySupportMessage(ctx context.Context, note domain.SupportNotification) error {
	return n.sendToKind(ctx, domain.TelegramRecipientSupport, formatSupportMessage(note, n.mgmtBase))
}

func (n Notifier) sendToKind(ctx context.Context, kind domain.TelegramRecipientKind, text string) error {
	if !n.client.Configured() {
		return ports.ErrNotificationNotConfigured
	}

	usernames, err := n.usernamesFor(ctx, kind)
	if err != nil {
		return err
	}
	bindings, err := n.bindingsFor(ctx, usernames)
	if err != nil {
		return err
	}

	route := domain.ResolveTelegramTargets(usernames, bindings, n.fallback)
	if len(route.ChatIDs) == 0 {
		if len(usernames) == 0 && n.fallback == "" {
			return ports.ErrNotificationNotConfigured
		}
		return nil
	}

	var sendErr error
	sent := 0
	for _, chatID := range route.ChatIDs {
		if err := n.client.SendMessage(ctx, chatID, text); err != nil {
			sendErr = errors.Join(sendErr, err)
			continue
		}
		sent++
	}
	if sent == 0 {
		return sendErr
	}
	return sendErr
}

func (n Notifier) usernamesFor(ctx context.Context, kind domain.TelegramRecipientKind) ([]string, error) {
	if n.recipients == nil {
		return nil, nil
	}
	settings, err := n.recipients.GetTelegramRecipients(ctx)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, nil
		}
		return nil, err
	}
	switch kind {
	case domain.TelegramRecipientSupport:
		return append([]string(nil), settings.SupportUsernames...), nil
	default:
		return append([]string(nil), settings.BookingUsernames...), nil
	}
}

func (n Notifier) bindingsFor(ctx context.Context, usernames []string) (map[string]string, error) {
	if n.chats == nil || len(usernames) == 0 {
		return map[string]string{}, nil
	}
	return n.chats.ListTelegramChatBindings(ctx, usernames)
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

func formatSupportMessage(note domain.SupportNotification, mgmtBase string) string {
	var b strings.Builder
	b.WriteString("<b>Сообщение в поддержку</b>\n")
	if note.Body != "" {
		b.WriteString(escapeHTML(note.Body))
		b.WriteString("\n")
	}
	if note.ThreadID != uuid.Nil {
		b.WriteString(fmt.Sprintf("Диалог: <code>%s</code>\n", note.ThreadID.String()))
	}
	b.WriteString("Ответ — в кабинете клиента (отдельного экрана поддержки в админке нет).")
	if mgmtBase != "" {
		b.WriteString(fmt.Sprintf("\n<a href=\"%s\">Кабинет управления</a>", strings.TrimSuffix(mgmtBase, "/bookings")))
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
