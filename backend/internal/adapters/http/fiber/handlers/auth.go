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
			Message: "Некорректные данные запроса",
		})
	}

	result, err := h.auth.Register(c.Context(), application.RegisterInput{
		Email:        req.Email,
		Phone:        req.Phone,
		Name:         req.Name,
		Password:     req.Password,
		PhoneCheckID: req.PhoneCheckID,
	})
	if err != nil {
		return respondError(c, err, MapError)
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
			Message: "Некорректные данные запроса",
		})
	}

	result, err := h.auth.Login(c.Context(), application.LoginInput{
		Login:    req.Login,
		Password: req.Password,
	})
	if err != nil {
		return respondError(c, err, MapError)
	}

	return c.JSON(dto.DataEnvelope[dto.AuthResponse]{
		Data: dto.AuthResponse{
			Token: result.Token,
			User:  dto.ToUserResponse(result.User),
		},
	})
}

func (h *Handler) AuthMethods(c *fiber.Ctx) error {
	methods := h.auth.AuthMethods()
	status := func(s application.AuthMethodStatus) dto.AuthMethodStatusResponse {
		return dto.AuthMethodStatusResponse{
			Available: s.Available,
			Message:   s.Message,
			Username:  s.Username,
		}
	}
	return c.JSON(dto.DataEnvelope[dto.AuthMethodsResponse]{
		Data: dto.AuthMethodsResponse{
			Password:  methods.Password,
			PhoneCall: status(methods.PhoneCall),
			Yandex:    status(methods.Yandex),
			VK:        status(methods.VK),
			Max:       status(methods.Max),
			Telegram:  status(methods.Telegram),
			Mail:      status(methods.Mail),
		},
	})
}


func (h *Handler) StartPhoneVerification(c *fiber.Ctx) error {
	var req dto.PhoneStartRequest
	if err := c.BodyParser(&req); err != nil {
		return writeAppError(c, &AppError{
			Status:  422,
			Code:    "VALIDATION_ERROR",
			Message: "Некорректные данные запроса",
		})
	}

	result, err := h.auth.StartPhoneVerification(c.Context(), req.Phone)
	if err != nil {
		return respondError(c, err, MapError)
	}

	return c.JSON(dto.DataEnvelope[dto.PhoneStartResponse]{
		Data: dto.PhoneStartResponse{
			CheckID:         result.CheckID,
			CallPhone:       result.CallPhone,
			CallPhonePretty: result.CallPhonePretty,
			ExpiresIn:       result.ExpiresIn,
		},
	})
}

func (h *Handler) PhoneVerificationStatus(c *fiber.Ctx) error {
	checkID := c.Query("check_id")
	if checkID == "" {
		var req dto.PhoneCheckRequest
		_ = c.BodyParser(&req)
		checkID = req.CheckID
	}

	result, err := h.auth.PhoneVerificationStatus(c.Context(), checkID)
	if err != nil {
		return respondError(c, err, MapError)
	}

	return c.JSON(dto.DataEnvelope[dto.PhoneStatusResponse]{
		Data: dto.PhoneStatusResponse{Status: result.Status},
	})
}

func (h *Handler) CompletePhoneLogin(c *fiber.Ctx) error {
	var req dto.PhoneCheckRequest
	if err := c.BodyParser(&req); err != nil {
		return writeAppError(c, &AppError{
			Status:  422,
			Code:    "VALIDATION_ERROR",
			Message: "Некорректные данные запроса",
		})
	}

	result, err := h.auth.CompletePhoneLogin(c.Context(), req.CheckID)
	if err != nil {
		return respondError(c, err, MapError)
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
			Message: "Нужно войти в аккаунт",
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
			Message: "Нужно войти в аккаунт",
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
