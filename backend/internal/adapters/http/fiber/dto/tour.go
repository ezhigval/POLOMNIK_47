package dto

import (
	"time"

	"polomnik/internal/domain"
)

type TourResponse struct {
	ID          string   `json:"id"`
	Slug        string   `json:"slug"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Price       int      `json:"price"`
	Currency    string   `json:"currency"`
	DateStart   string   `json:"date_start"`
	DateEnd     string   `json:"date_end"`
	SlotsTotal  int      `json:"slots_total"`
	SlotsLeft   int      `json:"slots_left"`
	Location    string   `json:"location"`
	Images      []string `json:"images"`
	IsHot       bool     `json:"is_hot"`
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
	OverbookingEnabled bool     `json:"overbooking_enabled"`
}

func ToTourResponse(tour domain.Tour) TourResponse {
	return TourResponse{
		ID:          tour.ID.String(),
		Slug:        tour.Slug,
		Title:       tour.Title,
		Description: tour.Description,
		Price:       tour.Price,
		Currency:    tour.Currency,
		DateStart:   formatDate(tour.DateStart),
		DateEnd:     formatDate(tour.DateEnd),
		SlotsTotal:  tour.SlotsTotal,
		SlotsLeft:   tour.SlotsLeft,
		Location:    tour.Location,
		Images:      tour.Images,
		IsHot:       tour.IsHot,
	}
}

func ToManagementTourResponse(tour domain.Tour) ManagementTourResponse {
	return ManagementTourResponse{
		TourResponse:       ToTourResponse(tour),
		IsActive:           tour.IsActive,
		OverbookingEnabled: tour.OverbookingEnabled,
	}
}

func formatDate(value time.Time) string {
	return value.UTC().Format("2006-01-02")
}
