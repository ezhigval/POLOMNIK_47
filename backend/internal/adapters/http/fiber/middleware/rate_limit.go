package middleware

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	memorylimit "polomnik/internal/adapters/ratelimit/memory"
	"polomnik/internal/ports"
)

func RateLimit(max int, window time.Duration) fiber.Handler {
	return RateLimitWithStore(nil, max, window)
}

func RateLimitWithStore(store ports.RateLimiter, max int, window time.Duration) fiber.Handler {
	fallback := memorylimit.New()
	if max < 1 {
		max = 1
	}
	if window <= 0 {
		window = time.Minute
	}

	return func(c *fiber.Ctx) error {
		key := c.IP() + ":" + c.Method() + ":" + c.Path()
		result, err := allow(c, store, fallback, key, max, window)
		if err != nil {
			result, err = fallback.Allow(c.Context(), key, max, window)
			if err != nil {
				return c.Next()
			}
		}
		if result.Allowed {
			return c.Next()
		}

		retryAfter := int(result.RetryAfter.Seconds())
		if retryAfter < 1 {
			retryAfter = 1
		}
		c.Set("Retry-After", strconv.Itoa(retryAfter))
		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
			"error": fiber.Map{
				"code":       "RATE_LIMITED",
				"message":    "Too many requests, try again later",
				"request_id": requestID(c),
			},
		})
	}
}

func allow(c *fiber.Ctx, store, fallback ports.RateLimiter, key string, max int, window time.Duration) (ports.RateLimitResult, error) {
	if store != nil {
		return store.Allow(c.Context(), key, max, window)
	}
	return fallback.Allow(c.Context(), key, max, window)
}

func requestID(c *fiber.Ctx) string {
	if value, ok := c.Locals("requestid").(string); ok {
		return value
	}
	return c.Get("X-Request-ID")
}
