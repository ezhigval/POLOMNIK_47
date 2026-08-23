package handlers

import (
	"github.com/gofiber/fiber/v2"

	"palomnik/internal/adapters/http/fiber/dto"
	appmiddleware "palomnik/internal/adapters/http/fiber/middleware"
	"palomnik/internal/application"
	"palomnik/internal/domain"
)

func (h *Handler) ListMyPhotos(c *fiber.Ctx) error {
	userID, ok := appmiddleware.UserIDFromContext(c)
	if !ok {
		return writeAppError(c, &AppError{Status: 401, Code: "UNAUTHORIZED", Message: "Нужно войти в аккаунт"})
	}
	items, err := h.photos.List(c.Context(), userID)
	if err != nil {
		return respondError(c, err, MapError)
	}
	return c.JSON(dto.DataEnvelope[[]dto.UserPhotoResponse]{
		Data: dto.ToUserPhotoResponses(items),
	})
}

func (h *Handler) CreateMyPhoto(c *fiber.Ctx) error {
	userID, ok := appmiddleware.UserIDFromContext(c)
	if !ok {
		return writeAppError(c, &AppError{Status: 401, Code: "UNAUTHORIZED", Message: "Нужно войти в аккаунт"})
	}

	var req dto.CreateUserPhotoRequest
	if err := c.BodyParser(&req); err != nil {
		return writeAppError(c, &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Некорректные данные запроса"})
	}
	if err := rejectHoneypot(req.Website); err != nil {
		return writeAppError(c, err)
	}
	if !req.ConsentPersonalData {
		return writeAppError(c, &AppError{
			Status:  422,
			Code:    "CONSENT_REQUIRED",
			Message: "Необходимо согласие на обработку персональных данных",
		})
	}

	created, err := h.photos.Create(c.Context(), userID, application.UserPhotoInput{
		URL:               req.URL,
		Caption:           req.Caption,
		AllowDistribution: req.AllowDistribution,
	})
	if err != nil {
		return respondError(c, err, MapError)
	}

	photoID := created.ID
	if _, err := h.consents.RecordConsent(c.Context(), application.RecordConsentInput{
		ConsentType: domain.ConsentTypePersonalData,
		UserID:      &userID,
		RequestID:   &photoID,
		IP:          c.IP(),
		UserAgent:   c.Get("User-Agent"),
	}); err != nil {
		return respondError(c, err, MapError)
	}
	if req.AllowDistribution {
		if _, err := h.consents.RecordConsent(c.Context(), application.RecordConsentInput{
			ConsentType: domain.ConsentTypeDistribution,
			UserID:      &userID,
			RequestID:   &photoID,
			IP:          c.IP(),
			UserAgent:   c.Get("User-Agent"),
		}); err != nil {
			return respondError(c, err, MapError)
		}
	}

	return c.Status(fiber.StatusCreated).JSON(dto.DataEnvelope[dto.UserPhotoResponse]{
		Data: dto.ToUserPhotoResponse(created),
	})
}

func (h *Handler) DeleteMyPhoto(c *fiber.Ctx) error {
	userID, ok := appmiddleware.UserIDFromContext(c)
	if !ok {
		return writeAppError(c, &AppError{Status: 401, Code: "UNAUTHORIZED", Message: "Нужно войти в аккаунт"})
	}
	id, err := parseUUID(c.Params("id"))
	if err != nil {
		return writeAppError(c, err)
	}
	if err := h.photos.Delete(c.Context(), userID, id); err != nil {
		return respondError(c, err, MapError)
	}
	return c.JSON(dto.DataEnvelope[map[string]bool]{Data: map[string]bool{"removed": true}})
}
