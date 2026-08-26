package domain

import "time"

// TourCatalogContext carries read-time rules for public tour presentation.
// Regular tours without fixed dates skip burning/hot logic.
type TourCatalogContext struct {
	Today                  time.Time
	HotTourDiscountPercent int
}

func (c TourCatalogContext) today() time.Time {
	if c.Today.IsZero() {
		return calendarDay(time.Now().UTC())
	}
	return calendarDay(c.Today)
}

func calendarDay(value time.Time) time.Time {
	return value.UTC().Truncate(24 * time.Hour)
}

// IsBurningOn is true when the tour departs today (date_start calendar day).
// Distinct from manual is_hot ("Популярный"). Regular tours never burn.
func (t Tour) IsBurningOn(today time.Time) bool {
	if t.IsRegular || t.DateStart.IsZero() {
		return false
	}
	return calendarDay(t.DateStart).Equal(calendarDay(today))
}

func (t Tour) IsBurningIn(catalog TourCatalogContext) bool {
	return t.IsBurningOn(catalog.today())
}

// ApplyPercentDiscount returns price minus percent (integer math, rounded down).
func ApplyPercentDiscount(price int, percent int) int {
	if price <= 0 || percent <= 0 {
		return price
	}
	if percent >= 100 {
		return 0
	}
	return price - (price*percent)/100
}

// UnitPriceIn returns the per-person price for catalog/booking on the given day.
func (t Tour) UnitPriceIn(catalog TourCatalogContext) int {
	if t.Price <= 0 {
		return 0
	}
	if t.IsBurningIn(catalog) && catalog.HotTourDiscountPercent > 0 {
		return ApplyPercentDiscount(t.Price, catalog.HotTourDiscountPercent)
	}
	return t.Price
}

// BookingTotalIn uses the effective unit price (burning discount when configured).
func (t Tour) BookingTotalIn(peopleCount int, catalog TourCatalogContext) int {
	unit := t.UnitPriceIn(catalog)
	if peopleCount <= 0 || unit <= 0 {
		return 0
	}
	return unit * peopleCount
}
