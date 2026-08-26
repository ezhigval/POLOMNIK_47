package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"palomnik/internal/adapters/http/fiber/dto"
	"palomnik/internal/application"
	"palomnik/internal/domain"
	"palomnik/internal/ports"
)

func (h *Handler) ManagementListTours(c *fiber.Ctx) error {
	list, err := h.tours.ListTours(c.Context(), ports.TourFilters{}, parsePagination(c))
	if err != nil {
		return respondError(c, err, MapError)
	}

	items := make([]dto.ManagementTourResponse, 0, len(list.Items))
	for _, tour := range list.Items {
		items = append(items, dto.ToManagementTourResponse(tour))
	}

	return c.JSON(dto.ListEnvelope[dto.ManagementTourResponse]{
		Data: items,
		Meta: list.Meta,
	})
}

func (h *Handler) ManagementGetTour(c *fiber.Ctx) error {
	id, err := parseUUID(c.Params("id"))
	if err != nil {
		return writeAppError(c, err)
	}

	tour, err := h.tours.GetTour(c.Context(), id)
	if err != nil {
		return respondError(c, err, MapError)
	}

	return c.JSON(dto.DataEnvelope[dto.ManagementTourResponse]{
		Data: dto.ToManagementTourResponse(tour),
	})
}

func (h *Handler) ManagementCreateTour(c *fiber.Ctx) error {
	var req dto.TourUpsertRequest
	if err := c.BodyParser(&req); err != nil {
		return writeAppError(c, &AppError{
			Status:  422,
			Code:    "VALIDATION_ERROR",
			Message: "Некорректные данные запроса",
		})
	}

	dateStart, dateEnd, err := parseTourSchedule(req)
	if err != nil {
		return writeAppError(c, err)
	}

	tour, err := h.tours.CreateTour(c.Context(), application.CreateTourInput{
		Slug:               req.Slug,
		Title:              req.Title,
		Description:        req.Description,
		Price:              req.Price,
		Currency:           req.Currency,
		DateStart:          dateStart,
		DateEnd:            dateEnd,
		SlotsTotal:         req.SlotsTotal,
		SlotsLeft:          req.SlotsLeft,
		Location:           req.Location,
		Images:             req.Images,
		IsActive:           req.IsActive,
		IsHot:              req.IsHot,
		IsRegular:          req.IsRegular,
		OverbookingEnabled: req.OverbookingEnabled,
	})
	if err != nil {
		return respondError(c, err, MapError)
	}

	return c.Status(fiber.StatusCreated).JSON(dto.DataEnvelope[dto.ManagementTourResponse]{
		Data: dto.ToManagementTourResponse(tour),
	})
}

func (h *Handler) ManagementUpdateTour(c *fiber.Ctx) error {
	id, err := parseUUID(c.Params("id"))
	if err != nil {
		return writeAppError(c, err)
	}

	var req dto.TourUpsertRequest
	if err := c.BodyParser(&req); err != nil {
		return writeAppError(c, &AppError{
			Status:  422,
			Code:    "VALIDATION_ERROR",
			Message: "Некорректные данные запроса",
		})
	}

	dateStart, dateEnd, err := parseTourSchedule(req)
	if err != nil {
		return writeAppError(c, err)
	}

	tour, err := h.tours.UpdateTour(c.Context(), id, application.UpdateTourInput{
		Slug:               req.Slug,
		Title:              req.Title,
		Description:        req.Description,
		Price:              req.Price,
		Currency:           req.Currency,
		DateStart:          dateStart,
		DateEnd:            dateEnd,
		SlotsTotal:         req.SlotsTotal,
		SlotsLeft:          req.SlotsLeft,
		Location:           req.Location,
		Images:             req.Images,
		IsActive:           req.IsActive,
		IsHot:              req.IsHot,
		IsRegular:          req.IsRegular,
		OverbookingEnabled: req.OverbookingEnabled,
	})
	if err != nil {
		return respondError(c, err, MapError)
	}

	return c.JSON(dto.DataEnvelope[dto.ManagementTourResponse]{
		Data: dto.ToManagementTourResponse(tour),
	})
}

func (h *Handler) ManagementDeleteTour(c *fiber.Ctx) error {
	id, err := parseUUID(c.Params("id"))
	if err != nil {
		return writeAppError(c, err)
	}

	if err := h.tours.DeleteTour(c.Context(), id); err != nil {
		return respondError(c, err, MapError)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) ManagementListBookings(c *fiber.Ctx) error {
	filters, err := parseManagementBookingFilters(c)
	if err != nil {
		return writeAppError(c, err)
	}

	if strings.EqualFold(strings.TrimSpace(c.Query("format")), "csv") {
		items, err := h.bookings.ListAllBookings(c.Context(), filters)
		if err != nil {
			return respondError(c, err, MapError)
		}
		body, err := managementBookingsCSV(items)
		if err != nil {
			return err
		}
		c.Set("Content-Type", "text/csv; charset=utf-8")
		c.Set("Content-Disposition", `attachment; filename="bookings.csv"`)
		return c.Send(body)
	}

	list, err := h.bookings.ListBookings(c.Context(), filters, parsePagination(c))
	if err != nil {
		return respondError(c, err, MapError)
	}

	items := make([]dto.ManagementBookingResponse, 0, len(list.Items))
	for _, booking := range list.Items {
		items = append(items, dto.ToManagementBookingResponse(booking))
	}

	return c.JSON(dto.ListEnvelope[dto.ManagementBookingResponse]{
		Data: items,
		Meta: list.Meta,
	})
}

func (h *Handler) ManagementGetBooking(c *fiber.Ctx) error {
	id, err := parseUUID(c.Params("id"))
	if err != nil {
		return writeAppError(c, err)
	}

	booking, err := h.bookings.GetBooking(c.Context(), id)
	if err != nil {
		return respondError(c, err, MapError)
	}

	return c.JSON(dto.DataEnvelope[dto.ManagementBookingResponse]{
		Data: dto.ToManagementBookingResponse(booking),
	})
}

func (h *Handler) ManagementUpdateBookingStatus(c *fiber.Ctx) error {
	id, err := parseUUID(c.Params("id"))
	if err != nil {
		return writeAppError(c, err)
	}

	var req dto.UpdateBookingStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return writeAppError(c, &AppError{
			Status:  422,
			Code:    "VALIDATION_ERROR",
			Message: "Некорректные данные запроса",
		})
	}

	status := domain.BookingStatus(req.Status)
	booking, err := h.bookings.UpdateBookingStatus(c.Context(), id, status)
	if err != nil {
		return respondError(c, err, MapError)
	}

	return c.JSON(dto.DataEnvelope[dto.ManagementBookingResponse]{
		Data: dto.ToManagementBookingResponse(booking),
	})
}

