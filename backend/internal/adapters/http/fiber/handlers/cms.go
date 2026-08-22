package handlers

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"polomnik/internal/adapters/http/fiber/dto"
	"polomnik/internal/application"
	"polomnik/internal/ports"
)

func (h *Handler) ListPublicCMSPages(c *fiber.Ctx) error {
	pages, err := h.cms.ListPages(c.Context(), true)
	if err != nil {
		return respondError(c, err, MapError)
	}
	items := make([]dto.CMSPageResponse, 0, len(pages))
	for _, page := range pages {
		items = append(items, dto.ToCMSPageResponse(page))
	}
	return c.JSON(dto.ListEnvelope[dto.CMSPageResponse]{
		Data: items,
		Meta: ports.PageMeta{Page: 1, Limit: len(items), Total: len(items), HasNext: false},
	})
}

func (h *Handler) GetPublicCMSPage(c *fiber.Ctx) error {
	slug := c.Params("slug")
	page, err := h.cms.GetPublishedPageBySlug(c.Context(), slug)
	if err != nil {
		return respondError(c, err, MapError)
	}
	return c.JSON(dto.DataEnvelope[dto.CMSPageResponse]{Data: dto.ToCMSPageResponse(page)})
}

func (h *Handler) ManagementListCMSPages(c *fiber.Ctx) error {
	pages, err := h.cms.ListPages(c.Context(), false)
	if err != nil {
		return respondError(c, err, MapError)
	}
	items := make([]dto.CMSPageResponse, 0, len(pages))
	for _, page := range pages {
		items = append(items, dto.ToCMSPageResponse(page))
	}
	return c.JSON(dto.ListEnvelope[dto.CMSPageResponse]{
		Data: items,
		Meta: ports.PageMeta{Page: 1, Limit: len(items), Total: len(items), HasNext: false},
	})
}

func (h *Handler) ManagementGetCMSPage(c *fiber.Ctx) error {
	id, err := parseUUID(c.Params("id"))
	if err != nil {
		return writeAppError(c, err)
	}
	page, err := h.cms.GetPage(c.Context(), id)
	if err != nil {
		return respondError(c, err, MapError)
	}
	return c.JSON(dto.DataEnvelope[dto.CMSPageResponse]{Data: dto.ToCMSPageResponse(page)})
}

func (h *Handler) ManagementBootstrapHomeCMSPage(c *fiber.Ctx) error {
	page, err := h.cms.BootstrapHomePage(c.Context())
	if err != nil {
		return respondError(c, err, MapError)
	}
	return c.Status(fiber.StatusCreated).JSON(dto.DataEnvelope[dto.CMSPageResponse]{Data: dto.ToCMSPageResponse(page)})
}

func (h *Handler) ManagementCreateCMSPage(c *fiber.Ctx) error {
	var req dto.CMSPageCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return writeAppError(c, &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Некорректные данные запроса"})
	}
	published := true
	if req.IsPublished != nil {
		published = *req.IsPublished
	}
	page, err := h.cms.CreatePage(c.Context(), application.CreatePageInput{
		Slug:        req.Slug,
		Title:       req.Title,
		Path:        req.Path,
		IsPublished: published,
	})
	if err != nil {
		return respondError(c, err, MapError)
	}
	return c.Status(fiber.StatusCreated).JSON(dto.DataEnvelope[dto.CMSPageResponse]{Data: dto.ToCMSPageResponse(page)})
}

func (h *Handler) ManagementUpdateCMSPage(c *fiber.Ctx) error {
	id, err := parseUUID(c.Params("id"))
	if err != nil {
		return writeAppError(c, err)
	}
	var req dto.CMSPageUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return writeAppError(c, &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Некорректные данные запроса"})
	}
	page, err := h.cms.UpdatePage(c.Context(), id, application.UpdatePageInput{
		Title:       req.Title,
		Path:        req.Path,
		IsPublished: req.IsPublished,
	})
	if err != nil {
		return respondError(c, err, MapError)
	}
	return c.JSON(dto.DataEnvelope[dto.CMSPageResponse]{Data: dto.ToCMSPageResponse(page)})
}

