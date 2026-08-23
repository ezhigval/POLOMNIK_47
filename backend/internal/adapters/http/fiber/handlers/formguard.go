package handlers

import (
	"github.com/gofiber/fiber/v2"

	"polomnik/internal/application"
)

func rejectHoneypot(value string) error {
	if application.HoneypotTriggered(value) {
		return &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Некорректные данные запроса"}
	}
	return nil
}

func (h *Handler) verifyCaptcha(c *fiber.Ctx, token string) error {
	if h == nil || h.captcha == nil || !h.captcha.Configured() {
		return nil
	}
	if err := h.captcha.Verify(c.Context(), token, c.IP()); err != nil {
		return &AppError{Status: 422, Code: "CAPTCHA_FAILED", Message: "Подтвердите, что вы не робот"}
	}
	return nil
}

func captchaDTO(h *Handler) dtoCaptchaStatus {
	if h == nil || h.captcha == nil || !h.captcha.Configured() {
		return dtoCaptchaStatus{Available: false, Provider: "noop"}
	}
	return dtoCaptchaStatus{
		Available: true,
		Provider:  h.captcha.Provider(),
		ClientKey: h.captcha.ClientKey(),
	}
}

type dtoCaptchaStatus struct {
	Available bool
	Provider  string
	ClientKey string
}
