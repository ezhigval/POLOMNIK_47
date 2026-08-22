package fiberhttp

import (
	"errors"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	"polomnik/internal/adapters/http/fiber/dto"
	"polomnik/internal/adapters/http/fiber/handlers"
	appmiddleware "polomnik/internal/adapters/http/fiber/middleware"
	"polomnik/internal/application"
	"polomnik/internal/config"
)

func NewRouter(cfg config.Config, log *slog.Logger, services Services, health HealthDeps) *fiber.App {
	bodyLimit := 1 * 1024 * 1024
	if cfg.UploadMaxBytes > 0 {
		// Multipart overhead on top of the file itself.
		bodyLimit = cfg.UploadMaxBytes + 512*1024
	}

	app := fiber.New(fiber.Config{
		AppName:      "polomnik-api",
		ErrorHandler: errorHandler(log),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
		BodyLimit:    bodyLimit,
	})

	app.Use(requestid.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.CORSAllowOrigins,
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, X-Admin-Token, X-Internal-Secret",
		AllowMethods: "GET,POST,PATCH,DELETE,OPTIONS",
	}))
	app.Use(appmiddleware.RequestLogger(log))

	if cfg.UploadDir != "" {
		if err := os.MkdirAll(cfg.UploadDir, 0o755); err != nil {
			log.Error("failed to create upload dir", slog.String("dir", cfg.UploadDir), slog.Any("error", err))
		} else {
			app.Static("/uploads", cfg.UploadDir)
		}
	}

	h := handlers.New(
		services.Tours,
		services.Bookings,
		services.Reviews,
		services.Integrations,
		services.Webhooks,
		services.Auth,
		services.Favorites,
		services.Support,
		services.CMS,
		services.News,
		services.Telegram,
		application.TelegramWebhookSecret(cfg.InternalAPISecret),
	)

	ready := handlers.NewReadinessChecker(health.PingDB, health.PingCache, health.CacheRequired)

	app.Get("/health", handlers.Health)
	app.Get("/health/ready", handlers.HealthReady(ready))

	v1 := app.Group("/api/v1")
	v1.Get("/health", handlers.Health)
	v1.Get("/health/ready", handlers.HealthReady(ready))

	authLimiter := appmiddleware.RateLimit(20, time.Minute)
	v1.Post("/auth/register", authLimiter, h.Register)
	v1.Post("/auth/login", authLimiter, h.Login)
	v1.Post("/auth/oauth", appmiddleware.InternalAuth(cfg.InternalAPISecret), authLimiter, h.OAuthLogin)

	v1.Get("/tours/popular", h.ListPopularTours)
	v1.Get("/tours/:id/reviews", h.ListTourReviews)
	v1.Get("/tours/:id", h.GetTour)
	v1.Get("/tours", h.ListTours)

	v1.Get("/reviews", h.ListReviews)
	v1.Get("/news", h.ListPublicNews)
	v1.Get("/pages", h.ListPublicCMSPages)
	v1.Get("/pages/:slug", h.GetPublicCMSPage)
	v1.Post("/bookings", appmiddleware.RateLimit(30, time.Minute), appmiddleware.OptionalUserAuth(services.Auth), h.CreateBooking)
	v1.Post("/webhooks/bitrix/deal", appmiddleware.RateLimit(60, time.Minute), h.BitrixDealWebhook)
	v1.Post("/webhooks/telegram", appmiddleware.RateLimit(120, time.Minute), h.TelegramWebhook)

	me := v1.Group("/me", appmiddleware.RequireUserAuth(services.Auth))
	me.Get("/", h.Me)
	me.Get("/bookings", h.MyBookings)
	me.Get("/favorites", h.ListFavoriteIDs)
	me.Post("/favorites/:tourId", h.AddFavorite)
	me.Delete("/favorites/:tourId", h.RemoveFavorite)
	me.Get("/support", h.GetSupportThread)
	me.Post("/support/messages", appmiddleware.RateLimit(30, time.Minute), h.SendSupportMessage)

	management := v1.Group("/management", appmiddleware.AdminAuth(cfg.AdminToken))
	management.Get("/tours", h.ManagementListTours)
	management.Post("/tours", h.ManagementCreateTour)
	management.Get("/tours/:id", h.ManagementGetTour)
	management.Patch("/tours/:id", h.ManagementUpdateTour)
	management.Delete("/tours/:id", h.ManagementDeleteTour)
	management.Post("/uploads", appmiddleware.RateLimit(30, time.Minute), handlers.UploadImage(cfg))

	management.Get("/bookings", h.ManagementListBookings)
	management.Get("/bookings/:id", h.ManagementGetBooking)
	management.Patch("/bookings/:id/status", h.ManagementUpdateBookingStatus)

	management.Get("/reviews", h.ManagementListReviews)
	management.Post("/reviews", h.ManagementCreateReview)
	management.Patch("/reviews/:id/approve", h.ManagementApproveReview)
	management.Patch("/reviews/:id/reject", h.ManagementRejectReview)
	management.Patch("/reviews/:id/reply", h.ManagementSetReviewReply)
	management.Delete("/reviews/:id", h.ManagementDeleteReview)

	management.Get("/news", h.ManagementListNews)
	management.Post("/news", h.ManagementCreateNews)
	management.Get("/news/:id", h.ManagementGetNews)
	management.Patch("/news/:id", h.ManagementUpdateNews)
	management.Delete("/news/:id", h.ManagementDeleteNews)

	management.Get("/telegram-settings", h.ManagementGetTelegramSettings)
	management.Patch("/telegram-settings", h.ManagementUpdateTelegramSettings)

	management.Get("/integration-references", h.ManagementListIntegrationReferences)
	management.Get("/outbox-events", h.ManagementListOutboxEvents)
	management.Get("/system-info", handlers.SystemInfo(cfg, services.Integrations))

	management.Get("/cms/templates", h.ManagementListCMSTemplates)
	management.Get("/cms/pages", h.ManagementListCMSPages)
	management.Post("/cms/pages/bootstrap-home", h.ManagementBootstrapHomeCMSPage)
	management.Post("/cms/pages", h.ManagementCreateCMSPage)
	management.Get("/cms/pages/:id", h.ManagementGetCMSPage)
	management.Patch("/cms/pages/:id", h.ManagementUpdateCMSPage)
	management.Delete("/cms/pages/:id", h.ManagementDeleteCMSPage)
	management.Post("/cms/pages/:id/blocks", h.ManagementCreateCMSBlock)
	management.Post("/cms/pages/:id/blocks/reorder", h.ManagementReorderCMSBlocks)
	management.Patch("/cms/blocks/:id", h.ManagementUpdateCMSBlock)
	management.Delete("/cms/blocks/:id", h.ManagementDeleteCMSBlock)

	return app
}

