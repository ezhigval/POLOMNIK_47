package handlers

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"palomnik/internal/adapters/http/fiber/dto"
	"palomnik/internal/application"
	"palomnik/internal/config"
	"palomnik/internal/ports"
)

func SystemInfo(cfg config.Config, integrations *application.IntegrationService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		info := dto.SystemInfo{
			CRMAdapter:          cfg.CRMAdapter,
			AccountingAdapter:   cfg.AccountingAdapter,
			NotificationAdapter: cfg.NotificationAdapter,
			TelegramConfigured:  strings.TrimSpace(cfg.TelegramBotToken) != "",
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

		if integrations != nil {
			summary, err := integrations.OutboxSummary(c.Context())
			if err != nil {
				return respondError(c, err, MapError)
			}
			info.Outbox = toOutboxSummaryDTO(summary)
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
