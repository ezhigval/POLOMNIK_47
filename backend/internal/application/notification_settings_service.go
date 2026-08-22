package application

import (
	"context"
	"errors"
	"strings"
	"time"

	"polomnik/internal/domain"
	"polomnik/internal/ports"
)

type NotificationRecipientStatus struct {
	Channel NotificationChannelView
	Address string
	Event   domain.NotificationEventKind
	Ready   bool
	Status  string
}

type NotificationChannelView struct {
	ID          domain.NotificationChannel
	Configured  bool
	Label       string
}

type NotificationEventBlock struct {
	Kind       domain.NotificationEventKind
	Title      string
	Recipients []NotificationRecipientStatus
}

type NotificationSettingsView struct {
	Routing  domain.NotificationRouting
	Channels []NotificationChannelView
	Events   []NotificationEventBlock
}

type NotificationSettingsService struct {
	routing    ports.NotificationRoutingRepository
	legacy     ports.TelegramRecipientsRepository
	chats      ports.TelegramChatMapRepository
	telegramOn bool
	maxOn      bool
}

func NewNotificationSettingsService(
	routing ports.NotificationRoutingRepository,
	legacy ports.TelegramRecipientsRepository,
	chats ports.TelegramChatMapRepository,
	telegramConfigured bool,
	maxConfigured bool,
) *NotificationSettingsService {
	return &NotificationSettingsService{
		routing:    routing,
		legacy:     legacy,
		chats:      chats,
		telegramOn: telegramConfigured,
		maxOn:      maxConfigured,
	}
}

func (s *NotificationSettingsService) Settings(ctx context.Context) (NotificationSettingsView, error) {
	stored, err := s.loadRouting(ctx)
	if err != nil {
		return NotificationSettingsView{}, err
	}

	usernames := domain.UniqueTelegramAddressesFromRouting(stored)
	bindings := map[string]string{}
	if s.chats != nil && len(usernames) > 0 {
		bindings, err = s.chats.ListTelegramChatBindings(ctx, usernames)
		if err != nil {
			return NotificationSettingsView{}, err
		}
	}

	channels := s.channelViews()
	events := make([]NotificationEventBlock, 0, len(domain.AllNotificationEventKinds()))
	for _, kind := range domain.AllNotificationEventKinds() {
		recipients := domain.RecipientsForEvent(stored, kind)
		statuses := make([]NotificationRecipientStatus, 0, len(recipients))
		for _, recipient := range recipients {
			status, ready := s.recipientStatus(recipient, bindings)
			statuses = append(statuses, NotificationRecipientStatus{
				Channel: NotificationChannelView{
					ID:         recipient.Channel,
					Configured: s.channelConfigured(recipient.Channel),
					Label:      channelLabel(recipient.Channel),
				},
				Address: recipient.Address,
				Event:   kind,
				Ready:   ready,
				Status:  status,
			})
		}
		events = append(events, NotificationEventBlock{
			Kind:       kind,
			Title:      eventTitle(kind),
			Recipients: statuses,
		})
	}

	return NotificationSettingsView{
		Routing:  stored,
		Channels: channels,
		Events:   events,
	}, nil
}

func (s *NotificationSettingsService) UpdateRouting(
	ctx context.Context,
	events map[domain.NotificationEventKind][]domain.NotificationRecipient,
) (NotificationSettingsView, error) {
	updated, err := domain.NewNotificationRouting(events, time.Time{})
	if err != nil {
		return NotificationSettingsView{}, err
	}
	if s == nil || s.routing == nil {
		return NotificationSettingsView{Routing: updated}, nil
	}
	if _, err := s.routing.UpsertNotificationRouting(ctx, updated); err != nil {
		return NotificationSettingsView{}, err
	}
	return s.Settings(ctx)
}

func (s *NotificationSettingsService) loadRouting(ctx context.Context) (domain.NotificationRouting, error) {
	if s == nil || s.routing == nil {
		return s.legacyRouting(ctx)
	}
	stored, err := s.routing.GetNotificationRouting(ctx)
	if err == nil {
		if routingHasAnyRecipient(stored) {
			return stored, nil
		}
		// Empty new table: try one-shot legacy fill so prod data is not lost before SQL migrate.
		if legacy, legacyErr := s.legacyRouting(ctx); legacyErr == nil && routingHasAnyRecipient(legacy) {
			if _, upsertErr := s.routing.UpsertNotificationRouting(ctx, legacy); upsertErr == nil {
				return legacy, nil
			}
			return legacy, nil
		}
		return stored, nil
	}
	if !errors.Is(err, domain.ErrNotFound) {
		return domain.NotificationRouting{}, err
	}
	return s.legacyRouting(ctx)
}

func (s *NotificationSettingsService) legacyRouting(ctx context.Context) (domain.NotificationRouting, error) {
	if s == nil || s.legacy == nil {
		return domain.EmptyNotificationRouting(time.Time{}), nil
	}
	legacy, err := s.legacy.GetTelegramRecipients(ctx)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.EmptyNotificationRouting(time.Time{}), nil
		}
		return domain.NotificationRouting{}, err
	}
	return domain.NotificationRoutingFromTelegramRecipients(legacy, legacy.UpdatedAt), nil
}

