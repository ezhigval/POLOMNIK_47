package bitrix

import (
	"context"
	"encoding/json"
	"fmt"

	"polomnik/internal/domain"
)

func (a CRMAdapter) findProductByOrigin(ctx context.Context, tourID string) (string, bool, error) {
	raw, err := a.client.Call(ctx, "crm.product.list", map[string]any{
		"filter": map[string]any{
			"ORIGINATOR_ID": a.originatorID(),
			"ORIGIN_ID":     tourID,
		},
		"select": []string{"ID"},
	})
	if err != nil {
		return "", false, err
	}
	return decodeFirstListID(raw)
}

func (a CRMAdapter) syncProduct(ctx context.Context, tour domain.Tour) (string, error) {
	productID, found, err := a.findProductByOrigin(ctx, tour.ID.String())
	if err != nil {
		return "", err
	}

	fields := productFields(tour, a.originatorID())
	if found {
		_, err = a.client.Call(ctx, "crm.product.update", map[string]any{
			"id":     productID,
			"fields": fields,
		})
		return productID, err
	}

	raw, err := a.client.Call(ctx, "crm.product.add", map[string]any{"fields": fields})
	if err != nil {
		return "", err
	}
	return decodeID(raw)
}

func productFields(tour domain.Tour, originatorID string) map[string]any {
	active := "N"
	if tour.IsActive {
		active = "Y"
	}
	fields := map[string]any{
		"NAME":          tour.Title,
		"DESCRIPTION":   tour.Description,
		"PRICE":         tour.Price,
		"CURRENCY_ID":   tour.Currency,
		"ORIGINATOR_ID": originatorID,
		"ORIGIN_ID":     tour.ID.String(),
		"ACTIVE":        active,
	}
	if tour.Location != "" {
		fields["DESCRIPTION"] = fmt.Sprintf("%s\n\nЛокация: %s", tour.Description, tour.Location)
	}
	return fields
}

func (a CRMAdapter) attachProductRows(ctx context.Context, dealID, productID string, booking domain.Booking, unitPrice int) error {
	_, err := a.client.Call(ctx, "crm.deal.productrows.set", map[string]any{
		"id": dealID,
		"rows": []map[string]any{
			{
				"PRODUCT_ID": productID,
				"PRICE":      unitPrice,
				"QUANTITY":   booking.PeopleCount,
			},
		},
	})
	return err
}

func decodeDealFields(raw json.RawMessage) (originID string, stageID string, err error) {
	var deal map[string]any
	if err := json.Unmarshal(raw, &deal); err != nil {
		return "", "", err
	}
	originID, _ = deal["ORIGIN_ID"].(string)
	stageID, _ = deal["STAGE_ID"].(string)
	return originID, stageID, nil
}
