package application

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"polomnik/internal/domain"
	"polomnik/internal/ports"
)

type FavoriteService struct {
	favorites ports.FavoriteRepository
	tours     ports.TourRepository
}

func NewFavoriteService(favorites ports.FavoriteRepository, tours ports.TourRepository) *FavoriteService {
	return &FavoriteService{favorites: favorites, tours: tours}
}

func (s *FavoriteService) ListFavoriteTourIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	return s.favorites.ListFavoriteTourIDs(ctx, userID)
}

func (s *FavoriteService) AddFavorite(ctx context.Context, userID, tourID uuid.UUID) error {
	if _, err := s.tours.GetTour(ctx, tourID); err != nil {
		return err
	}
	return s.favorites.AddFavorite(ctx, userID, tourID)
}

func (s *FavoriteService) RemoveFavorite(ctx context.Context, userID, tourID uuid.UUID) error {
	return s.favorites.RemoveFavorite(ctx, userID, tourID)
}

type SupportService struct {
	support ports.SupportRepository
}

func NewSupportService(support ports.SupportRepository) *SupportService {
	return &SupportService{support: support}
}

const supportWelcomeMessage = "Здравствуйте! Мы получили ваше сообщение. Менеджер ответит в рабочее время — обычно в течение нескольких часов."

func (s *SupportService) GetOrCreateThread(ctx context.Context, userID uuid.UUID) (domain.SupportThread, error) {
	thread, err := s.support.GetOpenThread(ctx, userID)
	if err == nil {
		return thread, nil
	}
	if !errors.Is(err, domain.ErrNotFound) {
		return domain.SupportThread{}, err
	}

	now := time.Now().UTC()
	thread = domain.SupportThread{
		ID:        uuid.New(),
		UserID:    userID,
		Subject:   "Обращение в поддержку",
		Status:    "open",
		CreatedAt: now,
		UpdatedAt: now,
	}
	return s.support.CreateThread(ctx, thread)
}

func (s *SupportService) ListMessages(ctx context.Context, userID, threadID uuid.UUID) ([]domain.SupportMessage, error) {
	thread, err := s.support.GetOpenThread(ctx, userID)
	if err != nil {
		return nil, err
	}
	if thread.ID != threadID {
		return nil, domain.ErrNotFound
	}
	return s.support.ListMessages(ctx, threadID)
}

func (s *SupportService) SendUserMessage(ctx context.Context, userID uuid.UUID, body string) ([]domain.SupportMessage, error) {
	thread, err := s.GetOrCreateThread(ctx, userID)
	if err != nil {
		return nil, err
	}

	message, err := domain.NewSupportMessage(domain.SupportMessage{
		ID:         uuid.New(),
		ThreadID:   thread.ID,
		SenderType: domain.SupportSenderUser,
		Body:       body,
	})
	if err != nil {
		return nil, err
	}

	if _, err := s.support.AddMessage(ctx, message); err != nil {
		return nil, err
	}
	_ = s.support.TouchThread(ctx, thread.ID)

	existing, err := s.support.ListMessages(ctx, thread.ID)
	if err != nil {
		return nil, err
	}
	if len(existing) == 1 {
		welcome, err := domain.NewSupportMessage(domain.SupportMessage{
			ID:         uuid.New(),
			ThreadID:   thread.ID,
			SenderType: domain.SupportSenderStaff,
			Body:       supportWelcomeMessage,
		})
		if err != nil {
			return nil, err
		}
		if _, err := s.support.AddMessage(ctx, welcome); err != nil {
			return nil, err
		}
	}

	return s.support.ListMessages(ctx, thread.ID)
}

func (s *SupportService) GetThread(ctx context.Context, userID uuid.UUID) (domain.SupportThread, []domain.SupportMessage, error) {
	thread, err := s.GetOrCreateThread(ctx, userID)
	if err != nil {
		return domain.SupportThread{}, nil, err
	}
	messages, err := s.support.ListMessages(ctx, thread.ID)
	if err != nil {
		return domain.SupportThread{}, nil, err
	}
	return thread, messages, nil
}
