package dto

import (
	"palomnik/internal/domain"
)

type LegalDocumentResponse struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	Version     string `json:"version"`
	Title       string `json:"title"`
	Content     string `json:"content,omitempty"`
	PublishedAt string `json:"published_at"`
	UpdatedAt   string `json:"updated_at"`
	IsActive    bool   `json:"is_active,omitempty"`
}

type LegalDocumentSummaryResponse struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	Version     string `json:"version"`
	Title       string `json:"title"`
	PublishedAt string `json:"published_at"`
	UpdatedAt   string `json:"updated_at"`
	IsActive    bool   `json:"is_active"`
}

type PublishLegalDocumentRequest struct {
	Type    string `json:"type"`
	Version string `json:"version"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type RecordConsentRequest struct {
	ConsentType string  `json:"consent_type"`
	RequestID   *string `json:"request_id,omitempty"`
}

type ConsentResponse struct {
	ID              string  `json:"id"`
	UserID          *string `json:"user_id,omitempty"`
	RequestID       *string `json:"request_id,omitempty"`
	ConsentType     string  `json:"consent_type"`
	DocumentID      string  `json:"document_id"`
	DocumentVersion string  `json:"document_version"`
	AcceptedAt      string  `json:"accepted_at"`
}

func ToLegalDocumentResponse(doc domain.LegalDocument, includeContent bool) LegalDocumentResponse {
	resp := LegalDocumentResponse{
		ID:          doc.ID.String(),
		Type:        string(doc.Type),
		Version:     doc.Version,
		Title:       doc.Title,
		PublishedAt: formatTime(doc.PublishedAt),
		UpdatedAt:   formatTime(doc.UpdatedAt),
		IsActive:    doc.IsActive,
	}
	if includeContent {
		resp.Content = doc.Content
	}
	return resp
}

func ToLegalDocumentSummaryResponse(doc domain.LegalDocument) LegalDocumentSummaryResponse {
	return LegalDocumentSummaryResponse{
		ID:          doc.ID.String(),
		Type:        string(doc.Type),
		Version:     doc.Version,
		Title:       doc.Title,
		PublishedAt: formatTime(doc.PublishedAt),
		UpdatedAt:   formatTime(doc.UpdatedAt),
		IsActive:    doc.IsActive,
	}
}

func ToConsentResponse(c domain.Consent) ConsentResponse {
	resp := ConsentResponse{
		ID:              c.ID.String(),
		ConsentType:     string(c.ConsentType),
		DocumentID:      c.DocumentID.String(),
		DocumentVersion: c.DocumentVersion,
		AcceptedAt:      formatTime(c.AcceptedAt),
	}
	if c.UserID != nil {
		s := c.UserID.String()
		resp.UserID = &s
	}
	if c.RequestID != nil {
		s := c.RequestID.String()
		resp.RequestID = &s
	}
	return resp
}
