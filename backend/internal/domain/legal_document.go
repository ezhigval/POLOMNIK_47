package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type LegalDocumentType string

const (
	LegalDocTypePrivacyPolicy LegalDocumentType = "privacy_policy"
	LegalDocTypePersonalData  LegalDocumentType = "personal_data"
	LegalDocTypeDistribution  LegalDocumentType = "distribution"
	LegalDocTypeMarketing     LegalDocumentType = "marketing"
	LegalDocTypeCookie        LegalDocumentType = "cookie"
	LegalDocTypeTerms         LegalDocumentType = "terms"
	LegalDocTypeOffer         LegalDocumentType = "offer"
)

var legalDocumentTypes = map[LegalDocumentType]struct{}{
	LegalDocTypePrivacyPolicy: {},
	LegalDocTypePersonalData:  {},
	LegalDocTypeDistribution:  {},
	LegalDocTypeMarketing:     {},
	LegalDocTypeCookie:        {},
	LegalDocTypeTerms:         {},
	LegalDocTypeOffer:         {},
}

func ParseLegalDocumentType(raw string) (LegalDocumentType, error) {
	t := LegalDocumentType(strings.TrimSpace(raw))
	if _, ok := legalDocumentTypes[t]; !ok {
		return "", ErrInvalidLegalDocumentType
	}
	return t, nil
}

type LegalDocument struct {
	ID          uuid.UUID
	Type        LegalDocumentType
	Version     string
	Title       string
	Content     string
	PublishedAt time.Time
	UpdatedAt   time.Time
	IsActive    bool
}

type NewLegalDocumentInput struct {
	ID          uuid.UUID
	Type        LegalDocumentType
	Version     string
	Title       string
	Content     string
	PublishedAt time.Time
	IsActive    bool
	Now         time.Time
}

func NewLegalDocument(input NewLegalDocumentInput) (LegalDocument, error) {
	if input.ID == uuid.Nil {
		return LegalDocument{}, ErrInvalidID
	}
	if _, ok := legalDocumentTypes[input.Type]; !ok {
		return LegalDocument{}, ErrInvalidLegalDocumentType
	}
	version := strings.TrimSpace(input.Version)
	title := strings.TrimSpace(input.Title)
	content := strings.TrimSpace(input.Content)
	if version == "" {
		return LegalDocument{}, ErrInvalidDocumentVersion
	}
	if title == "" {
		return LegalDocument{}, ErrInvalidTitle
	}
	if content == "" {
		return LegalDocument{}, ErrInvalidLegalDocumentContent
	}
	if input.PublishedAt.IsZero() {
		return LegalDocument{}, ErrInvalidPublishedAt
	}

	now := input.Now
	if now.IsZero() {
		now = time.Now().UTC()
	}

	return LegalDocument{
		ID:          input.ID,
		Type:        input.Type,
		Version:     version,
		Title:       title,
		Content:     content,
		PublishedAt: input.PublishedAt.UTC(),
		UpdatedAt:   now,
		IsActive:    input.IsActive,
	}, nil
}
