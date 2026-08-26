package dto

import (
	"time"

	"palomnik/internal/domain"
	"palomnik/internal/ports"
)

type CreateBookingRequest struct {
	TourID              string `json:"tour_id"`
	Name                string `json:"name"`
	Phone               string `json:"phone"`
	Email               string `json:"email"`
	PeopleCount         int    `json:"people_count"`
	Comment             string `json:"comment"`
	Website             string `json:"website"`
	CaptchaToken        string `json:"captcha_token"`
	ConsentPersonalData bool   `json:"consent_personal_data"`
	ConsentMarketing    bool   `json:"consent_marketing"`
}

type CreateBookingResponse struct {
	Status            string `json:"status"`
	BookingID         string `json:"booking_id"`
	BookingStatus     string `json:"booking_status"`
	PaymentStatus     string `json:"payment_status"`
	TotalPrice        int    `json:"total_price"`
	IntegrationStatus string `json:"integration_status"`
}

type ManagementBookingResponse struct {
	ID            string `json:"id"`
	TourID        string `json:"tour_id"`
	Name          string `json:"name"`
	Phone         string `json:"phone"`
	Email         string `json:"email"`
	PeopleCount   int    `json:"people_count"`
	Status        string `json:"status"`
	PaymentStatus string `json:"payment_status"`
	TotalPrice    int    `json:"total_price"`
	Comment       string `json:"comment"`
	Overbooked    bool   `json:"overbooked"`
	Source        string `json:"source"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

type UpdateBookingStatusRequest struct {
	Status string `json:"status"`
}

type UpdateBookingPaymentStatusRequest struct {
	PaymentStatus string `json:"payment_status"`
}

func ToCreateBookingResponse(booking domain.Booking, integrationStatus ports.IntegrationStatus) CreateBookingResponse {
	return CreateBookingResponse{
		Status:            "ok",
		BookingID:         booking.ID.String(),
		BookingStatus:     string(booking.Status),
		PaymentStatus:     string(booking.PaymentStatus),
		TotalPrice:        booking.TotalPrice,
		IntegrationStatus: string(integrationStatus),
	}
}

func ToManagementBookingResponse(booking domain.Booking) ManagementBookingResponse {
	return ManagementBookingResponse{
		ID:            booking.ID.String(),
		TourID:        booking.TourID.String(),
		Name:          booking.Name,
		Phone:         booking.Phone,
		Email:         booking.Email,
		PeopleCount:   booking.PeopleCount,
		Status:        string(booking.Status),
		PaymentStatus: string(booking.PaymentStatus),
		TotalPrice:    booking.TotalPrice,
		Comment:       booking.Comment,
		Overbooked:    booking.Overbooked,
		Source:        booking.Source,
		CreatedAt:     booking.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:     booking.UpdatedAt.UTC().Format(time.RFC3339),
	}
}
