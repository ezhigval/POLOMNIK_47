package memory

import (
	"context"
	"strings"

	"polomnik/internal/domain"
)

func (s *Store) GetTelegramRecipients(_ context.Context) (domain.TelegramRecipients, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.telegramRecipients == nil {
		return domain.TelegramRecipients{}, domain.ErrNotFound
	}
	return cloneTelegramRecipients(*s.telegramRecipients), nil
}

func (s *Store) UpsertTelegramRecipients(_ context.Context, settings domain.TelegramRecipients) (domain.TelegramRecipients, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	copy := cloneTelegramRecipients(settings)
	s.telegramRecipients = &copy
	return cloneTelegramRecipients(settings), nil
}

func (s *Store) UpsertTelegramChatBinding(_ context.Context, binding domain.TelegramChatBinding) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.telegramChats == nil {
		s.telegramChats = make(map[string]domain.TelegramChatBinding)
	}
	s.telegramChats[binding.Username] = binding
	return nil
}

func (s *Store) ListTelegramChatBindings(_ context.Context, usernames []string) (map[string]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make(map[string]string, len(usernames))
	for _, username := range usernames {
		key := strings.ToLower(strings.TrimSpace(username))
		if binding, ok := s.telegramChats[key]; ok {
			out[key] = binding.ChatID
		}
	}
	return out, nil
}

func cloneTelegramRecipients(settings domain.TelegramRecipients) domain.TelegramRecipients {
	settings.BookingUsernames = append([]string(nil), settings.BookingUsernames...)
	settings.SupportUsernames = append([]string(nil), settings.SupportUsernames...)
	return settings
}
