package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"palomnik/internal/adapters/http/fiber/dto"
	appmiddleware "palomnik/internal/adapters/http/fiber/middleware"
	"palomnik/internal/application"
	"palomnik/internal/domain"
	"palomnik/internal/ports"
)

func (h *Handler) ListPublicLegalDocuments(c *fiber.Ctx) error {
	docs, err := h.legal.ListActiveDocuments(c.Context())
	if err != nil {
		return respondError(c, err, MapError)
	}
	items := make([]dto.LegalDocumentSummaryResponse, 0, len(docs))
	for _, doc := range docs {
		items = append(items, dto.ToLegalDocumentSummaryResponse(doc))
	}
	return c.JSON(dto.DataEnvelope[[]dto.LegalDocumentSummaryResponse]{Data: items})
}

func (h *Handler) GetPublicLegalDocument(c *fiber.Ctx) error {
	docType, err := domain.ParseLegalDocumentType(c.Params("type"))
	if err != nil {
		return writeAppError(c, &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Некорректный тип документа"})
	}
	doc, err := h.legal.GetActiveByType(c.Context(), docType)
	if err != nil {
		return respondError(c, err, MapError)
	}
	return c.JSON(dto.DataEnvelope[dto.LegalDocumentResponse]{
		Data: dto.ToLegalDocumentResponse(doc, true),
	})
}

func (h *Handler) GetPublicLegalDocumentVersion(c *fiber.Ctx) error {
	docType, err := domain.ParseLegalDocumentType(c.Params("type"))
	if err != nil {
		return writeAppError(c, &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Некорректный тип документа"})
	}
	version := c.Params("version")
	doc, err := h.legal.GetByTypeAndVersion(c.Context(), docType, version)
	if err != nil {
		return respondError(c, err, MapError)
	}
	return c.JSON(dto.DataEnvelope[dto.LegalDocumentResponse]{
		Data: dto.ToLegalDocumentResponse(doc, true),
	})
}

func (h *Handler) RecordConsent(c *fiber.Ctx) error {
	var req dto.RecordConsentRequest
	if err := c.BodyParser(&req); err != nil {
		return writeAppError(c, &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Некорректные данные запроса"})
	}
	consentType, err := domain.ParseConsentType(req.ConsentType)
	if err != nil {
		return writeAppError(c, &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Некорректный тип согласия"})
	}

	var userID *uuid.UUID
	if id, ok := appmiddleware.UserIDFromContext(c); ok {
		userID = &id
	}

	var requestID *uuid.UUID
	if req.RequestID != nil && *req.RequestID != "" {
		id, err := uuid.Parse(*req.RequestID)
		if err != nil {
			return writeAppError(c, &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Некорректный идентификатор заявки"})
		}
		requestID = &id
	}

	consent, err := h.consents.RecordConsent(c.Context(), application.RecordConsentInput{
		ConsentType: consentType,
		UserID:      userID,
		RequestID:   requestID,
		IP:          c.IP(),
		UserAgent:   c.Get("User-Agent"),
	})
	if err != nil {
		return respondError(c, err, MapError)
	}
	return c.Status(fiber.StatusCreated).JSON(dto.DataEnvelope[dto.ConsentResponse]{
		Data: dto.ToConsentResponse(consent),
	})
}

func (h *Handler) ListMyConsents(c *fiber.Ctx) error {
	userID, ok := appmiddleware.UserIDFromContext(c)
	if !ok {
		return writeAppError(c, &AppError{Status: 401, Code: "UNAUTHORIZED", Message: "Нужно войти в аккаунт"})
	}
	pagination := parsePagination(c)
	list, err := h.consents.ListUserConsents(c.Context(), userID, pagination)
	if err != nil {
		return respondError(c, err, MapError)
	}
	items := make([]dto.ConsentResponse, 0, len(list.Items))
	for _, item := range list.Items {
		items = append(items, dto.ToConsentResponse(item))
	}
	return c.JSON(dto.ListEnvelope[dto.ConsentResponse]{Data: items, Meta: list.Meta})
}

func (h *Handler) ManagementListLegalDocuments(c *fiber.Ctx) error {
	var filters ports.LegalDocumentFilters
	filters.IncludeAll = true
	if docTypeRaw := c.Query("type"); docTypeRaw != "" {
		docType, err := domain.ParseLegalDocumentType(docTypeRaw)
		if err != nil {
			return writeAppError(c, &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Некорректный тип документа"})
		}
		filters.Type = docType
	}
	docs, err := h.legal.ListDocuments(c.Context(), filters)
	if err != nil {
		return respondError(c, err, MapError)
	}
	items := make([]dto.LegalDocumentSummaryResponse, 0, len(docs))
	for _, doc := range docs {
		items = append(items, dto.ToLegalDocumentSummaryResponse(doc))
	}
	return c.JSON(dto.DataEnvelope[[]dto.LegalDocumentSummaryResponse]{Data: items})
}

func (h *Handler) ManagementPublishLegalDocument(c *fiber.Ctx) error {
	var req dto.PublishLegalDocumentRequest
	if err := c.BodyParser(&req); err != nil {
		return writeAppError(c, &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Некорректные данные запроса"})
	}
	docType, err := domain.ParseLegalDocumentType(req.Type)
	if err != nil {
		return writeAppError(c, &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Некорректный тип документа"})
	}
	doc, err := h.legal.PublishNewVersion(c.Context(), application.PublishLegalDocumentInput{
		Type:    docType,
		Version: req.Version,
		Title:   req.Title,
		Content: req.Content,
	})
	if err != nil {
		return respondError(c, err, MapError)
	}
	return c.Status(fiber.StatusCreated).JSON(dto.DataEnvelope[dto.LegalDocumentResponse]{
		Data: dto.ToLegalDocumentResponse(doc, true),
	})
}

func (h *Handler) ManagementListConsents(c *fiber.Ctx) error {
	pagination := parsePagination(c)
	var filters ports.ConsentFilters
	if userIDRaw := c.Query("user_id"); userIDRaw != "" {
		id, err := uuid.Parse(userIDRaw)
		if err != nil {
			return writeAppError(c, &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Некорректный user_id"})
		}
		filters.UserID = &id
	}
	if consentTypeRaw := c.Query("consent_type"); consentTypeRaw != "" {
		ct, err := domain.ParseConsentType(consentTypeRaw)
		if err != nil {
			return writeAppError(c, &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Некорректный тип согласия"})
		}
		filters.ConsentType = ct
	}
	list, err := h.consents.ListManagementConsents(c.Context(), filters, pagination)
	if err != nil {
		return respondError(c, err, MapError)
	}
	items := make([]dto.ConsentResponse, 0, len(list.Items))
	for _, item := range list.Items {
		items = append(items, dto.ToConsentResponse(item))
	}
	return c.JSON(dto.ListEnvelope[dto.ConsentResponse]{Data: items, Meta: list.Meta})
}