func (h *Handler) ManagementDeleteCMSPage(c *fiber.Ctx) error {
	id, err := parseUUID(c.Params("id"))
	if err != nil {
		return writeAppError(c, err)
	}
	if err := h.cms.DeletePage(c.Context(), id); err != nil {
		return respondError(c, err, MapError)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) ManagementListCMSTemplates(c *fiber.Ctx) error {
	templates := application.BlockTemplates()
	items := make([]dto.CMSBlockTemplateResponse, 0, len(templates))
	for _, item := range templates {
		content, _ := item["content"].(json.RawMessage)
		items = append(items, dto.CMSBlockTemplateResponse{
			Type:    item["type"].(string),
			Label:   item["label"].(string),
			Content: content,
		})
	}
	return c.JSON(dto.DataEnvelope[[]dto.CMSBlockTemplateResponse]{Data: items})
}

func (h *Handler) ManagementCreateCMSBlock(c *fiber.Ctx) error {
	pageID, err := parseUUID(c.Params("id"))
	if err != nil {
		return writeAppError(c, err)
	}
	var req dto.CMSBlockCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return writeAppError(c, &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Некорректные данные запроса"})
	}
	visible := true
	if req.IsVisible != nil {
		visible = *req.IsVisible
	}
	block, err := h.cms.CreateBlock(c.Context(), pageID, application.CreateBlockInput{
		Type:      req.Type,
		Content:   req.Content,
		IsVisible: visible,
	})
	if err != nil {
		return respondError(c, err, MapError)
	}
	return c.Status(fiber.StatusCreated).JSON(dto.DataEnvelope[dto.CMSBlockResponse]{Data: dto.ToCMSBlockResponse(block)})
}

func (h *Handler) ManagementUpdateCMSBlock(c *fiber.Ctx) error {
	id, err := parseUUID(c.Params("id"))
	if err != nil {
		return writeAppError(c, err)
	}
	var req dto.CMSBlockUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return writeAppError(c, &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Некорректные данные запроса"})
	}
	block, err := h.cms.UpdateBlock(c.Context(), id, application.UpdateBlockInput{
		Content:   req.Content,
		IsVisible: req.IsVisible,
		SortOrder: req.SortOrder,
	})
	if err != nil {
		return respondError(c, err, MapError)
	}
	return c.JSON(dto.DataEnvelope[dto.CMSBlockResponse]{Data: dto.ToCMSBlockResponse(block)})
}

func (h *Handler) ManagementDeleteCMSBlock(c *fiber.Ctx) error {
	id, err := parseUUID(c.Params("id"))
	if err != nil {
		return writeAppError(c, err)
	}
	if err := h.cms.DeleteBlock(c.Context(), id); err != nil {
		return respondError(c, err, MapError)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) ManagementReorderCMSBlocks(c *fiber.Ctx) error {
	pageID, err := parseUUID(c.Params("id"))
	if err != nil {
		return writeAppError(c, err)
	}
	var req dto.CMSReorderRequest
	if err := c.BodyParser(&req); err != nil {
		return writeAppError(c, &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Некорректные данные запроса"})
	}
	ids := make([]uuid.UUID, 0, len(req.BlockIDs))
	for _, value := range req.BlockIDs {
		id, parseErr := uuid.Parse(value)
		if parseErr != nil {
			return writeAppError(c, &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Некорректный идентификатор блока"})
		}
		ids = append(ids, id)
	}
	if err := h.cms.ReorderBlocks(c.Context(), pageID, ids); err != nil {
		return respondError(c, err, MapError)
	}
	page, err := h.cms.GetPage(c.Context(), pageID)
	if err != nil {
		return respondError(c, err, MapError)
	}
	return c.JSON(dto.DataEnvelope[dto.CMSPageResponse]{Data: dto.ToCMSPageResponse(page)})
}