func (s *NotificationSettingsService) channelViews() []NotificationChannelView {
	return []NotificationChannelView{
		{ID: domain.NotificationChannelTelegram, Configured: s.telegramOn, Label: channelLabel(domain.NotificationChannelTelegram)},
		{ID: domain.NotificationChannelMax, Configured: s.maxOn, Label: channelLabel(domain.NotificationChannelMax)},
	}
}

func (s *NotificationSettingsService) channelConfigured(channel domain.NotificationChannel) bool {
	switch channel {
	case domain.NotificationChannelTelegram:
		return s.telegramOn
	case domain.NotificationChannelMax:
		return s.maxOn
	default:
		return false
	}
}

func (s *NotificationSettingsService) recipientStatus(
	recipient domain.NotificationRecipient,
	bindings map[string]string,
) (string, bool) {
	switch recipient.Channel {
	case domain.NotificationChannelTelegram:
		if !s.telegramOn {
			return "канал не подключён", false
		}
		if _, ok := bindings[recipient.Address]; ok {
			return "подключён", true
		}
		return "ждёт /start", false
	case domain.NotificationChannelMax:
		if !s.maxOn {
			return "канал не подключён", false
		}
		return "готов", true
	default:
		return "канал не подключён", false
	}
}

func routingHasAnyRecipient(routing domain.NotificationRouting) bool {
	for _, kind := range domain.AllNotificationEventKinds() {
		if len(domain.RecipientsForEvent(routing, kind)) > 0 {
			return true
		}
	}
	return false
}

func eventTitle(kind domain.NotificationEventKind) string {
	switch kind {
	case domain.NotificationEventBookingCreated:
		return "Новые заявки"
	case domain.NotificationEventBookingStatusChanged:
		return "Смена статуса заявки"
	case domain.NotificationEventSupportMessage:
		return "Сообщения в поддержку"
	default:
		return string(kind)
	}
}

func channelLabel(channel domain.NotificationChannel) string {
	switch channel {
	case domain.NotificationChannelTelegram:
		return "Telegram"
	case domain.NotificationChannelMax:
		return "Max"
	default:
		return strings.ToUpper(string(channel))
	}
}

// TelegramSettingsView is the legacy shape used by /management/telegram-settings.
type TelegramSettingsView struct {
	Recipients domain.TelegramRecipients
	Statuses   []TelegramRecipientStatusRow
}

type TelegramRecipientStatusRow struct {
	Username  string
	Kind      domain.TelegramRecipientKind
	ChatBound bool
}

func (s *NotificationSettingsService) TelegramSettings(ctx context.Context) (TelegramSettingsView, error) {
	view, err := s.Settings(ctx)
	if err != nil {
		return TelegramSettingsView{}, err
	}
	usernames := domain.UniqueTelegramAddressesFromRouting(view.Routing)
	bindings := map[string]string{}
	if s.chats != nil && len(usernames) > 0 {
		bindings, err = s.chats.ListTelegramChatBindings(ctx, usernames)
		if err != nil {
			return TelegramSettingsView{}, err
		}
	}
	return telegramViewFromNotification(view, bindings), nil
}

func (s *NotificationSettingsService) UpdateTelegramUsernames(
	ctx context.Context,
	bookingRaw, supportRaw string,
) (TelegramSettingsView, error) {
	legacy, err := domain.NewTelegramRecipients(bookingRaw, supportRaw, time.Time{})
	if err != nil {
		return TelegramSettingsView{}, err
	}
	routing := domain.NotificationRoutingFromTelegramRecipients(legacy, time.Time{})
	if _, err := s.UpdateRouting(ctx, routing.ByEvent); err != nil {
		return TelegramSettingsView{}, err
	}
	if s.legacy != nil {
		if _, err := s.legacy.UpsertTelegramRecipients(ctx, legacy); err != nil {
			return TelegramSettingsView{}, err
		}
	}
	return s.TelegramSettings(ctx)
}

func telegramViewFromNotification(view NotificationSettingsView, bindings map[string]string) TelegramSettingsView {
	booking := domain.TelegramUsernamesFromRouting(view.Routing, domain.NotificationEventBookingCreated)
	booking = domain.UniqueTelegramUsernames(
		booking,
		domain.TelegramUsernamesFromRouting(view.Routing, domain.NotificationEventBookingStatusChanged),
	)
	support := domain.TelegramUsernamesFromRouting(view.Routing, domain.NotificationEventSupportMessage)
	statuses := make([]TelegramRecipientStatusRow, 0)
	seen := map[string]struct{}{}
	for _, event := range view.Events {
		kind := domain.TelegramRecipientBooking
		if event.Kind == domain.NotificationEventSupportMessage {
			kind = domain.TelegramRecipientSupport
		}
		for _, item := range event.Recipients {
			if item.Channel.ID != domain.NotificationChannelTelegram {
				continue
			}
			key := string(kind) + ":" + item.Address
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			_, bound := bindings[item.Address]
			statuses = append(statuses, TelegramRecipientStatusRow{
				Username:  item.Address,
				Kind:      kind,
				ChatBound: bound,
			})
		}
	}
	return TelegramSettingsView{
		Recipients: domain.TelegramRecipients{
			BookingUsernames: booking,
			SupportUsernames: support,
			UpdatedAt:        view.Routing.UpdatedAt,
		},
		Statuses: statuses,
	}
}
