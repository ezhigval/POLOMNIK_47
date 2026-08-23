package handlers

import (
	"github.com/gofiber/fiber/v2"

	"palomnik/internal/adapters/http/fiber/dto"
	appmiddleware "palomnik/internal/adapters/http/fiber/middleware"
	"palomnik/internal/application"
)

func (h *Handler) ListMyPassengers(c *fiber.Ctx) error {
	userID, ok := appmiddleware.UserIDFromContext(c)
	if !ok {
		return writeAppError(c, &AppError{Status: 401, Code: "UNAUTHORIZED", Message: "Нужно войти в аккаунт"})
	}

	items, err := h.passengers.List(c.Context(), userID)
	if err != nil {
		return respondError(c, err, MapError)
	}

	return c.JSON(dto.DataEnvelope[[]dto.PassengerResponse]{
		Data: dto.ToPassengerResponses(items),
	})
}

func (h *Handler) CreateMyPassenger(c *fiber.Ctx) error {
	userID, ok := appmiddleware.UserIDFromContext(c)
	if !ok {
		return writeAppError(c, &AppError{Status: 401, Code: "UNAUTHORIZED", Message: "Нужно войти в аккаунт"})
	}

	input, err := parsePassengerRequest(c)
	if err != nil {
		return writeAppError(c, err)
	}

	created, err := h.passengers.Create(c.Context(), userID, input)
	if err != nil {
		return respondError(c, err, MapError)
	}

	return c.Status(fiber.StatusCreated).JSON(dto.DataEnvelope[dto.PassengerResponse]{
		Data: dto.ToPassengerResponse(created),
	})
}

func (h *Handler) UpdateMyPassenger(c *fiber.Ctx) error {
	userID, ok := appmiddleware.UserIDFromContext(c)
	if !ok {
		return writeAppError(c, &AppError{Status: 401, Code: "UNAUTHORIZED", Message: "Нужно войти в аккаунт"})
	}
	id, err := parseUUID(c.Params("id"))
	if err != nil {
		return writeAppError(c, err)
	}

	input, err := parsePassengerRequest(c)
	if err != nil {
		return writeAppError(c, err)
	}

	updated, err := h.passengers.Update(c.Context(), userID, id, input)
	if err != nil {
		return respondError(c, err, MapError)
	}

	return c.JSON(dto.DataEnvelope[dto.PassengerResponse]{
		Data: dto.ToPassengerResponse(updated),
	})
}

func (h *Handler) DeleteMyPassenger(c *fiber.Ctx) error {
	userID, ok := appmiddleware.UserIDFromContext(c)
	if !ok {
		return writeAppError(c, &AppError{Status: 401, Code: "UNAUTHORIZED", Message: "Нужно войти в аккаунт"})
	}
	id, err := parseUUID(c.Params("id"))
	if err != nil {
		return writeAppError(c, err)
	}

	if err := h.passengers.Delete(c.Context(), userID, id); err != nil {
		return respondError(c, err, MapError)
	}

	return c.JSON(dto.DataEnvelope[map[string]bool]{Data: map[string]bool{"removed": true}})
}

func parsePassengerRequest(c *fiber.Ctx) (application.PassengerInput, error) {
	var req dto.PassengerRequest
	if err := c.BodyParser(&req); err != nil {
		return application.PassengerInput{}, &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Некорректные данные запроса"}
	}
	if err := rejectHoneypot(req.Website); err != nil {
		return application.PassengerInput{}, err
	}
	birthDate, err := parseRequiredDate(req.BirthDate)
	if err != nil {
		return application.PassengerInput{}, err
	}
	return application.PassengerInput{
		Name:      req.Name,
		Phone:     req.Phone,
		BirthDate: birthDate,
		Passport:  req.Passport,
	}, nil
}
