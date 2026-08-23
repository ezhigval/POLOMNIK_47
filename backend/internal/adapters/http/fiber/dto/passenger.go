package dto

import (
	"time"

	"palomnik/internal/domain"
)

type PassengerRequest struct {
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	BirthDate string `json:"birth_date"`
	Passport  string `json:"passport"`
	Website   string `json:"website"`
}

type PassengerResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	BirthDate string    `json:"birth_date"`
	Passport  string    `json:"passport"`
	CreatedAt time.Time `json:"created_at"`
}

func ToPassengerResponse(passenger domain.Passenger) PassengerResponse {
	return PassengerResponse{
		ID:        passenger.ID.String(),
		Name:      passenger.Name,
		Phone:     passenger.Phone,
		BirthDate: passenger.BirthDate.Format("2006-01-02"),
		Passport:  passenger.Passport,
		CreatedAt: passenger.CreatedAt,
	}
}

func ToPassengerResponses(items []domain.Passenger) []PassengerResponse {
	out := make([]PassengerResponse, 0, len(items))
	for _, item := range items {
		out = append(out, ToPassengerResponse(item))
	}
	return out
}
