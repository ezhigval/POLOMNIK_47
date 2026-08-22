package fiberhttp

import (
	"polomnik/internal/application"
)

type Services struct {
	Tours        *application.TourService
	Bookings     *application.BookingService
	Reviews      *application.ReviewService
	Integrations *application.IntegrationService
	Webhooks     *application.WebhookService
	Auth         *application.AuthService
	Favorites    *application.FavoriteService
	Support      *application.SupportService
	CMS          *application.CMSService
	News         *application.NewsService
}
