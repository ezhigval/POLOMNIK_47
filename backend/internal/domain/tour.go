package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Tour struct {
	ID                 uuid.UUID
	Slug               string
	Title              string
	Description        string
	Price              int
	Currency           string
	DateStart          time.Time
	DateEnd            time.Time
	SlotsTotal         int
	SlotsLeft          int
	Location           string
	Images             []string
	IsActive           bool
	IsHot              bool
	OverbookingEnabled bool
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type NewTourInput struct {
	ID                 uuid.UUID
	Slug               string
	Title              string
	Description        string
	Price              int
	Currency           string
	DateStart          time.Time
	DateEnd            time.Time
	SlotsTotal         int
	SlotsLeft          int
	Location           string
	Images             []string
	IsActive           bool
	IsHot              bool
	OverbookingEnabled bool
	Now                time.Time
}

func NewTour(input NewTourInput) (Tour, error) {
	if input.ID == uuid.Nil {
		return Tour{}, ErrInvalidID
	}
	if strings.TrimSpace(input.Slug) == "" {
		return Tour{}, ErrInvalidSlug
	}
	if strings.TrimSpace(input.Title) == "" {
		return Tour{}, ErrInvalidTitle
	}
	if input.Price < 0 {
		return Tour{}, ErrInvalidPrice
	}
	if strings.TrimSpace(input.Currency) == "" {
		return Tour{}, ErrInvalidCurrency
	}
	if input.DateStart.IsZero() || input.DateEnd.IsZero() || input.DateStart.After(input.DateEnd) {
		return Tour{}, ErrInvalidDateRange
	}
	if err := validateSlots(input.SlotsTotal, input.SlotsLeft); err != nil {
		return Tour{}, err
	}

	now := input.Now
	if now.IsZero() {
		now = time.Now().UTC()
	}

	return Tour{
		ID:                 input.ID,
		Slug:               strings.TrimSpace(input.Slug),
		Title:              strings.TrimSpace(input.Title),
		Description:        strings.TrimSpace(input.Description),
		Price:              input.Price,
		Currency:           strings.TrimSpace(input.Currency),
		DateStart:          input.DateStart,
		DateEnd:            input.DateEnd,
		SlotsTotal:         input.SlotsTotal,
		SlotsLeft:          input.SlotsLeft,
		Location:           strings.TrimSpace(input.Location),
		Images:             stringSliceCopy(input.Images),
		IsActive:           input.IsActive,
		IsHot:              input.IsHot,
		OverbookingEnabled: input.OverbookingEnabled,
		CreatedAt:          now,
		UpdatedAt:          now,
	}, nil
}

func (t Tour) CanAcceptPeople(peopleCount int) bool {
	if peopleCount <= 0 {
		return false
	}
	return t.SlotsLeft >= peopleCount || t.OverbookingEnabled
}

func (t Tour) BookingWouldBeOverbooked(peopleCount int) bool {
	return peopleCount > 0 && peopleCount > t.SlotsLeft && t.OverbookingEnabled
}

func (t *Tour) ReserveSlots(peopleCount int) error {
	if peopleCount <= 0 {
		return ErrInvalidPeopleCount
	}
	if t.SlotsLeft < peopleCount {
		if t.OverbookingEnabled {
			t.UpdatedAt = time.Now().UTC()
			return nil
		}
		return ErrInsufficientSlots
	}
	t.SlotsLeft -= peopleCount
	t.UpdatedAt = time.Now().UTC()
	return nil
}

func (t *Tour) ReleaseSlots(peopleCount int) error {
	if peopleCount <= 0 {
		return ErrInvalidPeopleCount
	}
	t.SlotsLeft += peopleCount
	if t.SlotsLeft > t.SlotsTotal {
		t.SlotsLeft = t.SlotsTotal
	}
	t.UpdatedAt = time.Now().UTC()
	return nil
}

func validateSlots(total int, left int) error {
	if total < 0 || left < 0 || left > total {
		return ErrInvalidSlots
	}
	return nil
}

func stringSliceCopy(values []string) []string {
	if len(values) == 0 {
		return []string{}
	}
	out := make([]string, len(values))
	copy(out, values)
	return out
}
