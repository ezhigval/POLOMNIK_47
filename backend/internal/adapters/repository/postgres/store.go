package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	_ "github.com/jackc/pgx/v5/stdlib"

	"polomnik/internal/domain"
	"polomnik/internal/ports"
)

type Store struct {
	db *sql.DB
}

func Open(ctx context.Context, databaseURL string) (*Store, error) {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("open postgres: %w", err)
	}
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(time.Minute)

	return NewStore(db), nil
}

func (s *Store) Ping(ctx context.Context) error {
	if s == nil || s.db == nil {
		return fmt.Errorf("database is not configured")
	}
	return s.db.PingContext(ctx)
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) ListTours(ctx context.Context, filters ports.TourFilters, pagination ports.Pagination) (ports.TourList, error) {
	pagination = ports.NormalizePagination(pagination.Page, pagination.Limit)
	args := tourFilterArgs(filters)

	total, err := count(ctx, s.db, "SELECT COUNT(*) FROM tours WHERE "+tourWhereClause, args...)
	if err != nil {
		return ports.TourList{}, err
	}

	rows, err := s.db.QueryContext(ctx, `
SELECT id, slug, title, description, price, currency, date_start, date_end,
       slots_total, slots_left, location, images, is_active, is_hot,
       overbooking_enabled, created_at, updated_at
FROM tours
WHERE `+tourWhereClause+`
ORDER BY date_start ASC, id ASC
LIMIT $10 OFFSET $11
`, append(args, pagination.Limit, offset(pagination))...)
	if err != nil {
		return ports.TourList{}, fmt.Errorf("list tours: %w", err)
	}
	defer rows.Close()

	items, err := scanTours(rows)
	if err != nil {
		return ports.TourList{}, err
	}

	return ports.TourList{
		Items: items,
		Meta:  pageMeta(pagination, total),
	}, nil
}

func (s *Store) GetTour(ctx context.Context, id uuid.UUID) (domain.Tour, error) {
	row := s.db.QueryRowContext(ctx, `
SELECT id, slug, title, description, price, currency, date_start, date_end,
       slots_total, slots_left, location, images, is_active, is_hot,
       overbooking_enabled, created_at, updated_at
FROM tours
WHERE id = $1
`, id)
	return scanTour(row)
}

func (s *Store) CreateTour(ctx context.Context, tour domain.Tour) (domain.Tour, error) {
	row := s.db.QueryRowContext(ctx, `
INSERT INTO tours (
    id, slug, title, description, price, currency, date_start, date_end,
    slots_total, slots_left, location, images, is_active, is_hot,
    overbooking_enabled, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8,
    $9, $10, $11, $12, $13, $14, $15, $16, $17
)
RETURNING id, slug, title, description, price, currency, date_start, date_end,
          slots_total, slots_left, location, images, is_active, is_hot,
          overbooking_enabled, created_at, updated_at
`, tour.ID, tour.Slug, tour.Title, tour.Description, tour.Price, tour.Currency,
		tour.DateStart, tour.DateEnd, tour.SlotsTotal, tour.SlotsLeft, tour.Location,
		pq.Array(tour.Images), tour.IsActive, tour.IsHot, tour.OverbookingEnabled, tour.CreatedAt, tour.UpdatedAt)
	return scanTour(row)
}

func (s *Store) UpdateTour(ctx context.Context, tour domain.Tour) (domain.Tour, error) {
	row := s.db.QueryRowContext(ctx, `
UPDATE tours
SET slug = $2,
    title = $3,
    description = $4,
    price = $5,
    currency = $6,
    date_start = $7,
    date_end = $8,
    slots_total = $9,
    slots_left = $10,
    location = $11,
    images = $12,
    is_active = $13,
    is_hot = $14,
    overbooking_enabled = $15,
    updated_at = $16
WHERE id = $1
RETURNING id, slug, title, description, price, currency, date_start, date_end,
          slots_total, slots_left, location, images, is_active, is_hot,
          overbooking_enabled, created_at, updated_at
`, tour.ID, tour.Slug, tour.Title, tour.Description, tour.Price, tour.Currency,
		tour.DateStart, tour.DateEnd, tour.SlotsTotal, tour.SlotsLeft, tour.Location,
		pq.Array(tour.Images), tour.IsActive, tour.IsHot, tour.OverbookingEnabled, tour.UpdatedAt)
	return scanTour(row)
}

