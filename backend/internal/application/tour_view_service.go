package application

import (
	"context"

	"github.com/google/uuid"

	"palomnik/internal/ports"
)

const maxViewedTours = 20

type TourViewService struct {
	views ports.TourViewRepository
	tours ports.TourRepository
}

func NewTourViewService(views ports.TourViewRepository, tours ports.TourRepository) *TourViewService {
	return &TourViewService{views: views, tours: tours}
}

func (s *TourViewService) RecordView(ctx context.Context, userID, tourID uuid.UUID) error {
	if _, err := s.tours.GetTour(ctx, tourID); err != nil {
		return err
	}
	return s.views.RecordTourView(ctx, userID, tourID)
}

func (s *TourViewService) ListViewedTourIDs(ctx context.Context, userID uuid.UUID, limit int) ([]uuid.UUID, error) {
	if limit < 1 {
		limit = maxViewedTours
	}
	if limit > maxViewedTours {
		limit = maxViewedTours
	}
	return s.views.ListViewedTourIDs(ctx, userID, limit)
}

func (s *TourViewService) MergeViews(ctx context.Context, userID uuid.UUID, tourIDs []uuid.UUID) error {
	if len(tourIDs) == 0 {
		return nil
	}
	unique := make([]uuid.UUID, 0, len(tourIDs))
	seen := make(map[uuid.UUID]struct{}, len(tourIDs))
	for _, id := range tourIDs {
		if id == uuid.Nil {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		unique = append(unique, id)
		if len(unique) >= maxViewedTours {
			break
		}
	}
	return s.views.MergeTourViews(ctx, userID, unique)
}
