package application

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/uuid"

	"palomnik/internal/domain"
	"palomnik/internal/ports"
)

const botForbiddenText = "Недостаточно прав."

var threadUUIDPattern = regexp.MustCompile(`(?i)[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`)

func (s *TelegramService) handleStaffOrCommand(ctx context.Context, update TelegramInboundUpdate) (string, bool, error) {
	if s == nil {
		return "", false, nil
	}

	command := telegramCommand(update.Text)
	switch command {
	case "/reply":
		text, err := s.handleReplyCommand(ctx, update)
		return text, true, err
	case "/bookings":
		text, err := s.requirePermThen(ctx, update, domain.PermManageBookings, s.handleBookingsCommand)
		return text, true, err
	case "/booking":
		text, err := s.requirePermThen(ctx, update, domain.PermManageBookings, s.handleBookingCommand)
		return text, true, err
	case "/tours":
		text, err := s.requirePermThen(ctx, update, domain.PermManageTours, s.handleToursCommand)
		return text, true, err
	case "/tour":
		text, err := s.requirePermThen(ctx, update, domain.PermManageTours, s.handleTourCommand)
		return text, true, err
	}

	if threadID, ok := extractThreadID(update.ReplyToText); ok {
		text, err := s.handleThreadReply(ctx, update, threadID, strings.TrimSpace(update.Text))
		return text, true, err
	}
	return "", false, nil
}

func (s *TelegramService) requirePermThen(
	ctx context.Context,
	update TelegramInboundUpdate,
	perm domain.Permission,
	fn func(context.Context, TelegramInboundUpdate) (string, error),
) (string, error) {
	if !s.hasPermission(ctx, update, perm) {
		return botForbiddenText, nil
	}
	return fn(ctx, update)
}

func (s *TelegramService) hasPermission(ctx context.Context, update TelegramInboundUpdate, perm domain.Permission) bool {
	principal, err := s.principalForTelegram(ctx, update.UserID)
	if err != nil {
		return false
	}
	return s.adminRoles != nil && s.adminRoles.HasPermission(principal, perm)
}

func (s *TelegramService) principalForTelegram(ctx context.Context, telegramUserID string) (ManagementPrincipal, error) {
	if s.adminRoles == nil || s.users == nil {
		return ManagementPrincipal{}, domain.ErrForbidden
	}
	identity, err := s.users.GetIdentity(ctx, "telegram", strings.TrimSpace(telegramUserID))
	if err != nil {
		return ManagementPrincipal{}, domain.ErrForbidden
	}
	return s.adminRoles.PrincipalForUser(ctx, identity.UserID)
}

func (s *TelegramService) isSupportRecipient(ctx context.Context, username string) bool {
	if s.settings == nil {
		return false
	}
	normalized, err := domain.NormalizeTelegramUsername(username)
	if err != nil || normalized == "" {
		return false
	}
	view, err := s.settings.Settings(ctx)
	if err != nil {
		return false
	}
	for _, event := range view.Events {
		if event.Kind != domain.NotificationEventSupportMessage {
			continue
		}
		for _, item := range event.Recipients {
			if item.Channel.ID != domain.NotificationChannelTelegram {
				continue
			}
			addr, addrErr := domain.NormalizeTelegramUsername(item.Address)
			if addrErr == nil && addr == normalized {
				return true
			}
		}
	}
	return false
}

func (s *TelegramService) canReplyToThread(ctx context.Context, update TelegramInboundUpdate) bool {
	if s.hasPermission(ctx, update, domain.PermManageSupport) {
		return true
	}
	return s.isSupportRecipient(ctx, update.Username)
}

func (s *TelegramService) handleReplyCommand(ctx context.Context, update TelegramInboundUpdate) (string, error) {
	fields := strings.Fields(strings.TrimSpace(update.Text))
	if len(fields) < 3 {
		return "Формат: /reply <id диалога> текст", nil
	}
	threadID, ok := parseUUID(fields[1])
	if !ok {
		return "Некорректный id диалога.", nil
	}
	body := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(update.Text), fields[0]))
	body = strings.TrimSpace(strings.TrimPrefix(body, fields[1]))
	return s.handleThreadReply(ctx, update, threadID, body)
}