func (h *Handler) ManagementUpdateBookingPaymentStatus(c *fiber.Ctx) error {
	id, err := parseUUID(c.Params("id"))
	if err != nil {
		return writeAppError(c, err)
	}

	var req dto.UpdateBookingPaymentStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return writeAppError(c, &AppError{
			Status:  422,
			Code:    "VALIDATION_ERROR",
			Message: "Некорректные данные запроса",
		})
	}

	status := domain.PaymentStatus(req.PaymentStatus)
	booking, err := h.bookings.UpdateBookingPaymentStatus(c.Context(), id, status)
	if err != nil {
		return respondError(c, err, MapError)
	}

	return c.JSON(dto.DataEnvelope[dto.ManagementBookingResponse]{
		Data: dto.ToManagementBookingResponse(booking),
	})
}

func (h *Handler) ManagementListReviews(c *fiber.Ctx) error {
	list, err := h.reviews.ListReviews(c.Context(), ports.ReviewFilters{}, parsePagination(c))
	if err != nil {
		return respondError(c, err, MapError)
	}

	items := make([]dto.ManagementReviewResponse, 0, len(list.Items))
	for _, review := range list.Items {
		items = append(items, dto.ToManagementReviewResponse(review))
	}

	return c.JSON(dto.ListEnvelope[dto.ManagementReviewResponse]{
		Data: items,
		Meta: list.Meta,
	})
}

func (h *Handler) ManagementCreateReview(c *fiber.Ctx) error {
	var req dto.CreateReviewRequest
	if err := c.BodyParser(&req); err != nil {
		return writeAppError(c, &AppError{
			Status:  422,
			Code:    "VALIDATION_ERROR",
			Message: "Некорректные данные запроса",
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
		IsApproved:        req.IsApproved,
		AllowDistribution: req.AllowDistribution,
	})
	if err != nil {
		return respondError(c, err, MapError)
	}

	return c.Status(fiber.StatusCreated).JSON(dto.DataEnvelope[dto.ManagementReviewResponse]{
		Data: dto.ToManagementReviewResponse(review),
	})
}

func (h *Handler) ManagementApproveReview(c *fiber.Ctx) error {
	id, err := parseUUID(c.Params("id"))
	if err != nil {
		return writeAppError(c, err)
	}

	review, err := h.reviews.ApproveReview(c.Context(), id)
	if err != nil {
		return respondError(c, err, MapError)
	}

	return c.JSON(dto.DataEnvelope[dto.ManagementReviewResponse]{
		Data: dto.ToManagementReviewResponse(review),
	})
}

func (h *Handler) ManagementRejectReview(c *fiber.Ctx) error {
	id, err := parseUUID(c.Params("id"))
	if err != nil {
		return writeAppError(c, err)
	}

	review, err := h.reviews.RejectReview(c.Context(), id)
	if err != nil {
		return respondError(c, err, MapError)
	}

	return c.JSON(dto.DataEnvelope[dto.ManagementReviewResponse]{
		Data: dto.ToManagementReviewResponse(review),
	})
}

func (h *Handler) ManagementSetReviewReply(c *fiber.Ctx) error {
	id, err := parseUUID(c.Params("id"))
	if err != nil {
		return writeAppError(c, err)
	}

	var req dto.SetCompanyReplyRequest
	if err := c.BodyParser(&req); err != nil {
		return writeAppError(c, &AppError{
			Status:  422,
			Code:    "VALIDATION_ERROR",
			Message: "Некорректные данные запроса",
		})
	}

	review, err := h.reviews.SetCompanyReply(c.Context(), id, req.CompanyReply)
	if err != nil {
		return respondError(c, err, MapError)
	}

	return c.JSON(dto.DataEnvelope[dto.ManagementReviewResponse]{
		Data: dto.ToManagementReviewResponse(review),
	})
}

func (h *Handler) ManagementDeleteReview(c *fiber.Ctx) error {
	id, err := parseUUID(c.Params("id"))
	if err != nil {
		return writeAppError(c, err)
	}

	if err := h.reviews.DeleteReview(c.Context(), id); err != nil {
		return respondError(c, err, MapError)
	}

	return c.SendStatus(fiber.StatusNoContent)
}
