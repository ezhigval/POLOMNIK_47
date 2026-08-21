package bitrix

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"polomnik/internal/adapters/integration/stub"
	"polomnik/internal/config"
	"polomnik/internal/domain"
	"polomnik/internal/ports"
)

// CRMAdapter implements ports.CRMPort and ports.CRMInboundPort via Bitrix24 REST API.
type CRMAdapter struct {
	cfg    config.Config
	client Client
}

func NewCRMAdapter(cfg config.Config) CRMAdapter {
	return CRMAdapter{
		cfg:    cfg,
		client: NewClient(cfg.BitrixWebhookURL, cfg.IntegrationHTTPTimeout),
	}
}

var (
	_ ports.CRMPort        = CRMAdapter{}
	_ ports.CRMInboundPort = CRMAdapter{}
)

func (a CRMAdapter) configured() bool {
	return strings.TrimSpace(a.cfg.BitrixWebhookURL) != ""
}

func (a CRMAdapter) originatorID() string {
	if id := strings.TrimSpace(a.cfg.BitrixOriginatorID); id != "" {
		return id
	}
	return "polomnik"
}

func (a CRMAdapter) SyncTour(ctx context.Context, tour domain.Tour) (ports.IntegrationResult, error) {
	if !a.configured() {
		return stub.NotConfigured("bitrix24"), nil
	}

	productID, err := a.syncProduct(ctx, tour)
	if err != nil {
		return ports.IntegrationResult{}, err
	}

	return ports.IntegrationResult{
		Status:     ports.IntegrationStatusSynced,
		ExternalID: productID,
		Message:    "product synced in bitrix24 catalog",
	}, nil
}

func (a CRMAdapter) PullTours(context.Context) ([]domain.Tour, error) {
	return nil, nil
}

func (a CRMAdapter) PushBooking(ctx context.Context, booking domain.Booking) (ports.IntegrationResult, error) {
	if !a.configured() {
		return stub.NotConfigured("bitrix24"), nil
	}

	dealID, found, err := a.findDealByOrigin(ctx, booking.ID.String())
	if err != nil {
		return ports.IntegrationResult{}, err
	}

	contactID, err := a.ensureContact(ctx, booking)
	if err != nil {
		return ports.IntegrationResult{}, err
	}

	fields := a.dealFields(booking, contactID)

	if found {
		_, err = a.client.Call(ctx, "crm.deal.update", map[string]any{
			"id":     dealID,
			"fields": fields,
		})
	} else {
		var raw json.RawMessage
		raw, err = a.client.Call(ctx, "crm.deal.add", map[string]any{"fields": fields})
		if err != nil {
			return ports.IntegrationResult{}, err
		}
		dealID, err = decodeID(raw)
	}
	if err != nil {
		return ports.IntegrationResult{}, err
	}

	if productID, ok, err := a.findProductByOrigin(ctx, booking.TourID.String()); err != nil {
		return ports.IntegrationResult{}, err
	} else if ok {
		unitPrice := booking.TotalPrice
		if booking.PeopleCount > 0 {
			unitPrice = booking.TotalPrice / booking.PeopleCount
		}
		if err := a.attachProductRows(ctx, dealID, productID, booking, unitPrice); err != nil {
			return ports.IntegrationResult{}, err
		}
	}

	message := "deal created in bitrix24"
	if found {
		message = "deal updated in bitrix24"
	}

	return ports.IntegrationResult{
		Status:     ports.IntegrationStatusSynced,
		ExternalID: dealID,
		Message:    message,
	}, nil
}

func (a CRMAdapter) UpdateBookingStatus(ctx context.Context, booking domain.Booking) (ports.IntegrationResult, error) {
	if !a.configured() {
		return stub.NotConfigured("bitrix24"), nil
	}

	dealID, found, err := a.findDealByOrigin(ctx, booking.ID.String())
	if err != nil {
		return ports.IntegrationResult{}, err
	}
	if !found {
		return a.PushBooking(ctx, booking)
	}

	fields := map[string]any{}
	if stage := bookingStageID(a.cfg, booking.Status); stage != "" {
		fields["STAGE_ID"] = stage
	}
	if booking.Comment != "" {
		fields["COMMENTS"] = booking.Comment
	}

	if len(fields) > 0 {
		_, err = a.client.Call(ctx, "crm.deal.update", map[string]any{
			"id":     dealID,
			"fields": fields,
		})
		if err != nil {
			return ports.IntegrationResult{}, err
		}
	}

	return ports.IntegrationResult{
		Status:     ports.IntegrationStatusSynced,
		ExternalID: dealID,
		Message:    "deal status updated in bitrix24",
	}, nil
}

