package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"polomnik/internal/application"
)

type bitrixWebhookRequest struct {
	Event string `json:"event"`
	Data  struct {
		FIELDS struct {
			ID      string `json:"ID"`
			STAGEID string `json:"STAGE_ID"`
		} `json:"FIELDS"`
	} `json:"data"`
}

func (h *Handler) BitrixDealWebhook(c *fiber.Ctx) error {
	var body bitrixWebhookRequest
	if err := c.BodyParser(&body); err != nil {
		return &AppError{Status: fiber.StatusBadRequest, Code: "BAD_REQUEST", Message: "invalid webhook payload"}
	}

	dealID := strings.TrimSpace(body.Data.FIELDS.ID)
	if dealID == "" {
		dealID = strings.TrimSpace(c.Query("deal_id"))
	}

	err := h.webhooks.HandleBitrixDealUpdate(c.Context(), application.BitrixDealWebhookInput{
		Event:  body.Event,
		DealID: dealID,
		Token:  c.Query("token"),
	})
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}
