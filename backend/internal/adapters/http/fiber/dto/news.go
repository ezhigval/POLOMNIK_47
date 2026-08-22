package dto

import (
	"polomnik/internal/domain"
)

type NewsArticleResponse struct {
	ID          string `json:"id"`
	Slug        string `json:"slug"`
	Title       string `json:"title"`
	Excerpt     string `json:"excerpt"`
	Body        string `json:"body"`
	ImageURL    string `json:"image_url"`
	PublishedAt string `json:"published_at"`
	IsPublished bool   `json:"is_published"`
	SortOrder   int    `json:"sort_order"`
}

type NewsArticleUpsertRequest struct {
	Slug        string `json:"slug"`
	Title       string `json:"title"`
	Excerpt     string `json:"excerpt"`
	Body        string `json:"body"`
	ImageURL    string `json:"image_url"`
	PublishedAt string `json:"published_at"`
	IsPublished bool   `json:"is_published"`
	SortOrder   int    `json:"sort_order"`
}

func ToNewsArticleResponse(article domain.NewsArticle) NewsArticleResponse {
	return NewsArticleResponse{
		ID:          article.ID.String(),
		Slug:        article.Slug,
		Title:       article.Title,
		Excerpt:     article.Excerpt,
		Body:        article.Body,
		ImageURL:    article.ImageURL,
		PublishedAt: article.PublishedAt.UTC().Format("2006-01-02"),
		IsPublished: article.IsPublished,
		SortOrder:   article.SortOrder,
	}
}
