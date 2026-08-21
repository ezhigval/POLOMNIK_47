package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"polomnik/internal/adapters/http/fiber/dto"
	"polomnik/internal/ports"
)

func (h *Handler) ManagementListIntegrationReferences(c *fiber.Ctx) error {
	list, err := h.integrations.ListReferences(c.Context(), ports.IntegrationReferenceFilters{
		ExternalSystem:  strings.TrimSpace(c.Query("external_system")),
		LocalEntityType: strings.TrimSpace(c.Query("local_entity_type")),
		SyncStatus:      strings.TrimSpace(c.Query("sync_status")),
	}, parsePagination(c))
	if err != nil {
		return respondError(c, err, MapError)
	}

	items := make([]dto.IntegrationReferenceResponse, 0, len(list.Items))
	for _, ref := range list.Items {
		items = append(items, dto.ToIntegrationReferenceResponse(ref))
	}

	return c.JSON(dto.ListEnvelope[dto.IntegrationReferenceResponse]{
		Data: items,
		Meta: list.Meta,
	})
}

func (h *Handler) ManagementListOutboxEvents(c *fiber.Ctx) error {
	list, err := h.integrations.ListOutboxEvents(c.Context(), ports.OutboxFilters{
		Status:     strings.TrimSpace(c.Query("status")),
		EntityType: strings.TrimSpace(c.Query("entity_type")),
		EventType:  strings.TrimSpace(c.Query("event_type")),
	}, parsePagination(c))
	if err != nil {
		return respondError(c, err, MapError)
	}

	items := make([]dto.OutboxEventResponse, 0, len(list.Items))
	for _, event := range list.Items {
		items = append(items, dto.ToOutboxEventResponse(event))
	}

	return c.JSON(dto.ListEnvelope[dto.OutboxEventResponse]{
		Data: items,
		Meta: list.Meta,
	})
}
