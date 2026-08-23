package sitenews

import (
	"context"
	"strings"
	"testing"

	"palomnik/internal/adapters/repository/memory"
	"palomnik/internal/application"
	"palomnik/internal/ports"
)

func TestPublisherNotConfiguredWithoutNews(t *testing.T) {
	p := New(nil)
	if p.Configured() {
		t.Fatal("expected unconfigured")
	}
	if err := p.Publish(context.Background(), "site_news", ports.PublishContent{Title: "A", Body: "B"}); err != ports.ErrPublisherNotConfigured {
		t.Fatalf("got %v", err)
	}
}

func TestPublisherCreatesNewsFromContent(t *testing.T) {
	store := memory.NewStore()
	news := application.NewNewsService(store, nil)
	p := New(news)

	if err := p.Publish(context.Background(), "site_news", ports.PublishContent{
		Title: "Тихвинский путь",
		Body:  "Текст новости",
		URL:   "https://tikhvin-palomnik.ru/news/a",
	}); err != nil {
		t.Fatalf("publish: %v", err)
	}

	list, err := news.ListNews(context.Background(), ports.Pagination{Page: 1, Limit: 10})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(list.Items) != 1 {
		t.Fatalf("items %d", len(list.Items))
	}
	article := list.Items[0]
	if article.Title != "Тихвинский путь" {
		t.Fatalf("title %q", article.Title)
	}
	if !strings.Contains(article.Body, "Текст новости") || !strings.Contains(article.Body, "https://tikhvin-palomnik.ru/news/a") {
		t.Fatalf("body %q", article.Body)
	}
	if !article.IsPublished {
		t.Fatal("expected published row")
	}
	if article.Slug == "" {
		t.Fatal("expected slug")
	}
}

func TestPublisherRejectsOtherChannel(t *testing.T) {
	p := New(application.NewNewsService(memory.NewStore(), nil))
	if err := p.Publish(context.Background(), "vk_wall", ports.PublishContent{Title: "A", Body: "B"}); err != ports.ErrPublisherNotConfigured {
		t.Fatalf("got %v", err)
	}
}
