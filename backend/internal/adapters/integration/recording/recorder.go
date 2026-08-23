package recording

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"palomnik/internal/domain"
	"palomnik/internal/ports"
)

type CRM struct {
	inner    ports.CRMPort
	recorder Recorder
}

type CRMOptions struct {
	ExternalSystem string
}

func NewCRM(inner ports.CRMPort, refs ports.IntegrationReferenceRepository, outbox ports.OutboxRepository, opts CRMOptions) CRM {
	return CRM{
		inner: inner,
		recorder: Recorder{
			refs:           refs,
			outbox:         outbox,
			externalSystem: opts.ExternalSystem,
		},
	}
}

var _ ports.CRMPort = CRM{}

func (a CRM) SyncTour(ctx context.Context, tour domain.Tour) (ports.IntegrationResult, error) {
	result, err := a.inner.SyncTour(ctx, tour)
	a.recorder.record(ctx, domain.OutboxEventCRMSyncTour, domain.EntityTypeTour, tour.ID, domain.ExternalEntityTypeProduct, result, err)
	return result, err
}

func (a CRM) PullTours(ctx context.Context) ([]domain.Tour, error) {
	return a.inner.PullTours(ctx)
}

func (a CRM) PushBooking(ctx context.Context, booking domain.Booking) (ports.IntegrationResult, error) {
	result, err := a.inner.PushBooking(ctx, booking)
	a.recorder.record(ctx, domain.OutboxEventCRMPushBooking, domain.EntityTypeBooking, booking.ID, domain.ExternalEntityTypeDeal, result, err)
	return result, err
}

func (a CRM) UpdateBookingStatus(ctx context.Context, booking domain.Booking) (ports.IntegrationResult, error) {
	result, err := a.inner.UpdateBookingStatus(ctx, booking)
	a.recorder.record(ctx, domain.OutboxEventCRMUpdateBookingStatus, domain.EntityTypeBooking, booking.ID, domain.ExternalEntityTypeDeal, result, err)
	return result, err
}

func (a CRM) PushReview(ctx context.Context, review domain.Review) (ports.IntegrationResult, error) {
	result, err := a.inner.PushReview(ctx, review)
	a.recorder.record(ctx, domain.OutboxEventCRMPushReview, domain.EntityTypeReview, review.ID, domain.ExternalEntityTypeReview, result, err)
	return result, err
}

func (a CRM) PullReviews(ctx context.Context) ([]domain.Review, error) {
	return a.inner.PullReviews(ctx)
}

type Accounting struct {
	inner    ports.AccountingPort
	recorder Recorder
}

type AccountingOptions struct {
	ExternalSystem string
}

func NewAccounting(inner ports.AccountingPort, refs ports.IntegrationReferenceRepository, outbox ports.OutboxRepository, opts AccountingOptions) Accounting {
	return Accounting{
		inner: inner,
		recorder: Recorder{
			refs:           refs,
			outbox:         outbox,
			externalSystem: opts.ExternalSystem,
		},
	}
}

var _ ports.AccountingPort = Accounting{}

func (a Accounting) ExportBooking(ctx context.Context, booking domain.Booking) (ports.IntegrationResult, error) {
	result, err := a.inner.ExportBooking(ctx, booking)
	a.recorder.record(ctx, domain.OutboxEventAccountingExport, domain.EntityTypeBooking, booking.ID, domain.ExternalEntityTypeOrder, result, err)
	return result, err
}

func (a Accounting) SyncCounterparty(ctx context.Context, booking domain.Booking) (ports.IntegrationResult, error) {
	result, err := a.inner.SyncCounterparty(ctx, booking)
	a.recorder.record(ctx, domain.OutboxEventAccountingCounterparty, domain.EntityTypeBooking, booking.ID, "counterparty", result, err)
	return result, err
}

type Recorder struct {
	refs           ports.IntegrationReferenceRepository
	outbox         ports.OutboxRepository
	externalSystem string
}

func (r Recorder) record(
	ctx context.Context,
	eventType string,
	localEntityType string,
	localEntityID uuid.UUID,
	externalEntityType string,
	result ports.IntegrationResult,
	callErr error,
) {
	status := domain.SyncStatus(result.Status)
	lastError := result.Message
	if callErr != nil {
		status = domain.SyncStatusFailed
		lastError = callErr.Error()
	}

	externalID := result.ExternalID
	if externalID == "" {
		externalID = domain.PlaceholderExternalID(localEntityID)
	}

	if r.refs != nil {
		ref, err := domain.NewIntegrationReference(domain.NewIntegrationReferenceInput{
			ID:                 uuid.New(),
			LocalEntityType:    localEntityType,
			LocalEntityID:      localEntityID,
			ExternalSystem:     r.externalSystem,
			ExternalEntityType: externalEntityType,
			ExternalEntityID:   externalID,
			SyncStatus:         status,
			LastError:          lastError,
			Now:                time.Now().UTC(),
		})
		if err == nil {
			_, _ = r.refs.UpsertReference(ctx, ref)
		}
	}

	if r.outbox == nil || !shouldEnqueueOutbox(status) {
		return
	}

	payload, _ := json.Marshal(map[string]string{
		"external_system":      r.externalSystem,
		"external_entity_type": externalEntityType,
		"integration_status": string(status),
		"message":              lastError,
	})

	event, err := domain.NewOutboxEvent(domain.NewOutboxEventInput{
		ID:         uuid.New(),
		EventType:  eventType,
		EntityType: localEntityType,
		EntityID:   localEntityID,
		Payload:    payload,
		Status:     domain.OutboxStatusPending,
		LastError:  lastError,
		Now:        time.Now().UTC(),
	})
	if err != nil {
		return
	}

	_, _ = r.outbox.Enqueue(ctx, event)
}

func shouldEnqueueOutbox(status domain.SyncStatus) bool {
	return status == domain.SyncStatusPending || status == domain.SyncStatusFailed
}
