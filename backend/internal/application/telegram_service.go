package application

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"

	"polomnik/internal/domain"
	"polomnik/internal/ports"
)

func TelegramWebhookSecret(internalSecret string) string {
	if strings.TrimSpace(internalSecret) == "" {
		return ""
	}
	sum := sha256.Sum256([]byte("telegram-webhook:" + internalSecret))
	return hex.EncodeToString(sum[:])
}

type TelegramService struct {
	recipients ports.TelegramRecipientsRepository
	chats      ports.TelegramChatMapRepository
	sender     ports.TelegramBotSender
	fallback   string
}

func NewTelegramService(
	recipients ports.TelegramRecipientsRepository,
	chats ports.TelegramChatMapRepository,
	sender ports.TelegramBotSender,
	fallbackChatID string,
) *TelegramService {
	return &TelegramService{
		recipients: recipients,
		chats:      chats,
		sender:     sender,
		fallback:   strings.TrimSpace(fallbackChatID),
	}
}

type TelegramRecipientStatus struct {
	Username  string
	Kind      domain.TelegramRecipientKind
	ChatBound bool
}

type TelegramSettingsView struct {
	Recipients domain.TelegramRecipients
	Statuses   []TelegramRecipientStatus
}

type TelegramInboundUpdate struct {
	Username string
	ChatID   string
	Text     string
}

func (s *TelegramService) Settings(ctx context.Context) (TelegramSettingsView, error) {
	stored, err := s.getStored(ctx)
	if err != nil {
		return TelegramSettingsView{}, err
	}
	usernames := domain.UniqueTelegramUsernames(stored.BookingUsernames, stored.SupportUsernames)
	bindings := map[string]string{}
	if s.chats != nil && len(usernames) > 0 {
		bindings, err = s.chats.ListTelegramChatBindings(ctx, usernames)
		if err != nil {
			return TelegramSettingsView{}, err
		}
	}

	statuses := make([]TelegramRecipientStatus, 0, len(stored.BookingUsernames)+len(stored.SupportUsernames))
	for _, username := range stored.BookingUsernames {
		_, bound := bindings[username]
		statuses = append(statuses, TelegramRecipientStatus{
			Username:  username,
			Kind:      domain.TelegramRecipientBooking,
			ChatBound: bound,
		})
	}
	for _, username := range stored.SupportUsernames {
		_, bound := bindings[username]
		statuses = append(statuses, TelegramRecipientStatus{
			Username:  username,
			Kind:      domain.TelegramRecipientSupport,
			ChatBound: bound,
		})
	}
	return TelegramSettingsView{Recipients: stored, Statuses: statuses}, nil
}

func (s *TelegramService) UpdateRecipients(ctx context.Context, bookingRaw, supportRaw string) (TelegramSettingsView, error) {
	updated, err := domain.NewTelegramRecipients(bookingRaw, supportRaw, domain.TelegramRecipients{}.UpdatedAt)
	if err != nil {
		return TelegramSettingsView{}, err
	}
	if s == nil || s.recipients == nil {
		return TelegramSettingsView{Recipients: updated}, nil
	}
	if _, err := s.recipients.UpsertTelegramRecipients(ctx, updated); err != nil {
		return TelegramSettingsView{}, err
	}
	return s.Settings(ctx)
}

func (s *TelegramService) HandleInboundUpdate(ctx context.Context, update TelegramInboundUpdate) error {
	username := strings.TrimSpace(update.Username)
	chatID := strings.TrimSpace(update.ChatID)
	if username != "" && chatID != "" {
		binding, err := domain.NewTelegramChatBinding(username, chatID, domain.TelegramChatBinding{}.UpdatedAt)
		if err == nil && s.chats != nil {
			if err := s.chats.UpsertTelegramChatBinding(ctx, binding); err != nil {
				return err
			}
		}
	}

	if s.sender == nil || chatID == "" {
		return nil
	}

	switch telegramCommand(update.Text) {
	case "/start":
		return s.sender.SendMessage(ctx, chatID, "Готово. Бот запомнил ваш аккаунт и сможет присылать уведомления, если вас добавят в настройках.")
	case "/health":
		return s.sender.SendMessage(ctx, chatID, s.healthText(ctx))
	default:
		return nil
	}
}

func (s *TelegramService) healthText(ctx context.Context) string {
	view, err := s.Settings(ctx)
	if err != nil {
		return "Бот на связи. Список получателей сейчас недоступен."
	}

	var b strings.Builder
	b.WriteString("Бот на связи.\n")
	if len(view.Statuses) == 0 {
		if s.fallback != "" {
			b.WriteString("Списки в настройках пустые — используется запасной chat_id.")
		} else {
			b.WriteString("Списки получателей в настройках пустые.")
		}
		return b.String()
	}

	waiting := 0
	for _, item := range view.Statuses {
		if !item.ChatBound {
			waiting++
			b.WriteString("@")
			b.WriteString(item.Username)
			b.WriteString(" — ждёт /start\n")
		}
	}
	if waiting == 0 {
		b.WriteString("Все указанные получатели уже писали боту.")
	}
	return strings.TrimSpace(b.String())
}

func (s *TelegramService) getStored(ctx context.Context) (domain.TelegramRecipients, error) {
	if s == nil || s.recipients == nil {
		return domain.TelegramRecipients{}, nil
	}
	stored, err := s.recipients.GetTelegramRecipients(ctx)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return domain.TelegramRecipients{}, err
	}
	return stored, nil
}

func telegramCommand(text string) string {
	fields := strings.Fields(strings.TrimSpace(text))
	if len(fields) == 0 {
		return ""
	}
	command := fields[0]
	if i := strings.Index(command, "@"); i > 0 {
		command = command[:i]
	}
	return strings.ToLower(command)
}
