package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"polomnik/internal/adapters/http/fiber/dto"
	"polomnik/internal/ports"
)

func (h *Handler) ListReviews(c *fiber.Ctx) error {
	var tourID *uuid.UUID
	if value := c.Query("tour_id"); value != "" {
		parsed, err := parseUUID(value)
		if err != nil {
			return writeAppError(c, err.(*AppError))
		}
		tourID = &parsed
	}

	rating, err := parseOptionalInt(c.Query("rating"))
	if err != nil {
		return writeAppError(c, err.(*AppError))
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
