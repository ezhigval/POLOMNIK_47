package ports

import (
	"context"

	"github.com/google/uuid"

	"palomnik/internal/domain"
)

type CMSPageFilters struct {
	PublishedOnly bool
}

type CMSRepository interface {
	ListPages(ctx context.Context, filters CMSPageFilters) ([]domain.Page, error)
	GetPage(ctx context.Context, id uuid.UUID) (domain.Page, error)
	GetPageBySlug(ctx context.Context, slug string) (domain.Page, error)
	GetPageByPath(ctx context.Context, path string) (domain.Page, error)
	CreatePage(ctx context.Context, page domain.Page) (domain.Page, error)
	UpdatePage(ctx context.Context, page domain.Page) (domain.Page, error)
	DeletePage(ctx context.Context, id uuid.UUID) error

	ListBlocks(ctx context.Context, pageID uuid.UUID) ([]domain.Block, error)
	GetBlock(ctx context.Context, id uuid.UUID) (domain.Block, error)
	CreateBlock(ctx context.Context, block domain.Block) (domain.Block, error)
	UpdateBlock(ctx context.Context, block domain.Block) (domain.Block, error)
	DeleteBlock(ctx context.Context, id uuid.UUID) error
	ReorderBlocks(ctx context.Context, pageID uuid.UUID, orderedIDs []uuid.UUID) error
}
