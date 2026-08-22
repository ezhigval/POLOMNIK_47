package domain

import (
	"strings"
	"time"
	"unicode"
)

const (
	notificationRoutingRowID   int16 = 1
	maxRecipientsPerEvent            = 20
	maxNotificationAddressLen        = 64
)

type NotificationEventKind string

const (
	NotificationEventBookingCreated       NotificationEventKind = "booking_created"
	NotificationEventBookingStatusChanged NotificationEventKind = "booking_status_changed"
	NotificationEventSupportMessage       NotificationEventKind = "support_message"
)

type NotificationChannel string

const (
	NotificationChannelTelegram NotificationChannel = "telegram"
	NotificationChannelMax      NotificationChannel = "max"
)

type NotificationRecipient struct {
	Channel NotificationChannel
	Address string
}

type NotificationRouting struct {
	ByEvent   map[NotificationEventKind][]NotificationRecipient
	UpdatedAt time.Time
}

func NotificationRoutingID() int16 {
	return notificationRoutingRowID
}

func AllNotificationEventKinds() []NotificationEventKind {
	return []NotificationEventKind{
		NotificationEventBookingCreated,
		NotificationEventBookingStatusChanged,
		NotificationEventSupportMessage,
	}
}

func ValidNotificationEventKind(kind NotificationEventKind) bool {
	switch kind {
	case NotificationEventBookingCreated, NotificationEventBookingStatusChanged, NotificationEventSupportMessage:
		return true
	default:
		return false
	}
}

func ValidNotificationChannel(channel NotificationChannel) bool {
	switch channel {
	case NotificationChannelTelegram, NotificationChannelMax:
		return true
	default:
		return false
	}
}

func EmptyNotificationRouting(now time.Time) NotificationRouting {
	if now.IsZero() {
		now = time.Now().UTC()
	}
	byEvent := make(map[NotificationEventKind][]NotificationRecipient, 3)
	for _, kind := range AllNotificationEventKinds() {
		byEvent[kind] = nil
	}
	return NotificationRouting{ByEvent: byEvent, UpdatedAt: now}
}

func NewNotificationRecipient(channelRaw, addressRaw string) (NotificationRecipient, error) {
	channel := NotificationChannel(strings.ToLower(strings.TrimSpace(channelRaw)))
	if !ValidNotificationChannel(channel) {
		return NotificationRecipient{}, ErrInvalidNotificationChannel
	}
	address, err := NormalizeNotificationAddress(channel, addressRaw)
	if err != nil {
		return NotificationRecipient{}, err
	}
	if address == "" {
		return NotificationRecipient{}, ErrInvalidNotificationAddress
	}
	return NotificationRecipient{Channel: channel, Address: address}, nil
}

func NormalizeNotificationAddress(channel NotificationChannel, raw string) (string, error) {
	switch channel {
	case NotificationChannelTelegram:
		return NormalizeTelegramUsername(raw)
	case NotificationChannelMax:
		return NormalizeMaxAddress(raw)
	default:
		return "", ErrInvalidNotificationChannel
	}
}

func NormalizeMaxAddress(raw string) (string, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return "", nil
	}
	phone := NormalizePhone(value)
	digits := 0
	for _, r := range phone {
		if r >= '0' && r <= '9' {
			digits++
		}
	}
	if !strings.HasPrefix(phone, "+") || digits < 10 || digits > 15 {
		return "", ErrInvalidNotificationAddress
	}
	if len(phone) > maxNotificationAddressLen {
		return "", ErrInvalidNotificationAddress
	}
	return phone, nil
}

func NewNotificationRouting(
	events map[NotificationEventKind][]NotificationRecipient,
	now time.Time,
) (NotificationRouting, error) {
	routing := EmptyNotificationRouting(now)
	for _, kind := range AllNotificationEventKinds() {
		list := events[kind]
		normalized, err := normalizeRecipientList(list)
		if err != nil {
			return NotificationRouting{}, err
		}
		routing.ByEvent[kind] = normalized
	}
	// Reject unknown event keys that callers may have passed.
	for kind := range events {
		if !ValidNotificationEventKind(kind) {
			return NotificationRouting{}, ErrInvalidNotificationEvent
		}
	}
	return routing, nil
}

