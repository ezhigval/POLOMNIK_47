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
	Telegram       *application.TelegramService
	Notifications  *application.NotificationSettingsService
	SiteSettings   *application.SiteSettingsService
	AdminRoles     *application.AdminRoleService
	Captcha        ports.CaptchaPort
	Messenger      ports.MessengerPort
	Publisher      ports.PublisherPort
	WebhookGuard   *application.WebhookGuard
	RateLimiter    ports.RateLimiter
	Metrics        *middleware.RequestMetrics
	BackupLastPath string
}
