package postgres

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"testing"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	"palomnik/internal/domain"
	"palomnik/internal/ports"
)

func TestStoreImplementsRepositoryPorts(t *testing.T) {
	var store *Store

	var _ ports.TourRepository = store
	var _ ports.BookingRepository = store
	var _ ports.ReviewRepository = store
	var _ ports.IntegrationReferenceRepository = store
	var _ ports.OutboxRepository = store
}

func TestStoreIntegration(t *testing.T) {
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("set TEST_DATABASE_URL to run postgres repository integration tests")
	}

	ctx := context.Background()
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		t.Fatalf("open postgres: %v", err)
	}
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		t.Fatalf("set goose dialect: %v", err)
	}
	if err := goose.Up(db, "../../../../migrations"); err != nil {
		t.Fatalf("apply migrations: %v", err)
	}
	cleanupPostgres(t, db)

	store := NewStore(db)
	tour := mustTour(t, func(input *domain.NewTourInput) {
		input.ID = uuid.New()
		input.Slug = "postgres-integration"
		input.SlotsTotal = 5
		input.SlotsLeft = 5
	})

	createdTour, err := store.CreateTour(ctx, tour)
	if err != nil {
		t.Fatalf("create tour: %v", err)
	}

	if err := store.ReserveSlots(ctx, createdTour.ID, 2); err != nil {
		t.Fatalf("reserve slots: %v", err)
	}

	updatedTour, err := store.GetTour(ctx, createdTour.ID)
	if err != nil {
		t.Fatalf("get tour: %v", err)
	}
	if updatedTour.SlotsLeft != 3 {
		t.Fatalf("expected 3 slots left, got %d", updatedTour.SlotsLeft)
	}

	booking, err := domain.NewBooking(domain.NewBookingInput{
		ID:          uuid.New(),
		Tour:        updatedTour,
		Name:        "Иван Иванов",
		Phone:       "+79999999999",
		Email:       "mail@test.com",
		PeopleCount: 2,
		Source:      "web",
	})
	if err != nil {
		t.Fatalf("create domain booking: %v", err)
	}

	createdBooking, err := store.CreateBooking(ctx, booking)
	if err != nil {
		t.Fatalf("create booking: %v", err)
	}

	if _, err := store.UpdateBookingStatus(ctx, createdBooking.ID, domain.BookingStatusContacted); err != nil {
		t.Fatalf("update booking status: %v", err)
	}

	review, err := domain.NewReview(domain.NewReviewInput{
		ID:         uuid.New(),
		TourID:     updatedTour.ID,
		ClientName: "Мария",
		Rating:     5,
		Text:       "Хорошая поездка",
	})
	if err != nil {
		t.Fatalf("create domain review: %v", err)
	}

	createdReview, err := store.CreateReview(ctx, review)
	if err != nil {
		t.Fatalf("create review: %v", err)
	}
	approvedReview, err := store.ApproveReview(ctx, createdReview.ID)
	if err != nil {
		t.Fatalf("approve review: %v", err)
	}
	if !approvedReview.IsApproved {
		t.Fatal("expected review approved")
	}

	if _, err := store.GetTour(ctx, uuid.New()); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected not found for missing tour, got %v", err)
	}
}

func cleanupPostgres(t *testing.T, db *sql.DB) {
	t.Helper()
	t.Cleanup(func() {
		if _, err := db.Exec(`TRUNCATE outbox_events, integration_references, reviews, bookings, tours RESTART IDENTITY CASCADE`); err != nil {
			t.Fatalf("cleanup postgres: %v", err)
		}
	})
}

func mustTour(t *testing.T, mutators ...func(*domain.NewTourInput)) domain.Tour {
	t.Helper()
	input := domain.NewTourInput{
		ID:         uuid.New(),
		Slug:       "test-tour",
		Title:      "Test Tour",
		Price:      100,
		Currency:   "RUB",
		DateStart:  testDate(2026, 6, 1),
		DateEnd:    testDate(2026, 6, 2),
		SlotsTotal: 10,
		SlotsLeft:  10,
		IsActive:   true,
	}
	for _, mutate := range mutators {
		mutate(&input)
	}
	tour, err := domain.NewTour(input)
	if err != nil {
		t.Fatalf("create tour: %v", err)
	}
	return tour
}