func normalizeRecipientList(list []NotificationRecipient) ([]NotificationRecipient, error) {
	seen := make(map[string]struct{})
	out := make([]NotificationRecipient, 0, len(list))
	for _, item := range list {
		recipient, err := NewNotificationRecipient(string(item.Channel), item.Address)
		if err != nil {
			return nil, err
		}
		key := string(recipient.Channel) + ":" + recipient.Address
		if _, ok := seen[key]; ok {
			continue
		}
		if len(out) >= maxRecipientsPerEvent {
			return nil, ErrInvalidNotificationAddress
		}
		seen[key] = struct{}{}
		out = append(out, recipient)
	}
	return out, nil
}

func RecipientsForEvent(routing NotificationRouting, kind NotificationEventKind) []NotificationRecipient {
	if routing.ByEvent == nil {
		return nil
	}
	return append([]NotificationRecipient(nil), routing.ByEvent[kind]...)
}

// NotificationRoutingFromTelegramRecipients maps legacy username lists onto the new
// event→recipient model without dropping data: booking usernames feed both booking
// events; support usernames feed support messages.
func NotificationRoutingFromTelegramRecipients(legacy TelegramRecipients, now time.Time) NotificationRouting {
	routing := EmptyNotificationRouting(now)
	booking := telegramUsernamesToRecipients(legacy.BookingUsernames)
	support := telegramUsernamesToRecipients(legacy.SupportUsernames)
	routing.ByEvent[NotificationEventBookingCreated] = append([]NotificationRecipient(nil), booking...)
	routing.ByEvent[NotificationEventBookingStatusChanged] = append([]NotificationRecipient(nil), booking...)
	routing.ByEvent[NotificationEventSupportMessage] = support
	if !legacy.UpdatedAt.IsZero() {
		routing.UpdatedAt = legacy.UpdatedAt.UTC()
	}
	return routing
}

func telegramUsernamesToRecipients(usernames []string) []NotificationRecipient {
	out := make([]NotificationRecipient, 0, len(usernames))
	seen := make(map[string]struct{})
	for _, username := range usernames {
		normalized, err := NormalizeTelegramUsername(username)
		if err != nil || normalized == "" {
			continue
		}
		if _, ok := seen[normalized]; ok {
			continue
		}
		seen[normalized] = struct{}{}
		out = append(out, NotificationRecipient{
			Channel: NotificationChannelTelegram,
			Address: normalized,
		})
	}
	return out
}

func TelegramUsernamesFromRouting(routing NotificationRouting, kind NotificationEventKind) []string {
	var out []string
	seen := make(map[string]struct{})
	for _, recipient := range RecipientsForEvent(routing, kind) {
		if recipient.Channel != NotificationChannelTelegram {
			continue
		}
		if _, ok := seen[recipient.Address]; ok {
			continue
		}
		seen[recipient.Address] = struct{}{}
		out = append(out, recipient.Address)
	}
	return out
}

func UniqueTelegramAddressesFromRouting(routing NotificationRouting) []string {
	return UniqueTelegramUsernames(
		TelegramUsernamesFromRouting(routing, NotificationEventBookingCreated),
		TelegramUsernamesFromRouting(routing, NotificationEventBookingStatusChanged),
		TelegramUsernamesFromRouting(routing, NotificationEventSupportMessage),
	)
}

func ParseLooseAddressList(raw string) []string {
	fields := strings.FieldsFunc(raw, func(r rune) bool {
		return r == ',' || r == ';' || r == '\n' || r == '\r' || unicode.IsSpace(r)
	})
	out := make([]string, 0, len(fields))
	for _, field := range fields {
		value := strings.TrimSpace(field)
		if value == "" {
			continue
		}
		out = append(out, value)
	}
	return out
}
