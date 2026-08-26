package handlers

import (
	"time"

	"palomnik/internal/application"
	"palomnik/internal/domain"
	"palomnik/internal/ports"
)

type Handler struct {
	tours                 *application.TourService
	bookings              *application.BookingService
	reviews               *application.ReviewService
	integrations          *application.IntegrationService
	webhooks              *application.WebhookService
	auth                  *application.AuthService
	favorites             *application.FavoriteService
	passengers            *application.PassengerService
	support               *application.SupportService
	cms                   *application.CMSService
	news                  *application.NewsService
	smm                   *application.SMMService
	telegram              *application.TelegramService
	notifications         *application.NotificationSettingsService
	siteSettings          *application.SiteSettingsService
	adminRoles            *application.AdminRoleService
	legal                 *application.LegalDocumentService
	consents              *application.ConsentService
	photos                *application.UserPhotoService
	telegramWebhookSecret string
	captcha               ports.CaptchaPort
	webhookGuard          *application.WebhookGuard
	ai                    *application.AIFeaturesService
	newsEngagement        *application.NewsEngagementService
	tourCatalog           application.TourCatalog
}

func New(
	tours *application.TourService,
	bookings *application.BookingService,
	reviews *application.ReviewService,
	integrations *application.IntegrationService,
	webhooks *application.WebhookService,
	auth *application.AuthService,
	favorites *application.FavoriteService,
	passengers *application.PassengerService,
	support *application.SupportService,
	cms *application.CMSService,
	news *application.NewsService,
	smm *application.SMMService,
	telegram *application.TelegramService,
	notifications *application.NotificationSettingsService,
	siteSettings *application.SiteSettingsService,
	adminRoles *application.AdminRoleService,
	legal *application.LegalDocumentService,
	consents *application.ConsentService,
	photos *application.UserPhotoService,
	telegramWebhookSecret string,
	captcha ports.CaptchaPort,
	webhookGuard *application.WebhookGuard,
) *Handler {
	return &Handler{
		tours:                 tours,
		bookings:              bookings,
		reviews:               reviews,
		integrations:          integrations,
		webhooks:              webhooks,
		auth:                  auth,
		favorites:             favorites,
		passengers:            passengers,
		support:               support,
		cms:                   cms,
		news:                  news,
		smm:                   smm,
		telegram:              telegram,
		notifications:         notifications,
		siteSettings:          siteSettings,
		adminRoles:            adminRoles,
		legal:                 legal,
		consents:              consents,
		photos:                photos,
		telegramWebhookSecret: telegramWebhookSecret,
		captcha:               captcha,
		webhookGuard:          webhookGuard,
	}
}

func (h *Handler) WithAI(ai *application.AIFeaturesService) *Handler {
	if h != nil {
		h.ai = ai
	}
	return h
}

func (h *Handler) WithNewsEngagement(service *application.NewsEngagementService) *Handler {
	if h != nil {
		h.newsEngagement = service
	}
	return h
}

func (h *Handler) WithTourCatalog(catalog application.TourCatalog) *Handler {
	if h != nil {
		h.tourCatalog = catalog
	}
	return h
}

func (h *Handler) publicTourCatalog() domain.TourCatalogContext {
	if h == nil {
		return domain.TourCatalogContext{Today: time.Now().UTC()}
	}
	return h.tourCatalog.Context(time.Now().UTC())
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
