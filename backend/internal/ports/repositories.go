package ports

import (
	"context"
	"time"

	"github.com/google/uuid"

	"polomnik/internal/domain"
)

type TourFilters struct {
	DateFrom *time.Time
	DateTo   *time.Time
	PriceMin *int
	PriceMax *int
	Location string
	Query    string
	MinSlots *int
	IsActive *bool
	IsHot    *bool
}

type TourList struct {
	Items []domain.Tour
	Meta  PageMeta
}

type TourRepository interface {
	ListTours(ctx context.Context, filters TourFilters, pagination Pagination) (TourList, error)
	GetTour(ctx context.Context, id uuid.UUID) (domain.Tour, error)
	CreateTour(ctx context.Context, tour domain.Tour) (domain.Tour, error)
	UpdateTour(ctx context.Context, tour domain.Tour) (domain.Tour, error)
	DeleteTour(ctx context.Context, id uuid.UUID) error
	ReserveSlots(ctx context.Context, tourID uuid.UUID, peopleCount int) error
	ReleaseSlots(ctx context.Context, tourID uuid.UUID, peopleCount int) error
}

type BookingFilters struct {
	TourID *uuid.UUID
	UserID *uuid.UUID
	Status *domain.BookingStatus
	From   *time.Time
	To     *time.Time
}

type BookingList struct {
	Items []domain.Booking
	Meta  PageMeta
}

type BookingRepository interface {
	CreateBooking(ctx context.Context, booking domain.Booking) (domain.Booking, error)
	GetBooking(ctx context.Context, id uuid.UUID) (domain.Booking, error)
	ListBookings(ctx context.Context, filters BookingFilters, pagination Pagination) (BookingList, error)
	UpdateBookingStatus(ctx context.Context, id uuid.UUID, status domain.BookingStatus) (domain.Booking, error)
	MarkBookingOverbooked(ctx context.Context, id uuid.UUID) (domain.Booking, error)
}

type ReviewFilters struct {
	TourID     *uuid.UUID
	Rating     *int
	IsApproved *bool
}

type ReviewList struct {
	Items []domain.Review
	Meta  PageMeta
}

type ReviewRepository interface {
	ListReviews(ctx context.Context, filters ReviewFilters, pagination Pagination) (ReviewList, error)
	GetReview(ctx context.Context, id uuid.UUID) (domain.Review, error)
	CreateReview(ctx context.Context, review domain.Review) (domain.Review, error)
	ApproveReview(ctx context.Context, id uuid.UUID) (domain.Review, error)
	RejectReview(ctx context.Context, id uuid.UUID) (domain.Review, error)
	DeleteReview(ctx context.Context, id uuid.UUID) error
}

type UserRepository interface {
	CreateUser(ctx context.Context, user domain.User) (domain.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
	GetUserByPhone(ctx context.Context, phone string) (domain.User, error)
	GetUserByOAuth(ctx context.Context, provider, subject string) (domain.User, error)
	UpdateUserProfile(ctx context.Context, user domain.User) (domain.User, error)
}

type FavoriteRepository interface {
	AddFavorite(ctx context.Context, userID, tourID uuid.UUID) error
	RemoveFavorite(ctx context.Context, userID, tourID uuid.UUID) error
	ListFavoriteTourIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
	IsFavorite(ctx context.Context, userID, tourID uuid.UUID) (bool, error)
}

type SupportRepository interface {
	GetOpenThread(ctx context.Context, userID uuid.UUID) (domain.SupportThread, error)
	CreateThread(ctx context.Context, thread domain.SupportThread) (domain.SupportThread, error)
	ListMessages(ctx context.Context, threadID uuid.UUID) ([]domain.SupportMessage, error)
	AddMessage(ctx context.Context, message domain.SupportMessage) (domain.SupportMessage, error)
	TouchThread(ctx context.Context, threadID uuid.UUID) error
}
