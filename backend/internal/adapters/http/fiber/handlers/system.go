package handlers

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"palomnik/internal/adapters/http/fiber/dto"
	appmiddleware "palomnik/internal/adapters/http/fiber/middleware"
	"palomnik/internal/application"
	"palomnik/internal/config"
	"palomnik/internal/ports"
)

func SystemInfo(cfg config.Config, integrations *application.IntegrationService, metrics *appmiddleware.RequestMetrics, backupLastPath string, messenger ports.MessengerPort, publisher ports.PublisherPort, ai ports.AIPort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		info := dto.SystemInfo{
			CRMAdapter:          cfg.CRMAdapter,
			AccountingAdapter:   cfg.AccountingAdapter,
			NotificationAdapter: cfg.NotificationAdapter,
			MessengerAdapter:    cfg.MessengerAdapter,
			PublisherAdapter:    cfg.PublisherAdapter,
			AIAdapter:           cfg.AIAdapter,
			TelegramConfigured:  strings.TrimSpace(cfg.TelegramBotToken) != "",
			MessengerConfigured: messenger != nil && messenger.Configured(),
			PublisherConfigured: publisher != nil && publisher.Configured(),
			AIConfigured:        ai != nil && ai.Configured(),
			BitrixConfigured:    strings.TrimSpace(cfg.BitrixWebhookURL) != "",
			OneCConfigured:      strings.TrimSpace(cfg.OneCBaseURL) != "",
		}
		if info.CRMAdapter == "" {
			info.CRMAdapter = "noop"
		}
		if info.AccountingAdapter == "" {
			info.AccountingAdapter = "noop"
		}
		if info.NotificationAdapter == "" {
			info.NotificationAdapter = "noop"
		}
		if info.MessengerAdapter == "" {
			info.MessengerAdapter = "noop"
		}
		if info.PublisherAdapter == "" {
			info.PublisherAdapter = "noop"
		}
		if info.AIAdapter == "" {
			info.AIAdapter = "noop"
		}

		if integrations != nil {
			summary, err := integrations.OutboxSummary(c.Context())
			if err != nil {
				return respondError(c, err, MapError)
			}
			info.Outbox = toOutboxSummaryDTO(summary)
		}

		last, avg, count := metrics.Snapshot()
		info.Latency = dto.LatencyInfo{
			LastMS:   last.Milliseconds(),
			AvgMS:    avg.Milliseconds(),
			Requests: count,
		}

		backup := application.ReadBackupStatus(backupLastPath)
		info.LastBackup = dto.BackupInfo{
			File:         backup.File,
			Bytes:        backup.Bytes,
			Offsite:      backup.Offsite,
			OffsiteError: backup.OffsiteError,
		}
		if backup.At != nil {
			value := backup.At.UTC().Format(time.RFC3339)
			info.LastBackup.At = &value
		}

		return c.Status(fiber.StatusOK).JSON(dto.DataEnvelope[dto.SystemInfo]{Data: info})
	}
}

func toOutboxSummaryDTO(summary ports.OutboxSummary) dto.OutboxSummary {
	out := dto.OutboxSummary{
		Pending:           summary.Pending,
		Failed:            summary.Failed,
		Processed:         summary.Processed,
		LatestFailedError: summary.LatestFailedError,
	}
	if summary.OldestPendingAt != nil {
		value := summary.OldestPendingAt.UTC().Format(time.RFC3339)
		out.OldestPendingAt = &value
	}
	if summary.LatestFailedAt != nil {
		value := summary.LatestFailedAt.UTC().Format(time.RFC3339)
		out.LatestFailedAt = &value
	}
	return out
}
