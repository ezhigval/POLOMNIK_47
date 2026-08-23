package application

import (
	"context"
	"errors"
	"time"

	"palomnik/internal/domain"
	"palomnik/internal/ports"
)

type SiteSettingsService struct {
	repo     ports.SiteSettingsRepository
	defaults domain.SiteSettings
}

func NewSiteSettingsService(repo ports.SiteSettingsRepository, defaults domain.SiteSettings) *SiteSettingsService {
	return &SiteSettingsService{repo: repo, defaults: defaults}
}

func (s *SiteSettingsService) Public(ctx context.Context) (domain.SiteSettings, error) {
	stored, err := s.load(ctx)
	if err != nil {
		return domain.SiteSettings{}, err
	}
	return domain.MergeSiteSettings(s.defaults, stored), nil
}

func (s *SiteSettingsService) Settings(ctx context.Context) (domain.SiteSettings, error) {
	return s.Public(ctx)
}

func (s *SiteSettingsService) Update(ctx context.Context, input domain.SiteSettings) (domain.SiteSettings, error) {
	updated, err := domain.NewSiteSettings(input, time.Time{})
	if err != nil {
		return domain.SiteSettings{}, err
	}
	if s == nil || s.repo == nil {
		return domain.MergeSiteSettings(s.defaults, updated), nil
	}
	if _, err := s.repo.UpsertSiteSettings(ctx, updated); err != nil {
		return domain.SiteSettings{}, err
	}
	return s.Public(ctx)
}

func (s *SiteSettingsService) load(ctx context.Context) (domain.SiteSettings, error) {
	if s == nil || s.repo == nil {
		return domain.SiteSettings{}, nil
	}
	stored, err := s.repo.GetSiteSettings(ctx)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.SiteSettings{}, nil
		}
		return domain.SiteSettings{}, err
	}
	return stored, nil
}
