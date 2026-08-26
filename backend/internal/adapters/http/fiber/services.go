package fiberhttp

import (
	"palomnik/internal/adapters/http/fiber/middleware"
	"palomnik/internal/application"
	"palomnik/internal/ports"
)

type Services struct {
	Tours          *application.TourService
	Bookings       *application.BookingService
	Reviews        *application.ReviewService
	Integrations   *application.IntegrationService
	Webhooks       *application.WebhookService
	Auth           *application.AuthService
	Favorites      *application.FavoriteService
	Passengers     *application.PassengerService
	Support        *application.SupportService
	CMS            *application.CMSService
	News           *application.NewsService
	NewsEngagement *application.NewsEngagementService
	SMM            *application.SMMService
	Telegram       *application.TelegramService
	Notifications  *application.NotificationSettingsService
	SiteSettings   *application.SiteSettingsService
	AdminRoles     *application.AdminRoleService
	Legal          *application.LegalDocumentService
	Consents       *application.ConsentService
	Photos         *application.UserPhotoService
	Captcha        ports.CaptchaPort
	Messenger      ports.MessengerPort
	Publisher      ports.PublisherPort
	AI             ports.AIPort
	Payment        ports.PaymentPort
	AIFeatures     *application.AIFeaturesService
	WebhookGuard   *application.WebhookGuard
	RateLimiter    ports.RateLimiter
	Metrics        *middleware.RequestMetrics
	BackupLastPath string
	TourCatalog    application.TourCatalog
}
