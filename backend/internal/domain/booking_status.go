package domain

type BookingStatus string

const (
	BookingStatusNew       BookingStatus = "NEW"
	BookingStatusContacted BookingStatus = "CONTACTED"
	BookingStatusConfirmed BookingStatus = "CONFIRMED"
	BookingStatusCompleted BookingStatus = "COMPLETED"
	BookingStatusCancelled BookingStatus = "CANCELLED"
)

func (s BookingStatus) Valid() bool {
	switch s {
	case BookingStatusNew,
		BookingStatusContacted,
		BookingStatusConfirmed,
		BookingStatusCompleted,
		BookingStatusCancelled:
		return true
	default:
		return false
	}
}

func (s BookingStatus) CanTransitionTo(next BookingStatus) bool {
	if !s.Valid() || !next.Valid() {
		return false
	}
	if s == next {
		return true
	}
	if s == BookingStatusCompleted || s == BookingStatusCancelled {
		return false
	}
	if next == BookingStatusCancelled {
		return true
	}

	switch s {
	case BookingStatusNew:
		return next == BookingStatusContacted
	case BookingStatusContacted:
		return next == BookingStatusConfirmed
	case BookingStatusConfirmed:
		return next == BookingStatusCompleted
	default:
		return false
	}
}
