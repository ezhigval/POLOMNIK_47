package dto

import (
	"time"

	"palomnik/internal/domain"
)

type TourResponse struct {
	ID          string   `json:"id"`
	Slug        string   `json:"slug"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Price       *int     `json:"price"`
	Currency    string   `json:"currency"`
	DateStart   *string  `json:"date_start"`
	DateEnd     *string  `json:"date_end"`
	SlotsTotal  int      `json:"slots_total"`
	SlotsLeft   int      `json:"slots_left"`
	Location    string   `json:"location"`
	Images      []string `json:"images"`
	IsHot       bool     `json:"is_hot"`
	IsRegular   bool     `json:"is_regular"`
}

type ManagementTourResponse struct {
	TourResponse
	IsActive           bool `json:"is_active"`
	OverbookingEnabled bool `json:"overbooking_enabled"`
}

type TourUpsertRequest struct {
	Slug               string   `json:"slug"`
	Title              string   `json:"title"`
	Description        string   `json:"description"`
	Price              int      `json:"price"`
	Currency           string   `json:"currency"`
	DateStart          string   `json:"date_start"`
	DateEnd            string   `json:"date_end"`
	SlotsTotal         int      `json:"slots_total"`
	SlotsLeft          int      `json:"slots_left"`
	Location           string   `json:"location"`
	Images             []string `json:"images"`
	IsActive           bool     `json:"is_active"`
	IsHot              bool     `json:"is_hot"`
	IsRegular          bool     `json:"is_regular"`
	OverbookingEnabled bool     `json:"overbooking_enabled"`
}

func ToTourResponse(tour domain.Tour) TourResponse {
	resp := TourResponse{
		ID:          tour.ID.String(),
		Slug:        tour.Slug,
		Title:       tour.Title,
		Description: tour.Description,
		Currency:    tour.Currency,
		SlotsTotal:  tour.SlotsTotal,
		SlotsLeft:   tour.SlotsLeft,
		Location:    tour.Location,
		Images:      tour.Images,
		IsHot:       tour.IsHot,
		IsRegular:   tour.IsRegular,
	}
	if !tour.IsRegular {
		price := tour.Price
		resp.Price = &price
		resp.DateStart = formatDatePtr(tour.DateStart)
		resp.DateEnd = formatDatePtr(tour.DateEnd)
	}
	return resp
}

func ToManagementTourResponse(tour domain.Tour) ManagementTourResponse {
	return ManagementTourResponse{
		TourResponse:       ToTourResponse(tour),
		IsActive:           tour.IsActive,
		OverbookingEnabled: tour.OverbookingEnabled,
	}
}

func formatDatePtr(value time.Time) *string {
	if value.IsZero() {
		return nil
	}
	formatted := value.UTC().Format("2006-01-02")
	return &formatted
}
