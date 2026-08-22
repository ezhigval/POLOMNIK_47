package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"

	cachenoop "polomnik/internal/adapters/cache/noop"
	rediscache "polomnik/internal/adapters/cache/redis"
	fiberhttp "polomnik/internal/adapters/http/fiber"
	"polomnik/internal/adapters/integration"
	"polomnik/internal/adapters/notification"
	"polomnik/internal/adapters/repository/memory"
	"polomnik/internal/adapters/repository/postgres"
	"polomnik/internal/application"
	"polomnik/internal/config"
	"polomnik/internal/logger"
	"polomnik/internal/ports"
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
	notifier := notification.New(cfg, outboxRepo)
	bookingService := application.NewBookingService(
		bookingRepo,
		tourRepo,
		crm,
		integration.NewAccounting(cfg, integrationRefs, outboxRepo),
		notifier,
		txManager(tourRepo),
	)

	services := fiberhttp.Services{
		Tours: application.NewTourService(
			tourRepo,
			cache,
			crm,
		),
		Bookings: bookingService,
		Reviews: application.NewReviewService(
			reviewRepo,
			tourRepo,
			crm,
		),
		Integrations: application.NewIntegrationService(integrationRefs, outboxRepo),
		Webhooks: application.NewWebhookService(
			bookingService,
			integration.NewCRMInbound(cfg),
			cfg.BitrixInboundToken,
			cfg.IsProduction() || strings.EqualFold(cfg.CRMAdapter, "bitrix"),
		),
		Auth:      application.NewAuthService(userRepo, bookingRepo, cfg.JWTSecret, cfg.JWTTokenTTL),
		Favorites: application.NewFavoriteService(favoriteRepo, tourRepo),
		Support:   application.NewSupportService(supportRepo),
		CMS:       application.NewCMSService(cmsRepo),
		News:      application.NewNewsService(newsRepo),
	}

	health := fiberhttp.HealthDeps{
		CacheRequired: cfg.RedisURL != "",
	}
	if pgStore != nil {
		health.PingDB = pgStore.Ping
	}
	if redisCache != nil {
		health.PingCache = redisCache.Ping
	}

	app := fiberhttp.NewRouter(cfg, log, services, health)

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
		return nil, nil, func() {}, err
	}

	return cache, cache, func() { _ = cache.Close() }, nil
}

func txManager(repo ports.TourRepository) ports.TransactionManager {
	if tx, ok := repo.(ports.TransactionManager); ok {
		return tx
	}
	return nil
}