func (s *TelegramService) handleThreadReply(ctx context.Context, update TelegramInboundUpdate, threadID uuid.UUID, body string) (string, error) {
	if s.support == nil {
		return "", nil
	}
	if !s.canReplyToThread(ctx, update) {
		return botForbiddenText, nil
	}
	body = strings.TrimSpace(body)
	if body == "" {
		return "Введите текст ответа.", nil
	}
	if _, _, err := s.support.SendStaffMessage(ctx, threadID, body); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return "Диалог не найден.", nil
		}
		return "", err
	}
	return "Ответ записан в диалог.", nil
}

func (s *TelegramService) handleBookingsCommand(ctx context.Context, update TelegramInboundUpdate) (string, error) {
	if s.bookings == nil {
		return "", nil
	}
	list, err := s.bookings.ListBookings(ctx, ports.BookingFilters{}, ports.NormalizePagination(1, 10))
	if err != nil {
		return "", err
	}
	if len(list.Items) == 0 {
		return "Заявок нет.", nil
	}
	var b strings.Builder
	b.WriteString("Заявки:\n")
	for _, item := range list.Items {
		fmt.Fprintf(&b, "%s  %s  %s  %s  %d ₽\n", item.ID, item.Status, item.Name, domain.MaskPhone(item.Phone), item.TotalPrice)
	}
	return strings.TrimSpace(b.String()), nil
}

func (s *TelegramService) handleBookingCommand(ctx context.Context, update TelegramInboundUpdate) (string, error) {
	if s.bookings == nil {
		return "", nil
	}
	fields := strings.Fields(strings.TrimSpace(update.Text))
	if len(fields) < 2 {
		return "Формат: /booking <id>", nil
	}
	id, ok := parseUUID(fields[1])
	if !ok {
		return "Некорректный id заявки.", nil
	}
	booking, err := s.bookings.GetBooking(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return "Заявка не найдена.", nil
		}
		return "", err
	}
	var b strings.Builder
	fmt.Fprintf(&b, "Заявка %s\n", booking.ID)
	fmt.Fprintf(&b, "Статус: %s\n", booking.Status)
	fmt.Fprintf(&b, "Клиент: %s\n", booking.Name)
	fmt.Fprintf(&b, "Телефон: %s\n", domain.MaskPhone(booking.Phone))
	if booking.Email != "" {
		fmt.Fprintf(&b, "Email: %s\n", booking.Email)
	}
	fmt.Fprintf(&b, "Участников: %d\n", booking.PeopleCount)
	fmt.Fprintf(&b, "Сумма: %d ₽\n", booking.TotalPrice)
	fmt.Fprintf(&b, "Тур: %s\n", booking.TourID)
	if booking.Comment != "" {
		fmt.Fprintf(&b, "Комментарий: %s\n", booking.Comment)
	}
	return strings.TrimSpace(b.String()), nil
}

func (s *TelegramService) handleToursCommand(ctx context.Context, update TelegramInboundUpdate) (string, error) {
	if s.tours == nil {
		return "", nil
	}
	list, err := s.tours.ListTours(ctx, ports.TourFilters{}, ports.NormalizePagination(1, 10))
	if err != nil {
		return "", err
	}
	if len(list.Items) == 0 {
		return "Туров нет.", nil
	}
	var b strings.Builder
	b.WriteString("Туры:\n")
	for _, tour := range list.Items {
		active := "выкл"
		if tour.IsActive {
			active = "вкл"
		}
		fmt.Fprintf(&b, "%s  %s  %d ₽  мест %d/%d  %s\n", tour.ID, tour.Slug, tour.Price, tour.SlotsLeft, tour.SlotsTotal, active)
	}
	return strings.TrimSpace(b.String()), nil
}

