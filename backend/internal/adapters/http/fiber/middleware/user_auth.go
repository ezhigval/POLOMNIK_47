package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"polomnik/internal/adapters/http/fiber/dto"
	"polomnik/internal/application"
)

const UserIDLocalKey = "userID"

func OptionalUserAuth(auth *application.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, ok := parseBearerUserID(c, auth)
		if ok {
			c.Locals(UserIDLocalKey, userID)
		}
		return c.Next()
	}
}

func RequireUserAuth(auth *application.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, ok := parseBearerUserID(c, auth)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorEnvelope{
				Error: dto.ErrorBody{
					Code:    "UNAUTHORIZED",
					Message: "Нужно войти в аккаунт",
				},
			})
		}
		c.Locals(UserIDLocalKey, userID)
		return c.Next()
	}
}

func UserIDFromContext(c *fiber.Ctx) (uuid.UUID, bool) {
	value := c.Locals(UserIDLocalKey)
	userID, ok := value.(uuid.UUID)
	return userID, ok && userID != uuid.Nil
}

func parseBearerUserID(c *fiber.Ctx, auth *application.AuthService) (uuid.UUID, bool) {
	if auth == nil {
		return uuid.Nil, false
	}

	header := strings.TrimSpace(c.Get("Authorization"))
	if header == "" {
		return uuid.Nil, false
	}

	const prefix = "Bearer "
	if !strings.HasPrefix(header, prefix) {
		return uuid.Nil, false
	}

	userID, err := auth.ParseToken(strings.TrimSpace(header[len(prefix):]))
	if err != nil {
		return uuid.Nil, false
	}

	return userID, true
}
