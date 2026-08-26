package application

import (
	"context"
	"time"

	"palomnik/internal/domain"
	"palomnik/internal/ports"
)

const tourLifecyclePageSize = 100

// TourCatalog holds env-driven rules for public tour presentation.
type TourCatalog struct {
	HotTourDiscountPercent int
}

func NewTourCatalog(hotTourDiscountPercent int) TourCatalog {
	if hotTourDiscountPercent < 0 {
		hotTourDiscountPercent = 0
	}
	if hotTourDiscountPercent > 100 {
		hotTourDiscountPercent = 100
	}
	return TourCatalog{HotTourDiscountPercent: hotTourDiscountPercent}
}

func (c TourCatalog) Context(now time.Time) domain.TourCatalogContext {
	if now.IsZero() {
		now = time.Now().UTC()
	}
	return domain.TourCatalogContext{
		Today:                  now,
		HotTourDiscountPercent: c.HotTourDiscountPercent,
	}
}

type TourLifecycleService struct {
	tours         *TourService
	notifications ports.NotificationPort
}

func NewTourLifecycleService(tours *TourService, notifications ports.NotificationPort) *TourLifecycleService {
	return &TourLifecycleService{tours: tours, notifications: notifications}
}

// HideExpiredActiveTours deactivates dated tours whose schedule ended before today
// and notifies managers. Regular tours and tours without dates are skipped.
func (s *TourLifecycleService) HideExpiredActiveTours(ctx context.Context, now time.Time) (int, error) {
	if s == nil || s.tours == nil {
		return 0, nil
	}
	if now.IsZero() {
		now = time.Now().UTC()
	}

	active := true
	hidden := 0
	page := 1
	for {
		list, err := s.tours.ListTours(ctx, ports.TourFilters{IsActive: &active}, ports.Pagination{
			Page:  page,
			Limit: tourLifecyclePageSize,
		})
		if err != nil {
			return hidden, err
		}

		for _, tour := range list.Items {
			if !tour.ScheduleEndedOn(now) {
				continue
			}
			inactive := false
			if _, err := s.tours.PatchTourOps(ctx, tour.ID, TourOpsPatch{IsActive: &inactive}); err != nil {
				return hidden, err
			}
			if s.notifications != nil {
				_ = s.notifications.NotifyTourHidden(ctx, tour)
			}
			hidden++
		}

		if len(list.Items) < tourLifecyclePageSize {
			break
		}
		page++
	}

	return hidden, nil
}
