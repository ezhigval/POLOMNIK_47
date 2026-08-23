package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type ConsentType string

const (
	ConsentTypePersonalData ConsentType = "personal_data"
	ConsentTypeMarketing    ConsentType = "marketing"
	ConsentTypeDistribution ConsentType = "distribution"
	ConsentTypeCookieAll    ConsentType = "cookie_all"
	ConsentTypeCookieEssential ConsentType = "cookie_essential"
	ConsentTypeCookieReject ConsentType = "cookie_reject"
)

var consentTypes = map[ConsentType]struct{}{
	ConsentTypePersonalData:    {},
	ConsentTypeMarketing:       {},
	ConsentTypeDistribution:    {},
	ConsentTypeCookieAll:       {},
	ConsentTypeCookieEssential: {},
	ConsentTypeCookieReject:    {},
}

func ParseConsentType(raw string) (ConsentType, error) {
	t := ConsentType(strings.TrimSpace(raw))
	if _, ok := consentTypes[t]; !ok {
		return "", ErrInvalidConsentType
	}
	return t, nil
}

type Consent struct {
	ID              uuid.UUID
	UserID          *uuid.UUID
	RequestID       *uuid.UUID
	ConsentType     ConsentType
	DocumentID      uuid.UUID
	DocumentVersion string
	AcceptedAt      time.Time
	IP              string
	UserAgent       string
}

type NewConsentInput struct {
	ID              uuid.UUID
	UserID          *uuid.UUID
	RequestID       *uuid.UUID
	ConsentType     ConsentType
	DocumentID      uuid.UUID
	DocumentVersion string
	AcceptedAt      time.Time
	IP              string
	UserAgent       string
}

func NewConsent(input NewConsentInput) (Consent, error) {
	if input.ID == uuid.Nil {
		return Consent{}, ErrInvalidID
	}
	if _, ok := consentTypes[input.ConsentType]; !ok {
		return Consent{}, ErrInvalidConsentType
	}
	if input.DocumentID == uuid.Nil {
		return Consent{}, ErrInvalidID
	}
	version := strings.TrimSpace(input.DocumentVersion)
	if version == "" {
		return Consent{}, ErrInvalidDocumentVersion
	}
	if input.AcceptedAt.IsZero() {
		return Consent{}, ErrInvalidAcceptedAt
	}

	return Consent{
		ID:              input.ID,
		UserID:          input.UserID,
		RequestID:       input.RequestID,
		ConsentType:     input.ConsentType,
		DocumentID:      input.DocumentID,
		DocumentVersion: version,
		AcceptedAt:      input.AcceptedAt.UTC(),
		IP:              strings.TrimSpace(input.IP),
		UserAgent:       strings.TrimSpace(input.UserAgent),
	}, nil
}
