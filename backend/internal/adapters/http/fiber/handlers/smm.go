package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"palomnik/internal/adapters/http/fiber/dto"
	"palomnik/internal/application"
)

func (h *Handler) ManagementListSMM(c *fiber.Ctx) error {
	if h.smm == nil {
		return writeAppError(c, &AppError{Status: 503, Code: "SERVICE_UNAVAILABLE", Message: "SMM не настроен"})
	}
	list, err := h.smm.ListPosts(c.Context(), parsePagination(c))
	if err != nil {
		return respondError(c, err, MapError)
	}
	items := make([]dto.SMMPostResponse, 0, len(list.Items))
	for _, post := range list.Items {
		items = append(items, dto.ToSMMPostResponse(post))
	}
	return c.JSON(dto.ListEnvelope[dto.SMMPostResponse]{Data: items, Meta: list.Meta})
}

func (h *Handler) ManagementCreateSMM(c *fiber.Ctx) error {
	if h.smm == nil {
		return writeAppError(c, &AppError{Status: 503, Code: "SERVICE_UNAVAILABLE", Message: "SMM не настроен"})
	}
	var req dto.SMMPostCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return writeAppError(c, &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Некорректные данные запроса"})
	}
	publishAt, err := time.Parse(time.RFC3339, req.PublishAt)
	if err != nil {
		return writeAppError(c, &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Укажите дату и время слота"})
	}
	post, err := h.smm.CreatePost(c.Context(), application.CreateSMMPostInput{
		Title:     req.Title,
		Body:      req.Body,
		URL:       req.URL,
		PublishAt: publishAt,
		Channels:  req.Channels,
	})
	if err != nil {
		return respondError(c, err, MapError)
	}
	return c.Status(fiber.StatusCreated).JSON(dto.DataEnvelope[dto.SMMPostResponse]{Data: dto.ToSMMPostResponse(post)})
}

func (h *Handler) ManagementGetSMM(c *fiber.Ctx) error {
	if h.smm == nil {
		return writeAppError(c, &AppError{Status: 503, Code: "SERVICE_UNAVAILABLE", Message: "SMM не настроен"})
	}
	id, err := parseUUID(c.Params("id"))
	if err != nil {
		return writeAppError(c, err)
	}
	post, err := h.smm.GetPost(c.Context(), id)
	if err != nil {
		return respondError(c, err, MapError)
	}
	return c.JSON(dto.DataEnvelope[dto.SMMPostResponse]{Data: dto.ToSMMPostResponse(post)})
}

func (h *Handler) ManagementPublishSMM(c *fiber.Ctx) error {
	if h.smm == nil {
		return writeAppError(c, &AppError{Status: 503, Code: "SERVICE_UNAVAILABLE", Message: "SMM не настроен"})
	}
	id, err := parseUUID(c.Params("id"))
	if err != nil {
		return writeAppError(c, err)
	}
	post, err := h.smm.PublishPost(c.Context(), id)
	if err != nil {
		return respondError(c, err, MapError)
	}
	return c.JSON(dto.DataEnvelope[dto.SMMPostResponse]{Data: dto.ToSMMPostResponse(post)})
}

func (h *Handler) ManagementDeleteSMM(c *fiber.Ctx) error {
	if h.smm == nil {
		return writeAppError(c, &AppError{Status: 503, Code: "SERVICE_UNAVAILABLE", Message: "SMM не настроен"})
	}
	id, err := parseUUID(c.Params("id"))
	if err != nil {
		return writeAppError(c, err)
	}
	if err := h.smm.DeletePost(c.Context(), id); err != nil {
		return respondError(c, err, MapError)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) GetPublicNewsBySlug(c *fiber.Ctx) error {
	article, err := h.news.GetPublishedNewsBySlug(c.Context(), c.Params("slug"))
	if err != nil {
		return respondError(c, err, MapError)
	}
	return c.JSON(dto.DataEnvelope[dto.NewsArticleResponse]{Data: dto.ToNewsArticleResponse(article)})
}