func (s *TelegramService) handleTourCommand(ctx context.Context, update TelegramInboundUpdate) (string, error) {
	if s.tours == nil {
		return "", nil
	}
	fields := strings.Fields(strings.TrimSpace(update.Text))
	if len(fields) < 2 {
		return "Формат: /tour <id|slug> [slots|price|on|off] …", nil
	}
	tour, err := s.tours.FindTour(ctx, fields[1])
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return "Тур не найден.", nil
		}
		return "", err
	}
	if len(fields) == 2 {
		return formatTourCard(tour), nil
	}

	var patch TourOpsPatch
	switch strings.ToLower(fields[2]) {
	case "on":
		active := true
		patch.IsActive = &active
	case "off":
		active := false
		patch.IsActive = &active
	case "price":
		if len(fields) < 4 {
			return "Формат: /tour <id|slug> price <число>", nil
		}
		price, convErr := strconv.Atoi(fields[3])
		if convErr != nil {
			return "Некорректная цена.", nil
		}
		patch.Price = &price
	case "slots":
		if len(fields) < 4 {
			return "Формат: /tour <id|slug> slots <свободно> [всего]", nil
		}
		left, convErr := strconv.Atoi(fields[3])
		if convErr != nil {
			return "Некорректное число мест.", nil
		}
		patch.SlotsLeft = &left
		if len(fields) >= 5 {
			total, totalErr := strconv.Atoi(fields[4])
			if totalErr != nil {
				return "Некорректное число мест.", nil
			}
			patch.SlotsTotal = &total
		}
	default:
		return "Формат: /tour <id|slug> [slots|price|on|off] …", nil
	}

	saved, err := s.tours.PatchTourOps(ctx, tour.ID, patch)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidSlots) || errors.Is(err, domain.ErrInvalidPrice) {
			return "Некорректные данные тура.", nil
		}
		return "", err
	}
	return "Сохранено.\n" + formatTourCard(saved), nil
}

func formatTourCard(tour domain.Tour) string {
	active := "выкл"
	if tour.IsActive {
		active = "вкл"
	}
	return fmt.Sprintf(
		"Тур %s\n%s\nЦена: %d ₽\nМест: %d/%d\n%s",
		tour.ID,
		tour.Slug,
		tour.Price,
		tour.SlotsLeft,
		tour.SlotsTotal,
		active,
	)
}

func (s *TelegramService) helpText(ctx context.Context, update TelegramInboundUpdate) string {
	var b strings.Builder
	b.WriteString("/start — запомнить чат для уведомлений\n")
	b.WriteString("/health — статус получателей\n")
	if s.canReplyToThread(ctx, update) {
		b.WriteString("/reply <id> текст — ответ в диалог поддержки\n")
		b.WriteString("или ответьте реплаем на уведомление с id диалога\n")
	}
	if s.hasPermission(ctx, update, domain.PermManageBookings) {
		b.WriteString("/bookings — список заявок\n")
		b.WriteString("/booking <id> — заявка\n")
	}
	if s.hasPermission(ctx, update, domain.PermManageTours) {
		b.WriteString("/tours — список туров\n")
		b.WriteString("/tour <id|slug> — тур\n")
		b.WriteString("/tour <id|slug> slots <свободно> [всего]\n")
		b.WriteString("/tour <id|slug> price <число>\n")
		b.WriteString("/tour <id|slug> on|off\n")
	}
	return strings.TrimSpace(b.String())
}

func extractThreadID(text string) (uuid.UUID, bool) {
	match := threadUUIDPattern.FindString(text)
	if match == "" {
		return uuid.Nil, false
	}
	id, err := uuid.Parse(match)
	if err != nil {
		return uuid.Nil, false
	}
	return id, true
}

func parseUUID(raw string) (uuid.UUID, bool) {
	id, err := uuid.Parse(strings.TrimSpace(raw))
	if err != nil || id == uuid.Nil {
		return uuid.Nil, false
	}
	return id, true
}
