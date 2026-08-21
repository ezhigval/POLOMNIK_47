package middleware

import (
	"github.com/gofiber/fiber/v2"

	"polomnik/internal/adapters/http/fiber/dto"
)

func InternalAuth(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if secret == "" {
			return c.Status(fiber.StatusServiceUnavailable).JSON(dto.ErrorEnvelope{
				Error: dto.ErrorBody{
					Code:    "SERVICE_UNAVAILABLE",
					Message: "Internal API is not configured",
				},
			})
		}

		if c.Get("X-Internal-Secret") != secret {
			return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorEnvelope{
				Error: dto.ErrorBody{
					Code:    "UNAUTHORIZED",
					Message: "Invalid internal secret",
				},
			})
		}

		return c.Next()
	}
}
