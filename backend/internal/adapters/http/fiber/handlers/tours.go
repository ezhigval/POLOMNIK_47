package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"palomnik/internal/adapters/http/fiber/dto"
	"palomnik/internal/ports"
)

func (h *Handler) ListTours(c *fiber.Ctx) error {
	dateFrom, err := parseOptionalDate(c.Query("date_from"))
	if err != nil {
		return writeAppError(c, err)
	}
	dateTo, err := parseOptionalDate(c.Query("date_to"))
	if err != nil {
		return writeAppError(c, err)
	}
	priceMin, err := parseOptionalInt(c.Query("price_min"))
	if err != nil {
		return writeAppError(c, err)
	}
	priceMax, err := parseOptionalInt(c.Query("price_max"))
	if err != nil {
		return writeAppError(c, err)
	}
	minSlots, err := parseOptionalInt(c.Query("min_slots"))
	if err != nil {
		return writeAppError(c, err)
	}

	list, err := h.tours.ListPublicToursCached(c.Context(), ports.TourFilters{
		DateFrom: dateFrom,
		DateTo:   dateTo,
		PriceMin: priceMin,
		PriceMax: priceMax,
		Location: c.Query("location"),
		Query:    c.Query("q"),
		MinSlots: minSlots,
		IsHot:    parseOptionalBool(c.Query("is_hot")),
	}, parsePagination(c))
	if err != nil {
		return respondError(c, err, MapError)
	}

	items := make([]dto.TourResponse, 0, len(list.Items))
	for _, tour := range list.Items {
		items = append(items, dto.ToTourResponse(tour))
	}

	return c.JSON(dto.ListEnvelope[dto.TourResponse]{
		Data: items,
		Meta: list.Meta,
	})
}

func (h *Handler) GetTour(c *fiber.Ctx) error {
	tour, err := h.tours.GetPublicTourByKeyCached(c.Context(), c.Params("id"))
	if err != nil {
		return respondError(c, err, MapError)
	}

	return c.JSON(dto.DataEnvelope[dto.TourResponse]{
		Data: dto.ToTourResponse(tour),
	})
}

func (h *Handler) ListPopularTours(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit"))
	if limit < 1 {
		limit = 10
	}

	items, err := h.tours.ListPopularToursCached(c.Context(), limit)
	if err != nil {
		return respondError(c, err, MapError)
	}

	response := make([]dto.TourResponse, 0, len(items))
	for _, tour := range items {
		response = append(response, dto.ToTourResponse(tour))
	}

	return c.JSON(dto.DataEnvelope[[]dto.TourResponse]{
		Data: response,
	})
}

func (h *Handler) ListTourReviews(c *fiber.Ctx) error {
	tourID, err := parseUUID(c.Params("id"))
	if err != nil {
		return writeAppError(c, err)
	}

	if _, err := h.tours.GetPublicTourCached(c.Context(), tourID); err != nil {
		return respondError(c, err, MapError)
	}

	list, err := h.reviews.ListPublicReviews(c.Context(), ports.ReviewFilters{
		TourID: &tourID,
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
