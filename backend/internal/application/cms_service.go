package application

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"

	"polomnik/internal/domain"
	"polomnik/internal/ports"
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
	Title       *string
	Path        *string
	IsPublished *bool
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

func (s *CMSService) GetPublishedPageByPath(ctx context.Context, path string) (domain.Page, error) {
	page, err := s.cms.GetPageByPath(ctx, path)
	if err != nil {
		return domain.Page{}, err
	}
	if !page.IsPublished {
		return domain.Page{}, errNotPublished()
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

func errNotPublished() error {
	return domain.ErrNotFound
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
		return json.RawMessage(`{"eyebrow":"","title":"Заголовок","subtitle":"","primaryCta":"Выбрать тур","primaryHref":"/search","secondaryCta":"","secondaryHref":"/#how-it-works","stats":[]}`)
	case domain.BlockTypeAbout:
		return json.RawMessage(`{"eyebrow":"О службе","title":"О нас","paragraphs":[""],"highlights":[],"showContacts":true}`)
	case domain.BlockTypeWhyUs:
		return json.RawMessage(`{"eyebrow":"Почему мы","title":"Почему нам доверяют","description":"","items":[]}`)
	case domain.BlockTypeHowItWorks:
		return json.RawMessage(`{"eyebrow":"Просто","title":"Как записаться","description":"","steps":[],"ctaLabel":"Записаться","ctaHref":"/search"}`)
	case domain.BlockTypeFAQ:
		return json.RawMessage(`{"eyebrow":"Вопросы","title":"Частые вопросы","description":"","items":[]}`)
	case domain.BlockTypeCTA:
		return json.RawMessage(`{"title":"Готовы отправиться?","subtitle":"","button":"Смотреть туры","href":"/search"}`)
	case domain.BlockTypeRichText:
		return json.RawMessage(`{"eyebrow":"","title":"Текст","body":""}`)
	case domain.BlockTypePopularDestinations, domain.BlockTypeTestimonials:
		return json.RawMessage(`{}`)
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

type homeBlockSpec struct {
	Type    string
	Content json.RawMessage
}

func HomePageBlockSpecs() []homeBlockSpec {
	return []homeBlockSpec{
		{
			Type: domain.BlockTypeHero,
			Content: json.RawMessage(`{
				"eyebrow":"Паломнические туры по России",
				"title":"Путь к святыням — с заботой и сопровождением",
				"subtitle":"Организуем поездки в монастыри и святые места: продуманная программа, комфортный транспорт и духовник на маршруте. Вы выбираете тур — мы берём на себя всё остальное.",
				"primaryCta":"Выбрать тур",
				"primaryHref":"/search",
				"secondaryCta":"Как записаться",
				"secondaryHref":"/#how-it-works",
				"stats":[
					{"value":"12+","label":"лет опыта"},
					{"value":"3 000+","label":"паломников в год"},
					{"value":"40+","label":"направлений"},
					{"value":"4.9","label":"средняя оценка"}
				]
			}`),
		},
		{Type: domain.BlockTypePopularDestinations, Content: json.RawMessage(`{}`)},
		{
			Type: domain.BlockTypeAbout,
			Content: json.RawMessage(`{
				"eyebrow":"О службе",
				"title":"Паломничество без лишней суеты",
				"paragraphs":[
					"Мы организуем паломнические туры по России: от классических направлений — Оптина, Дивеево, Валаам — до сезонных поездок в монастыри и святые места Северо-Запада.",
					"Наша задача — создать спокойное пространство для паломничества: комфортный транспорт, проверенное размещение, понятная программа и человек на связи, если возникнут вопросы.",
					"Запись через сайт простая: вы оставляете заявку, менеджер перезванивает, уточняет детали и подтверждает участие."
				],
				"highlights":[
					"Групповые поездки с продуманной программой",
					"Сопровождение духовника на маршруте",
					"Помощь менеджера на каждом этапе — от заявки до возвращения",
					"Выезд из Санкт-Петербурга и других городов по согласованию"
				],
				"showContacts":true
			}`),
		},
		{
			Type: domain.BlockTypeWhyUs,
			Content: json.RawMessage(`{
				"eyebrow":"Почему мы",
				"title":"Почему нам доверяют",
				"description":"",
				"items":[
					{"title":"Продуманная программа","description":"Маршрут, проживание и питание согласованы заранее. Вы знаете, куда едете и что вас ждёт — без сюрпризов в дороге.","icon":"route"},
					{"title":"Сопровождение духовника","description":"На каждой поездке — священнослужитель, который помогает сосредоточиться на главном: молитве, покаянии и встрече со святынями.","icon":"cross"},
					{"title":"Комфорт и безопасность","description":"Современный транспорт, проверенные гостиницы, страховка и круглосуточная связь с координатором группы.","icon":"shield"},
					{"title":"Прозрачная стоимость","description":"Цена на сайте — за человека, без скрытых доплат. Менеджер заранее расскажет, что входит в тур.","icon":"wallet"}
				]
			}`),
		},
		{
			Type: domain.BlockTypeHowItWorks,
			Content: json.RawMessage(`{
				"eyebrow":"Просто",
				"title":"Как записаться",
				"description":"Три шага — и вы в списке участников. Никаких личных кабинетов и сложных форм.",
				"steps":[
					{"title":"Выберите тур","description":"Изучите даты, программу и стоимость. Фильтры помогут найти подходящее направление."},
					{"title":"Оставьте заявку","description":"Укажите имя, телефон и количество участников — без регистрации и оплаты на сайте."},
					{"title":"Подтверждение","description":"Менеджер перезвонит, ответит на вопросы и подтвердит ваше участие в группе."}
				],
				"ctaLabel":"Записаться",
				"ctaHref":"/search"
			}`),
		},
		{Type: domain.BlockTypeTestimonials, Content: json.RawMessage(`{}`)},
		{
			Type: domain.BlockTypeFAQ,
			Content: json.RawMessage(`{
				"eyebrow":"Вопросы",
				"title":"Частые вопросы",
				"description":"Не нашли ответ? Позвоните или оставьте заявку — менеджер поможет.",
				"items":[
					{"question":"Нужна ли предоплата при записи?","answer":"Нет. Вы оставляете заявку на сайте, менеджер связывается с вами, уточняет детали и рассказывает о порядке бронирования и оплаты."},
					{"question":"Можно ли поехать одному или только группой?","answer":"Можно и одному, и с семьёй. Мы формируем группы из паломников с разным составом — от одиночных участников до больших семей."},
					{"question":"Что входит в стоимость тура?","answer":"Обычно — трансфер, проживание, питание, сопровождение духовника и организационные расходы. Точный состав указан в описании каждого тура."},
					{"question":"Есть ли ограничения по возрасту или здоровью?","answer":"Большинство туров подходят для взрослых и детей от 7 лет. Если у вас есть особые потребности — напишите в комментарии к заявке, мы подберём подходящий маршрут."},
					{"question":"Как отменить или перенести участие?","answer":"Свяжитесь с менеджером по телефону. Условия возврата зависят от срока до выезда — мы объясним их при подтверждении брони."}
				]
			}`),
		},
		{
			Type: domain.BlockTypeCTA,
			Content: json.RawMessage(`{
				"title":"Готовы отправиться в путь?",
				"subtitle":"Оставьте заявку — менеджер перезвонит в течение рабочего дня и ответит на все вопросы.",
				"button":"Смотреть туры",
				"href":"/search"
			}`),
		},
	}
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
		domain.BlockTypePopularDestinations,
		domain.BlockTypeTestimonials,
	}
	labels := map[string]string{
		domain.BlockTypeHero:                "Hero / шапка",
		domain.BlockTypeAbout:               "О службе",
		domain.BlockTypeWhyUs:               "Почему мы",
		domain.BlockTypeHowItWorks:          "Как записаться",
		domain.BlockTypeFAQ:                 "Частые вопросы",
		domain.BlockTypeCTA:                 "CTA-баннер",
		domain.BlockTypeRichText:            "Текстовый блок",
		domain.BlockTypePopularDestinations: "Популярные направления (виджет)",
		domain.BlockTypeTestimonials:        "Отзывы (виджет)",
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
