package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"palomnik/internal/adapters/http/fiber/dto"
	"palomnik/internal/application"
)

func (h *Handler) ManagementGetTelegramSettings(c *fiber.Ctx) error {
	if h.telegram == nil {
		return writeAppError(c, &AppError{Status: 503, Code: "SERVICE_UNAVAILABLE", Message: "Telegram не настроен"})
	}
	settings, err := h.telegram.Settings(c.Context())
	if err != nil {
		return respondError(c, err, MapError)
	}
	return c.JSON(dto.DataEnvelope[dto.TelegramSettingsResponse]{Data: dto.ToTelegramSettings(settings)})
}

func (h *Handler) ManagementUpdateTelegramSettings(c *fiber.Ctx) error {
	if h.telegram == nil {
		return writeAppError(c, &AppError{Status: 503, Code: "SERVICE_UNAVAILABLE", Message: "Telegram не настроен"})
	}
	var req dto.TelegramSettingsUpsertRequest
	if err := c.BodyParser(&req); err != nil {
		return writeAppError(c, &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Некорректные данные запроса"})
	}
	settings, err := h.telegram.UpdateRecipients(c.Context(), req.BookingUsernames, req.SupportUsernames)
	if err != nil {
		return respondError(c, err, MapError)
	}
	return c.JSON(dto.DataEnvelope[dto.TelegramSettingsResponse]{Data: dto.ToTelegramSettings(settings)})
}

func (h *Handler) TelegramWebhook(c *fiber.Ctx) error {
	if h.telegram == nil {
		return c.SendStatus(fiber.StatusOK)
	}
	if secret := strings.TrimSpace(h.telegramWebhookSecret); secret != "" {
		got := strings.TrimSpace(c.Get("X-Telegram-Bot-Api-Secret-Token"))
		if !application.SecretEqual(got, secret) {
			return writeAppError(c, &AppError{Status: 401, Code: "UNAUTHORIZED", Message: "Некорректный webhook secret"})
		}
	}

	var req dto.TelegramUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return writeAppError(c, &AppError{Status: 400, Code: "BAD_REQUEST", Message: "Некорректное обновление Telegram"})
	}
	updateID := strings.TrimSpace(req.UpdateID.String())
	if h.webhookGuard != nil && updateID != "" && h.webhookGuard.AlreadyProcessed(c.Context(), "telegram", updateID) {
		return c.SendStatus(fiber.StatusOK)
	}
	if req.Message == nil {
		if h.webhookGuard != nil && updateID != "" {
			h.webhookGuard.Remember(c.Context(), "telegram", updateID)
		}
		return c.SendStatus(fiber.StatusOK)
	}

	err := h.telegram.HandleInboundUpdate(c.Context(), application.TelegramInboundUpdate{
		Username: req.Message.Username(),
		ChatID:   req.Message.ChatID(),
		Text:     req.Message.Text,
	})
	if err != nil {
		return respondError(c, err, MapError)
	}
	if h.webhookGuard != nil && updateID != "" {
		h.webhookGuard.Remember(c.Context(), "telegram", updateID)
	}
	return c.SendStatus(fiber.StatusOK)
}
