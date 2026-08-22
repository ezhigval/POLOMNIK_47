package application

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"polomnik/internal/domain"
	"polomnik/internal/ports"
)

// NotificationDispatcher routes domain notification events to channel adapters.
// Adapters only send; routing and Telegram username→chat_id resolution live here.
type NotificationDispatcher struct {
	routing  ports.NotificationRoutingRepository
	legacy   ports.TelegramRecipientsRepository
	chats    ports.TelegramChatMapRepository
	telegram ports.TelegramBotSender
	max      ports.MaxMessageSender
	fallback string
	mgmtBase string
	settings *NotificationSettingsService
}

func NewNotificationDispatcher(
	routing ports.NotificationRoutingRepository,
	legacy ports.TelegramRecipientsRepository,
	chats ports.TelegramChatMapRepository,
	telegram ports.TelegramBotSender,
	max ports.MaxMessageSender,
	fallbackChatID string,
	mgmtBaseURL string,
) *NotificationDispatcher {
	return &NotificationDispatcher{
		routing:  routing,
		legacy:   legacy,
		chats:    chats,
		telegram: telegram,
		max:      max,
		fallback: strings.TrimSpace(fallbackChatID),
		mgmtBase: strings.TrimRight(strings.TrimSpace(mgmtBaseURL), "/"),
		settings: NewNotificationSettingsService(routing, legacy, chats, telegram != nil, max != nil && max.Configured()),
	}
}

var _ ports.NotificationPort = (*NotificationDispatcher)(nil)

func (d *NotificationDispatcher) NotifyNewBooking(ctx context.Context, booking domain.Booking, tour domain.Tour) error {
	return d.dispatch(ctx, domain.NotificationEventBookingCreated, formatNewBookingMessage(booking, tour, d.mgmtBase))
}

func (d *NotificationDispatcher) NotifyBookingStatusChanged(
	ctx context.Context,
	booking domain.Booking,
	tour domain.Tour,
	previousStatus domain.BookingStatus,
) error {
	if booking.Status == previousStatus {
		return nil
	}
	return d.dispatch(ctx, domain.NotificationEventBookingStatusChanged, formatStatusChangedMessage(booking, tour, previousStatus, d.mgmtBase))
}

func (d *NotificationDispatcher) NotifySupportMessage(ctx context.Context, note domain.SupportNotification) error {
	return d.dispatch(ctx, domain.NotificationEventSupportMessage, formatSupportMessage(note, d.mgmtBase))
}

func (d *NotificationDispatcher) dispatch(ctx context.Context, kind domain.NotificationEventKind, text string) error {
	routing, err := d.loadRouting(ctx)
	if err != nil {
		return err
	}
	recipients := domain.RecipientsForEvent(routing, kind)
	if len(recipients) == 0 {
		return d.sendTelegramFallback(ctx, text)
	}

	usernames := domain.TelegramUsernamesFromRouting(domain.NotificationRouting{
		ByEvent: map[domain.NotificationEventKind][]domain.NotificationRecipient{kind: recipients},
	}, kind)
	bindings := map[string]string{}
	if d.chats != nil && len(usernames) > 0 {
		bindings, err = d.chats.ListTelegramChatBindings(ctx, usernames)
		if err != nil {
			return err
		}
	}

	var sendErr error
	sent := 0
	for _, recipient := range recipients {
		switch recipient.Channel {
		case domain.NotificationChannelTelegram:
			chatID := strings.TrimSpace(bindings[recipient.Address])
			if chatID == "" {
				continue
			}
			if d.telegram == nil {
				continue
			}
			if err := d.telegram.SendMessage(ctx, chatID, text); err != nil {
				sendErr = errors.Join(sendErr, err)
				continue
			}
			sent++
		case domain.NotificationChannelMax:
			if d.max == nil || !d.max.Configured() {
				continue
			}
			if err := d.max.SendMessage(ctx, recipient.Address, text); err != nil {
				if errors.Is(err, ports.ErrNotificationNotConfigured) {
					continue
				}
				sendErr = errors.Join(sendErr, err)
				continue
			}
			sent++
		default:
			continue
		}
	}

	if sent == 0 {
		if sendErr != nil {
			return sendErr
		}
		// All recipients skipped (unbound / channel off) — not a hard failure.
		return nil
	}
	return sendErr
}

func (d *NotificationDispatcher) sendTelegramFallback(ctx context.Context, text string) error {
	if d.fallback == "" || d.telegram == nil {
		return ports.ErrNotificationNotConfigured
	}
	if err := d.telegram.SendMessage(ctx, d.fallback, text); err != nil {
		return err
	}
	return nil
}

func (d *NotificationDispatcher) loadRouting(ctx context.Context) (domain.NotificationRouting, error) {
	if d.settings != nil {
		return d.settings.loadRouting(ctx)
	}
	return domain.EmptyNotificationRouting(domain.NotificationRouting{}.UpdatedAt), nil
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
	b.WriteString(fmt.Sprintf("%s → %s\n\n", formatBookingStatus(previousStatus), formatBookingStatus(booking.Status)))
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
	b.WriteString(fmt.Sprintf("Сумма: %d ₽\n", booking.TotalPrice))
	b.WriteString(fmt.Sprintf("Статус: %s\n", formatBookingStatus(booking.Status)))
	if booking.Comment != "" {
		b.WriteString(fmt.Sprintf("Комментарий: %s\n", escapeHTML(booking.Comment)))
	}
	if booking.Overbooked {
		b.WriteString("⚠ Overbooking\n")
	}
	b.WriteString(fmt.Sprintf("ID: <code>%s</code>", booking.ID.String()))
}

func formatBookingStatus(status domain.BookingStatus) string {
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
