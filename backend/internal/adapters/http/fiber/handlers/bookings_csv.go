package handlers

import (
	"bytes"
	"encoding/csv"
	"strconv"
	"time"

	"palomnik/internal/domain"
)

func managementBookingsCSV(items []domain.Booking) ([]byte, error) {
	var buf bytes.Buffer
	buf.Write([]byte{0xEF, 0xBB, 0xBF})
	w := csv.NewWriter(&buf)
	if err := w.Write([]string{
		"id",
		"created_at",
		"status",
		"payment_status",
		"name",
		"phone",
		"email",
		"tour_id",
		"people_count",
		"total_price",
		"comment",
		"overbooked",
		"source",
	}); err != nil {
		return nil, err
	}
	for _, booking := range items {
		if err := w.Write([]string{
			booking.ID.String(),
			booking.CreatedAt.UTC().Format(time.RFC3339),
			string(booking.Status),
			string(booking.PaymentStatus),
			booking.Name,
			booking.Phone,
			booking.Email,
			booking.TourID.String(),
			strconv.Itoa(booking.PeopleCount),
			strconv.Itoa(booking.TotalPrice),
			booking.Comment,
			strconv.FormatBool(booking.Overbooked),
			booking.Source,
		}); err != nil {
			return nil, err
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
