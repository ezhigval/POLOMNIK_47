package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"palomnik/internal/adapters/ai"
	cachenoop "palomnik/internal/adapters/cache/noop"
	rediscache "palomnik/internal/adapters/cache/redis"
	"palomnik/internal/adapters/captcha"
	fiberhttp "palomnik/internal/adapters/http/fiber"
	appmiddleware "palomnik/internal/adapters/http/fiber/middleware"
	"palomnik/internal/adapters/integration"
	"palomnik/internal/adapters/mail"
	"palomnik/internal/adapters/messenger"
	"palomnik/internal/adapters/notification"
	telegramnotify "palomnik/internal/adapters/notification/telegram"
	"palomnik/internal/adapters/payment"
	"palomnik/internal/adapters/phone"
	"palomnik/internal/adapters/publisher"
	memorylimit "palomnik/internal/adapters/ratelimit/memory"
	"palomnik/internal/adapters/repository/memory"
	"palomnik/internal/adapters/repository/postgres"
	"palomnik/internal/application"
	"palomnik/internal/config"
	"palomnik/internal/domain"
	"palomnik/internal/legal/operator"
	"palomnik/internal/logger"
	"palomnik/internal/ports"
)

func main() {
	os.Exit(run())
}

func run() int {
	cfg := config.Load()
	log := logger.New(cfg.LogLevel)

	if err := config.Validate(cfg); err != nil {
		log.Error("invalid configuration", slog.Any("error", err))
		return 1
	}

	tourRepo, bookingRepo, reviewRepo, integrationRefs, outboxRepo, userRepo, favoriteRepo, supportRepo, cmsRepo, newsRepo, pgStore, cleanupStore, err := openRepositories(context.Background(), cfg, log)
	if err != nil {
		log.Error("failed to open storage", slog.Any("error", err))
		return 1
	}
	defer cleanupStore()

	cache, redisCache, cleanupCache, err := openCache(cfg, log)
	if err != nil {
		log.Error("failed to open cache", slog.Any("error", err))
		return 1
	}
	defer cleanupCache()

	crm := integration.NewCRM(cfg, integrationRefs, outboxRepo)
	telegramDeps := notification.Deps{}
	if routing, ok := tourRepo.(ports.NotificationRoutingRepository); ok {
		telegramDeps.Routing = routing
	}
	if recipients, ok := tourRepo.(ports.TelegramRecipientsRepository); ok {
		telegramDeps.Legacy = recipients
	}
	if chats, ok := tourRepo.(ports.TelegramChatMapRepository); ok {
		telegramDeps.Chats = chats
	}
	notifier := notification.New(cfg, outboxRepo, telegramDeps)
	bookingService := application.NewBookingService(
		bookingRepo,
		tourRepo,
		crm,
		integration.NewAccounting(cfg, integrationRefs, outboxRepo),
		notifier,
		txManager(tourRepo),
	)

	var siteSettingsRepo ports.SiteSettingsRepository
	if repo, ok := tourRepo.(ports.SiteSettingsRepository); ok {
		siteSettingsRepo = repo
	}
	var adminRoleRepo ports.AdminRoleRepository
	if repo, ok := tourRepo.(ports.AdminRoleRepository); ok {
		adminRoleRepo = repo
	}

	notificationSettings := application.NewNotificationSettingsService(
		telegramDeps.Routing,
		telegramDeps.Legacy,
		telegramDeps.Chats,
		cfg.NotificationAdapter == "telegram" && cfg.TelegramBotToken != "",
		cfg.MaxBotToken != "",
	)
	siteDefaults := domain.SiteSettings{}
	siteSettings := application.NewSiteSettingsService(siteSettingsRepo, siteDefaults)
	adminRoles := application.NewAdminRoleService(adminRoleRepo, userRepo, cfg.AdminToken, cfg.JWTSecret)
	newsService := application.NewNewsService(newsRepo, cache)
	publisherPort := publisher.New(cfg, newsService)
	smmService := application.NewSMMService(smmRepo(tourRepo), publisherPort)
	tourService := application.NewTourService(tourRepo, cache, crm)
	messengerPort := messenger.New(cfg)
	supportService := application.NewSupportService(supportRepo, notifier).WithClientMessenger(messengerPort, userRepo)
	integrationService := application.NewIntegrationService(integrationRefs, outboxRepo)
	metrics := appmiddleware.NewRequestMetrics()
	diskPath := cfg.UploadDir
	if diskPath == "" {
		diskPath = "/"
	}
	var pingDB func(context.Context) error
	if pgStore != nil {
		pingDB = pgStore.Ping
	}
	aiPort := ai.New(cfg)
	aiFeatures := application.NewAIFeaturesService(
		aiPort,
		tourService,
		bookingService,
		supportService,
		integrationService,
		&application.Watchdog{
			Outbox:     outboxRepo,
			BackupPath: cfg.BackupLastPath,
			DiskPath:   diskPath,
			PingDB:     pingDB,
			Status5xx:  metrics.Status5xx,
		},
	)
	telegramService := application.NewTelegramService(
		notificationSettings,
		telegramDeps.Chats,
		telegramnotify.NewClient(cfg),
		cfg.TelegramChatID,
	).WithBot(application.TelegramBotDeps{
		Support:    supportService,
		Bookings:   bookingService,
		Tours:      tourService,
		AdminRoles: adminRoles,
		Users:      userRepo,
	})

	opConfig := operator.Default()
	var legalRepo ports.LegalDocumentRepository
	var consentRepo ports.ConsentRepository
	if repo, ok := tourRepo.(ports.LegalDocumentRepository); ok {
		legalRepo = repo
	}
	if repo, ok := tourRepo.(ports.ConsentRepository); ok {
		consentRepo = repo
	}
	legalService := application.NewLegalDocumentService(legalRepo, opConfig)
	consentService := application.NewConsentService(consentRepo, legalRepo)
	if err := legalService.BootstrapInitialDocuments(context.Background()); err != nil {
		log.Error("failed to bootstrap legal documents", slog.Any("error", err))
	}

	services := fiberhttp.Services{
		Tours:    tourService,
		Bookings: bookingService,
		Reviews: application.NewReviewService(
			reviewRepo,
			tourRepo,
			crm,
		),
		Integrations: integrationService,
		Webhooks: application.NewWebhookService(
			bookingService,
			integration.NewCRMInbound(cfg),
			cfg.BitrixInboundToken,
			cfg.IsProduction() || strings.EqualFold(cfg.CRMAdapter, "bitrix"),
		),
		Auth: application.NewAuthService(
			userRepo,
			bookingRepo,
			phone.New(cfg),
			mail.New(cfg),
			application.SocialAuthConfig{
				YandexConfigured:    strings.TrimSpace(cfg.YandexOAuthClientID) != "" && strings.TrimSpace(cfg.YandexOAuthClientSecret) != "",
				VKConfigured:        strings.TrimSpace(cfg.VKOAuthClientID) != "" && strings.TrimSpace(cfg.VKOAuthClientSecret) != "",
				MaxConfigured:       strings.TrimSpace(cfg.MaxOAuthClientID) != "" && strings.TrimSpace(cfg.MaxOAuthClientSecret) != "",
				TelegramConfigured:  cfg.EffectiveTelegramLoginBotToken() != "" && cfg.EffectiveTelegramLoginBotUsername() != "",
				TelegramBotUsername: cfg.EffectiveTelegramLoginBotUsername(),
			},
			cfg.JWTSecret,
			cfg.JWTTokenTTL,
			cfg.PublicSiteURL,
			txManager(tourRepo),
		),
		Favorites:      application.NewFavoriteService(favoriteRepo, tourRepo),
		Passengers:     application.NewPassengerService(passengerRepo(tourRepo)),
		Photos:         application.NewUserPhotoService(photoRepo(tourRepo)),
		Support:        supportService,
		CMS:            application.NewCMSService(cmsRepo),
		News:           newsService,
		SMM:            smmService,
		Telegram:       telegramService,
		Notifications:  notificationSettings,
		SiteSettings:   siteSettings,
		AdminRoles:     adminRoles,
		Legal:          legalService,
		Consents:       consentService,
		Captcha:        captcha.New(cfg),
		Messenger:      messengerPort,
		Publisher:      publisherPort,
		AI:             aiPort,
		Payment:        payment.New(cfg),
		AIFeatures:     aiFeatures,
		WebhookGuard:   application.NewWebhookGuard(cache),
		RateLimiter:    rateLimiterFromCache(redisCache),
		Metrics:        metrics,
		BackupLastPath: cfg.BackupLastPath,
		TourCatalog:    application.NewTourCatalog(),
	}
	if engagementRepo, ok := tourRepo.(ports.NewsEngagementRepository); ok {
		services.NewsEngagement = application.NewNewsEngagementService(engagementRepo, newsRepo, userRepo)
	}

	health := fiberhttp.HealthDeps{
		CacheRequired: false,
	}
	if pgStore != nil {
		health.PingDB = pgStore.Ping
	}
	if redisCache != nil {
		health.PingCache = redisCache.Ping
	}

	app := fiberhttp.NewRouter(cfg, log, services, health)

	if cfg.NotificationAdapter == "telegram" && cfg.TelegramBotToken != "" {
		webhookURL := cfg.EffectiveTelegramWebhookURL()
		if webhookURL != "" {
			client := telegramnotify.NewClient(cfg)
			if err := client.SetWebhook(context.Background(), webhookURL, application.TelegramWebhookSecret(cfg.InternalAPISecret)); err != nil {
				log.Error("telegram setWebhook failed", slog.Any("error", err))
			} else {
				log.Info("telegram webhook registered")
			}
		}
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		log.Info("starting http server", slog.String("addr", cfg.HTTPAddr))
		errCh <- app.Listen(cfg.HTTPAddr)
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()

		if err := app.ShutdownWithContext(shutdownCtx); err != nil {
			log.Error("http server shutdown failed", slog.Any("error", err))
			return 1
		}

		log.Info("http server stopped")
		return 0
	case err := <-errCh:
		log.Error("http server failed", slog.Any("error", err))
		return 1
	}
}