func (a CRMAdapter) PushReview(ctx context.Context, review domain.Review) (ports.IntegrationResult, error) {
	if !a.configured() {
		return stub.NotConfigured("bitrix24"), nil
	}

	productID, found, err := a.findProductByOrigin(ctx, review.TourID.String())
	if err != nil {
		return ports.IntegrationResult{}, err
	}
	if !found {
		return ports.IntegrationResult{
			Status:  ports.IntegrationStatusSynced,
			Message: "review skipped: tour product not found in bitrix24",
		}, nil
	}

	comment := fmt.Sprintf("Отзыв %d/5 от %s: %s", review.Rating, review.ClientName, review.Text)
	raw, err := a.client.Call(ctx, "crm.timeline.comment.add", map[string]any{
		"fields": map[string]any{
			"ENTITY_ID":   productID,
			"ENTITY_TYPE": "product",
			"COMMENT":     comment,
		},
	})
	if err != nil {
		return ports.IntegrationResult{}, err
	}

	commentID, err := decodeID(raw)
	if err != nil {
		return ports.IntegrationResult{}, err
	}

	return ports.IntegrationResult{
		Status:     ports.IntegrationStatusSynced,
		ExternalID: commentID,
		Message:    "review posted to bitrix24 product timeline",
	}, nil
}

func (a CRMAdapter) PullReviews(context.Context) ([]domain.Review, error) {
	return nil, nil
}

func (a CRMAdapter) ResolveBookingStatusFromDeal(ctx context.Context, dealID string) (uuid.UUID, domain.BookingStatus, error) {
	if !a.configured() {
		return uuid.Nil, "", fmt.Errorf("bitrix24 is not configured")
	}
	if strings.TrimSpace(dealID) == "" {
		return uuid.Nil, "", fmt.Errorf("deal id is required")
	}

	raw, err := a.client.Call(ctx, "crm.deal.get", map[string]any{"id": dealID})
	if err != nil {
		return uuid.Nil, "", err
	}

	originID, stageID, err := decodeDealFields(raw)
	if err != nil {
		return uuid.Nil, "", err
	}
	if originID == "" {
		return uuid.Nil, "", domain.ErrNotFound
	}

	bookingID, err := uuid.Parse(originID)
	if err != nil {
		return uuid.Nil, "", fmt.Errorf("invalid booking origin id: %w", err)
	}

	status, ok := bookingStatusFromStage(a.cfg, stageID)
	if !ok {
		return bookingID, "", fmt.Errorf("unknown bitrix stage: %s", stageID)
	}

	return bookingID, status, nil
}

func (a CRMAdapter) findDealByOrigin(ctx context.Context, originID string) (string, bool, error) {
	raw, err := a.client.Call(ctx, "crm.deal.list", map[string]any{
		"filter": map[string]any{
			"ORIGINATOR_ID": a.originatorID(),
			"ORIGIN_ID":     originID,
		},
		"select": []string{"ID"},
	})
	if err != nil {
		return "", false, err
	}
	return decodeFirstListID(raw)
}

func (a CRMAdapter) ensureContact(ctx context.Context, booking domain.Booking) (string, error) {
	if booking.Phone == "" {
		return "", fmt.Errorf("booking phone is required for bitrix contact")
	}

	raw, err := a.client.Call(ctx, "crm.duplicate.findbycomm", map[string]any{
		"entity_type": "CONTACT",
		"type":        "PHONE",
		"values":      []string{booking.Phone},
	})
	if err != nil {
		return "", err
	}

	var duplicates map[string][]any
	if len(raw) > 0 && string(raw) != "null" {
		if err := json.Unmarshal(raw, &duplicates); err != nil {
			return "", err
		}
	}
	if ids, ok := duplicates["CONTACT"]; ok && len(ids) > 0 {
		return fmt.Sprintf("%v", ids[0]), nil
	}

	firstName, lastName := splitName(booking.Name)
	fields := map[string]any{
		"NAME":      firstName,
		"LAST_NAME": lastName,
		"PHONE":     []map[string]string{{"VALUE": booking.Phone, "VALUE_TYPE": "WORK"}},
	}
	if booking.Email != "" {
		fields["EMAIL"] = []map[string]string{{"VALUE": booking.Email, "VALUE_TYPE": "WORK"}}
	}

	raw, err = a.client.Call(ctx, "crm.contact.add", map[string]any{"fields": fields})
	if err != nil {
		return "", err
	}
	return decodeID(raw)
}

func (a CRMAdapter) dealFields(booking domain.Booking, contactID string) map[string]any {
	fields := map[string]any{
		"TITLE":              fmt.Sprintf("Заявка: %s", booking.Name),
		"CONTACT_ID":         contactID,
		"OPPORTUNITY":        booking.TotalPrice,
		"CURRENCY_ID":        "RUB",
		"ORIGINATOR_ID":      a.originatorID(),
		"ORIGIN_ID":          booking.ID.String(),
		"SOURCE_ID":          "WEB",
		"SOURCE_DESCRIPTION": booking.Source,
		"COMMENTS":           booking.Comment,
	}

	if a.cfg.BitrixDealCategoryID > 0 {
		fields["CATEGORY_ID"] = a.cfg.BitrixDealCategoryID
	}
	if stage := bookingStageID(a.cfg, booking.Status); stage != "" {
		fields["STAGE_ID"] = stage
	}

	if uf := strings.TrimSpace(a.cfg.BitrixUFTourID); uf != "" {
		fields[uf] = booking.TourID.String()
	}
	if uf := strings.TrimSpace(a.cfg.BitrixUFPeopleCount); uf != "" {
		fields[uf] = booking.PeopleCount
	}

	return fields
}
