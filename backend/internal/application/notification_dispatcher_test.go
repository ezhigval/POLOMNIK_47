package application

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"

	"palomnik/internal/domain"
)

func TestFormatPassengersForMessengerMasksPassportAndPhone(t *testing.T) {
	passenger, err := domain.NewPassenger(domain.NewPassengerInput{
		ID:        uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"),
		UserID:    uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"),
		Name:      "Иван Иванов",
		Phone:     "+79001234567",
		BirthDate: time.Date(1990, 5, 15, 0, 0, 0, 0, time.UTC),
		Passport:  "4010 123456",
		Now:       time.Date(2026, 8, 23, 0, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("passenger: %v", err)
	}

	text := formatPassengersForMessenger([]domain.Passenger{passenger})
	if !strings.Contains(text, "Иван Иванов") {
		t.Fatalf("name missing: %q", text)
	}
	for _, leak := range []string{"+79001234567", "4010 123456", "4010", "1990-05-15"} {
		if strings.Contains(text, leak) {
			t.Fatalf("leaked %q in %q", leak, text)
		}
	}
}
