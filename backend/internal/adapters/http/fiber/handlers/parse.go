package handlers

import (
	"errors"
	"net/mail"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"palomnik/internal/ports"
)

func parsePagination(c *fiber.Ctx) ports.Pagination {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	return ports.NormalizePagination(page, limit)
}

func parseOptionalDate(value string) (*time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}
	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		return nil, &AppError{
			Status:  422,
			Code:    "VALIDATION_ERROR",
			Message: "Некорректная дата, ожидается ГГГГ-ММ-ДД",
		}
	}
	return &parsed, nil
}

func parseRequiredDate(value string) (time.Time, error) {
	parsed, err := parseOptionalDate(value)
	if err != nil {
		return time.Time{}, err
	}
	if parsed == nil {
		return time.Time{}, &AppError{
			Status:  422,
			Code:    "VALIDATION_ERROR",
			Message: "Укажите дату",
		}
	}
	return *parsed, nil
}

func parseUUID(value string) (uuid.UUID, error) {
	id, err := uuid.Parse(strings.TrimSpace(value))
	if err != nil {
		return uuid.Nil, &AppError{
			Status:  422,
			Code:    "VALIDATION_ERROR",
			Message: "Некорректный идентификатор",
		}
	}
	return id, nil
}

func parseOptionalInt(value string) (*int, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return nil, &AppError{
			Status:  422,
			Code:    "VALIDATION_ERROR",
			Message: "Некорректное число",
		}
	}
	return &parsed, nil
}

func parseOptionalBool(value string) *bool {
	value = strings.TrimSpace(strings.ToLower(value))
	if value == "" {
		return nil
	}
	parsed := value == "true" || value == "1"
	return &parsed
}

func validateEmail(email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return nil
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return &AppError{
			Status:  422,
			Code:    "VALIDATION_ERROR",
			Message: "Некорректный email",
		}
	}
	return nil
}

func respondError(c *fiber.Ctx, err error, mapper func(error) *AppError) error {
	if appErr := asAppError(err); appErr != nil {
		return writeAppError(c, appErr)
	}
	if appErr := mapper(err); appErr != nil {
		return writeAppError(c, appErr)
	}
	return err
}

func asAppError(err error) *AppError {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	return nil
}

func writeAppError(c *fiber.Ctx, err error) error {
	appErr := asAppError(err)
	if appErr == nil {
		return c.Status(500).JSON(fiber.Map{
			"error": fiber.Map{
				"code":       "INTERNAL_ERROR",
				"message":    "Внутренняя ошибка сервера",
				"request_id": requestIDFrom(c),
			},
		})
	}
	return c.Status(appErr.Status).JSON(fiber.Map{
		"error": fiber.Map{
			"code":       appErr.Code,
			"message":    appErr.Message,
			"request_id": requestIDFrom(c),
		},
	})
}

func requestIDFrom(c *fiber.Ctx) string {
	if value, ok := c.Locals("requestid").(string); ok && value != "" {
		return value
	}
	return c.Get("X-Request-ID")
}
