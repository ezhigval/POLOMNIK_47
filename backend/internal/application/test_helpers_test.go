package application

import (
	"time"

	"github.com/google/uuid"

	"polomnik/internal/domain"
)

func testUUID(value string) uuid.UUID {
	id, err := uuid.Parse(value)
	if err != nil {
		panic(err)
	}
	return id
}

func testDate(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

func testTour(mutators ...func(*domain.NewTourInput)) domain.Tour {
	input := domain.NewTourInput{
		ID:         testUUID("11111111-1111-1111-1111-111111111111"),
		Slug:       "test-tour",
		Title:      "Test Tour",
		Price:      10000,
		Currency:   "RUB",
		DateStart:  testDate(2026, 8, 1),
		DateEnd:    testDate(2026, 8, 5),
		SlotsTotal: 10,
		SlotsLeft:  10,
		Location:   "Moscow",
		IsActive:   true,
		Now:        testDate(2026, 5, 20),
	}
	for _, mutate := range mutators {
		mutate(&input)
	}

	tour, err := domain.NewTour(input)
	if err != nil {
		panic(err)
	}
	return tour
}
