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

	"palomnik/internal/adapters/http/fiber/dto"
	"palomnik/internal/adapters/http/fiber/handlers"
	appmiddleware "palomnik/internal/adapters/http/fiber/middleware"
	"palomnik/internal/application"
	"palomnik/internal/config"
	"palomnik/internal/domain"
)

func NewRouter(cfg config.Config, log *slog.Logger, services Services, health HealthDeps) *fiber.App {
	bodyLimit := 1 * 1024 * 1024
	if cfg.UploadMaxBytes > 0 {
		bodyLimit = cfg.UploadMaxBytes + 512*1024
	}

	app := fiber.New(fiber.Config{
		AppName:      "palomnik-api",
		ErrorHandler: errorHandler(log),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
		BodyLimit:    bodyLimit,
	})

	metrics := services.Metrics
	if metrics == nil {
		metrics = appmiddleware.NewRequestMetrics()
	}

	app.Use(requestid.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:  cfg.CORSAllowOrigins,
		AllowHeaders:  "Origin, Content-Type, Accept, Authorization, X-Admin-Token, X-Admin-Session, X-Internal-Secret",
		AllowMethods:  "GET,POST,PATCH,DELETE,OPTIONS",
		ExposeHeaders: "X-Request-ID, Retry-After",
	}))
	app.Use(appmiddleware.RequestLogger(log, metrics))

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
		services.Passengers,
		services.Support,
		services.CMS,
		services.News,
		services.SMM,
		services.Telegram,
		services.Notifications,
		services.SiteSettings,
		services.AdminRoles,
		services.Legal,
		services.Consents,
		services.Photos,
		application.TelegramWebhookSecret(cfg.InternalAPISecret),
		services.Captcha,
		services.WebhookGuard,
	).WithAI(services.AIFeatures).WithNewsEngagement(services.NewsEngagement).WithTourCatalog(services.TourCatalog)

	ready := handlers.NewReadinessChecker(health.PingDB, health.PingCache, false)

	app.Get("/health", handlers.Health)
	app.Get("/health/ready", handlers.HealthReady(ready))

	v1 := app.Group("/api/v1")
	v1.Get("/health", handlers.Health)
	v1.Get("/health/ready", handlers.HealthReady(ready))

	authLimiter := appmiddleware.RateLimitWithStore(services.RateLimiter, 20, time.Minute)
	bookingLimiter := appmiddleware.RateLimitWithStore(services.RateLimiter, 30, time.Minute)
	supportLimiter := appmiddleware.RateLimitWithStore(services.RateLimiter, 30, time.Minute)
	webhookLimiter := appmiddleware.RateLimitWithStore(services.RateLimiter, 60, time.Minute)
	telegramLimiter := appmiddleware.RateLimitWithStore(services.RateLimiter, 120, time.Minute)
	v1.Get("/auth/methods", h.AuthMethods)
	v1.Post("/auth/register", authLimiter, h.Register)
	v1.Post("/auth/login", authLimiter, h.Login)
	v1.Post("/auth/forgot-password", authLimiter, h.ForgotPassword)
	v1.Post("/auth/reset-password", authLimiter, h.ResetPassword)
	v1.Post("/auth/phone/start", authLimiter, h.StartPhoneVerification)
	v1.Get("/auth/phone/status", authLimiter, h.PhoneVerificationStatus)
	v1.Post("/auth/phone/complete", authLimiter, h.CompletePhoneLogin)
	v1.Post("/auth/oauth", appmiddleware.InternalAuth(cfg.InternalAPISecret), authLimiter, h.OAuthLogin)

	v1.Get("/tours/popular", h.ListPopularTours)
	v1.Get("/tours/:id/reviews", h.ListTourReviews)
	v1.Get("/tours/:id/recommendations", h.ListTourRecommendations)
	v1.Get("/tours/:id", h.GetTour)
	v1.Get("/tours", h.ListTours)

	v1.Get("/reviews", h.ListReviews)
	v1.Post("/reviews", appmiddleware.OptionalUserAuth(services.Auth), h.CreatePublicReview)
	v1.Get("/news", h.ListPublicNews)
	v1.Get("/news/:slug/likes", h.GetNewsLikeState)
	v1.Post("/news/:slug/likes", h.ToggleNewsLike)
	v1.Get("/news/:slug/comments", h.ListNewsComments)
	v1.Post("/news/:slug/comments", appmiddleware.RequireUserAuth(services.Auth), h.AddNewsComment)
	v1.Get("/news/:slug", h.GetPublicNewsBySlug)
	v1.Get("/pages", h.ListPublicCMSPages)
	v1.Get("/pages/:slug", h.GetPublicCMSPage)
	v1.Get("/site-settings", h.GetPublicSiteSettings)
	v1.Get("/legal/documents", h.ListPublicLegalDocuments)
	v1.Get("/legal/documents/:type/download", h.DownloadPublicLegalDocument)
	v1.Get("/legal/documents/:type", h.GetPublicLegalDocument)
	v1.Get("/legal/documents/:type/versions/:version", h.GetPublicLegalDocumentVersion)
	v1.Post("/consents", appmiddleware.OptionalUserAuth(services.Auth), h.RecordConsent)
	v1.Post("/bookings", bookingLimiter, appmiddleware.OptionalUserAuth(services.Auth), h.CreateBooking)
	v1.Post("/webhooks/bitrix/deal", webhookLimiter, h.BitrixDealWebhook)
	v1.Post("/webhooks/telegram", telegramLimiter, h.TelegramWebhook)

	me := v1.Group("/me", appmiddleware.RequireUserAuth(services.Auth))
	me.Get("/", h.Me)
	me.Patch("/", authLimiter, h.UpdateMe)
	me.Get("/identities", h.MyIdentities)
	me.Get("/passengers", h.ListMyPassengers)
	me.Post("/passengers", authLimiter, h.CreateMyPassenger)
	me.Patch("/passengers/:id", authLimiter, h.UpdateMyPassenger)
	me.Delete("/passengers/:id", h.DeleteMyPassenger)
	me.Get("/bookings", h.MyBookings)
	me.Get("/favorites", h.ListFavoriteIDs)
	me.Post("/favorites/:tourId", h.AddFavorite)
	me.Delete("/favorites/:tourId", h.RemoveFavorite)
	me.Get("/support", h.GetSupportThread)
	me.Post("/support/messages", supportLimiter, h.SendSupportMessage)
	me.Get("/consents", h.ListMyConsents)
	me.Get("/photos", h.ListMyPhotos)
	me.Post("/photos", authLimiter, h.CreateMyPhoto)
	me.Delete("/photos/:id", h.DeleteMyPhoto)
	me.Post("/uploads", appmiddleware.RateLimit(30, time.Minute), handlers.UploadImage(cfg))

	adminAuth := appmiddleware.AdminAuth(services.AdminRoles, cfg.AdminToken)
	require := func(perm domain.Permission) fiber.Handler {
		return appmiddleware.RequirePermission(services.AdminRoles, perm)
	}

	v1.Post("/management/auth/login", authLimiter, h.ManagementLogin)

	management := v1.Group("/management", adminAuth)
	management.Get("/session", h.ManagementSession)

	management.Get("/tours", require(domain.PermManageTours), h.ManagementListTours)
	management.Post("/tours", require(domain.PermManageTours), h.ManagementCreateTour)
	management.Get("/tours/:id", require(domain.PermManageTours), h.ManagementGetTour)
	management.Patch("/tours/:id", require(domain.PermManageTours), h.ManagementUpdateTour)
	management.Delete("/tours/:id", require(domain.PermManageTours), h.ManagementDeleteTour)
	management.Post("/uploads", require(domain.PermManageTours), appmiddleware.RateLimit(30, time.Minute), handlers.UploadImage(cfg))

	management.Get("/bookings", require(domain.PermManageBookings), h.ManagementListBookings)
	management.Get("/bookings/:id", require(domain.PermManageBookings), h.ManagementGetBooking)
	management.Patch("/bookings/:id/status", require(domain.PermManageBookings), h.ManagementUpdateBookingStatus)
	management.Patch("/bookings/:id/payment-status", require(domain.PermManageBookings), h.ManagementUpdateBookingPaymentStatus)

	management.Get("/support", require(domain.PermManageSupport), h.ManagementListSupportThreads)
	management.Get("/support/:id", require(domain.PermManageSupport), h.ManagementGetSupportThread)
	management.Post("/support/:id/draft", require(domain.PermManageSupport), appmiddleware.RateLimit(20, time.Minute), h.ManagementSupportDraft)
	management.Post("/support/:id/messages", require(domain.PermManageSupport), appmiddleware.RateLimit(60, time.Minute), h.ManagementSendSupportMessage)
	management.Get("/ai/metrics-digest", require(domain.PermViewStats), h.ManagementMetricsDigest)
	management.Get("/watchdog", require(domain.PermViewStats), h.ManagementWatchdog)

	management.Get("/reviews", require(domain.PermManageContent), h.ManagementListReviews)
	management.Post("/reviews", require(domain.PermManageContent), h.ManagementCreateReview)
	management.Patch("/reviews/:id/approve", require(domain.PermManageContent), h.ManagementApproveReview)
	management.Patch("/reviews/:id/reject", require(domain.PermManageContent), h.ManagementRejectReview)
	management.Patch("/reviews/:id/reply", require(domain.PermManageContent), h.ManagementSetReviewReply)
	management.Delete("/reviews/:id", require(domain.PermManageContent), h.ManagementDeleteReview)

	management.Get("/news", require(domain.PermManageContent), h.ManagementListNews)
	management.Post("/news", require(domain.PermManageContent), h.ManagementCreateNews)
	management.Get("/news/:id", require(domain.PermManageContent), h.ManagementGetNews)
	management.Patch("/news/:id/pin", require(domain.PermManageContent), h.ManagementSetNewsPinned)
	management.Patch("/news/:id", require(domain.PermManageContent), h.ManagementUpdateNews)
	management.Delete("/news/:id", require(domain.PermManageContent), h.ManagementDeleteNews)

	management.Get("/smm", require(domain.PermManageContent), h.ManagementListSMM)
	management.Post("/smm", require(domain.PermManageContent), h.ManagementCreateSMM)
	management.Get("/smm/:id", require(domain.PermManageContent), h.ManagementGetSMM)
	management.Post("/smm/:id/publish", require(domain.PermManageContent), h.ManagementPublishSMM)
	management.Delete("/smm/:id", require(domain.PermManageContent), h.ManagementDeleteSMM)

	management.Get("/telegram-settings", require(domain.PermManageRecipients), h.ManagementGetTelegramSettings)
	management.Patch("/telegram-settings", require(domain.PermManageRecipients), h.ManagementUpdateTelegramSettings)
	management.Get("/notification-settings", require(domain.PermManageRecipients), h.ManagementGetNotificationSettings)
	management.Patch("/notification-settings", require(domain.PermManageRecipients), h.ManagementUpdateNotificationSettings)

	management.Get("/site-settings", require(domain.PermManageSettingsSite), h.ManagementGetSiteSettings)
	management.Patch("/site-settings", require(domain.PermManageSettingsSite), h.ManagementUpdateSiteSettings)

	management.Get("/roles", require(domain.PermManageRoles), h.ManagementListRoles)
	management.Get("/roles/templates", require(domain.PermManageRoles), h.ManagementListRoleTemplates)
	management.Post("/roles", require(domain.PermManageRoles), h.ManagementCreateRole)
	management.Patch("/roles/:id", require(domain.PermManageRoles), h.ManagementUpdateRole)
	management.Delete("/roles/:id", require(domain.PermManageRoles), h.ManagementDeleteRole)
	management.Get("/roles/:id/assignments", require(domain.PermManageRoles), h.ManagementListRoleAssignments)
	management.Post("/roles/:id/assignments", require(domain.PermManageRoles), h.ManagementAssignRoleUser)
	management.Delete("/roles/:id/assignments/:userId", require(domain.PermManageRoles), h.ManagementUnassignRoleUser)

	management.Get("/integration-references", require(domain.PermManageIntegrations), h.ManagementListIntegrationReferences)
	management.Get("/outbox-events", require(domain.PermManageIntegrations), h.ManagementListOutboxEvents)
	management.Get("/system-info", require(domain.PermViewStats), handlers.SystemInfo(cfg, services.Integrations, metrics, services.BackupLastPath, services.Messenger, services.Publisher, services.AI, services.Payment))

	management.Get("/cms/templates", require(domain.PermManageContent), h.ManagementListCMSTemplates)
	management.Get("/cms/pages", require(domain.PermManageContent), h.ManagementListCMSPages)
	management.Post("/cms/pages/bootstrap-home", require(domain.PermManageContent), h.ManagementBootstrapHomeCMSPage)
	management.Get("/cms/pages/:id", require(domain.PermManageContent), h.ManagementGetCMSPage)
	management.Patch("/cms/pages/:id", require(domain.PermManageContent), h.ManagementUpdateCMSPage)
	management.Post("/cms/pages/:id/blocks", require(domain.PermManageContent), h.ManagementCreateCMSBlock)
	management.Post("/cms/pages/:id/blocks/reorder", require(domain.PermManageContent), h.ManagementReorderCMSBlocks)
	management.Patch("/cms/blocks/:id", require(domain.PermManageContent), h.ManagementUpdateCMSBlock)
	management.Delete("/cms/blocks/:id", require(domain.PermManageContent), h.ManagementDeleteCMSBlock)

	management.Get("/legal/documents", require(domain.PermManageContent), h.ManagementListLegalDocuments)
	management.Post("/legal/documents", require(domain.PermManageContent), h.ManagementPublishLegalDocument)
	management.Get("/consents", require(domain.PermManageContent), h.ManagementListConsents)

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
			log.Error("http error", slog.Int("status", status), slog.String("request_id", requestIDFrom(c)), slog.Any("error", err))
		}

		return c.Status(status).JSON(dto.ErrorEnvelope{
			Error: dto.ErrorBody{
				Code:      code,
				Message:   message,
				RequestID: requestIDFrom(c),
			},
		})
	}
}

func requestIDFrom(c *fiber.Ctx) string {
	if value, ok := c.Locals("requestid").(string); ok && value != "" {
		return value
	}
	return c.Get("X-Request-ID")
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
