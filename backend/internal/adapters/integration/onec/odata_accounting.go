package onec

import (
	"context"
	"fmt"
	"strings"
	"time"

	"polomnik/internal/adapters/integration/stub"
	"polomnik/internal/config"
	"polomnik/internal/domain"
	"polomnik/internal/ports"
)

type ODataAccountingAdapter struct {
	cfg    config.Config
	client ODataClient
}

func NewODataAccountingAdapter(cfg config.Config) ODataAccountingAdapter {
	return ODataAccountingAdapter{
		cfg: cfg,
		client: NewODataClient(
			cfg.OneCBaseURL,
			cfg.OneCODataPath,
			cfg.OneCUsername,
			cfg.OneCPassword,
			cfg.IntegrationHTTPTimeout,
		),
	}
}

var _ ports.AccountingPort = ODataAccountingAdapter{}

func (a ODataAccountingAdapter) configured() bool {
	return strings.TrimSpace(a.cfg.OneCBaseURL) != ""
}

func (a ODataAccountingAdapter) counterpartyEntity() string {
	if entity := strings.TrimSpace(a.cfg.OneCODataCounterparty); entity != "" {
		return entity
	}
	return "Catalog_Контрагенты"
}

func (a ODataAccountingAdapter) orderEntity() string {
	if entity := strings.TrimSpace(a.cfg.OneCODataSalesOrder); entity != "" {
		return entity
	}
	return "Document_ЗаказКлиента"
}

func (a ODataAccountingAdapter) ExportBooking(ctx context.Context, booking domain.Booking) (ports.IntegrationResult, error) {
	if !a.configured() {
		return stub.NotConfigured("1c"), nil
	}
	if !shouldExportBooking(booking.Status) {
		return ports.IntegrationResult{
			Status:  ports.IntegrationStatusSynced,
			Message: "booking export skipped until confirmed",
		}, nil
	}

	marker := bookingMarker(booking.ID.String())
	raw, err := a.client.Get(ctx, a.orderEntity(), fmt.Sprintf("Comment eq '%s'", marker))
	if err != nil {
		return ports.IntegrationResult{}, err
	}
	if refKey, found, err := firstODataRefKey(raw); err != nil {
		return ports.IntegrationResult{}, err
	} else if found {
		return ports.IntegrationResult{
			Status:     ports.IntegrationStatusSynced,
			ExternalID: refKey,
			Message:    "sales order already exists in 1c",
		}, nil
	}

	counterpartyID, err := a.ensureCounterparty(ctx, booking)
	if err != nil {
		return ports.IntegrationResult{}, err
	}

	refKey, err := a.client.Post(ctx, a.orderEntity(), map[string]any{
		"Date":            time.Now().UTC().Format(time.RFC3339),
		"Comment":         marker,
		"Контрагент_Key":  counterpartyID,
		"Комментарий":     fmt.Sprintf("Palomnik tour %s, %d чел.", booking.TourID, booking.PeopleCount),
	})
	if err != nil {
		return ports.IntegrationResult{}, err
	}

	return ports.IntegrationResult{
		Status:     ports.IntegrationStatusSynced,
		ExternalID: refKey,
		Message:    "sales order created in 1c odata",
	}, nil
}

func (a ODataAccountingAdapter) SyncCounterparty(ctx context.Context, booking domain.Booking) (ports.IntegrationResult, error) {
	if !a.configured() {
		return stub.NotConfigured("1c"), nil
	}

	refKey, err := a.ensureCounterparty(ctx, booking)
	if err != nil {
		return ports.IntegrationResult{}, err
	}

	return ports.IntegrationResult{
		Status:     ports.IntegrationStatusSynced,
		ExternalID: refKey,
		Message:    "counterparty synced in 1c odata",
	}, nil
}

func (a ODataAccountingAdapter) ensureCounterparty(ctx context.Context, booking domain.Booking) (string, error) {
	filter := fmt.Sprintf("Description eq '%s'", escapeODataString(booking.Phone))
	if booking.Phone == "" {
		filter = fmt.Sprintf("Description eq '%s'", escapeODataString(booking.Name))
	}

	raw, err := a.client.Get(ctx, a.counterpartyEntity(), filter)
	if err != nil {
		return "", err
	}
	if refKey, found, err := firstODataRefKey(raw); err != nil {
		return "", err
	} else if found {
		return refKey, nil
	}

	return a.client.Post(ctx, a.counterpartyEntity(), map[string]any{
		"Description": booking.Name,
		"Комментарий": fmt.Sprintf("phone: %s email: %s", booking.Phone, booking.Email),
	})
}

func escapeODataString(value string) string {
	return strings.ReplaceAll(value, "'", "''")
}

func shouldExportBooking(status domain.BookingStatus) bool {
	return status == domain.BookingStatusConfirmed || status == domain.BookingStatusCompleted
}
