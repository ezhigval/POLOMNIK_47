package handlers

import (
	"github.com/gofiber/fiber/v2"

	"palomnik/internal/adapters/http/fiber/dto"
	"palomnik/internal/ports"
)

func (h *Handler) ListTourRecommendations(c *fiber.Ctx) error {
	id, err := parseUUID(c.Params("id"))
	if err != nil {
		return writeAppError(c, err)
	}
	if h.ai == nil {
		return c.JSON(dto.ListEnvelope[dto.TourResponse]{
			Data: []dto.TourResponse{},
			Meta: ports.PageMeta{Page: 1, Limit: 0, Total: 0},
		})
	}
	tours, err := h.ai.RecommendPublishedTours(c.Context(), id, 3)
	if err != nil {
		return respondError(c, err, MapError)
	}
	items := dto.ToTourList(tours)
	return c.JSON(dto.ListEnvelope[dto.TourResponse]{
		Data: items,
		Meta: ports.PageMeta{Page: 1, Limit: len(items), Total: len(items)},
	})
}

func (h *Handler) ManagementSupportDraft(c *fiber.Ctx) error {
	id, err := parseUUID(c.Params("id"))
	if err != nil {
		return writeAppError(c, err)
	}
	if h.ai == nil {
		return c.JSON(dto.DataEnvelope[dto.SupportDraftResponse]{
			Data: dto.SupportDraftResponse{
				Configured: false,
				Escalate:   true,
				Note:       "Черновик для менеджера. Клиенту не отправляется. Ответить должен человек. Цены и богословие не выдумывать.",
			},
		})
	}
	draft, err := h.ai.SupportDraft(c.Context(), id)
	if err != nil {
		return respondError(c, err, MapError)
	}
	return c.JSON(dto.DataEnvelope[dto.SupportDraftResponse]{Data: dto.ToSupportDraftResponse(draft)})
}

func (h *Handler) ManagementMetricsDigest(c *fiber.Ctx) error {
	if h.ai == nil {
		return c.JSON(dto.DataEnvelope[dto.MetricsDigestResponse]{
			Data: dto.MetricsDigestResponse{Configured: false, BookingsByStatus: map[string]int{}},
		})
	}
	digest, err := h.ai.MetricsDigest(c.Context())
	if err != nil {
		return respondError(c, err, MapError)
	}
	return c.JSON(dto.DataEnvelope[dto.MetricsDigestResponse]{Data: dto.ToMetricsDigestResponse(digest)})
}

func (h *Handler) ManagementWatchdog(c *fiber.Ctx) error {
	if h.ai == nil {
		return c.JSON(dto.DataEnvelope[dto.WatchdogResponse]{
			Data: dto.WatchdogResponse{RestartAttempted: false, Database: "unknown"},
		})
	}
	return c.JSON(dto.DataEnvelope[dto.WatchdogResponse]{
		Data: dto.ToWatchdogResponse(h.ai.WatchdogReport(c.Context())),
	})
}