func openRepositories(
	ctx context.Context,
	cfg config.Config,
	log *slog.Logger,
) (ports.TourRepository, ports.BookingRepository, ports.ReviewRepository, ports.IntegrationReferenceRepository, ports.OutboxRepository, ports.UserRepository, ports.FavoriteRepository, ports.SupportRepository, ports.CMSRepository, ports.NewsRepository, *postgres.Store, func(), error) {
	if cfg.DatabaseURL != "" {
		store, err := postgres.Open(ctx, cfg.DatabaseURL)
		if err != nil {
			return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, func() {}, err
		}
		cleanup := func() { _ = store.Close() }
		return store, store, store, store, store, store, store, store, store, store, store, cleanup, nil
	}

	log.Warn("DATABASE_URL is empty, using in-memory storage")
	store := memory.NewStore()
	return store, store, store, store, store, store, store, store, store, store, nil, func() {}, nil
}

func openCache(cfg config.Config, log *slog.Logger) (ports.CachePort, *rediscache.Cache, func(), error) {
	if cfg.RedisURL == "" {
		log.Info("REDIS_URL is empty, cache disabled")
		return cachenoop.New(), nil, func() {}, nil
	}

	cache, err := rediscache.New(cfg.RedisURL)
	if err != nil {
		log.Warn("redis unavailable, continuing without cache", slog.Any("error", err))
		return cachenoop.New(), nil, func() {}, nil
	}

	return cache, cache, func() { _ = cache.Close() }, nil
}

func rateLimiterFromCache(redisCache *rediscache.Cache) ports.RateLimiter {
	if redisCache != nil {
		return redisCache
	}
	return memorylimit.New()
}

func txManager(repo ports.TourRepository) ports.TransactionManager {
	if tx, ok := repo.(ports.TransactionManager); ok {
		return tx
	}
	return nil
}

func passengerRepo(repo ports.TourRepository) ports.PassengerRepository {
	if passengers, ok := repo.(ports.PassengerRepository); ok {
		return passengers
	}
	return nil
}

func photoRepo(repo ports.TourRepository) ports.UserPhotoRepository {
	if photos, ok := repo.(ports.UserPhotoRepository); ok {
		return photos
	}
	return nil
}

func smmRepo(repo ports.TourRepository) ports.SMMPostRepository {
	if posts, ok := repo.(ports.SMMPostRepository); ok {
		return posts
	}
	return nil
}
