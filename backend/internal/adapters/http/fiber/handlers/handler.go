package handlers

import (
	"polomnik/internal/application"
)

type Handler struct {
	tours        *application.TourService
	bookings     *application.BookingService
	reviews      *application.ReviewService
	integrations *application.IntegrationService
	webhooks     *application.WebhookService
	auth         *application.AuthService
	favorites    *application.FavoriteService
	support      *application.SupportService
	cms          *application.CMSService
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
) *Handler {
	return &Handler{
		tours:        tours,
		bookings:     bookings,
		reviews:      reviews,
		integrations: integrations,
		webhooks:     webhooks,
		auth:         auth,
		favorites:    favorites,
		support:      support,
		cms:          cms,
	}
}
