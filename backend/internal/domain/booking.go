package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Booking struct {
	ID          uuid.UUID
	TourID      uuid.UUID
	UserID      *uuid.UUID
	Name        string
	Phone       string
	Email       string
	PeopleCount int
	Status        BookingStatus
	PaymentStatus PaymentStatus
	TotalPrice    int
	Comment     string
	Overbooked  bool
	Source      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type NewBookingInput struct {
	ID                     uuid.UUID
	Tour                   Tour
	UserID                 *uuid.UUID
	Name                   string
	Phone                  string
	Email                  string
	PeopleCount            int
	Comment                string
	Source                 string
	HotTourDiscountPercent int
	Now                    time.Time
}

func NewBooking(input NewBookingInput) (Booking, error) {
	if input.ID == uuid.Nil {
		return Booking{}, ErrInvalidID
	}
	if input.Tour.ID == uuid.Nil {
		return Booking{}, ErrInvalidID
	}
	if strings.TrimSpace(input.Name) == "" {
		return Booking{}, ErrInvalidContactName
	}
	if strings.TrimSpace(input.Phone) == "" {
		return Booking{}, ErrInvalidPhone
	}
	if input.PeopleCount <= 0 {
		return Booking{}, ErrInvalidPeopleCount
	}
	if !input.Tour.CanAcceptPeople(input.PeopleCount) {
		return Booking{}, ErrInsufficientSlots
	}

	now := input.Now
	if now.IsZero() {
		now = time.Now().UTC()
	}

	catalog := TourCatalogContext{
		Today:                  now,
		HotTourDiscountPercent: input.HotTourDiscountPercent,
	}
	totalPrice := input.Tour.BookingTotalIn(input.PeopleCount, catalog)
	if totalPrice < 0 {
		return Booking{}, ErrInvalidTotalPrice
	}

	return Booking{
		ID:            input.ID,
		TourID:        input.Tour.ID,
		UserID:        input.UserID,
		Name:          strings.TrimSpace(input.Name),
		Phone:         strings.TrimSpace(input.Phone),
		Email:         strings.TrimSpace(input.Email),
		PeopleCount:   input.PeopleCount,
		Status:        BookingStatusNew,
		PaymentStatus: DefaultPaymentStatus(totalPrice),
		TotalPrice:    totalPrice,
		Comment:     strings.TrimSpace(input.Comment),
		Overbooked:  input.Tour.BookingWouldBeOverbooked(input.PeopleCount),
		Source:      strings.TrimSpace(input.Source),
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func (b *Booking) ChangeStatus(next BookingStatus) error {
	if !next.Valid() {
		return ErrInvalidBookingStatus
	}
	if !b.Status.CanTransitionTo(next) {
		return ErrInvalidStatusTransition
	}
	b.Status = next
	b.UpdatedAt = time.Now().UTC()
	return nil
}

func (b Booking) IsTerminal() bool {
	return b.Status == BookingStatusCompleted || b.Status == BookingStatusCancelled
}

func (b *Booking) ChangePaymentStatus(next PaymentStatus) error {
	if !next.Valid() {
		return ErrInvalidPaymentStatus
	}
	b.PaymentStatus = next
	b.UpdatedAt = time.Now().UTC()
	return nil
}
