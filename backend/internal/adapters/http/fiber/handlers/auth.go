package handlers

import (
	"github.com/gofiber/fiber/v2"

	"polomnik/internal/adapters/http/fiber/dto"
	appmiddleware "polomnik/internal/adapters/http/fiber/middleware"
	"polomnik/internal/application"
)

func (h *Handler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return writeAppError(c, &AppError{
			Status:  422,
			Code:    "VALIDATION_ERROR",
			Message: "Invalid request body",
		})
	}

	result, err := h.auth.Register(c.Context(), application.RegisterInput{
		Email:    req.Email,
		Phone:    req.Phone,
		Name:     req.Name,
		Password: req.Password,
	})
	if err != nil {
		return respondError(c, err, MapAuthError)
	}

	return c.Status(fiber.StatusCreated).JSON(dto.DataEnvelope[dto.AuthResponse]{
		Data: dto.AuthResponse{
			Token: result.Token,
			User:  dto.ToUserResponse(result.User),
		},
	})
}

func (h *Handler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return writeAppError(c, &AppError{
			Status:  422,
			Code:    "VALIDATION_ERROR",
			Message: "Invalid request body",
		})
	}

	result, err := h.auth.Login(c.Context(), application.LoginInput{
		Login:    req.Login,
		Password: req.Password,
	})
	if err != nil {
		return respondError(c, err, MapAuthError)
	}

	return c.JSON(dto.DataEnvelope[dto.AuthResponse]{
		Data: dto.AuthResponse{
			Token: result.Token,
			User:  dto.ToUserResponse(result.User),
		},
	})
}

func (h *Handler) Me(c *fiber.Ctx) error {
	userID, ok := appmiddleware.UserIDFromContext(c)
	if !ok {
		return writeAppError(c, &AppError{
			Status:  401,
			Code:    "UNAUTHORIZED",
			Message: "Authentication required",
		})
	}

	user, err := h.auth.GetUser(c.Context(), userID)
	if err != nil {
		return respondError(c, err, MapError)
	}

	return c.JSON(dto.DataEnvelope[dto.UserResponse]{
		Data: dto.ToUserResponse(user),
	})
}

func (h *Handler) MyBookings(c *fiber.Ctx) error {
	userID, ok := appmiddleware.UserIDFromContext(c)
	if !ok {
		return writeAppError(c, &AppError{
			Status:  401,
			Code:    "UNAUTHORIZED",
			Message: "Authentication required",
		})
	}

	list, err := h.auth.ListMyBookings(c.Context(), userID, parsePagination(c))
	if err != nil {
		return respondError(c, err, MapError)
	}

	return c.JSON(dto.ListEnvelope[dto.MyBookingResponse]{
		Data: dto.ToMyBookingResponses(list.Items),
		Meta: list.Meta,
	})
}

func MapAuthError(err error) *AppError {
	return MapError(err)
}
