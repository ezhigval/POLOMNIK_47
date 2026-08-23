package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"palomnik/internal/adapters/http/fiber/dto"
	appmiddleware "palomnik/internal/adapters/http/fiber/middleware"
	"palomnik/internal/application"
)

func (h *Handler) CreateBooking(c *fiber.Ctx) error {
	var req dto.CreateBookingRequest
	if err := c.BodyParser(&req); err != nil {
		return writeAppError(c, &AppError{
			Status:  422,
			Code:    "VALIDATION_ERROR",
			Message: "Некорректные данные запроса",
		})
	}

	tourID, err := parseUUID(req.TourID)
	if err != nil {
		return writeAppError(c, err)
	}
	if err := validateEmail(req.Email); err != nil {
		return writeAppError(c, err)
	}

	result, err := h.bookings.CreateBooking(c.Context(), application.CreateBookingInput{
		TourID:      tourID,
		UserID:      optionalUserID(c),
		Name:        req.Name,
		Phone:       req.Phone,
		Email:       req.Email,
		PeopleCount: req.PeopleCount,
		Comment:     req.Comment,
		Source:      "web",
	})
	if err != nil {
		return respondError(c, err, MapBookingError)
	}

	return c.Status(fiber.StatusCreated).JSON(dto.DataEnvelope[dto.CreateBookingResponse]{
		Data: dto.ToCreateBookingResponse(result.Booking, result.IntegrationStatus),
	})
}

func optionalUserID(c *fiber.Ctx) *uuid.UUID {
	userID, ok := appmiddleware.UserIDFromContext(c)
	if !ok {
		return nil
	}
	return dto.UserIDPtr(userID)
}
