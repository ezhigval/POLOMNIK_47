package dto

import (
	"time"

	"github.com/google/uuid"

	"palomnik/internal/domain"
)

type RegisterRequest struct {
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Name         string `json:"name"`
	Password     string `json:"password"`
	PhoneCheckID string `json:"phone_check_id"`
	Website      string `json:"website"`
	CaptchaToken string `json:"captcha_token"`
}

type LoginRequest struct {
	Login        string `json:"login"`
	Password     string `json:"password"`
	Website      string `json:"website"`
	CaptchaToken string `json:"captcha_token"`
}

type ForgotPasswordRequest struct {
	Email        string `json:"email"`
	Website      string `json:"website"`
	CaptchaToken string `json:"captcha_token"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

type PhoneStartRequest struct {
	Phone string `json:"phone"`
}

type PhoneCheckRequest struct {
	CheckID string `json:"check_id"`
}

type AuthMethodsResponse struct {
	Password  bool                     `json:"password"`
	PhoneCall AuthMethodStatusResponse `json:"phone_call"`
	Yandex    AuthMethodStatusResponse `json:"yandex"`
	VK        AuthMethodStatusResponse `json:"vk"`
	Max       AuthMethodStatusResponse `json:"max"`
	Telegram  AuthMethodStatusResponse `json:"telegram"`
	Mail      AuthMethodStatusResponse `json:"mail"`
	Captcha   CaptchaStatusResponse    `json:"captcha"`
}

type CaptchaStatusResponse struct {
	Available bool   `json:"available"`
	Provider  string `json:"provider"`
	ClientKey string `json:"client_key,omitempty"`
}

type AuthMethodStatusResponse struct {
	Available bool   `json:"available"`
	Message   string `json:"message,omitempty"`
	Username  string `json:"username,omitempty"`
}

type PhoneStartResponse struct {
	CheckID         string `json:"check_id"`
	CallPhone       string `json:"call_phone"`
	CallPhonePretty string `json:"call_phone_pretty"`
	ExpiresIn       int    `json:"expires_in"`
}

type PhoneStatusResponse struct {
	Status string `json:"status"`
}

type AuthResponse struct {
	Token      string       `json:"token"`
	User       UserResponse `json:"user"`
	Linked     bool         `json:"linked,omitempty"`
	Merged     bool         `json:"merged,omitempty"`
	KeptFields []string     `json:"kept_fields,omitempty"`
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

func ToOAuthAuthResponse(token string, user domain.User, linked, merged bool, conflicts []domain.ProfileConflict) AuthResponse {
	resp := AuthResponse{
		Token:  token,
		User:   ToUserResponse(user),
		Linked: linked,
		Merged: merged,
	}
	if len(conflicts) == 0 {
		return resp
	}
	resp.KeptFields = make([]string, 0, len(conflicts))
	for _, conflict := range conflicts {
		if conflict.Field != "" {
			resp.KeptFields = append(resp.KeptFields, conflict.Field)
		}
	}
	return resp
}

type UserIdentityResponse struct {
	Provider  string    `json:"provider"`
	Subject   string    `json:"subject"`
	CreatedAt time.Time `json:"created_at"`
}

func ToUserIdentityResponses(identities []domain.UserIdentity) []UserIdentityResponse {
	items := make([]UserIdentityResponse, 0, len(identities))
	for _, identity := range identities {
		items = append(items, UserIdentityResponse{
			Provider:  identity.Provider,
			Subject:   identity.Subject,
			CreatedAt: identity.CreatedAt,
		})
	}
	return items
}

type MyBookingResponse struct {
	ID          string    `json:"id"`
	TourID      string    `json:"tour_id"`
	Name        string    `json:"name"`
	Phone       string    `json:"phone"`
	Email       string    `json:"email"`
	PeopleCount int       `json:"people_count"`
	Status      string    `json:"status"`
	TotalPrice  int       `json:"total_price"`
	Comment     string    `json:"comment"`
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
