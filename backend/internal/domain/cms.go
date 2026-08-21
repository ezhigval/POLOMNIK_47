package domain

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	BlockTypeHero                 = "hero"
	BlockTypeAbout                = "about"
	BlockTypeWhyUs                = "why_us"
	BlockTypeHowItWorks           = "how_it_works"
	BlockTypeFAQ                  = "faq"
	BlockTypeCTA                  = "cta"
	BlockTypeRichText             = "rich_text"
	BlockTypePopularDestinations  = "popular_destinations"
	BlockTypeTestimonials         = "testimonials"

	PageSlugHome = "home"
	PagePathHome = "/"
)

type Page struct {
	ID          uuid.UUID
	Slug        string
	Title       string
	Path        string
	IsPublished bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Blocks      []Block
}

type Block struct {
	ID        uuid.UUID
	PageID    uuid.UUID
	Type      string
	SortOrder int
	Content   json.RawMessage
	IsVisible bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type NewPageInput struct {
	ID          uuid.UUID
	Slug        string
	Title       string
	Path        string
	IsPublished bool
	Now         time.Time
}

type NewBlockInput struct {
	ID        uuid.UUID
	PageID    uuid.UUID
	Type      string
	SortOrder int
	Content   json.RawMessage
	IsVisible bool
	Now       time.Time
}

func NewPage(input NewPageInput) (Page, error) {
	slug := normalizeSlug(input.Slug)
	title := strings.TrimSpace(input.Title)
	path := normalizePath(input.Path)
	if input.ID == uuid.Nil {
		return Page{}, ErrInvalidID
	}
	if slug == "" {
		return Page{}, ErrInvalidSlug
	}
	if title == "" {
		return Page{}, ErrInvalidTitle
	}
	if path == "" {
		return Page{}, ErrInvalidPath
	}
	now := input.Now
	if now.IsZero() {
		now = time.Now().UTC()
	}
	return Page{
		ID:          input.ID,
		Slug:        slug,
		Title:       title,
		Path:        path,
		IsPublished: input.IsPublished,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func NewBlock(input NewBlockInput) (Block, error) {
	if input.ID == uuid.Nil || input.PageID == uuid.Nil {
		return Block{}, ErrInvalidID
	}
	blockType := strings.TrimSpace(input.Type)
	if !IsValidBlockType(blockType) {
		return Block{}, ErrInvalidBlockType
	}
	content := input.Content
	if len(content) == 0 {
		content = json.RawMessage(`{}`)
	}
	now := input.Now
	if now.IsZero() {
		now = time.Now().UTC()
	}
	return Block{
		ID:        input.ID,
		PageID:    input.PageID,
		Type:      blockType,
		SortOrder: input.SortOrder,
		Content:   content,
		IsVisible: input.IsVisible,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func IsValidBlockType(blockType string) bool {
	switch blockType {
	case BlockTypeHero, BlockTypeAbout, BlockTypeWhyUs, BlockTypeHowItWorks, BlockTypeFAQ, BlockTypeCTA, BlockTypeRichText, BlockTypePopularDestinations, BlockTypeTestimonials:
		return true
	default:
		return false
	}
}

func normalizeSlug(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	value = strings.ReplaceAll(value, " ", "-")
	return value
}

func normalizePath(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	if !strings.HasPrefix(value, "/") {
		value = "/" + value
	}
	if len(value) > 1 {
		value = strings.TrimRight(value, "/")
	}
	return value
}
