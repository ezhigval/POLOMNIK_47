package handlers

import (
	"github.com/gofiber/fiber/v2"

	"palomnik/internal/adapters/http/fiber/dto"
	appmiddleware "palomnik/internal/adapters/http/fiber/middleware"
)

const visitorIDCookie = "palomnik_visitor"

func visitorIDFromRequest(c *fiber.Ctx) string {
	return c.Cookies(visitorIDCookie)
}

func (h *Handler) GetNewsLikeState(c *fiber.Ctx) error {
	if h.newsEngagement == nil {
		return writeAppError(c, &AppError{Status: 503, Code: "UNAVAILABLE", Message: "Лайки недоступны"})
	}
	state, err := h.newsEngagement.GetLikeState(c.Context(), c.Params("slug"), visitorIDFromRequest(c))
	if err != nil {
		return respondError(c, err, MapError)
	}
	return c.JSON(dto.DataEnvelope[dto.NewsLikeStateResponse]{
		Data: dto.ToNewsLikeStateResponse(state.Count, state.LikedByYou),
	})
}

func (h *Handler) ToggleNewsLike(c *fiber.Ctx) error {
	if h.newsEngagement == nil {
		return writeAppError(c, &AppError{Status: 503, Code: "UNAVAILABLE", Message: "Лайки недоступны"})
	}
	visitorID := visitorIDFromRequest(c)
	if visitorID == "" {
		return writeAppError(c, &AppError{Status: 400, Code: "VISITOR_REQUIRED", Message: "Нужен идентификатор посетителя"})
	}
	state, err := h.newsEngagement.ToggleLike(c.Context(), c.Params("slug"), visitorID)
	if err != nil {
		return respondError(c, err, MapError)
	}
	return c.JSON(dto.DataEnvelope[dto.NewsLikeStateResponse]{
		Data: dto.ToNewsLikeStateResponse(state.Count, state.LikedByYou),
	})
}

func (h *Handler) ListNewsComments(c *fiber.Ctx) error {
	if h.newsEngagement == nil {
		return writeAppError(c, &AppError{Status: 503, Code: "UNAVAILABLE", Message: "Комментарии недоступны"})
	}
	pagination := parsePagination(c)
	list, err := h.newsEngagement.ListComments(c.Context(), c.Params("slug"), pagination)
	if err != nil {
		return respondError(c, err, MapError)
	}
	items := make([]dto.NewsCommentResponse, 0, len(list.Items))
	for _, comment := range list.Items {
		items = append(items, dto.ToNewsCommentResponse(comment))
	}
	return c.JSON(dto.ListEnvelope[dto.NewsCommentResponse]{
		Data: items,
		Meta: list.Meta,
	})
}

func (h *Handler) AddNewsComment(c *fiber.Ctx) error {
	if h.newsEngagement == nil {
		return writeAppError(c, &AppError{Status: 503, Code: "UNAVAILABLE", Message: "Комментарии недоступны"})
	}
	userID, ok := appmiddleware.UserIDFromContext(c)
	if !ok {
		return writeAppError(c, &AppError{Status: 401, Code: "UNAUTHORIZED", Message: "Войдите, чтобы комментировать"})
	}
	var req dto.AddNewsCommentRequest
	if err := c.BodyParser(&req); err != nil {
		return writeAppError(c, &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Некорректные данные запроса"})
	}
	comment, err := h.newsEngagement.AddComment(c.Context(), c.Params("slug"), userID, req.Body)
	if err != nil {
		return respondError(c, err, MapError)
	}
	return c.Status(fiber.StatusCreated).JSON(dto.DataEnvelope[dto.NewsCommentResponse]{
		Data: dto.ToNewsCommentResponse(comment),
	})
}
