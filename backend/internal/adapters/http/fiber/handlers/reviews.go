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

func (h *Handler) ListReviews(c *fiber.Ctx) error {
	var tourID *uuid.UUID
	if value := c.Query("tour_id"); value != "" {
		parsed, err := parseUUID(value)
		if err != nil {
			return writeAppError(c, err)
		}
		tourID = &parsed
	}

	rating, err := parseOptionalInt(c.Query("rating"))
	if err != nil {
		return writeAppError(c, err)
	}

	list, err := h.reviews.ListPublicReviews(c.Context(), ports.ReviewFilters{
		TourID: tourID,
		Rating: rating,
	}, parsePagination(c))
	if err != nil {
		return respondError(c, err, MapError)
	}

	items := make([]dto.ReviewResponse, 0, len(list.Items))
	for _, review := range list.Items {
		items = append(items, dto.ToReviewResponse(review))
	}

	return c.JSON(dto.ListEnvelope[dto.ReviewResponse]{
		Data: items,
		Meta: list.Meta,
	})
}

func (h *Handler) CreatePublicReview(c *fiber.Ctx) error {
	var req dto.CreateReviewRequest
	if err := c.BodyParser(&req); err != nil {
		return writeAppError(c, &AppError{
			Status:  422,
			Code:    "VALIDATION_ERROR",
			Message: "Некорректные данные запроса",
		})
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

	tourID, err := parseUUID(req.TourID)
	if err != nil {
		return writeAppError(c, err)
	}

	review, err := h.reviews.CreateReview(c.Context(), application.CreateReviewInput{
		TourID:            tourID,
		ClientName:        req.ClientName,
		Rating:            req.Rating,
		Text:              req.Text,
		IsApproved:        false,
		AllowDistribution: req.AllowDistribution,
	})
	if err != nil {
		return respondError(c, err, MapError)
	}

	var userID *uuid.UUID
	if id, ok := appmiddleware.UserIDFromContext(c); ok {
		userID = &id
	}
	if _, err := h.consents.RecordConsent(c.Context(), application.RecordConsentInput{
		ConsentType: domain.ConsentTypePersonalData,
		UserID:      userID,
		IP:          c.IP(),
		UserAgent:   c.Get("User-Agent"),
	}); err != nil {
		return respondError(c, err, MapError)
	}
	if req.AllowDistribution {
		if _, err := h.consents.RecordConsent(c.Context(), application.RecordConsentInput{
			ConsentType: domain.ConsentTypeDistribution,
			UserID:      userID,
			IP:          c.IP(),
			UserAgent:   c.Get("User-Agent"),
		}); err != nil {
			return respondError(c, err, MapError)
		}
	}

	return c.Status(fiber.StatusCreated).JSON(dto.DataEnvelope[dto.ReviewResponse]{
		Data: dto.ToReviewResponse(review),
	})
}
