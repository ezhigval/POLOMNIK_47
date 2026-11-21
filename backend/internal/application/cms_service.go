package application

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"palomnik/internal/domain"
	"palomnik/internal/ports"
)

type CMSService struct {
	cms ports.CMSRepository
}

func NewCMSService(cms ports.CMSRepository) *CMSService {
	return &CMSService{cms: cms}
}

type CreatePageInput struct {
	Slug        string
	Title       string
	Path        string
	IsPublished bool
}

type UpdatePageInput struct {
	Title           *string
	Path            *string
	MetaTitle       *string
	MetaDescription *string
	IsPublished     *bool
}

type CreateBlockInput struct {
	Type      string
	Content   json.RawMessage
	IsVisible bool
}

type UpdateBlockInput struct {
	Content   *json.RawMessage
	IsVisible *bool
	SortOrder *int
}

func (s *CMSService) ListPages(ctx context.Context, publishedOnly bool) ([]domain.Page, error) {
	return s.cms.ListPages(ctx, ports.CMSPageFilters{PublishedOnly: publishedOnly})
}

func (s *CMSService) GetPage(ctx context.Context, id uuid.UUID) (domain.Page, error) {
	page, err := s.cms.GetPage(ctx, id)
	if err != nil {
		return domain.Page{}, err
	}
	blocks, err := s.cms.ListBlocks(ctx, page.ID)
	if err != nil {
		return domain.Page{}, err
	}
	page.Blocks = blocks
	return page, nil
}

func (s *CMSService) GetPublishedPageBySlug(ctx context.Context, slug string) (domain.Page, error) {
	page, err := s.cms.GetPageBySlug(ctx, slug)
	if err != nil {
		return domain.Page{}, err
	}
	if !page.IsPublished {
		return domain.Page{}, domain.ErrNotFound
	}
	blocks, err := s.cms.ListBlocks(ctx, page.ID)
	if err != nil {
		return domain.Page{}, err
	}
	visible := make([]domain.Block, 0, len(blocks))
	for _, block := range blocks {
		if block.IsVisible {
			visible = append(visible, block)
		}
	}
	page.Blocks = visible
	return page, nil
}

func (s *CMSService) CreatePage(ctx context.Context, input CreatePageInput) (domain.Page, error) {
	page, err := domain.NewPage(domain.NewPageInput{
		ID:          uuid.New(),
		Slug:        input.Slug,
		Title:       input.Title,
		Path:        input.Path,
		IsPublished: input.IsPublished,
		Now:         time.Now().UTC(),
	})
	if err != nil {
		return domain.Page{}, err
	}
	created, err := s.cms.CreatePage(ctx, page)
	if err != nil {
		return domain.Page{}, err
	}
	created.Blocks = []domain.Block{}
	return created, nil
}

func (s *CMSService) UpdatePage(ctx context.Context, id uuid.UUID, input UpdatePageInput) (domain.Page, error) {
	page, err := s.cms.GetPage(ctx, id)
	if err != nil {
		return domain.Page{}, err
	}
	if input.Title != nil {
		title := *input.Title
		if title == "" {
			return domain.Page{}, domain.ErrInvalidTitle
		}
		page.Title = title
	}
	if input.Path != nil {
		pathPage, err := domain.NewPage(domain.NewPageInput{
			ID:          page.ID,
			Slug:        page.Slug,
			Title:       page.Title,
			Path:        *input.Path,
			IsPublished: page.IsPublished,
			Now:         page.CreatedAt,
		})
		if err != nil {
			return domain.Page{}, err
		}
		page.Path = pathPage.Path
	}
	if input.IsPublished != nil {
		page.IsPublished = *input.IsPublished
	}
	if input.MetaTitle != nil {
		page.MetaTitle = strings.TrimSpace(*input.MetaTitle)
	}
	if input.MetaDescription != nil {
		page.MetaDescription = strings.TrimSpace(*input.MetaDescription)
	}
	page.UpdatedAt = time.Now().UTC()
	updated, err := s.cms.UpdatePage(ctx, page)
	if err != nil {
		return domain.Page{}, err
	}
	blocks, err := s.cms.ListBlocks(ctx, updated.ID)
	if err != nil {
		return domain.Page{}, err
	}
	updated.Blocks = blocks
	return updated, nil
}

func (s *CMSService) DeletePage(ctx context.Context, id uuid.UUID) error {
	return s.cms.DeletePage(ctx, id)
}

