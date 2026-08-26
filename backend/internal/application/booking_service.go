package application

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"

	"palomnik/internal/domain"
	"palomnik/internal/ports"
)

type BookingService struct {
	bookings      ports.BookingRepository
	tours         ports.TourRepository
	crm           ports.CRMPort
	accounting    ports.AccountingPort
	notifications ports.NotificationPort
	tx            ports.TransactionManager
}

func NewBookingService(
	bookings ports.BookingRepository,
	tours ports.TourRepository,
	crm ports.CRMPort,
	accounting ports.AccountingPort,
	notifications ports.NotificationPort,
	tx ports.TransactionManager,
) *BookingService {
	return &BookingService{
		bookings:      bookings,
		tours:         tours,
		crm:           crm,
		accounting:    accounting,
		notifications: notifications,
		tx:            tx,
	}
}

type CreateBookingInput struct {
	TourID      uuid.UUID
	UserID      *uuid.UUID
	Name        string
	Phone       string
	Email       string
	PeopleCount int
	Comment     string
	Source      string
}

type CreateBookingResult struct {
	Booking           domain.Booking
	IntegrationStatus ports.IntegrationStatus
}

func (s *BookingService) CreateBooking(ctx context.Context, input CreateBookingInput) (CreateBookingResult, error) {
	tour, err := s.tours.GetTour(ctx, input.TourID)
	if err != nil {
		return CreateBookingResult{}, err
	}
	if !tour.IsActive {
		return CreateBookingResult{}, ErrTourInactive
	}
	if tourEnded(tour) {
		return CreateBookingResult{}, ErrTourExpired
	}

	source := strings.TrimSpace(input.Source)
	if source == "" {
		source = "web"
	}

	booking, err := domain.NewBooking(domain.NewBookingInput{
		ID:          uuid.New(),
		Tour:        tour,
		UserID:      input.UserID,
		Name:        input.Name,
		Phone:       input.Phone,
		Email:       input.Email,
		PeopleCount: input.PeopleCount,
		Comment:     input.Comment,
		Source:      source,
	})
	if err != nil {
		return CreateBookingResult{}, err
	}

	var created domain.Booking
	err = s.runInTx(ctx, func(ctx context.Context) error {
		if err := s.tours.ReserveSlots(ctx, tour.ID, booking.PeopleCount); err != nil {
			return err
		}

		created, err = s.bookings.CreateBooking(ctx, booking)
		if err != nil {
			return err
		}

		if created.Overbooked {
			created, err = s.bookings.MarkBookingOverbooked(ctx, created.ID)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return CreateBookingResult{}, err
	}

	result, err := s.crm.PushBooking(ctx, created)
	integrationStatus := resolveIntegrationStatus(result, err)
	runAccountingSync(s.accounting, ctx, created)
	_ = s.notifications.NotifyNewBooking(ctx, created, tour)

	return CreateBookingResult{
		Booking:           created,
		IntegrationStatus: integrationStatus,
	}, nil
}

func (s *BookingService) runInTx(ctx context.Context, fn func(context.Context) error) error {
	if s.tx == nil {
		return fn(ctx)
	}
	return s.tx.WithinTransaction(ctx, fn)
}

func (s *BookingService) GetBooking(ctx context.Context, id uuid.UUID) (domain.Booking, error) {
	return s.bookings.GetBooking(ctx, id)
}

func (s *BookingService) ListBookings(ctx context.Context, filters ports.BookingFilters, pagination ports.Pagination) (ports.BookingList, error) {
	return s.bookings.ListBookings(ctx, filters, pagination)
}

func (s *BookingService) ListAllBookings(ctx context.Context, filters ports.BookingFilters) ([]domain.Booking, error) {
	var all []domain.Booking
	for page := 1; page <= 1000; page++ {
		list, err := s.bookings.ListBookings(ctx, filters, ports.NormalizePagination(page, 100))
		if err != nil {
			return nil, err
		}
		all = append(all, list.Items...)
		if !list.Meta.HasNext {
			return all, nil
		}
	}
	return all, nil
}

func (s *BookingService) UpdateBookingStatus(ctx context.Context, id uuid.UUID, status domain.BookingStatus) (domain.Booking, error) {
	existing, err := s.bookings.GetBooking(ctx, id)
	if err != nil {
		return domain.Booking{}, err
	}
	if existing.Status == status {
		return existing, nil
	}

	updated, err := s.bookings.UpdateBookingStatus(ctx, id, status)
	if err != nil {
		return domain.Booking{}, err
	}

	if status == domain.BookingStatusCancelled && !existing.IsTerminal() {
		if err := s.tours.ReleaseSlots(ctx, existing.TourID, existing.PeopleCount); err != nil {
			return domain.Booking{}, err
		}
	}

	_, _ = s.crm.UpdateBookingStatus(ctx, updated)
	runAccountingSync(s.accounting, ctx, updated)

	tour, tourErr := s.tours.GetTour(ctx, updated.TourID)
	if tourErr != nil {
		tour = domain.Tour{ID: updated.TourID}
	}
	_ = s.notifications.NotifyBookingStatusChanged(ctx, updated, tour, existing.Status)

	return updated, nil
}

func tourEnded(tour domain.Tour) bool {
	return tour.ScheduleEndedOn(time.Now())
}
