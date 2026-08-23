package dto

import (
	"encoding/json"

	"palomnik/internal/application"
	"palomnik/internal/domain"
)

type TelegramSettingsResponse struct {
	BookingUsernames string                    `json:"booking_usernames"`
	SupportUsernames string                    `json:"support_usernames"`
	Recipients       []TelegramRecipientStatus `json:"recipients"`
}

type TelegramRecipientStatus struct {
	Username  string `json:"username"`
	Kind      string `json:"kind"`
	ChatBound bool   `json:"chat_bound"`
	Status    string `json:"status"`
}

type TelegramSettingsUpsertRequest struct {
	BookingUsernames string `json:"booking_usernames"`
	SupportUsernames string `json:"support_usernames"`
}

func ToTelegramSettings(view application.TelegramSettingsView) TelegramSettingsResponse {
	recipients := make([]TelegramRecipientStatus, 0, len(view.Statuses))
	for _, item := range view.Statuses {
		status := "ждёт /start"
		if item.ChatBound {
			status = "подключён"
		}
		recipients = append(recipients, TelegramRecipientStatus{
			Username:  item.Username,
			Kind:      string(item.Kind),
			ChatBound: item.ChatBound,
			Status:    status,
		})
	}
	return TelegramSettingsResponse{
		BookingUsernames: domain.FormatTelegramUsernameList(view.Recipients.BookingUsernames),
		SupportUsernames: domain.FormatTelegramUsernameList(view.Recipients.SupportUsernames),
		Recipients:       recipients,
	}
}

type TelegramUpdateRequest struct {
	Message *TelegramMessage `json:"message"`
}

type TelegramMessage struct {
	From *TelegramUser `json:"from"`
	Chat *TelegramChat `json:"chat"`
	Text string        `json:"text"`
}

type TelegramUser struct {
	Username string `json:"username"`
}

type TelegramChat struct {
	ID json.Number `json:"id"`
}

func (m *TelegramMessage) ChatID() string {
	if m == nil || m.Chat == nil {
		return ""
	}
	return m.Chat.ID.String()
}

func (m *TelegramMessage) Username() string {
	if m == nil || m.From == nil {
		return ""
	}
	return m.From.Username
}
