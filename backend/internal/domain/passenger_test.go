package domain

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewPassengerRequiresNamePhoneBirthPassport(t *testing.T) {
	now := time.Date(2026, 8, 23, 12, 0, 0, 0, time.UTC)
	valid := NewPassengerInput{
		ID:        uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"),
		UserID:    uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"),
		Name:      "Иван Иванов",
		Phone:     "89001234567",
		BirthDate: time.Date(1990, 5, 15, 10, 0, 0, 0, time.UTC),
		Passport:  "4010 123456",
		Now:       now,
	}

	got, err := NewPassenger(valid)
	if err != nil {
		t.Fatalf("new passenger: %v", err)
	}
	if got.Phone != "+79001234567" {
		t.Fatalf("phone not normalized: %q", got.Phone)
	}
	if got.BirthDate != time.Date(1990, 5, 15, 0, 0, 0, 0, time.UTC) {
		t.Fatalf("birth date not truncated: %v", got.BirthDate)
	}

	emptyName := valid
	emptyName.Name = "  "
	if _, err := NewPassenger(emptyName); err != ErrInvalidContactName {
		t.Fatalf("empty name: %v", err)
	}

	emptyPhone := valid
	emptyPhone.Phone = ""
	if _, err := NewPassenger(emptyPhone); err != ErrInvalidPhone {
		t.Fatalf("empty phone: %v", err)
	}

	future := valid
	future.BirthDate = now.Add(24 * time.Hour)
	if _, err := NewPassenger(future); err != ErrInvalidBirthDate {
		t.Fatalf("future birth: %v", err)
	}

	emptyPassport := valid
	emptyPassport.Passport = "  "
	if _, err := NewPassenger(emptyPassport); err != ErrInvalidPassport {
		t.Fatalf("empty passport: %v", err)
	}
}

func TestPassengerMessengerLineMasksPersonalData(t *testing.T) {
	passenger, err := NewPassenger(NewPassengerInput{
		ID:        uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"),
		UserID:    uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"),
		Name:      "Иван Иванов",
		Phone:     "+79001234567",
		BirthDate: time.Date(1990, 5, 15, 0, 0, 0, 0, time.UTC),
		Passport:  "4010 123456",
		Now:       time.Date(2026, 8, 23, 0, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("new passenger: %v", err)
	}

	line := passenger.MessengerLine()
	if !strings.Contains(line, "Иван Иванов") {
		t.Fatalf("name should stay readable: %q", line)
	}
	leaks := []string{passenger.Phone, "+79001234567", "4010", "123456", "1990", "15.05", "1990-05-15"}
	for _, leak := range leaks {
		if strings.Contains(line, leak) {
			t.Fatalf("messenger line leaked %q: %q", leak, line)
		}
	}
	if strings.Contains(line, "СНИЛС") || strings.Contains(line, "снилс") {
		t.Fatalf("SNILS must not appear: %q", line)
	}
}
