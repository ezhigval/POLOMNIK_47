package application

import (
	"context"
	"fmt"
	"strings"

	"polomnik/internal/ports"
)

type WebhookService struct {
	bookings     *BookingService
	inbound      ports.CRMInboundPort
	inboundToken string
	requireToken bool
}

func NewWebhookService(bookings *BookingService, inbound ports.CRMInboundPort, inboundToken string, requireToken bool) *WebhookService {
	return &WebhookService{
		bookings:     bookings,
		inbound:      inbound,
		inboundToken: inboundToken,
		requireToken: requireToken,
	}
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
	expected := strings.TrimSpace(s.inboundToken)
	if expected == "" {
		if s.requireToken {
			return ErrUnauthorized
		}
		return nil
	}
	if token != expected {
		return ErrUnauthorized
	}
	return nil
}
