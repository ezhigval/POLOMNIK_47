package application

import (
	"context"
	"time"

	"github.com/google/uuid"

	"palomnik/internal/domain"
	"palomnik/internal/ports"
)

type PassengerService struct {
	passengers ports.PassengerRepository
}

func NewPassengerService(passengers ports.PassengerRepository) *PassengerService {
	return &PassengerService{passengers: passengers}
}

type PassengerInput struct {
	Name      string
	Phone     string
	BirthDate time.Time
	Passport  string
}

func (s *PassengerService) List(ctx context.Context, userID uuid.UUID) ([]domain.Passenger, error) {
	if s == nil || s.passengers == nil {
		return nil, domain.ErrNotFound
	}
	if userID == uuid.Nil {
		return nil, domain.ErrInvalidID
	}
	return s.passengers.ListPassengers(ctx, userID)
}

func (s *PassengerService) Create(ctx context.Context, userID uuid.UUID, input PassengerInput) (domain.Passenger, error) {
	if s == nil || s.passengers == nil {
		return domain.Passenger{}, domain.ErrNotFound
	}
	passenger, err := domain.NewPassenger(domain.NewPassengerInput{
		ID:        uuid.New(),
		UserID:    userID,
		Name:      input.Name,
		Phone:     input.Phone,
		BirthDate: input.BirthDate,
		Passport:  input.Passport,
	})
	if err != nil {
		return domain.Passenger{}, err
	}
	return s.passengers.CreatePassenger(ctx, passenger)
}

func (s *PassengerService) Update(ctx context.Context, userID, id uuid.UUID, input PassengerInput) (domain.Passenger, error) {
	if s == nil || s.passengers == nil {
		return domain.Passenger{}, domain.ErrNotFound
	}
	existing, err := s.passengers.GetPassenger(ctx, userID, id)
	if err != nil {
		return domain.Passenger{}, err
	}
	updated, err := existing.Update(input.Name, input.Phone, input.BirthDate, input.Passport, time.Time{})
	if err != nil {
		return domain.Passenger{}, err
	}
	return s.passengers.UpdatePassenger(ctx, updated)
}

func (s *PassengerService) Delete(ctx context.Context, userID, id uuid.UUID) error {
	if s == nil || s.passengers == nil {
		return domain.ErrNotFound
	}
	if userID == uuid.Nil || id == uuid.Nil {
		return domain.ErrInvalidID
	}
	return s.passengers.DeletePassenger(ctx, userID, id)
}
