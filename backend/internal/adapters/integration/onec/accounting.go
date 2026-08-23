package onec

import (
	"context"
	"strings"

	"palomnik/internal/adapters/integration/stub"
	"palomnik/internal/config"
	"palomnik/internal/domain"
	"palomnik/internal/ports"
)

type bookingPayload struct {
	BookingID   string `json:"booking_id"`
	TourID      string `json:"tour_id"`
	Name        string `json:"name"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	PeopleCount int    `json:"people_count"`
	TotalPrice  int    `json:"total_price"`
	Status      string `json:"status"`
	Comment     string `json:"comment"`
	Source      string `json:"source"`
	Overbooked  bool   `json:"overbooked"`
}

type counterpartyPayload struct {
	BookingID string `json:"booking_id"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
}

// AccountingAdapter implements ports.AccountingPort via HTTP JSON or 1C OData.
type AccountingAdapter struct {
	inner ports.AccountingPort
}

func NewAccountingAdapter(cfg config.Config) AccountingAdapter {
	switch strings.ToLower(strings.TrimSpace(cfg.OneCProtocol)) {
	case "odata":
		return AccountingAdapter{inner: NewODataAccountingAdapter(cfg)}
	default:
		return AccountingAdapter{inner: NewHTTPAccountingAdapter(cfg)}
	}
}

var _ ports.AccountingPort = AccountingAdapter{}

func (a AccountingAdapter) ExportBooking(ctx context.Context, booking domain.Booking) (ports.IntegrationResult, error) {
	return a.inner.ExportBooking(ctx, booking)
}

func (a AccountingAdapter) SyncCounterparty(ctx context.Context, booking domain.Booking) (ports.IntegrationResult, error) {
	return a.inner.SyncCounterparty(ctx, booking)
}

// HTTPAccountingAdapter uses custom HTTP JSON endpoints (see docs/ONEC_INTEGRATOR_TZ.md).
type HTTPAccountingAdapter struct {
	cfg    config.Config
	client Client
}

func NewHTTPAccountingAdapter(cfg config.Config) HTTPAccountingAdapter {
	return HTTPAccountingAdapter{
		cfg: cfg,
		client: NewClient(
			cfg.OneCBaseURL,
			cfg.OneCUsername,
			cfg.OneCPassword,
			cfg.OneCAPIKey,
			cfg.IntegrationHTTPTimeout,
		),
	}
}

var _ ports.AccountingPort = HTTPAccountingAdapter{}

func (a HTTPAccountingAdapter) configured() bool {
	return strings.TrimSpace(a.cfg.OneCBaseURL) != ""
}

func (a HTTPAccountingAdapter) exportPath() string {
	if path := strings.TrimSpace(a.cfg.OneCExportBookingPath); path != "" {
		return path
	}
	return "/hs/palomnik/booking"
}

func (a HTTPAccountingAdapter) counterpartyPath() string {
	if path := strings.TrimSpace(a.cfg.OneCCounterpartyPath); path != "" {
		return path
	}
	return "/hs/palomnik/counterparty"
}

func (a HTTPAccountingAdapter) ExportBooking(ctx context.Context, booking domain.Booking) (ports.IntegrationResult, error) {
	if !a.configured() {
		return stub.NotConfigured("1c"), nil
	}
	if !shouldExportBooking(booking.Status) {
		return ports.IntegrationResult{
			Status:  ports.IntegrationStatusSynced,
			Message: "booking export skipped until confirmed",
		}, nil
	}

	resp, err := a.client.PostJSON(ctx, a.exportPath(), bookingPayload{
		BookingID:   booking.ID.String(),
		TourID:      booking.TourID.String(),
		Name:        booking.Name,
		Phone:       booking.Phone,
		Email:       booking.Email,
		PeopleCount: booking.PeopleCount,
		TotalPrice:  booking.TotalPrice,
		Status:      string(booking.Status),
		Comment:     booking.Comment,
		Source:      booking.Source,
		Overbooked:  booking.Overbooked,
	})
	if err != nil {
		return ports.IntegrationResult{}, err
	}

	return ports.IntegrationResult{
		Status:     ports.IntegrationStatusSynced,
		ExternalID: resp.ExternalID,
		Message:    resp.Message,
	}, nil
}

func (a HTTPAccountingAdapter) SyncCounterparty(ctx context.Context, booking domain.Booking) (ports.IntegrationResult, error) {
	if !a.configured() {
		return stub.NotConfigured("1c"), nil
	}

	resp, err := a.client.PostJSON(ctx, a.counterpartyPath(), counterpartyPayload{
		BookingID: booking.ID.String(),
		Name:      booking.Name,
		Phone:     booking.Phone,
		Email:     booking.Email,
	})
	if err != nil {
		return ports.IntegrationResult{}, err
	}

	return ports.IntegrationResult{
		Status:     ports.IntegrationStatusSynced,
		ExternalID: resp.ExternalID,
		Message:    resp.Message,
	}, nil
}
