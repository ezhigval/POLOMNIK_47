package application

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
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
	settings *NotificationSettingsService
	chats    ports.TelegramChatMapRepository
	sender   ports.TelegramBotSender
	fallback string
}

func NewTelegramService(
	settings *NotificationSettingsService,
	chats ports.TelegramChatMapRepository,
	sender ports.TelegramBotSender,
	fallbackChatID string,
) *TelegramService {
	return &TelegramService{
		settings: settings,
		chats:    chats,
		sender:   sender,
		fallback: strings.TrimSpace(fallbackChatID),
	}
}

// NewTelegramServiceFromRepos wires legacy telegram username lists through the
// notification routing settings service (memory/postgres stores implement both).
func NewTelegramServiceFromRepos(
	recipients ports.TelegramRecipientsRepository,
	chats ports.TelegramChatMapRepository,
	sender ports.TelegramBotSender,
	fallbackChatID string,
) *TelegramService {
	var routing ports.NotificationRoutingRepository
	if r, ok := recipients.(ports.NotificationRoutingRepository); ok {
		routing = r
	}
	telegramOn := sender != nil
	settings := NewNotificationSettingsService(routing, recipients, chats, telegramOn, false)
	return NewTelegramService(settings, chats, sender, fallbackChatID)
}

func (s *TelegramService) Settings(ctx context.Context) (TelegramSettingsView, error) {
	if s == nil || s.settings == nil {
		return TelegramSettingsView{}, nil
	}
	return s.settings.TelegramSettings(ctx)
}

func (s *TelegramService) UpdateRecipients(ctx context.Context, bookingRaw, supportRaw string) (TelegramSettingsView, error) {
	if s == nil || s.settings == nil {
		return TelegramSettingsView{}, nil
	}
	return s.settings.UpdateTelegramUsernames(ctx, bookingRaw, supportRaw)
}

func (s *TelegramService) NotificationSettings() *NotificationSettingsService {
	if s == nil {
		return nil
	}
	return s.settings
}

type TelegramInboundUpdate struct {
	Username string
	ChatID   string
	Text     string
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
	if s.settings == nil {
		return "Бот на связи. Список получателей сейчас недоступен."
	}
	view, err := s.settings.Settings(ctx)
	if err != nil {
		return "Бот на связи. Список получателей сейчас недоступен."
	}

	var b strings.Builder
	b.WriteString("Бот на связи.\n")
	waiting := 0
	totalTG := 0
	for _, event := range view.Events {
		for _, item := range event.Recipients {
			if item.Channel.ID != domain.NotificationChannelTelegram {
				continue
			}
			totalTG++
			if !item.Ready {
				waiting++
				b.WriteString("@")
				b.WriteString(item.Address)
				b.WriteString(" — ждёт /start\n")
			}
		}
	}
	if totalTG == 0 {
		if s.fallback != "" {
			b.WriteString("Списки в настройках пустые — используется запасной chat_id.")
		} else {
			b.WriteString("Списки получателей в настройках пустые.")
		}
		return strings.TrimSpace(b.String())
	}
	if waiting == 0 {
		b.WriteString("Все указанные Telegram-получатели уже писали боту.")
	}
	return strings.TrimSpace(b.String())
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
