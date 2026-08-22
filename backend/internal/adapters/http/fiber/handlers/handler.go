package handlers

import (
	"polomnik/internal/application"
)

type Handler struct {
	tours                 *application.TourService
	bookings              *application.BookingService
	reviews               *application.ReviewService
	integrations          *application.IntegrationService
	webhooks              *application.WebhookService
	auth                  *application.AuthService
	favorites             *application.FavoriteService
	support               *application.SupportService
	cms                   *application.CMSService
	news                  *application.NewsService
	telegram              *application.TelegramService
	notifications         *application.NotificationSettingsService
	siteSettings          *application.SiteSettingsService
	adminRoles            *application.AdminRoleService
	telegramWebhookSecret string
}

func New(
	tours *application.TourService,
	bookings *application.BookingService,
	reviews *application.ReviewService,
	integrations *application.IntegrationService,
	webhooks *application.WebhookService,
	auth *application.AuthService,
	favorites *application.FavoriteService,
	support *application.SupportService,
	cms *application.CMSService,
	news *application.NewsService,
	telegram *application.TelegramService,
	notifications *application.NotificationSettingsService,
	siteSettings *application.SiteSettingsService,
	adminRoles *application.AdminRoleService,
	telegramWebhookSecret string,
) *Handler {
	return &Handler{
		tours:                 tours,
		bookings:              bookings,
		reviews:               reviews,
		integrations:          integrations,
		webhooks:              webhooks,
		auth:                  auth,
		favorites:             favorites,
		support:               support,
		cms:                   cms,
		news:                  news,
		telegram:              telegram,
		notifications:         notifications,
		siteSettings:          siteSettings,
		adminRoles:            adminRoles,
		telegramWebhookSecret: telegramWebhookSecret,
	}
}

func (h *Handler) notificationSettings() *application.NotificationSettingsService {
	if h == nil {
		return nil
	}
	if h.notifications != nil {
		return h.notifications
	}
	if h.telegram != nil {
		return h.telegram.NotificationSettings()
	}
	return nil
}
