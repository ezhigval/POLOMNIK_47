package memory

import (
	"context"
	"strings"
	"sync"

	"github.com/google/uuid"

	"polomnik/internal/domain"
	"polomnik/internal/ports"
)

type Store struct {
	mu                 sync.RWMutex
	tours              map[uuid.UUID]domain.Tour
	bookings           map[uuid.UUID]domain.Booking
	reviews            map[uuid.UUID]domain.Review
	users              map[uuid.UUID]domain.User
	favorites          map[string]domain.Favorite
	supportThreads     map[uuid.UUID]domain.SupportThread
	supportMessages    map[uuid.UUID]domain.SupportMessage
	integrationRefs    map[string]domain.IntegrationReference
	outboxEvents       map[uuid.UUID]domain.OutboxEvent
	cmsPages           map[uuid.UUID]domain.Page
	cmsBlocks          map[uuid.UUID]domain.Block
	news               map[uuid.UUID]domain.NewsArticle
	telegramRecipients *domain.TelegramRecipients
	telegramChats      map[string]domain.TelegramChatBinding
}

func (s *Store) WithinTransaction(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

func NewStore() *Store {
	return &Store{
		tours:           make(map[uuid.UUID]domain.Tour),
		bookings:        make(map[uuid.UUID]domain.Booking),
		reviews:         make(map[uuid.UUID]domain.Review),
		users:           make(map[uuid.UUID]domain.User),
		favorites:       make(map[string]domain.Favorite),
		supportThreads:  make(map[uuid.UUID]domain.SupportThread),
		supportMessages: make(map[uuid.UUID]domain.SupportMessage),
		integrationRefs: make(map[string]domain.IntegrationReference),
		outboxEvents:    make(map[uuid.UUID]domain.OutboxEvent),
		cmsPages:        make(map[uuid.UUID]domain.Page),
		cmsBlocks:       make(map[uuid.UUID]domain.Block),
		news:            make(map[uuid.UUID]domain.NewsArticle),
		telegramChats:   make(map[string]domain.TelegramChatBinding),
	}
}

func (s *Store) ListTours(_ context.Context, filters ports.TourFilters, pagination ports.Pagination) (ports.TourList, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	items := make([]domain.Tour, 0, len(s.tours))
	for _, tour := range s.tours {
		if !matchesTourFilters(tour, filters) {
			continue
		}
		items = append(items, cloneTour(tour))
	}

	pageItems, meta := page(items, pagination)
	return ports.TourList{Items: pageItems, Meta: meta}, nil
}

func (s *Store) GetTour(_ context.Context, id uuid.UUID) (domain.Tour, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tour, ok := s.tours[id]
	if !ok {
		return domain.Tour{}, domain.ErrNotFound
	}
	return cloneTour(tour), nil
}

func (s *Store) CreateTour(_ context.Context, tour domain.Tour) (domain.Tour, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.tours[tour.ID] = cloneTour(tour)
	return cloneTour(tour), nil
}

func (s *Store) UpdateTour(_ context.Context, tour domain.Tour) (domain.Tour, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.tours[tour.ID]; !ok {
		return domain.Tour{}, domain.ErrNotFound
	}
	s.tours[tour.ID] = cloneTour(tour)
	return cloneTour(tour), nil
}

func (s *Store) DeleteTour(_ context.Context, id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.tours[id]; !ok {
		return domain.ErrNotFound
	}
	delete(s.tours, id)
	return nil
}

func (s *Store) ReserveSlots(_ context.Context, tourID uuid.UUID, peopleCount int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	tour, ok := s.tours[tourID]
	if !ok {
		return domain.ErrNotFound
	}
	if err := tour.ReserveSlots(peopleCount); err != nil {
		return err
	}
	s.tours[tourID] = tour
	return nil
}

func (s *Store) ReleaseSlots(_ context.Context, tourID uuid.UUID, peopleCount int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	tour, ok := s.tours[tourID]
	if !ok {
		return domain.ErrNotFound
	}
	if err := tour.ReleaseSlots(peopleCount); err != nil {
		return err
	}
	s.tours[tourID] = tour
	return nil
}

func (s *Store) CreateBooking(_ context.Context, booking domain.Booking) (domain.Booking, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.bookings[booking.ID] = booking
	return booking, nil
}

func (s *Store) GetBooking(_ context.Context, id uuid.UUID) (domain.Booking, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	booking, ok := s.bookings[id]
	if !ok {
		return domain.Booking{}, domain.ErrNotFound
	}
	return booking, nil
}

func (s *Store) ListBookings(_ context.Context, filters ports.BookingFilters, pagination ports.Pagination) (ports.BookingList, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	items := make([]domain.Booking, 0, len(s.bookings))
	for _, booking := range s.bookings {
		if !matchesBookingFilters(booking, filters) {
			continue
		}
		items = append(items, booking)
	}

	pageItems, meta := page(items, pagination)
	return ports.BookingList{Items: pageItems, Meta: meta}, nil
}

func (s *Store) UpdateBookingStatus(_ context.Context, id uuid.UUID, status domain.BookingStatus) (domain.Booking, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	booking, ok := s.bookings[id]
	if !ok {
		return domain.Booking{}, domain.ErrNotFound
	}
	if err := booking.ChangeStatus(status); err != nil {
		return domain.Booking{}, err
	}
	s.bookings[id] = booking
	return booking, nil
}

func (s *Store) MarkBookingOverbooked(_ context.Context, id uuid.UUID) (domain.Booking, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	booking, ok := s.bookings[id]
	if !ok {
		return domain.Booking{}, domain.ErrNotFound
	}
	booking.Overbooked = true
	s.bookings[id] = booking
	return booking, nil
}

func (s *Store) ListReviews(_ context.Context, filters ports.ReviewFilters, pagination ports.Pagination) (ports.ReviewList, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	items := make([]domain.Review, 0, len(s.reviews))
	for _, review := range s.reviews {
		if !matchesReviewFilters(review, filters) {
			continue
		}
		items = append(items, review)
	}

	pageItems, meta := page(items, pagination)
	return ports.ReviewList{Items: pageItems, Meta: meta}, nil
}

func (s *Store) GetReview(_ context.Context, id uuid.UUID) (domain.Review, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	review, ok := s.reviews[id]
	if !ok {
		return domain.Review{}, domain.ErrNotFound
	}
	return review, nil
}

func (s *Store) CreateReview(_ context.Context, review domain.Review) (domain.Review, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.reviews[review.ID] = review
	return review, nil
}

func (s *Store) ApproveReview(_ context.Context, id uuid.UUID) (domain.Review, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	review, ok := s.reviews[id]
	if !ok {
		return domain.Review{}, domain.ErrNotFound
	}
	review.Approve()
	s.reviews[id] = review
	return review, nil
}

func (s *Store) RejectReview(_ context.Context, id uuid.UUID) (domain.Review, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	review, ok := s.reviews[id]
	if !ok {
		return domain.Review{}, domain.ErrNotFound
	}
	review.Reject()
	s.reviews[id] = review
	return review, nil
}

func (s *Store) DeleteReview(_ context.Context, id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.reviews[id]; !ok {
		return domain.ErrNotFound
	}
	delete(s.reviews, id)
	return nil
}

func (s *Store) UpdateReview(_ context.Context, review domain.Review) (domain.Review, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.reviews[review.ID]; !ok {
		return domain.Review{}, domain.ErrNotFound
	}
	s.reviews[review.ID] = review
	return review, nil
}

func matchesTourFilters(tour domain.Tour, filters ports.TourFilters) bool {
	if filters.IsActive != nil && tour.IsActive != *filters.IsActive {
		return false
	}
	if filters.IsHot != nil && tour.IsHot != *filters.IsHot {
		return false
	}
	if filters.PriceMin != nil && tour.Price < *filters.PriceMin {
		return false
	}
	if filters.PriceMax != nil && tour.Price > *filters.PriceMax {
		return false
	}
	if filters.DateFrom != nil && tour.DateEnd.Before(*filters.DateFrom) {
		return false
	}
	if filters.DateTo != nil && tour.DateStart.After(*filters.DateTo) {
		return false
	}
	if strings.TrimSpace(filters.Location) != "" {
		needle := strings.ToLower(strings.TrimSpace(filters.Location))
		if !strings.Contains(strings.ToLower(tour.Location), needle) {
			return false
		}
	}
	if strings.TrimSpace(filters.Query) != "" {
		needle := strings.ToLower(strings.TrimSpace(filters.Query))
		haystack := strings.ToLower(tour.Title + " " + tour.Location + " " + tour.Slug)
		if !strings.Contains(haystack, needle) {
			return false
		}
	}
	if filters.MinSlots != nil && tour.SlotsLeft < *filters.MinSlots {
		return false
	}
	return true
}

func matchesBookingFilters(booking domain.Booking, filters ports.BookingFilters) bool {
	if filters.TourID != nil && booking.TourID != *filters.TourID {
		return false
	}
	if filters.Status != nil && booking.Status != *filters.Status {
		return false
	}
	if filters.From != nil && booking.CreatedAt.Before(*filters.From) {
		return false
	}
	if filters.To != nil && booking.CreatedAt.After(*filters.To) {
		return false
	}
	if filters.UserID != nil {
		if booking.UserID == nil || *booking.UserID != *filters.UserID {
			return false
		}
	}
	return true
}

func matchesReviewFilters(review domain.Review, filters ports.ReviewFilters) bool {
	if filters.TourID != nil && review.TourID != *filters.TourID {
		return false
	}
	if filters.Rating != nil && review.Rating != *filters.Rating {
		return false
	}
	if filters.IsApproved != nil && review.IsApproved != *filters.IsApproved {
		return false
	}
	return true
}

func page[T any](items []T, pagination ports.Pagination) ([]T, ports.PageMeta) {
	pagination = ports.NormalizePagination(pagination.Page, pagination.Limit)
	total := len(items)
	start := (pagination.Page - 1) * pagination.Limit
	if start >= total {
		return []T{}, ports.PageMeta{
			Page:    pagination.Page,
			Limit:   pagination.Limit,
			Total:   total,
			HasNext: false,
		}
	}

	end := start + pagination.Limit
	if end > total {
		end = total
	}

	return items[start:end], ports.PageMeta{
		Page:    pagination.Page,
		Limit:   pagination.Limit,
		Total:   total,
		HasNext: end < total,
	}
}

func cloneTour(tour domain.Tour) domain.Tour {
	tour.Images = append([]string(nil), tour.Images...)
	return tour
}
