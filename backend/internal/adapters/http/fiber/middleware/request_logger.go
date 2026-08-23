package middleware

import (
	"log/slog"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func RequestLogger(log *slog.Logger, metrics *RequestMetrics) fiber.Handler {
	return func(c *fiber.Ctx) error {
		startedAt := time.Now()

		err := c.Next()

		duration := time.Since(startedAt)
		metrics.Observe(duration)

		status := c.Response().StatusCode()
		if err != nil {
			status = fiber.StatusInternalServerError
			if fiberErr, ok := err.(*fiber.Error); ok {
				status = fiberErr.Code
			}
		}

		c.Set("X-Request-ID", requestID(c))

		log.Info(
			"http request",
			slog.String("method", c.Method()),
			slog.String("path", sanitizePath(c.Path())),
			slog.Int("status", status),
			slog.Duration("duration", duration),
			slog.String("request_id", requestID(c)),
		)

		return err
	}
}

func sanitizePath(path string) string {
	lower := strings.ToLower(path)
	if strings.Contains(lower, "password") || strings.Contains(lower, "token") {
		return path
	}
	return path
}
