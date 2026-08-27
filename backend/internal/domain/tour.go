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
	IsHot               bool
	IsRegular           bool
	OverbookingEnabled  bool
	HotDiscountPercent  int
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
	IsHot               bool
	IsRegular           bool
	OverbookingEnabled  bool
	HotDiscountPercent  int
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
	if input.HotDiscountPercent < 0 || input.HotDiscountPercent > 100 {
		return Tour{}, ErrInvalidHotDiscountPercent
	}
	if strings.TrimSpace(input.Currency) == "" {
		return Tour{}, ErrInvalidCurrency
	}
	if input.IsRegular {
		input.DateStart = time.Time{}
		input.DateEnd = time.Time{}
	} else if input.DateStart.IsZero() || input.DateEnd.IsZero() || input.DateStart.After(input.DateEnd) {
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
		IsRegular:          input.IsRegular,
		OverbookingEnabled: input.OverbookingEnabled,
		HotDiscountPercent: input.HotDiscountPercent,
		CreatedAt:          now,
		UpdatedAt:          now,
	}, nil
}

func (t Tour) CanAcceptPeople(peopleCount int) bool {
	if peopleCount <= 0 {
		return false
	}
	return t.RemainingSlots() >= peopleCount || t.OverbookingEnabled
}

func (t Tour) BookingWouldBeOverbooked(peopleCount int) bool {
	return peopleCount > 0 && peopleCount > t.SlotsLeft && t.OverbookingEnabled
}

func (t Tour) HasPublicSchedule() bool {
	return !t.IsRegular && !t.DateStart.IsZero() && !t.DateEnd.IsZero()
}

func (t Tour) BookingTotal(peopleCount int) int {
	return t.BookingTotalIn(peopleCount, TourCatalogContext{})
}

func (t Tour) HasPublicPrice() bool {
	return t.Price > 0
}

func (t Tour) ScheduleEndedOn(today time.Time) bool {
	if t.IsRegular || t.DateEnd.IsZero() {
		return false
	}
	day := today.UTC().Truncate(24 * time.Hour)
	end := t.DateEnd.UTC().Truncate(24 * time.Hour)
	return end.Before(day)
}

// CompareToursByCatalog orders dated tours first (by date_start, then id), then regular tours (by id).
func CompareToursByCatalog(a, b Tour) int {
	if a.IsRegular != b.IsRegular {
		if a.IsRegular {
			return 1
		}
		return -1
	}
	if !a.DateStart.Equal(b.DateStart) {
		if a.DateStart.IsZero() {
			return 1
		}
		if b.DateStart.IsZero() {
			return -1
		}
		if a.DateStart.Before(b.DateStart) {
			return -1
		}
		return 1
	}
	if a.ID == b.ID {
		return 0
	}
	if a.ID.String() < b.ID.String() {
		return -1
	}
	return 1
}

func (t *Tour) ReserveSlots(peopleCount int) error {
	if peopleCount <= 0 {
		return ErrInvalidPeopleCount
	}
	t.openRegularGroupIfNeeded()
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

// RemainingSlots is free seats on the public site.
// Regular tours often have only "всего мест" filled in admin (slots_left stays 0);
// that must not look like a sold-out departure.
func (t Tour) RemainingSlots() int {
	if t.IsRegular && t.SlotsTotal > 0 && t.SlotsLeft == 0 {
		return t.SlotsTotal
	}
	return t.SlotsLeft
}

func (t *Tour) openRegularGroupIfNeeded() {
	if t.IsRegular && t.SlotsTotal > 0 && t.SlotsLeft == 0 {
		t.SlotsLeft = t.SlotsTotal
	}
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
