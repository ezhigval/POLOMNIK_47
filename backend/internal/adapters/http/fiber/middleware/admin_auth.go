package middleware

import (
	"github.com/gofiber/fiber/v2"

	"polomnik/internal/adapters/http/fiber/dto"
)

func AdminAuth(adminToken string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if adminToken == "" {
			return c.Status(fiber.StatusServiceUnavailable).JSON(dto.ErrorEnvelope{
				Error: dto.ErrorBody{
					Code:    "SERVICE_UNAVAILABLE",
					Message: "Management API is not configured",
				},
			})
		}

		if c.Get("X-Admin-Token") != adminToken {
			return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorEnvelope{
				Error: dto.ErrorBody{
					Code:    "UNAUTHORIZED",
					Message: "Invalid admin token",
				},
			})
		}

		return c.Next()
	}
}
