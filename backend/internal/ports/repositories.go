package ports

import (
	"context"
	"time"

	"github.com/google/uuid"

	"palomnik/internal/domain"
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
	TourID            *uuid.UUID
	Rating            *int
	IsApproved        *bool
	AllowDistribution *bool
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
	UpdateReview(ctx context.Context, review domain.Review) (domain.Review, error)
	DeleteReview(ctx context.Context, id uuid.UUID) error
}

type UserRepository interface {
	CreateUser(ctx context.Context, user domain.User) (domain.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
	GetUserByPhone(ctx context.Context, phone string) (domain.User, error)
	GetUserByOAuth(ctx context.Context, provider, subject string) (domain.User, error)
	GetIdentity(ctx context.Context, provider, subject string) (domain.UserIdentity, error)
	ListIdentities(ctx context.Context, userID uuid.UUID) ([]domain.UserIdentity, error)
	AddIdentity(ctx context.Context, identity domain.UserIdentity) error
	UpdateUserProfile(ctx context.Context, user domain.User) (domain.User, error)
	UpdateUserPassword(ctx context.Context, userID uuid.UUID, passwordHash string) error
	MergeAccountInto(ctx context.Context, targetID, sourceID uuid.UUID) error
}

type PassengerRepository interface {
	ListPassengers(ctx context.Context, userID uuid.UUID) ([]domain.Passenger, error)
	GetPassenger(ctx context.Context, userID, id uuid.UUID) (domain.Passenger, error)
	CreatePassenger(ctx context.Context, passenger domain.Passenger) (domain.Passenger, error)
	UpdatePassenger(ctx context.Context, passenger domain.Passenger) (domain.Passenger, error)
	DeletePassenger(ctx context.Context, userID, id uuid.UUID) error
}

type UserPhotoRepository interface {
	ListUserPhotos(ctx context.Context, userID uuid.UUID) ([]domain.UserPhoto, error)
	CreateUserPhoto(ctx context.Context, photo domain.UserPhoto) (domain.UserPhoto, error)
	DeleteUserPhoto(ctx context.Context, userID, id uuid.UUID) error
}

type FavoriteRepository interface {
	AddFavorite(ctx context.Context, userID, tourID uuid.UUID) error
	RemoveFavorite(ctx context.Context, userID, tourID uuid.UUID) error
	ListFavoriteTourIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
	IsFavorite(ctx context.Context, userID, tourID uuid.UUID) (bool, error)
}

type SupportRepository interface {
	GetOpenThread(ctx context.Context, userID uuid.UUID) (domain.SupportThread, error)
	GetThreadByID(ctx context.Context, threadID uuid.UUID) (domain.SupportThread, error)
	ListThreads(ctx context.Context) ([]domain.SupportThread, error)
	CreateThread(ctx context.Context, thread domain.SupportThread) (domain.SupportThread, error)
	ListMessages(ctx context.Context, threadID uuid.UUID) ([]domain.SupportMessage, error)
	AddMessage(ctx context.Context, message domain.SupportMessage) (domain.SupportMessage, error)
	TouchThread(ctx context.Context, threadID uuid.UUID) error
}

type NewsFilters struct {
	PublishedOnly bool
}

type NewsList struct {
	Items []domain.NewsArticle
	Meta  PageMeta
}

type NewsRepository interface {
	ListNews(ctx context.Context, filters NewsFilters, pagination Pagination) (NewsList, error)
	GetNews(ctx context.Context, id uuid.UUID) (domain.NewsArticle, error)
	GetNewsBySlug(ctx context.Context, slug string) (domain.NewsArticle, error)
	CreateNews(ctx context.Context, article domain.NewsArticle) (domain.NewsArticle, error)
	UpdateNews(ctx context.Context, article domain.NewsArticle) (domain.NewsArticle, error)
	DeleteNews(ctx context.Context, id uuid.UUID) error
}

type SMMPostList struct {
	Items []domain.SMMPost
	Meta  PageMeta
}

type SMMPostRepository interface {
	CreateSMMPost(ctx context.Context, post domain.SMMPost) (domain.SMMPost, error)
	GetSMMPost(ctx context.Context, id uuid.UUID) (domain.SMMPost, error)
	ListSMMPosts(ctx context.Context, pagination Pagination) (SMMPostList, error)
	ListDueSMMPosts(ctx context.Context, now time.Time) ([]domain.SMMPost, error)
	SaveSMMPost(ctx context.Context, post domain.SMMPost) (domain.SMMPost, error)
	DeleteSMMPost(ctx context.Context, id uuid.UUID) error
}

type TelegramRecipientsRepository interface {
	GetTelegramRecipients(ctx context.Context) (domain.TelegramRecipients, error)
	UpsertTelegramRecipients(ctx context.Context, settings domain.TelegramRecipients) (domain.TelegramRecipients, error)
}

type TelegramChatMapRepository interface {
	UpsertTelegramChatBinding(ctx context.Context, binding domain.TelegramChatBinding) error
	ListTelegramChatBindings(ctx context.Context, usernames []string) (map[string]string, error)
}

type NotificationRoutingRepository interface {
	GetNotificationRouting(ctx context.Context) (domain.NotificationRouting, error)
	UpsertNotificationRouting(ctx context.Context, routing domain.NotificationRouting) (domain.NotificationRouting, error)
}

type SiteSettingsRepository interface {
	GetSiteSettings(ctx context.Context) (domain.SiteSettings, error)
	UpsertSiteSettings(ctx context.Context, settings domain.SiteSettings) (domain.SiteSettings, error)
}

type AdminRoleRepository interface {
	ListAdminRoles(ctx context.Context) ([]domain.AdminRole, error)
	GetAdminRole(ctx context.Context, id uuid.UUID) (domain.AdminRole, error)
	GetAdminRoleByName(ctx context.Context, name string) (domain.AdminRole, error)
	CreateAdminRole(ctx context.Context, role domain.AdminRole) (domain.AdminRole, error)
	UpdateAdminRole(ctx context.Context, role domain.AdminRole) (domain.AdminRole, error)
	DeleteAdminRole(ctx context.Context, id uuid.UUID) error
	AssignUserToRole(ctx context.Context, assignment domain.AdminRoleAssignment) error
	UnassignUserFromRole(ctx context.Context, roleID, userID uuid.UUID) error
	ListRoleAssignments(ctx context.Context, roleID uuid.UUID) ([]domain.AdminRoleAssignment, error)
}