func (s *CMSService) CreateBlock(ctx context.Context, pageID uuid.UUID, input CreateBlockInput) (domain.Block, error) {
	if _, err := s.cms.GetPage(ctx, pageID); err != nil {
		return domain.Block{}, err
	}
	existing, err := s.cms.ListBlocks(ctx, pageID)
	if err != nil {
		return domain.Block{}, err
	}
	content := input.Content
	if len(content) == 0 {
		content = DefaultBlockContent(input.Type)
	}
	block, err := domain.NewBlock(domain.NewBlockInput{
		ID:        uuid.New(),
		PageID:    pageID,
		Type:      input.Type,
		SortOrder: len(existing),
		Content:   content,
		IsVisible: input.IsVisible,
		Now:       time.Now().UTC(),
	})
	if err != nil {
		return domain.Block{}, err
	}
	return s.cms.CreateBlock(ctx, block)
}

func (s *CMSService) UpdateBlock(ctx context.Context, id uuid.UUID, input UpdateBlockInput) (domain.Block, error) {
	block, err := s.cms.GetBlock(ctx, id)
	if err != nil {
		return domain.Block{}, err
	}
	if input.Content != nil {
		block.Content = *input.Content
		if len(block.Content) == 0 {
			block.Content = json.RawMessage(`{}`)
		}
	}
	if input.IsVisible != nil {
		block.IsVisible = *input.IsVisible
	}
	if input.SortOrder != nil {
		block.SortOrder = *input.SortOrder
	}
	block.UpdatedAt = time.Now().UTC()
	return s.cms.UpdateBlock(ctx, block)
}

func (s *CMSService) DeleteBlock(ctx context.Context, id uuid.UUID) error {
	return s.cms.DeleteBlock(ctx, id)
}

func (s *CMSService) ReorderBlocks(ctx context.Context, pageID uuid.UUID, orderedIDs []uuid.UUID) error {
	if _, err := s.cms.GetPage(ctx, pageID); err != nil {
		return err
	}
	return s.cms.ReorderBlocks(ctx, pageID, orderedIDs)
}

func DefaultBlockContent(blockType string) json.RawMessage {
	switch blockType {
	case domain.BlockTypeHero:
		return json.RawMessage(`{"eyebrow":"","title":"Заголовок","subtitle":""}`)
	case domain.BlockTypeAbout:
		return json.RawMessage(`{"eyebrow":"О службе","title":"О нас","paragraphs":[""],"highlights":[],"stats":[],"showContacts":true}`)
	case domain.BlockTypeWhyUs:
		return json.RawMessage(`{"eyebrow":"Почему мы","title":"Паломничество без лишних забот","description":"Мы не просто везём вас в монастырь — мы создаём возможность для тишины, молитвы и встречи с Богом.","items":[]}`)
	case domain.BlockTypeHowItWorks:
		return json.RawMessage(`{"eyebrow":"Просто","title":"Как записаться","description":"","steps":[],"ctaLabel":"Записаться","ctaHref":"/search"}`)
	case domain.BlockTypeFAQ:
		return json.RawMessage(`{"eyebrow":"Вопросы","title":"Частые вопросы","description":"","items":[]}`)
	case domain.BlockTypeCTA:
		return json.RawMessage(`{"title":"Готовы отправиться?","subtitle":"","button":"Смотреть туры","href":"/search"}`)
	case domain.BlockTypeRichText:
		return json.RawMessage(`{"eyebrow":"","title":"Текст","body":""}`)
	case domain.BlockTypePopularDestinations, domain.BlockTypeTestimonials, domain.BlockTypeFeaturedRoute:
		return json.RawMessage(`{}`)
	case domain.BlockTypeInformationHero:
		return json.RawMessage(`{"eyebrow":"Паломническая служба Тихвинской епархии","title":"Дорога к святыням начинается с доброй встречи","subtitle":"Рассказываем о монастырях, святынях и самой службе. Помогаем подготовиться к паломничеству и выбрать подходящую поездку.","primaryCta":"Смотреть поездки","primaryHref":"/search","secondaryCta":"О службе","secondaryHref":"#service"}`)
	case domain.BlockTypeHolyPlaces:
		return json.RawMessage(`{"eyebrow":"Места паломничества","title":"Монастыри и святыни","description":"История, расположение и полезная информация для паломников.","items":[]}`)
	case domain.BlockTypeServiceProfile:
		return json.RawMessage(`{"eyebrow":"О службе","title":"Под покровом Божией Матери «Тихвинская»","paragraphs":["Паломническая служба Тихвинской епархии организует поездки к православным святыням и помогает паломникам подготовиться к дороге."],"links":[{"label":"Сайт Тихвинской епархии","href":"https://www.tikhvin-eparhia.ru/"},{"label":"Храм Покрова на Боровой","href":"https://church-pokrov.ru/"}]}`)
	case domain.BlockTypeServiceTeam:
		return json.RawMessage(`{"eyebrow":"Люди службы","title":"Кто сопровождает паломников","description":"Руководители и сотрудники службы.","members":[]}`)
	case domain.BlockTypeLegalDonation:
		return json.RawMessage(`{"eyebrow":"Официально","title":"Документы и реквизиты","legalName":"","legalAddress":"","inn":"","kpp":"","ogrn":"","donationRecipient":"","donationBankName":"","donationAccount":"","donationCorrespondentAccount":"","donationBik":"","donationPurpose":"","donationUrl":""}`)
	case domain.BlockTypeLaborTrips:
		return json.RawMessage(`{"eyebrow":"Трудовое послушание","title":"Поездки с добрым делом","description":"Информация о трудовых поездках паломнической службы.","items":[]}`)
	default:
		return json.RawMessage(`{}`)
	}
}

