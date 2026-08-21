package dto

import (
	"time"

	"github.com/google/uuid"

	"polomnik/internal/domain"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func ToUserResponse(user domain.User) UserResponse {
	return UserResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		Phone:     user.Phone,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
	}
}

type MyBookingResponse struct {
	ID          string `json:"id"`
	TourID      string `json:"tour_id"`
	Name        string `json:"name"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	PeopleCount int    `json:"people_count"`
	Status      string `json:"status"`
	TotalPrice  int    `json:"total_price"`
	Comment     string `json:"comment"`
	CreatedAt   time.Time `json:"created_at"`
}

func ToMyBookingResponse(booking domain.Booking) MyBookingResponse {
	return MyBookingResponse{
		ID:          booking.ID.String(),
		TourID:      booking.TourID.String(),
		Name:        booking.Name,
		Phone:       booking.Phone,
		Email:       booking.Email,
		PeopleCount: booking.PeopleCount,
		Status:      string(booking.Status),
		TotalPrice:  booking.TotalPrice,
		Comment:     booking.Comment,
		CreatedAt:   booking.CreatedAt,
	}
}

func ToMyBookingResponses(bookings []domain.Booking) []MyBookingResponse {
	items := make([]MyBookingResponse, 0, len(bookings))
	for _, booking := range bookings {
		items = append(items, ToMyBookingResponse(booking))
	}
	return items
}

func UserIDPtr(id uuid.UUID) *uuid.UUID {
	if id == uuid.Nil {
		return nil
	}
	return &id
}
