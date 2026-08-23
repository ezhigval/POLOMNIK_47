package memory

import (
	"context"
	"encoding/json"
	"sort"

	"github.com/google/uuid"

	"palomnik/internal/domain"
	"palomnik/internal/ports"
)

func (s *Store) ListPages(_ context.Context, filters ports.CMSPageFilters) ([]domain.Page, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	pages := make([]domain.Page, 0, len(s.cmsPages))
	for _, page := range s.cmsPages {
		if filters.PublishedOnly && !page.IsPublished {
			continue
		}
		pages = append(pages, cloneCMSPage(page))
	}
	sort.Slice(pages, func(i, j int) bool {
		return pages[i].Title < pages[j].Title
	})
	return pages, nil
}

func (s *Store) GetPage(_ context.Context, id uuid.UUID) (domain.Page, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	page, ok := s.cmsPages[id]
	if !ok {
		return domain.Page{}, domain.ErrNotFound
	}
	return cloneCMSPage(page), nil
}

func (s *Store) GetPageBySlug(_ context.Context, slug string) (domain.Page, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, page := range s.cmsPages {
		if page.Slug == slug {
			return cloneCMSPage(page), nil
		}
	}
	return domain.Page{}, domain.ErrNotFound
}

func (s *Store) GetPageByPath(_ context.Context, path string) (domain.Page, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, page := range s.cmsPages {
		if page.Path == path {
			return cloneCMSPage(page), nil
		}
	}
	return domain.Page{}, domain.ErrNotFound
}

func (s *Store) CreatePage(_ context.Context, page domain.Page) (domain.Page, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, existing := range s.cmsPages {
		if existing.Slug == page.Slug {
			return domain.Page{}, domain.ErrDuplicateSlug
		}
		if existing.Path == page.Path {
			return domain.Page{}, domain.ErrDuplicatePath
		}
	}
	s.cmsPages[page.ID] = cloneCMSPage(page)
	return cloneCMSPage(page), nil
}

func (s *Store) UpdatePage(_ context.Context, page domain.Page) (domain.Page, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.cmsPages[page.ID]; !ok {
		return domain.Page{}, domain.ErrNotFound
	}
	for _, existing := range s.cmsPages {
		if existing.ID == page.ID {
			continue
		}
		if existing.Path == page.Path {
			return domain.Page{}, domain.ErrDuplicatePath
		}
	}
	s.cmsPages[page.ID] = cloneCMSPage(page)
	return cloneCMSPage(page), nil
}

func (s *Store) DeletePage(_ context.Context, id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.cmsPages[id]; !ok {
		return domain.ErrNotFound
	}
	delete(s.cmsPages, id)
	for blockID, block := range s.cmsBlocks {
		if block.PageID == id {
			delete(s.cmsBlocks, blockID)
		}
	}
	return nil
}

func (s *Store) ListBlocks(_ context.Context, pageID uuid.UUID) ([]domain.Block, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	blocks := make([]domain.Block, 0)
	for _, block := range s.cmsBlocks {
		if block.PageID == pageID {
			blocks = append(blocks, cloneCMSBlock(block))
		}
	}
	sort.Slice(blocks, func(i, j int) bool {
		if blocks[i].SortOrder == blocks[j].SortOrder {
			return blocks[i].CreatedAt.Before(blocks[j].CreatedAt)
		}
		return blocks[i].SortOrder < blocks[j].SortOrder
	})
	return blocks, nil
}

func (s *Store) GetBlock(_ context.Context, id uuid.UUID) (domain.Block, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	block, ok := s.cmsBlocks[id]
	if !ok {
		return domain.Block{}, domain.ErrNotFound
	}
	return cloneCMSBlock(block), nil
}

func (s *Store) CreateBlock(_ context.Context, block domain.Block) (domain.Block, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.cmsPages[block.PageID]; !ok {
		return domain.Block{}, domain.ErrNotFound
	}
	s.cmsBlocks[block.ID] = cloneCMSBlock(block)
	return cloneCMSBlock(block), nil
}

func (s *Store) UpdateBlock(_ context.Context, block domain.Block) (domain.Block, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.cmsBlocks[block.ID]; !ok {
		return domain.Block{}, domain.ErrNotFound
	}
	s.cmsBlocks[block.ID] = cloneCMSBlock(block)
	return cloneCMSBlock(block), nil
}

func (s *Store) DeleteBlock(_ context.Context, id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.cmsBlocks[id]; !ok {
		return domain.ErrNotFound
	}
	delete(s.cmsBlocks, id)
	return nil
}

func (s *Store) ReorderBlocks(_ context.Context, pageID uuid.UUID, orderedIDs []uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.cmsPages[pageID]; !ok {
		return domain.ErrNotFound
	}
	for index, id := range orderedIDs {
		block, ok := s.cmsBlocks[id]
		if !ok || block.PageID != pageID {
			return domain.ErrNotFound
		}
		block.SortOrder = index
		s.cmsBlocks[id] = block
	}
	return nil
}

func cloneCMSPage(page domain.Page) domain.Page {
	cloned := page
	cloned.Blocks = nil
	return cloned
}

func cloneCMSBlock(block domain.Block) domain.Block {
	cloned := block
	if len(block.Content) > 0 {
		cloned.Content = append(json.RawMessage(nil), block.Content...)
	}
	return cloned
}