func (s *CMSService) BootstrapHomePage(ctx context.Context) (domain.Page, error) {
	if _, err := s.cms.GetPageBySlug(ctx, domain.PageSlugHome); err == nil {
		return domain.Page{}, domain.ErrDuplicateSlug
	} else if !errors.Is(err, domain.ErrNotFound) {
		return domain.Page{}, err
	}

	page, err := s.CreatePage(ctx, CreatePageInput{
		Slug:        domain.PageSlugHome,
		Title:       "Главная",
		Path:        domain.PagePathHome,
		IsPublished: true,
	})
	if err != nil {
		return domain.Page{}, err
	}

	for _, spec := range HomePageBlockSpecs() {
		if _, err := s.CreateBlock(ctx, page.ID, CreateBlockInput{
			Type:      spec.Type,
			Content:   spec.Content,
			IsVisible: true,
		}); err != nil {
			return domain.Page{}, err
		}
	}

	return s.GetPage(ctx, page.ID)
}

func BlockTemplates() []map[string]any {
	types := []string{
		domain.BlockTypeHero,
		domain.BlockTypeAbout,
		domain.BlockTypeWhyUs,
		domain.BlockTypeHowItWorks,
		domain.BlockTypeFAQ,
		domain.BlockTypeCTA,
		domain.BlockTypeRichText,
		domain.BlockTypeFeaturedRoute,
		domain.BlockTypePopularDestinations,
		domain.BlockTypeTestimonials,
		domain.BlockTypeInformationHero,
		domain.BlockTypeHolyPlaces,
		domain.BlockTypeServiceProfile,
		domain.BlockTypeServiceTeam,
		domain.BlockTypeLegalDonation,
		domain.BlockTypeLaborTrips,
	}
	labels := map[string]string{
		domain.BlockTypeHero:                "Шапка главной",
		domain.BlockTypeAbout:               "О службе",
		domain.BlockTypeWhyUs:               "Почему мы",
		domain.BlockTypeHowItWorks:          "Как записаться",
		domain.BlockTypeFAQ:                 "Частые вопросы",
		domain.BlockTypeCTA:                 "Баннер-призыв",
		domain.BlockTypeRichText:            "Текстовый блок",
		domain.BlockTypeFeaturedRoute:       "Тихвинский путь (главный маршрут)",
		domain.BlockTypePopularDestinations: "Популярные направления (виджет)",
		domain.BlockTypeTestimonials:        "Отзывы (виджет)",
		domain.BlockTypeInformationHero:     "Инфо-режим: приветственный баннер",
		domain.BlockTypeHolyPlaces:          "Инфо-режим: монастыри и святыни",
		domain.BlockTypeServiceProfile:      "Инфо-режим: о службе и полезные ссылки",
		domain.BlockTypeServiceTeam:         "Инфо-режим: руководители и сотрудники",
		domain.BlockTypeLegalDonation:       "Инфо-режим: юридические данные и пожертвования",
		domain.BlockTypeLaborTrips:          "Инфо-режим: трудовые поездки",
	}
	out := make([]map[string]any, 0, len(types))
	for _, t := range types {
		out = append(out, map[string]any{
			"type":    t,
			"label":   labels[t],
			"content": json.RawMessage(DefaultBlockContent(t)),
		})
	}
	return out
}
