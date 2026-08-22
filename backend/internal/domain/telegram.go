package domain

import (
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
)

const telegramRecipientsRowID int16 = 1

const maxTelegramRecipients = 20

type TelegramRecipientKind string

const (
	TelegramRecipientBooking TelegramRecipientKind = "booking"
	TelegramRecipientSupport TelegramRecipientKind = "support"
)

type TelegramRecipients struct {
	BookingUsernames []string
	SupportUsernames []string
	UpdatedAt        time.Time
}

type TelegramChatBinding struct {
	Username  string
	ChatID    string
	UpdatedAt time.Time
}

type SupportNotification struct {
	ThreadID  uuid.UUID
	MessageID uuid.UUID
	UserID    uuid.UUID
	Body      string
}

type TelegramRoute struct {
	ChatIDs []string
	Waiting []string
}

func TelegramRecipientsID() int16 {
	return telegramRecipientsRowID
}

func NewTelegramRecipients(bookingRaw, supportRaw string, now time.Time) (TelegramRecipients, error) {
	booking, err := ParseTelegramUsernameList(bookingRaw)
	if err != nil {
		return TelegramRecipients{}, err
	}
	support, err := ParseTelegramUsernameList(supportRaw)
	if err != nil {
		return TelegramRecipients{}, err
	}
	if now.IsZero() {
		now = time.Now().UTC()
	}
	return TelegramRecipients{
		BookingUsernames: booking,
		SupportUsernames: support,
		UpdatedAt:        now,
	}, nil
}

func ParseTelegramUsernameList(raw string) ([]string, error) {
	seen := make(map[string]struct{})
	var usernames []string

	fields := strings.FieldsFunc(raw, func(r rune) bool {
		return r == ',' || r == ';' || r == '\n' || r == '\r' || unicode.IsSpace(r)
	})
	for _, field := range fields {
		username, err := NormalizeTelegramUsername(field)
		if err != nil {
			return nil, err
		}
		if username == "" {
			continue
		}
		if _, ok := seen[username]; ok {
			continue
		}
		if len(usernames) >= maxTelegramRecipients {
			return nil, ErrInvalidTelegramUsername
		}
		seen[username] = struct{}{}
		usernames = append(usernames, username)
	}
	return usernames, nil
}

func NormalizeTelegramUsername(raw string) (string, error) {
	value := strings.TrimSpace(raw)
	value = strings.TrimPrefix(value, "@")
	value = strings.TrimSpace(value)
	if value == "" {
		return "", nil
	}
	value = strings.ToLower(value)
	if !validTelegramUsername(value) {
		return "", ErrInvalidTelegramUsername
	}
	return value, nil
}

func validTelegramUsername(value string) bool {
	if len(value) < 5 || len(value) > 32 {
		return false
	}
	if value[0] < 'a' || value[0] > 'z' {
		return false
	}
	for i := 1; i < len(value); i++ {
		c := value[i]
		if (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '_' {
			continue
		}
		return false
	}
	return true
}

func FormatTelegramUsernameList(usernames []string) string {
	return strings.Join(usernames, ", ")
}

func NewTelegramChatBinding(username, chatID string, now time.Time) (TelegramChatBinding, error) {
	normalized, err := NormalizeTelegramUsername(username)
	if err != nil {
		return TelegramChatBinding{}, err
	}
	if normalized == "" {
		return TelegramChatBinding{}, ErrInvalidTelegramUsername
	}
	chatID = strings.TrimSpace(chatID)
	if chatID == "" {
		return TelegramChatBinding{}, ErrInvalidTelegramChatID
	}
	if now.IsZero() {
		now = time.Now().UTC()
	}
	return TelegramChatBinding{
		Username:  normalized,
		ChatID:    chatID,
		UpdatedAt: now,
	}, nil
}

func ResolveTelegramTargets(usernames []string, bindings map[string]string, fallbackChatID string) TelegramRoute {
	if len(usernames) == 0 {
		if fallback := strings.TrimSpace(fallbackChatID); fallback != "" {
			return TelegramRoute{ChatIDs: []string{fallback}}
		}
		return TelegramRoute{}
	}

	seen := make(map[string]struct{})
	route := TelegramRoute{}
	for _, username := range usernames {
		key := strings.ToLower(strings.TrimSpace(username))
		chatID := strings.TrimSpace(bindings[key])
		if chatID == "" {
			route.Waiting = append(route.Waiting, key)
			continue
		}
		if _, ok := seen[chatID]; ok {
			continue
		}
		seen[chatID] = struct{}{}
		route.ChatIDs = append(route.ChatIDs, chatID)
	}
	return route
}

func UniqueTelegramUsernames(lists ...[]string) []string {
	seen := make(map[string]struct{})
	var out []string
	for _, list := range lists {
		for _, username := range list {
			key := strings.ToLower(strings.TrimSpace(username))
			if key == "" {
				continue
			}
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			out = append(out, key)
		}
	}
	return out
}
