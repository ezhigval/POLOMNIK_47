package memory

import (
	"context"

	"polomnik/internal/domain"
)

func (s *Store) GetNotificationRouting(_ context.Context) (domain.NotificationRouting, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.notificationRouting == nil {
		if s.telegramRecipients != nil {
			return domain.NotificationRoutingFromTelegramRecipients(*s.telegramRecipients, s.telegramRecipients.UpdatedAt), nil
		}
		return domain.NotificationRouting{}, domain.ErrNotFound
	}
	return cloneNotificationRouting(*s.notificationRouting), nil
}

func (s *Store) UpsertNotificationRouting(_ context.Context, routing domain.NotificationRouting) (domain.NotificationRouting, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	copy := cloneNotificationRouting(routing)
	s.notificationRouting = &copy
	return cloneNotificationRouting(routing), nil
}

func (s *Store) GetSiteSettings(_ context.Context) (domain.SiteSettings, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.siteSettings == nil {
		return domain.SiteSettings{}, domain.ErrNotFound
	}
	return cloneSiteSettings(*s.siteSettings), nil
}

func (s *Store) UpsertSiteSettings(_ context.Context, settings domain.SiteSettings) (domain.SiteSettings, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	copy := cloneSiteSettings(settings)
	s.siteSettings = &copy
	return cloneSiteSettings(settings), nil
}

func cloneNotificationRouting(routing domain.NotificationRouting) domain.NotificationRouting {
	out := domain.EmptyNotificationRouting(routing.UpdatedAt)
	out.UpdatedAt = routing.UpdatedAt
	if routing.ByEvent == nil {
		return out
	}
	for kind, list := range routing.ByEvent {
		out.ByEvent[kind] = append([]domain.NotificationRecipient(nil), list...)
	}
	return out
}

func cloneSiteSettings(settings domain.SiteSettings) domain.SiteSettings {
	settings.MailForwardTo = append([]string(nil), settings.MailForwardTo...)
	return settings
}
