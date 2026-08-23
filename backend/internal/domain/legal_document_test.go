package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func validLegalDocumentInput(mutators ...func(*NewLegalDocumentInput)) NewLegalDocumentInput {
	input := NewLegalDocumentInput{
		ID:          uuid.New(),
		Type:        LegalDocTypePersonalData,
		Version:     "1.0",
		Title:       "Согласие",
		Content:     "<p>Текст</p>",
		PublishedAt: time.Now().UTC(),
		IsActive:    true,
		Now:         time.Now().UTC(),
	}
	for _, mutate := range mutators {
		mutate(&input)
	}
	return input
}

func TestNewLegalDocument(t *testing.T) {
	doc, err := NewLegalDocument(validLegalDocumentInput())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if doc.Version != "1.0" {
		t.Fatalf("expected version 1.0, got %s", doc.Version)
	}
}

func TestNewLegalDocumentInvalidType(t *testing.T) {
	_, err := NewLegalDocument(validLegalDocumentInput(func(input *NewLegalDocumentInput) {
		input.Type = "invalid"
	}))
	if err != ErrInvalidLegalDocumentType {
		t.Fatalf("expected ErrInvalidLegalDocumentType, got %v", err)
	}
}

func TestNewConsent(t *testing.T) {
	userID := uuid.New()
	docID := uuid.New()
	consent, err := NewConsent(NewConsentInput{
		ID:              uuid.New(),
		UserID:          &userID,
		ConsentType:     ConsentTypePersonalData,
		DocumentID:      docID,
		DocumentVersion: "1.0",
		AcceptedAt:      time.Now().UTC(),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if consent.DocumentVersion != "1.0" {
		t.Fatalf("expected version 1.0, got %s", consent.DocumentVersion)
	}
}

func TestNewConsentRequiresVersion(t *testing.T) {
	_, err := NewConsent(NewConsentInput{
		ID:          uuid.New(),
		ConsentType: ConsentTypePersonalData,
		DocumentID:  uuid.New(),
		AcceptedAt:  time.Now().UTC(),
	})
	if err != ErrInvalidDocumentVersion {
		t.Fatalf("expected ErrInvalidDocumentVersion, got %v", err)
	}
}
