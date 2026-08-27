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
	IsHot              bool `json:"is_hot"`
	IsBurning          bool `json:"is_burning"`
	IsRegular          bool `json:"is_regular"`
	OverbookingEnabled bool `json:"overbooking_enabled"`
	OriginalPrice      *int `json:"original_price,omitempty"`
}

type ManagementTourResponse struct {
	TourResponse
	IsActive           bool `json:"is_active"`
	OverbookingEnabled bool `json:"overbooking_enabled"`
	HotDiscountPercent int  `json:"hot_discount_percent"`
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
	HotDiscountPercent int      `json:"hot_discount_percent"`
}

func ToTourResponse(tour domain.Tour) TourResponse {
	return ToPublicTourResponse(tour, domain.TourCatalogContext{})
}

func ToPublicTourResponse(tour domain.Tour, catalog domain.TourCatalogContext) TourResponse {
	resp := TourResponse{
		ID:          tour.ID.String(),
		Slug:        tour.Slug,
		Title:       tour.Title,
		Description: tour.Description,
		Currency:    tour.Currency,
		SlotsTotal:  tour.SlotsTotal,
		SlotsLeft:   tour.RemainingSlots(),
		Location:    tour.Location,
		Images:      tour.Images,
		IsHot:              tour.IsHot,
		IsBurning:          tour.IsBurningIn(catalog),
		IsRegular:          tour.IsRegular,
		OverbookingEnabled: tour.OverbookingEnabled,
	}
	if tour.HasPublicPrice() {
		unit := tour.UnitPriceIn(catalog)
		price := unit
		resp.Price = &price
		if resp.IsBurning && tour.HotDiscountPercent > 0 && unit < tour.Price {
			original := tour.Price
			resp.OriginalPrice = &original
		}
	}
	if !tour.IsRegular {
		resp.DateStart = formatDatePtr(tour.DateStart)
		resp.DateEnd = formatDatePtr(tour.DateEnd)
	}
	return resp
}

func ToManagementTourResponse(tour domain.Tour) ManagementTourResponse {
	resp := ToTourResponse(tour)
	resp.SlotsLeft = tour.SlotsLeft
	return ManagementTourResponse{
		TourResponse:       resp,
		IsActive:           tour.IsActive,
		OverbookingEnabled: tour.OverbookingEnabled,
		HotDiscountPercent: tour.HotDiscountPercent,
	}
}

func formatDatePtr(value time.Time) *string {
	if value.IsZero() {
		return nil
	}
	formatted := value.UTC().Format("2006-01-02")
	return &formatted
}