func (s *Store) DeleteTour(ctx context.Context, id uuid.UUID) error {
	result, err := s.db.ExecContext(ctx, `DELETE FROM tours WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete tour: %w", err)
	}
	return requireAffected(result)
}

func (s *Store) ReserveSlots(ctx context.Context, tourID uuid.UUID, peopleCount int) error {
	if peopleCount <= 0 {
		return domain.ErrInvalidPeopleCount
	}

	result, err := s.db.ExecContext(ctx, `
UPDATE tours
SET slots_left = slots_left - $2,
    updated_at = NOW()
WHERE id = $1 AND slots_left >= $2
`, tourID, peopleCount)
	if err != nil {
		return fmt.Errorf("reserve slots: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("reserve slots rows affected: %w", err)
	}
	if affected > 0 {
		return nil
	}

	tour, err := s.GetTour(ctx, tourID)
	if err != nil {
		return err
	}
	if tour.OverbookingEnabled {
		_, err := s.db.ExecContext(ctx, `UPDATE tours SET updated_at = NOW() WHERE id = $1`, tourID)
		if err != nil {
			return fmt.Errorf("touch overbooked tour: %w", err)
		}
		return nil
	}

	return domain.ErrInsufficientSlots
}

func (s *Store) ReleaseSlots(ctx context.Context, tourID uuid.UUID, peopleCount int) error {
	if peopleCount <= 0 {
		return domain.ErrInvalidPeopleCount
	}

	result, err := s.db.ExecContext(ctx, `
UPDATE tours
SET slots_left = LEAST(slots_total, slots_left + $2),
    updated_at = NOW()
WHERE id = $1
`, tourID, peopleCount)
	if err != nil {
		return fmt.Errorf("release slots: %w", err)
	}
	return requireAffected(result)
}

func (s *Store) CreateBooking(ctx context.Context, booking domain.Booking) (domain.Booking, error) {
	row := s.db.QueryRowContext(ctx, `
INSERT INTO bookings (
    id, tour_id, user_id, name, phone, email, people_count, status, total_price,
    comment, overbooked, source, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9,
    $10, $11, $12, $13, $14
)
RETURNING id, tour_id, user_id, name, phone, email, people_count, status, total_price,
          comment, overbooked, source, created_at, updated_at
`, booking.ID, booking.TourID, uuidValue(booking.UserID), booking.Name, booking.Phone, booking.Email,
		booking.PeopleCount, booking.Status, booking.TotalPrice, booking.Comment,
		booking.Overbooked, booking.Source, booking.CreatedAt, booking.UpdatedAt)
	return scanBooking(row)
}

func (s *Store) GetBooking(ctx context.Context, id uuid.UUID) (domain.Booking, error) {
	row := s.db.QueryRowContext(ctx, `
SELECT id, tour_id, user_id, name, phone, email, people_count, status, total_price,
       comment, overbooked, source, created_at, updated_at
FROM bookings
WHERE id = $1
`, id)
	return scanBooking(row)
}

func (s *Store) ListBookings(ctx context.Context, filters ports.BookingFilters, pagination ports.Pagination) (ports.BookingList, error) {
	pagination = ports.NormalizePagination(pagination.Page, pagination.Limit)
	args := bookingFilterArgs(filters)

	total, err := count(ctx, s.db, "SELECT COUNT(*) FROM bookings WHERE "+bookingWhereClause, args...)
	if err != nil {
		return ports.BookingList{}, err
	}

	rows, err := s.db.QueryContext(ctx, `
SELECT id, tour_id, user_id, name, phone, email, people_count, status, total_price,
       comment, overbooked, source, created_at, updated_at
FROM bookings
WHERE `+bookingWhereClause+`
ORDER BY created_at DESC, id ASC
LIMIT $6 OFFSET $7
`, append(args, pagination.Limit, offset(pagination))...)
	if err != nil {
		return ports.BookingList{}, fmt.Errorf("list bookings: %w", err)
	}
	defer rows.Close()

	items, err := scanBookings(rows)
	if err != nil {
		return ports.BookingList{}, err
	}

	return ports.BookingList{Items: items, Meta: pageMeta(pagination, total)}, nil
}

func (s *Store) UpdateBookingStatus(ctx context.Context, id uuid.UUID, status domain.BookingStatus) (domain.Booking, error) {
	booking, err := s.GetBooking(ctx, id)
	if err != nil {
		return domain.Booking{}, err
	}
	if err := booking.ChangeStatus(status); err != nil {
		return domain.Booking{}, err
	}

	row := s.db.QueryRowContext(ctx, `
UPDATE bookings
SET status = $2,
    updated_at = $3
WHERE id = $1
RETURNING id, tour_id, user_id, name, phone, email, people_count, status, total_price,
          comment, overbooked, source, created_at, updated_at
`, id, booking.Status, booking.UpdatedAt)
	return scanBooking(row)
}

func (s *Store) MarkBookingOverbooked(ctx context.Context, id uuid.UUID) (domain.Booking, error) {
	row := s.db.QueryRowContext(ctx, `
UPDATE bookings
SET overbooked = TRUE,
    updated_at = NOW()
WHERE id = $1
RETURNING id, tour_id, user_id, name, phone, email, people_count, status, total_price,
          comment, overbooked, source, created_at, updated_at
`, id)
	return scanBooking(row)
}

func (s *Store) ListReviews(ctx context.Context, filters ports.ReviewFilters, pagination ports.Pagination) (ports.ReviewList, error) {
	pagination = ports.NormalizePagination(pagination.Page, pagination.Limit)
	args := reviewFilterArgs(filters)

	total, err := count(ctx, s.db, "SELECT COUNT(*) FROM reviews WHERE "+reviewWhereClause, args...)
	if err != nil {
		return ports.ReviewList{}, err
	}

	rows, err := s.db.QueryContext(ctx, `
SELECT id, tour_id, client_name, rating, text, is_approved, created_at, updated_at
FROM reviews
WHERE `+reviewWhereClause+`
ORDER BY created_at DESC, id ASC
LIMIT $4 OFFSET $5
`, append(args, pagination.Limit, offset(pagination))...)
	if err != nil {
		return ports.ReviewList{}, fmt.Errorf("list reviews: %w", err)
	}
	defer rows.Close()

	items, err := scanReviews(rows)
	if err != nil {
		return ports.ReviewList{}, err
	}

	return ports.ReviewList{Items: items, Meta: pageMeta(pagination, total)}, nil
}

func (s *Store) GetReview(ctx context.Context, id uuid.UUID) (domain.Review, error) {
	row := s.db.QueryRowContext(ctx, `
SELECT id, tour_id, client_name, rating, text, is_approved, created_at, updated_at
FROM reviews
WHERE id = $1
`, id)
	return scanReview(row)
}

func (s *Store) CreateReview(ctx context.Context, review domain.Review) (domain.Review, error) {
	row := s.db.QueryRowContext(ctx, `
INSERT INTO reviews (
    id, tour_id, client_name, rating, text, is_approved, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING id, tour_id, client_name, rating, text, is_approved, created_at, updated_at
`, review.ID, review.TourID, review.ClientName, review.Rating, review.Text,
		review.IsApproved, review.CreatedAt, review.UpdatedAt)
	return scanReview(row)
}

func (s *Store) ApproveReview(ctx context.Context, id uuid.UUID) (domain.Review, error) {
	return s.setReviewApproval(ctx, id, true)
}

func (s *Store) RejectReview(ctx context.Context, id uuid.UUID) (domain.Review, error) {
	return s.setReviewApproval(ctx, id, false)
}

func (s *Store) DeleteReview(ctx context.Context, id uuid.UUID) error {
	result, err := s.db.ExecContext(ctx, `DELETE FROM reviews WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete review: %w", err)
	}
	return requireAffected(result)
}

func (s *Store) setReviewApproval(ctx context.Context, id uuid.UUID, approved bool) (domain.Review, error) {
	row := s.db.QueryRowContext(ctx, `
UPDATE reviews
SET is_approved = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING id, tour_id, client_name, rating, text, is_approved, created_at, updated_at
`, id, approved)
	return scanReview(row)
}

const tourWhereClause = `
($1::date IS NULL OR date_end >= $1) AND
($2::date IS NULL OR date_start <= $2) AND
($3::integer IS NULL OR price >= $3) AND
($4::integer IS NULL OR price <= $4) AND
($5::text = '' OR LOWER(location) LIKE '%' || LOWER($5) || '%') AND
($6::boolean IS NULL OR is_active = $6) AND
($7::boolean IS NULL OR is_hot = $7) AND
($8::text = '' OR (
    LOWER(title) LIKE '%' || LOWER($8) || '%' OR
    LOWER(location) LIKE '%' || LOWER($8) || '%' OR
    LOWER(slug) LIKE '%' || LOWER($8) || '%'
)) AND
($9::integer IS NULL OR slots_left >= $9)
`

func tourFilterArgs(filters ports.TourFilters) []any {
	return []any{
		dateValue(filters.DateFrom),
		dateValue(filters.DateTo),
		intValue(filters.PriceMin),
		intValue(filters.PriceMax),
		strings.TrimSpace(filters.Location),
		boolValue(filters.IsActive),
		boolValue(filters.IsHot),
		strings.TrimSpace(filters.Query),
		intValue(filters.MinSlots),
	}
}

const bookingWhereClause = `
($1::uuid IS NULL OR tour_id = $1) AND
($2::text IS NULL OR status = $2) AND
($3::timestamptz IS NULL OR created_at >= $3) AND
($4::timestamptz IS NULL OR created_at <= $4) AND
($5::uuid IS NULL OR user_id = $5)
`

func bookingFilterArgs(filters ports.BookingFilters) []any {
	var status any
	if filters.Status != nil {
		status = string(*filters.Status)
	}
	return []any{
		uuidValue(filters.TourID),
		status,
		timeValue(filters.From),
		timeValue(filters.To),
		uuidValue(filters.UserID),
	}
}

const reviewWhereClause = `
($1::uuid IS NULL OR tour_id = $1) AND
($2::integer IS NULL OR rating = $2) AND
($3::boolean IS NULL OR is_approved = $3)
`

func reviewFilterArgs(filters ports.ReviewFilters) []any {
	return []any{
		uuidValue(filters.TourID),
		intValue(filters.Rating),
		boolValue(filters.IsApproved),
	}
}

type scanner interface {
	Scan(dest ...any) error
}

func scanTour(row scanner) (domain.Tour, error) {
	var tour domain.Tour
	var images pq.StringArray
	err := row.Scan(
		&tour.ID,
		&tour.Slug,
		&tour.Title,
		&tour.Description,
		&tour.Price,
		&tour.Currency,
		&tour.DateStart,
		&tour.DateEnd,
		&tour.SlotsTotal,
		&tour.SlotsLeft,
		&tour.Location,
		&images,
		&tour.IsActive,
		&tour.IsHot,
		&tour.OverbookingEnabled,
		&tour.CreatedAt,
		&tour.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Tour{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.Tour{}, fmt.Errorf("scan tour: %w", err)
	}
	tour.Images = []string(images)
	return tour, nil
}

func scanTours(rows *sql.Rows) ([]domain.Tour, error) {
	var tours []domain.Tour
	for rows.Next() {
		tour, err := scanTour(rows)
		if err != nil {
			return nil, err
		}
		tours = append(tours, tour)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate tours: %w", err)
	}
	return tours, nil
}

func scanBooking(row scanner) (domain.Booking, error) {
	var booking domain.Booking
	var userID sql.NullString
	err := row.Scan(
		&booking.ID,
		&booking.TourID,
		&userID,
		&booking.Name,
		&booking.Phone,
		&booking.Email,
		&booking.PeopleCount,
		&booking.Status,
		&booking.TotalPrice,
		&booking.Comment,
		&booking.Overbooked,
		&booking.Source,
		&booking.CreatedAt,
		&booking.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Booking{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.Booking{}, fmt.Errorf("scan booking: %w", err)
	}
	if userID.Valid {
		parsed, parseErr := uuid.Parse(userID.String)
		if parseErr == nil {
			booking.UserID = &parsed
		}
	}
	return booking, nil
}

func scanBookings(rows *sql.Rows) ([]domain.Booking, error) {
	var bookings []domain.Booking
	for rows.Next() {
		booking, err := scanBooking(rows)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, booking)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate bookings: %w", err)
	}
	return bookings, nil
}

func scanReview(row scanner) (domain.Review, error) {
	var review domain.Review
	err := row.Scan(
		&review.ID,
		&review.TourID,
		&review.ClientName,
		&review.Rating,
		&review.Text,
		&review.IsApproved,
		&review.CreatedAt,
		&review.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Review{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.Review{}, fmt.Errorf("scan review: %w", err)
	}
	return review, nil
}

func scanReviews(rows *sql.Rows) ([]domain.Review, error) {
	var reviews []domain.Review
	for rows.Next() {
		review, err := scanReview(rows)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate reviews: %w", err)
	}
	return reviews, nil
}

func count(ctx context.Context, db *sql.DB, query string, args ...any) (int, error) {
	var total int
	if err := db.QueryRowContext(ctx, query, args...).Scan(&total); err != nil {
		return 0, fmt.Errorf("count rows: %w", err)
	}
	return total, nil
}

func pageMeta(pagination ports.Pagination, total int) ports.PageMeta {
	return ports.PageMeta{
		Page:    pagination.Page,
		Limit:   pagination.Limit,
		Total:   total,
		HasNext: pagination.Page*pagination.Limit < total,
	}
}

func offset(pagination ports.Pagination) int {
	return (pagination.Page - 1) * pagination.Limit
}

func requireAffected(result sql.Result) error {
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if affected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func dateValue(value *time.Time) any {
	if value == nil {
		return nil
	}
	return value.Format(time.DateOnly)
}

func timeValue(value *time.Time) any {
	if value == nil {
		return nil
	}
	return *value
}

func intValue(value *int) any {
	if value == nil {
		return nil
	}
	return *value
}

func boolValue(value *bool) any {
	if value == nil {
		return nil
	}
	return *value
}

func uuidValue(value *uuid.UUID) any {
	if value == nil {
		return nil
	}
	return *value
}
