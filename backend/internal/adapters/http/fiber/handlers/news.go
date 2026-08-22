package handlers

import (
	"github.com/gofiber/fiber/v2"

	"polomnik/internal/adapters/http/fiber/dto"
	"polomnik/internal/application"
)

func (h *Handler) ListPublicNews(c *fiber.Ctx) error {
	list, err := h.news.ListPublicNews(c.Context(), parsePagination(c))
	if err != nil {
		return respondError(c, err, MapError)
	}

	items := make([]dto.NewsArticleResponse, 0, len(list.Items))
	for _, article := range list.Items {
		items = append(items, dto.ToNewsArticleResponse(article))
	}

	return c.JSON(dto.ListEnvelope[dto.NewsArticleResponse]{
		Data: items,
		Meta: list.Meta,
	})
}

func (h *Handler) ManagementListNews(c *fiber.Ctx) error {
	list, err := h.news.ListNews(c.Context(), parsePagination(c))
	if err != nil {
		return respondError(c, err, MapError)
	}

	items := make([]dto.NewsArticleResponse, 0, len(list.Items))
	for _, article := range list.Items {
		items = append(items, dto.ToNewsArticleResponse(article))
	}

	return c.JSON(dto.ListEnvelope[dto.NewsArticleResponse]{
		Data: items,
		Meta: list.Meta,
	})
}

func (h *Handler) ManagementGetNews(c *fiber.Ctx) error {
	id, err := parseUUID(c.Params("id"))
	if err != nil {
		return writeAppError(c, err)
	}

	article, err := h.news.GetNews(c.Context(), id)
	if err != nil {
		return respondError(c, err, MapError)
	}

	return c.JSON(dto.DataEnvelope[dto.NewsArticleResponse]{
		Data: dto.ToNewsArticleResponse(article),
	})
}

func (h *Handler) ManagementCreateNews(c *fiber.Ctx) error {
	input, appErr := parseNewsArticleInput(c)
	if appErr != nil {
		return writeAppError(c, appErr)
	}

	article, err := h.news.CreateNews(c.Context(), input)
	if err != nil {
		return respondError(c, err, MapError)
	}

	return c.Status(fiber.StatusCreated).JSON(dto.DataEnvelope[dto.NewsArticleResponse]{
		Data: dto.ToNewsArticleResponse(article),
	})
}

func (h *Handler) ManagementUpdateNews(c *fiber.Ctx) error {
	id, err := parseUUID(c.Params("id"))
	if err != nil {
		return writeAppError(c, err)
	}

	input, appErr := parseNewsArticleInput(c)
	if appErr != nil {
		return writeAppError(c, appErr)
	}

	article, err := h.news.UpdateNews(c.Context(), id, input)
	if err != nil {
		return respondError(c, err, MapError)
	}

	return c.JSON(dto.DataEnvelope[dto.NewsArticleResponse]{
		Data: dto.ToNewsArticleResponse(article),
	})
}

func (h *Handler) ManagementDeleteNews(c *fiber.Ctx) error {
	id, err := parseUUID(c.Params("id"))
	if err != nil {
		return writeAppError(c, err)
	}

	if err := h.news.DeleteNews(c.Context(), id); err != nil {
		return respondError(c, err, MapError)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func parseNewsArticleInput(c *fiber.Ctx) (application.NewsArticleInput, *AppError) {
	var req dto.NewsArticleUpsertRequest
	if err := c.BodyParser(&req); err != nil {
		return application.NewsArticleInput{}, &AppError{
			Status:  422,
			Code:    "VALIDATION_ERROR",
			Message: "Некорректные данные запроса",
		}
	}

	publishedAt, err := parseRequiredDate(req.PublishedAt)
	if err != nil {
		return application.NewsArticleInput{}, asAppError(err)
	}

	return application.NewsArticleInput{
		Slug:        req.Slug,
		Title:       req.Title,
		Excerpt:     req.Excerpt,
		Body:        req.Body,
		ImageURL:    req.ImageURL,
		PublishedAt: publishedAt,
		IsPublished: req.IsPublished,
		SortOrder:   req.SortOrder,
	}, nil
}
