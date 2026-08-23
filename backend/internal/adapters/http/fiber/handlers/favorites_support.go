package handlers

import (
	"github.com/gofiber/fiber/v2"

	"palomnik/internal/adapters/http/fiber/dto"
	appmiddleware "palomnik/internal/adapters/http/fiber/middleware"
	"palomnik/internal/application"
	"palomnik/internal/ports"
)

func (h *Handler) OAuthLogin(c *fiber.Ctx) error {
	var req dto.OAuthLoginRequest
	if err := c.BodyParser(&req); err != nil {
		return writeAppError(c, &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Некорректные данные запроса"})
	}

	result, err := h.auth.OAuthLogin(c.Context(), application.OAuthLoginInput{
		Provider: req.Provider,
		Subject:  req.Subject,
		Email:    req.Email,
		Name:     req.Name,
	})
	if err != nil {
		return respondError(c, err, MapError)
	}

	return c.JSON(dto.DataEnvelope[dto.AuthResponse]{
		Data: dto.AuthResponse{Token: result.Token, User: dto.ToUserResponse(result.User)},
	})
}

func (h *Handler) ListFavoriteIDs(c *fiber.Ctx) error {
	userID, ok := appmiddleware.UserIDFromContext(c)
	if !ok {
		return writeAppError(c, &AppError{Status: 401, Code: "UNAUTHORIZED", Message: "Нужно войти в аккаунт"})
	}

	ids, err := h.favorites.ListFavoriteTourIDs(c.Context(), userID)
	if err != nil {
		return respondError(c, err, MapError)
	}

	response := make([]string, 0, len(ids))
	for _, id := range ids {
		response = append(response, id.String())
	}
	return c.JSON(dto.DataEnvelope[[]string]{Data: response})
}

func (h *Handler) AddFavorite(c *fiber.Ctx) error {
	userID, ok := appmiddleware.UserIDFromContext(c)
	if !ok {
		return writeAppError(c, &AppError{Status: 401, Code: "UNAUTHORIZED", Message: "Нужно войти в аккаунт"})
	}
	tourID, err := parseUUID(c.Params("tourId"))
	if err != nil {
		return writeAppError(c, err)
	}

	if err := h.favorites.AddFavorite(c.Context(), userID, tourID); err != nil {
		return respondError(c, err, MapError)
	}
	return c.Status(fiber.StatusCreated).JSON(dto.DataEnvelope[map[string]bool]{Data: map[string]bool{"saved": true}})
}

func (h *Handler) RemoveFavorite(c *fiber.Ctx) error {
	userID, ok := appmiddleware.UserIDFromContext(c)
	if !ok {
		return writeAppError(c, &AppError{Status: 401, Code: "UNAUTHORIZED", Message: "Нужно войти в аккаунт"})
	}
	tourID, err := parseUUID(c.Params("tourId"))
	if err != nil {
		return writeAppError(c, err)
	}

	if err := h.favorites.RemoveFavorite(c.Context(), userID, tourID); err != nil {
		return respondError(c, err, MapError)
	}
	return c.JSON(dto.DataEnvelope[map[string]bool]{Data: map[string]bool{"removed": true}})
}

func (h *Handler) GetSupportThread(c *fiber.Ctx) error {
	userID, ok := appmiddleware.UserIDFromContext(c)
	if !ok {
		return writeAppError(c, &AppError{Status: 401, Code: "UNAUTHORIZED", Message: "Нужно войти в аккаунт"})
	}

	thread, messages, err := h.support.GetThread(c.Context(), userID)
	if err != nil {
		return respondError(c, err, MapError)
	}

	return c.JSON(dto.DataEnvelope[dto.SupportThreadResponse]{
		Data: dto.ToSupportThreadResponse(thread, messages),
	})
}

func (h *Handler) SendSupportMessage(c *fiber.Ctx) error {
	userID, ok := appmiddleware.UserIDFromContext(c)
	if !ok {
		return writeAppError(c, &AppError{Status: 401, Code: "UNAUTHORIZED", Message: "Нужно войти в аккаунт"})
	}

	var req dto.SendSupportMessageRequest
	if err := c.BodyParser(&req); err != nil {
		return writeAppError(c, &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Некорректные данные запроса"})
	}

	messages, err := h.support.SendUserMessage(c.Context(), userID, req.Body)
	if err != nil {
		return respondError(c, err, MapError)
	}

	thread, err := h.support.GetOrCreateThread(c.Context(), userID)
	if err != nil {
		return respondError(c, err, MapError)
	}

	return c.Status(fiber.StatusCreated).JSON(dto.DataEnvelope[dto.SupportThreadResponse]{
		Data: dto.ToSupportThreadResponse(thread, messages),
	})
}

func (h *Handler) ManagementListSupportThreads(c *fiber.Ctx) error {
	threads, err := h.support.ListThreads(c.Context())
	if err != nil {
		return respondError(c, err, MapError)
	}

	items := make([]dto.SupportThreadResponse, 0, len(threads))
	for _, thread := range threads {
		items = append(items, dto.ToManagementSupportThreadSummary(thread))
	}

	return c.JSON(dto.ListEnvelope[dto.SupportThreadResponse]{
		Data: items,
		Meta: ports.PageMeta{Page: 1, Limit: len(items), Total: len(items), HasNext: false},
	})
}

func (h *Handler) ManagementGetSupportThread(c *fiber.Ctx) error {
	id, err := parseUUID(c.Params("id"))
	if err != nil {
		return writeAppError(c, err)
	}

	thread, messages, err := h.support.GetThreadByID(c.Context(), id)
	if err != nil {
		return respondError(c, err, MapError)
	}

	return c.JSON(dto.DataEnvelope[dto.SupportThreadResponse]{
		Data: dto.ToManagementSupportThreadResponse(thread, messages),
	})
}

func (h *Handler) ManagementSendSupportMessage(c *fiber.Ctx) error {
	id, err := parseUUID(c.Params("id"))
	if err != nil {
		return writeAppError(c, err)
	}

	var req dto.SendSupportMessageRequest
	if err := c.BodyParser(&req); err != nil {
		return writeAppError(c, &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Некорректные данные запроса"})
	}

	thread, messages, err := h.support.SendStaffMessage(c.Context(), id, req.Body)
	if err != nil {
		return respondError(c, err, MapError)
	}

	return c.Status(fiber.StatusCreated).JSON(dto.DataEnvelope[dto.SupportThreadResponse]{
		Data: dto.ToManagementSupportThreadResponse(thread, messages),
	})
}
