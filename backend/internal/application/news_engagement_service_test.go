package application

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"palomnik/internal/adapters/repository/memory"
	"palomnik/internal/domain"
	"palomnik/internal/ports"
)

func TestNewsEngagementToggleLike(t *testing.T) {
	ctx := context.Background()
	store := memory.NewStore()
	service := NewNewsEngagementService(store, store, store)

	article, err := domain.NewNewsArticle(domain.NewNewsArticleInput{
		ID:          uuid.New(),
		Slug:        "test-news",
		Title:       "Test",
		Excerpt:     "Excerpt",
		Body:        "Body text",
		PublishedAt: time.Now().UTC(),
		IsPublished: true,
	})
	if err != nil {
		t.Fatalf("new article: %v", err)
	}
	if _, err := store.CreateNews(ctx, article); err != nil {
		t.Fatalf("create news: %v", err)
	}

	state, err := service.ToggleLike(ctx, "test-news", "visitor-1")
	if err != nil {
		t.Fatalf("toggle like: %v", err)
	}
	if !state.LikedByYou || state.Count != 1 {
		t.Fatalf("expected liked with count 1, got %+v", state)
	}

	state, err = service.ToggleLike(ctx, "test-news", "visitor-1")
	if err != nil {
		t.Fatalf("toggle unlike: %v", err)
	}
	if state.LikedByYou || state.Count != 0 {
		t.Fatalf("expected unliked with count 0, got %+v", state)
	}
}

func TestNewsEngagementAddCommentRequiresUser(t *testing.T) {
	ctx := context.Background()
	store := memory.NewStore()
	service := NewNewsEngagementService(store, store, store)

	article, err := domain.NewNewsArticle(domain.NewNewsArticleInput{
		ID:          uuid.New(),
		Slug:        "comment-news",
		Title:       "Test",
		Excerpt:     "Excerpt",
		Body:        "Body text",
		PublishedAt: time.Now().UTC(),
		IsPublished: true,
	})
	if err != nil {
		t.Fatalf("new article: %v", err)
	}
	if _, err := store.CreateNews(ctx, article); err != nil {
		t.Fatalf("create news: %v", err)
	}

	userID := uuid.New()
	user := domain.User{
		ID:        userID,
		Email:     "user@example.com",
		Name:      "Anna",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	if _, err := store.CreateUser(ctx, user); err != nil {
		t.Fatalf("create user: %v", err)
	}

	comment, err := service.AddComment(ctx, "comment-news", userID, "  Hello world  ")
	if err != nil {
		t.Fatalf("add comment: %v", err)
	}
	if comment.Author != "Anna" || comment.Body != "Hello world" {
		t.Fatalf("unexpected comment: %+v", comment)
	}

	list, err := service.ListComments(ctx, "comment-news", ports.Pagination{Page: 1, Limit: 20})
	if err != nil {
		t.Fatalf("list comments: %v", err)
	}
	if len(list.Items) != 1 || list.Meta.Total != 1 {
		t.Fatalf("expected 1 comment, got %+v", list)
	}
}

func TestNewsEngagementUnpublishedNewsNotFound(t *testing.T) {
	ctx := context.Background()
	store := memory.NewStore()
	service := NewNewsEngagementService(store, store, store)

	article, err := domain.NewNewsArticle(domain.NewNewsArticleInput{
		ID:          uuid.New(),
		Slug:        "draft-news",
		Title:       "Draft",
		Excerpt:     "Excerpt",
		Body:        "Body",
		PublishedAt: time.Now().UTC(),
		IsPublished: false,
	})
	if err != nil {
		t.Fatalf("new article: %v", err)
	}
	if _, err := store.CreateNews(ctx, article); err != nil {
		t.Fatalf("create news: %v", err)
	}

	_, err = service.GetLikeState(ctx, "draft-news", "visitor-1")
	if err == nil || err != domain.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
