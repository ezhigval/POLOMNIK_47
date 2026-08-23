package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"palomnik/internal/adapters/integration"
	"palomnik/internal/adapters/notification"
	"palomnik/internal/adapters/repository/memory"
	"palomnik/internal/adapters/repository/postgres"
	"palomnik/internal/application"
	"palomnik/internal/config"
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

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	tourRepo, bookingRepo, reviewRepo, integrationRefs, outboxRepo, cleanupStore, err := openRepositories(ctx, cfg, log)
	if err != nil {
		log.Error("failed to open storage", slog.Any("error", err))
		return 1
	}
	defer cleanupStore()

	crm := integration.NewCRM(cfg, integrationRefs, outboxRepo)
	accounting := integration.NewAccounting(cfg, integrationRefs, outboxRepo)
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
	notifications := notification.Inner(cfg, telegramDeps)
	worker := application.NewOutboxWorker(
		outboxRepo,
		tourRepo,
		bookingRepo,
		reviewRepo,
		crm,
		accounting,
		notifications,
		cfg.OutboxWorkerMaxAttempts,
	)

	log.Info(
		"outbox worker started",
		slog.Duration("poll_interval", cfg.OutboxWorkerPollInterval),
		slog.Int("batch_size", cfg.OutboxWorkerBatchSize),
		slog.Int("max_attempts", cfg.OutboxWorkerMaxAttempts),
		slog.String("heartbeat_path", cfg.WorkerHeartbeatPath),
	)

	ticker := time.NewTicker(cfg.OutboxWorkerPollInterval)
	defer ticker.Stop()

	for {
		processed, err := worker.ProcessBatch(ctx, cfg.OutboxWorkerBatchSize)
		if err != nil {
			log.Error("outbox batch failed", slog.Any("error", err))
		} else {
			if processed > 0 {
				log.Info("outbox batch processed", slog.Int("count", processed))
			}
			if hbErr := writeHeartbeat(cfg.WorkerHeartbeatPath); hbErr != nil {
				log.Error("failed to write worker heartbeat", slog.Any("error", hbErr))
			}
		}

		select {
		case <-ctx.Done():
			log.Info("outbox worker stopped")
			return 0
		case <-ticker.C:
		}
	}
}

func writeHeartbeat(path string) error {
	if path == "" {
		return nil
	}
	dir := filepath.Dir(path)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	return os.WriteFile(path, []byte(time.Now().UTC().Format(time.RFC3339)+"\n"), 0o644)
}

func openRepositories(
	ctx context.Context,
	cfg config.Config,
	log *slog.Logger,
) (ports.TourRepository, ports.BookingRepository, ports.ReviewRepository, ports.IntegrationReferenceRepository, ports.OutboxRepository, func(), error) {
	if cfg.DatabaseURL != "" {
		store, err := postgres.Open(ctx, cfg.DatabaseURL)
		if err != nil {
			return nil, nil, nil, nil, nil, func() {}, err
		}
		cleanup := func() { _ = store.Close() }
		return store, store, store, store, store, cleanup, nil
	}

	log.Warn("DATABASE_URL is empty, using in-memory storage")
	store := memory.NewStore()
	return store, store, store, store, store, func() {}, nil
}
