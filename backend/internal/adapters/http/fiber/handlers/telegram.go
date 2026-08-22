package handlers

import (
	"crypto/subtle"
	"strings"

	"github.com/gofiber/fiber/v2"

	"polomnik/internal/adapters/http/fiber/dto"
	"polomnik/internal/application"
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
		if subtle.ConstantTimeCompare([]byte(got), []byte(secret)) != 1 {
			return writeAppError(c, &AppError{Status: 401, Code: "UNAUTHORIZED", Message: "Некорректный webhook secret"})
		}
	}

	var req dto.TelegramUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return writeAppError(c, &AppError{Status: 400, Code: "BAD_REQUEST", Message: "Некорректное обновление Telegram"})
	}
	if req.Message == nil {
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
	return c.SendStatus(fiber.StatusOK)
}
