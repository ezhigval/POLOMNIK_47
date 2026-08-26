package application

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"

	"palomnik/internal/domain"
	"palomnik/internal/ports"
)

type TourService struct {
	tours ports.TourRepository
	cache ports.CachePort
	crm   ports.CRMPort
}

func NewTourService(tours ports.TourRepository, cache ports.CachePort, crm ports.CRMPort) *TourService {
	return &TourService{tours: tours, cache: cache, crm: crm}
}

type CreateTourInput struct {
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
	IsRegular          bool
	OverbookingEnabled bool
	HotDiscountPercent int
}

type UpdateTourInput = CreateTourInput

func (s *TourService) ListPublicTours(ctx context.Context, filters ports.TourFilters, pagination ports.Pagination) (ports.TourList, error) {
	active := true
	filters.IsActive = &active
	return s.tours.ListTours(ctx, filters, pagination)
}

func (s *TourService) GetPublicTour(ctx context.Context, id uuid.UUID) (domain.Tour, error) {
	tour, err := s.tours.GetTour(ctx, id)
	if err != nil {
		return domain.Tour{}, err
	}
	if !tour.IsActive {
		return domain.Tour{}, domain.ErrNotFound
	}
	return tour, nil
}

func (s *TourService) FindPublicTour(ctx context.Context, key string) (domain.Tour, error) {
	tour, err := s.FindTour(ctx, key)
	if err != nil {
		return domain.Tour{}, err
	}
	if !tour.IsActive {
		return domain.Tour{}, domain.ErrNotFound
	}
	return tour, nil
}

func (s *TourService) ListPopularTours(ctx context.Context, limit int) ([]domain.Tour, error) {
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	active := true
	hot := true
	list, err := s.tours.ListTours(ctx, ports.TourFilters{
		IsActive: &active,
		IsHot:    &hot,
	}, ports.Pagination{Page: 1, Limit: limit})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (s *TourService) ListTours(ctx context.Context, filters ports.TourFilters, pagination ports.Pagination) (ports.TourList, error) {
	return s.tours.ListTours(ctx, filters, pagination)
}

func (s *TourService) GetTour(ctx context.Context, id uuid.UUID) (domain.Tour, error) {
	return s.tours.GetTour(ctx, id)
}

func (s *TourService) CreateTour(ctx context.Context, input CreateTourInput) (domain.Tour, error) {
	tour, err := domain.NewTour(domain.NewTourInput{
		ID:                 uuid.New(),
		Slug:               input.Slug,
		Title:              input.Title,
		Description:        input.Description,
		Price:              input.Price,
		Currency:           input.Currency,
		DateStart:          input.DateStart,
		DateEnd:            input.DateEnd,
		SlotsTotal:         input.SlotsTotal,
		SlotsLeft:          input.SlotsLeft,
		Location:           input.Location,
		Images:             input.Images,
		IsActive:           input.IsActive,
		IsHot:              input.IsHot,
		IsRegular:          input.IsRegular,
		OverbookingEnabled: input.OverbookingEnabled,
		HotDiscountPercent: input.HotDiscountPercent,
	})
	if err != nil {
		return domain.Tour{}, err
	}
	created, err := s.tours.CreateTour(ctx, tour)
	if err != nil {
		return domain.Tour{}, err
	}
	s.invalidateTourCache(ctx)
	_, _ = s.crm.SyncTour(ctx, created)
	return created, nil
}

func (s *TourService) UpdateTour(ctx context.Context, id uuid.UUID, input UpdateTourInput) (domain.Tour, error) {
	existing, err := s.tours.GetTour(ctx, id)
	if err != nil {
		return domain.Tour{}, err
	}

	updated, err := domain.NewTour(domain.NewTourInput{
		ID:                 existing.ID,
		Slug:               input.Slug,
		Title:              input.Title,
		Description:        input.Description,
		Price:              input.Price,
		Currency:           input.Currency,
		DateStart:          input.DateStart,
		DateEnd:            input.DateEnd,
		SlotsTotal:         input.SlotsTotal,
		SlotsLeft:          input.SlotsLeft,
		Location:           input.Location,
		Images:             input.Images,
		IsActive:           input.IsActive,
		IsHot:              input.IsHot,
		IsRegular:          input.IsRegular,
		OverbookingEnabled: input.OverbookingEnabled,
		HotDiscountPercent: input.HotDiscountPercent,
		Now:                existing.CreatedAt,
	})
	if err != nil {
		return domain.Tour{}, err
	}
	updated.CreatedAt = existing.CreatedAt
	saved, err := s.tours.UpdateTour(ctx, updated)
	if err != nil {
		return domain.Tour{}, err
	}
	s.invalidateTourCache(ctx)
	_, _ = s.crm.SyncTour(ctx, saved)
	return saved, nil
}

// TourOpsPatch changes only slots, price or is_active — the same fields HTTP already updates.
type TourOpsPatch struct {
	SlotsLeft  *int
	SlotsTotal *int
	Price      *int
	IsActive   *bool
}

func (s *TourService) FindTour(ctx context.Context, key string) (domain.Tour, error) {
	key = strings.TrimSpace(key)
	if key == "" {
		return domain.Tour{}, domain.ErrNotFound
	}
	if id, err := uuid.Parse(key); err == nil {
		return s.GetTour(ctx, id)
	}
	return s.tours.GetTourBySlug(ctx, key)
}

func (s *TourService) PatchTourOps(ctx context.Context, id uuid.UUID, patch TourOpsPatch) (domain.Tour, error) {
	existing, err := s.GetTour(ctx, id)
	if err != nil {
		return domain.Tour{}, err
	}
	input := UpdateTourInput{
		Slug:               existing.Slug,
		Title:              existing.Title,
		Description:        existing.Description,
		Price:              existing.Price,
		Currency:           existing.Currency,
		DateStart:          existing.DateStart,
		DateEnd:            existing.DateEnd,
		SlotsTotal:         existing.SlotsTotal,
		SlotsLeft:          existing.SlotsLeft,
		Location:           existing.Location,
		Images:             append([]string(nil), existing.Images...),
		IsActive:           existing.IsActive,
		IsHot:              existing.IsHot,
		IsRegular:          existing.IsRegular,
		OverbookingEnabled: existing.OverbookingEnabled,
	}
	if patch.SlotsTotal != nil {
		input.SlotsTotal = *patch.SlotsTotal
	}
	if patch.SlotsLeft != nil {
		input.SlotsLeft = *patch.SlotsLeft
	}
	if patch.Price != nil {
		input.Price = *patch.Price
	}
	if patch.IsActive != nil {
		input.IsActive = *patch.IsActive
	}
	return s.UpdateTour(ctx, id, input)
}

func (s *TourService) DeleteTour(ctx context.Context, id uuid.UUID) error {
	if err := s.tours.DeleteTour(ctx, id); err != nil {
		return err
	}
	s.invalidateTourCache(ctx)
	return nil
}