func errorHandler(log *slog.Logger) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		var appErr *handlers.AppError
		if errors.As(err, &appErr) {
			return c.Status(appErr.Status).JSON(dto.ErrorEnvelope{
				Error: dto.ErrorBody{
					Code:    appErr.Code,
					Message: appErr.Message,
				},
			})
		}

		if mapped := handlers.MapError(err); mapped != nil {
			return c.Status(mapped.Status).JSON(dto.ErrorEnvelope{
				Error: dto.ErrorBody{
					Code:    mapped.Code,
					Message: mapped.Message,
				},
			})
		}

		status := fiber.StatusInternalServerError
		code := "INTERNAL_ERROR"
		message := "Внутренняя ошибка сервера"

		if fiberErr, ok := err.(*fiber.Error); ok {
			status = fiberErr.Code
			code = errorCodeForStatus(status)
			message = http.StatusText(status)
			if fiberErr.Message != "" && status < fiber.StatusInternalServerError {
				message = fiberErr.Message
			}
		}

		if status >= fiber.StatusInternalServerError {
			log.Error("http error", slog.Int("status", status), slog.Any("error", err))
		}

		return c.Status(status).JSON(dto.ErrorEnvelope{
			Error: dto.ErrorBody{
				Code:    code,
				Message: message,
			},
		})
	}
}

func errorCodeForStatus(status int) string {
	switch status {
	case fiber.StatusBadRequest:
		return "BAD_REQUEST"
	case fiber.StatusUnauthorized:
		return "UNAUTHORIZED"
	case fiber.StatusForbidden:
		return "FORBIDDEN"
	case fiber.StatusNotFound:
		return "NOT_FOUND"
	case fiber.StatusConflict:
		return "CONFLICT"
	case fiber.StatusUnprocessableEntity:
		return "VALIDATION_ERROR"
	case fiber.StatusTooManyRequests:
		return "RATE_LIMITED"
	case fiber.StatusRequestEntityTooLarge:
		return "FILE_TOO_LARGE"
	case fiber.StatusServiceUnavailable:
		return "SERVICE_UNAVAILABLE"
	default:
		return "INTERNAL_ERROR"
	}
}
