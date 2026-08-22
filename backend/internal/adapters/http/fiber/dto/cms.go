package dto

import (
	"encoding/json"
	"time"

	"polomnik/internal/domain"
)

type CMSPageResponse struct {
	ID              string             `json:"id"`
	Slug            string             `json:"slug"`
	Title           string             `json:"title"`
	Path            string             `json:"path"`
	MetaTitle       string             `json:"meta_title"`
	MetaDescription string             `json:"meta_description"`
	IsPublished     bool               `json:"is_published"`
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at"`
	Blocks          []CMSBlockResponse `json:"blocks,omitempty"`
}

type CMSBlockResponse struct {
	ID        string          `json:"id"`
	PageID    string          `json:"page_id"`
	Type      string          `json:"type"`
	SortOrder int             `json:"sort_order"`
	Content   json.RawMessage `json:"content"`
	IsVisible bool            `json:"is_visible"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type CMSPageCreateRequest struct {
	Slug        string `json:"slug"`
	Title       string `json:"title"`
	Path        string `json:"path"`
	IsPublished *bool  `json:"is_published"`
}

type CMSPageUpdateRequest struct {
	Title           *string `json:"title"`
	Path            *string `json:"path"`
	MetaTitle       *string `json:"meta_title"`
	MetaDescription *string `json:"meta_description"`
	IsPublished     *bool   `json:"is_published"`
}

type CMSBlockCreateRequest struct {
	Type      string          `json:"type"`
	Content   json.RawMessage `json:"content"`
	IsVisible *bool           `json:"is_visible"`
}

type CMSBlockUpdateRequest struct {
	Content   *json.RawMessage `json:"content"`
	IsVisible *bool            `json:"is_visible"`
	SortOrder *int             `json:"sort_order"`
}

type CMSReorderRequest struct {
	BlockIDs []string `json:"block_ids"`
}

type CMSBlockTemplateResponse struct {
	Type    string          `json:"type"`
	Label   string          `json:"label"`
	Content json.RawMessage `json:"content"`
}

func ToCMSPageResponse(page domain.Page) CMSPageResponse {
	blocks := make([]CMSBlockResponse, 0, len(page.Blocks))
	for _, block := range page.Blocks {
		blocks = append(blocks, ToCMSBlockResponse(block))
	}
	return CMSPageResponse{
		ID:              page.ID.String(),
		Slug:            page.Slug,
		Title:           page.Title,
		Path:            page.Path,
		MetaTitle:       page.MetaTitle,
		MetaDescription: page.MetaDescription,
		IsPublished:     page.IsPublished,
		CreatedAt:       page.CreatedAt,
		UpdatedAt:       page.UpdatedAt,
		Blocks:          blocks,
	}
}

func ToCMSBlockResponse(block domain.Block) CMSBlockResponse {
	content := block.Content
	if len(content) == 0 {
		content = json.RawMessage(`{}`)
	}
	return CMSBlockResponse{
		ID:        block.ID.String(),
		PageID:    block.PageID.String(),
		Type:      block.Type,
		SortOrder: block.SortOrder,
		Content:   content,
		IsVisible: block.IsVisible,
		CreatedAt: block.CreatedAt,
		UpdatedAt: block.UpdatedAt,
	}
}
