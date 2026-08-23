package sitenews

import (
	"context"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"

	"palomnik/internal/application"
	"palomnik/internal/ports"
)

const excerptMaxRunes = 180

type Publisher struct {
	news *application.NewsService
}

func New(news *application.NewsService) Publisher {
	return Publisher{news: news}
}

var _ ports.PublisherPort = Publisher{}

func (p Publisher) Configured() bool {
	return p.news != nil
}

func (p Publisher) Publish(ctx context.Context, channel string, content ports.PublishContent) error {
	if !p.Configured() {
		return ports.ErrPublisherNotConfigured
	}
	if ch := strings.ToLower(strings.TrimSpace(channel)); ch != "" && ch != "site_news" {
		return ports.ErrPublisherNotConfigured
	}

	title := strings.TrimSpace(content.Title)
	body := strings.TrimSpace(content.Body)
	link := strings.TrimSpace(content.URL)
	if body == "" {
		body = link
	} else if link != "" {
		body = body + "\n\n" + link
	}
	if title == "" || body == "" {
		return ports.ErrPublisherNotConfigured
	}

	_, err := p.news.CreateNews(ctx, application.NewsArticleInput{
		Slug:        articleSlug(title),
		Title:       title,
		Excerpt:     excerptFrom(body, title),
		Body:        body,
		PublishedAt: time.Now().UTC(),
		IsPublished: true,
	})
	return err
}

func articleSlug(title string) string {
	slug := strings.ToLower(strings.TrimSpace(title))
	slug = strings.ReplaceAll(slug, " ", "-")
	if slug == "" {
		slug = "news"
	}
	return slug + "-" + uuid.NewString()[:8]
}

func excerptFrom(body, title string) string {
	text := strings.TrimSpace(body)
	if text == "" {
		text = strings.TrimSpace(title)
	}
	if utf8.RuneCountInString(text) <= excerptMaxRunes {
		return text
	}
	runes := []rune(text)
	return strings.TrimSpace(string(runes[:excerptMaxRunes])) + "…"
}
