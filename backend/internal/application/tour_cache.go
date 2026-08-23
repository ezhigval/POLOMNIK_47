package application

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"palomnik/internal/domain"
	"palomnik/internal/ports"
)

const (
	toursNamespaceKey = "tours:namespace"
	defaultCacheTTL   = 5 * time.Minute
)

func (s *TourService) ListPublicToursCached(ctx context.Context, filters ports.TourFilters, pagination ports.Pagination) (ports.TourList, error) {
	if s.cache == nil {
		return s.ListPublicTours(ctx, filters, pagination)
	}

	key, err := s.cacheKey(ctx, "list:public", tourListCacheID(filters, pagination))
	if err != nil {
		return s.ListPublicTours(ctx, filters, pagination)
	}

	var cached ports.TourList
	if ok, _ := s.readCache(ctx, key, &cached); ok {
		return cached, nil
	}

	list, err := s.ListPublicTours(ctx, filters, pagination)
	if err != nil {
		return ports.TourList{}, err
	}

	_ = s.writeCache(ctx, key, list)
	return list, nil
}

func (s *TourService) GetPublicTourCached(ctx context.Context, id uuid.UUID) (domain.Tour, error) {
	if s.cache == nil {
		return s.GetPublicTour(ctx, id)
	}

	key, err := s.cacheKey(ctx, "detail:public", id.String())
	if err != nil {
		return s.GetPublicTour(ctx, id)
	}

	var cached domain.Tour
	if ok, _ := s.readCache(ctx, key, &cached); ok {
		return cached, nil
	}

	tour, err := s.GetPublicTour(ctx, id)
	if err != nil {
		return domain.Tour{}, err
	}

	_ = s.writeCache(ctx, key, tour)
	return tour, nil
}

func (s *TourService) ListPopularToursCached(ctx context.Context, limit int) ([]domain.Tour, error) {
	if s.cache == nil {
		return s.ListPopularTours(ctx, limit)
	}

	key, err := s.cacheKey(ctx, "popular", fmt.Sprintf("limit=%d", limit))
	if err != nil {
		return s.ListPopularTours(ctx, limit)
	}

	var cached []domain.Tour
	if ok, _ := s.readCache(ctx, key, &cached); ok {
		return cached, nil
	}

	items, err := s.ListPopularTours(ctx, limit)
	if err != nil {
		return nil, err
	}

	_ = s.writeCache(ctx, key, items)
	return items, nil
}

func (s *TourService) invalidateTourCache(ctx context.Context) {
	if s.cache == nil {
		return
	}
	_ = s.cache.Set(ctx, toursNamespaceKey, []byte(uuid.NewString()), defaultCacheTTL)
}

func (s *TourService) cacheNamespace(ctx context.Context) (string, error) {
	value, err := s.cache.Get(ctx, toursNamespaceKey)
	if err == nil {
		return string(value), nil
	}
	if err != nil && !errors.Is(err, ports.ErrCacheMiss) {
		return "", err
	}
	namespace := uuid.NewString()
	_ = s.cache.Set(ctx, toursNamespaceKey, []byte(namespace), defaultCacheTTL)
	return namespace, nil
}

func (s *TourService) cacheKey(ctx context.Context, parts ...string) (string, error) {
	namespace, err := s.cacheNamespace(ctx)
	if err != nil {
		return "", err
	}

	raw := namespace
	for _, part := range parts {
		raw += ":" + part
	}
	sum := sha256.Sum256([]byte(raw))
	return "tours:" + hex.EncodeToString(sum[:]), nil
}

func tourListCacheID(filters ports.TourFilters, pagination ports.Pagination) string {
	payload, _ := json.Marshal(struct {
		Filters    ports.TourFilters
		Pagination ports.Pagination
	}{
		Filters:    filters,
		Pagination: pagination,
	})
	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:])
}

func (s *TourService) readCache(ctx context.Context, key string, target any) (bool, error) {
	value, err := s.cache.Get(ctx, key)
	if err != nil {
		return false, nil
	}
	if err := json.Unmarshal(value, target); err != nil {
		return false, nil
	}
	return true, nil
}

func (s *TourService) writeCache(ctx context.Context, key string, value any) error {
	payload, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return s.cache.Set(ctx, key, payload, defaultCacheTTL)
}
