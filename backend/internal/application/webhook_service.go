package application

import (
	"context"
	"fmt"
	"strings"

	"polomnik/internal/config"
	"polomnik/internal/ports"
)

type WebhookService struct {
	bookings *BookingService
	inbound  ports.CRMInboundPort
	cfg      config.Config
}

func NewWebhookService(bookings *BookingService, inbound ports.CRMInboundPort, cfg config.Config) *WebhookService {
	return &WebhookService{bookings: bookings, inbound: inbound, cfg: cfg}
}

type BitrixDealWebhookInput struct {
	Event  string
	DealID string
	Token  string
}

func (s *WebhookService) HandleBitrixDealUpdate(ctx context.Context, input BitrixDealWebhookInput) error {
	if err := s.validateInboundToken(input.Token); err != nil {
		return err
	}
	if input.DealID == "" {
		return fmt.Errorf("deal id is required")
	}

	bookingID, status, err := s.inbound.ResolveBookingStatusFromDeal(ctx, input.DealID)
	if err != nil {
		return err
	}

	existing, err := s.bookings.GetBooking(ctx, bookingID)
	if err != nil {
		return err
	}
	if existing.Status == status {
		return nil
	}

	_, err = s.bookings.UpdateBookingStatus(ctx, bookingID, status)
	return err
}

func (s *WebhookService) validateInboundToken(token string) error {
	expected := s.cfg.BitrixInboundToken
	if expected == "" {
		if s.cfg.IsProduction() || strings.EqualFold(s.cfg.CRMAdapter, "bitrix") {
			return ErrUnauthorized
		}
		return nil
	}
	if token != expected {
		return ErrUnauthorized
	}
	return nil
}
