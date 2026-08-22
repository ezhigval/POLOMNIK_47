package telegram

import (
	"fmt"
	"strings"

	"github.com/google/uuid"

	"polomnik/internal/domain"
)

func FormatNewBookingMessage(booking domain.Booking, tour domain.Tour, mgmtBase string) string {
	var b strings.Builder
	b.WriteString("<b>Новая заявка</b>\n")
	appendBookingDetails(&b, booking, tour)
	if mgmtBase != "" {
		b.WriteString(fmt.Sprintf("\n<a href=\"%s\">Открыть заявки</a>", mgmtBase))
	}
	return b.String()
}

func FormatStatusChangedMessage(
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

func FormatSupportMessage(note domain.SupportNotification, mgmtBase string) string {
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
