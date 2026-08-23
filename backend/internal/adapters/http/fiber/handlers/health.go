package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"

	"palomnik/internal/adapters/http/fiber/dto"
)

type HealthResponse struct {
	Status    string            `json:"status"`
	Service   string            `json:"service"`
	Timestamp string            `json:"timestamp"`
	Checks    map[string]string `json:"checks,omitempty"`
}

func Health(c *fiber.Ctx) error {
	return c.JSON(dto.DataEnvelope[HealthResponse]{
		Data: HealthResponse{
			Status:    "ok",
			Service:   "palomnik-api",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
	})
}

type readinessChecker struct {
	pingDB        func(context.Context) error
	pingCache     func(context.Context) error
	cacheRequired bool
}

func HealthReady(deps readinessChecker) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
		defer cancel()

		checks := map[string]string{}
		ready := true

		if deps.pingDB != nil {
			if err := deps.pingDB(ctx); err != nil {
				checks["database"] = "error"
				ready = false
			} else {
				checks["database"] = "ok"
			}
		}

		if deps.cacheRequired {
			if deps.pingCache == nil {
				checks["cache"] = "error"
				ready = false
			} else if err := deps.pingCache(ctx); err != nil {
				checks["cache"] = "error"
				ready = false
			} else {
				checks["cache"] = "ok"
			}
		}

		status := "ok"
		statusCode := fiber.StatusOK
		if !ready {
			status = "degraded"
			statusCode = fiber.StatusServiceUnavailable
		}

		return c.Status(statusCode).JSON(dto.DataEnvelope[HealthResponse]{
			Data: HealthResponse{
				Status:    status,
				Service:   "palomnik-api",
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				Checks:    checks,
			},
		})
	}
}

func NewReadinessChecker(pingDB, pingCache func(context.Context) error, cacheRequired bool) readinessChecker {
	return readinessChecker{
		pingDB:        pingDB,
		pingCache:     pingCache,
		cacheRequired: cacheRequired,
	}
}
