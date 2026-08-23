package domain

import (
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

const maxPassportLength = 120

type Passenger struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Name      string
	Phone     string
	BirthDate time.Time
	Passport  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type NewPassengerInput struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Name      string
	Phone     string
	BirthDate time.Time
	Passport  string
	Now       time.Time
}

func NewPassenger(input NewPassengerInput) (Passenger, error) {
	if input.ID == uuid.Nil || input.UserID == uuid.Nil {
		return Passenger{}, ErrInvalidID
	}

	now := input.Now
	if now.IsZero() {
		now = time.Now().UTC()
	}

	name, phone, birthDate, passport, err := normalizePassengerFields(input.Name, input.Phone, input.BirthDate, input.Passport, now)
	if err != nil {
		return Passenger{}, err
	}

	return Passenger{
		ID:        input.ID,
		UserID:    input.UserID,
		Name:      name,
		Phone:     phone,
		BirthDate: birthDate,
		Passport:  passport,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (p Passenger) Update(name, phone string, birthDate time.Time, passport string, now time.Time) (Passenger, error) {
	if now.IsZero() {
		now = time.Now().UTC()
	}
	normalizedName, normalizedPhone, normalizedBirth, normalizedPassport, err := normalizePassengerFields(name, phone, birthDate, passport, now)
	if err != nil {
		return Passenger{}, err
	}
	p.Name = normalizedName
	p.Phone = normalizedPhone
	p.BirthDate = normalizedBirth
	p.Passport = normalizedPassport
	p.UpdatedAt = now
	return p, nil
}

func normalizePassengerFields(name, phone string, birthDate time.Time, passport string, now time.Time) (string, string, time.Time, string, error) {
	name = strings.TrimSpace(name)
	if name == "" || utf8.RuneCountInString(name) > maxNameLength {
		return "", "", time.Time{}, "", ErrInvalidContactName
	}

	phone = NormalizePhone(phone)
	if phone == "" {
		return "", "", time.Time{}, "", ErrInvalidPhone
	}

	birthDate = calendarDateUTC(birthDate)
	if birthDate.IsZero() || birthDate.After(calendarDateUTC(now)) {
		return "", "", time.Time{}, "", ErrInvalidBirthDate
	}

	passport = strings.TrimSpace(passport)
	if passport == "" || utf8.RuneCountInString(passport) > maxPassportLength {
		return "", "", time.Time{}, "", ErrInvalidPassport
	}

	return name, phone, birthDate, passport, nil
}

func calendarDateUTC(value time.Time) time.Time {
	if value.IsZero() {
		return time.Time{}
	}
	value = value.UTC()
	return time.Date(value.Year(), value.Month(), value.Day(), 0, 0, 0, 0, time.UTC)
}

func MaskPhone(phone string) string {
	return maskKeepLastRunes(phone, 2)
}

func MaskPassport(passport string) string {
	return maskKeepLastRunes(passport, 2)
}

func MaskBirthDate(date time.Time) string {
	if date.IsZero() {
		return ""
	}
	return "**.**.****"
}

func maskKeepLastRunes(value string, keep int) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	runes := []rune(value)
	if keep < 0 {
		keep = 0
	}
	if len(runes) <= keep {
		return strings.Repeat("*", len(runes))
	}
	return strings.Repeat("*", len(runes)-keep) + string(runes[len(runes)-keep:])
}

// MessengerLine is the only passenger representation allowed in Telegram/Max text:
// name stays readable; phone, date of birth and passport are masked. No SNILS.
func (p Passenger) MessengerLine() string {
	parts := []string{p.Name}
	if phone := MaskPhone(p.Phone); phone != "" {
		parts = append(parts, "тел. "+phone)
	}
	if dob := MaskBirthDate(p.BirthDate); dob != "" {
		parts = append(parts, "ДР "+dob)
	}
	if passport := MaskPassport(p.Passport); passport != "" {
		parts = append(parts, "паспорт "+passport)
	}
	return strings.Join(parts, ", ")
}
