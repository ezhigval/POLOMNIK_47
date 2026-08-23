package dto

import (
	"time"

	"palomnik/internal/domain"
)

type SMMChannelResultResponse struct {
	Channel     string `json:"channel"`
	OK          bool   `json:"ok"`
	Error       string `json:"error,omitempty"`
	AttemptedAt string `json:"attempted_at,omitempty"`
}

type SMMPostResponse struct {
	ID          string                     `json:"id"`
	Title       string                     `json:"title"`
	Body        string                     `json:"body"`
	URL         string                     `json:"url"`
	PublishAt   string                     `json:"publish_at"`
	Channels    []string                   `json:"channels"`
	PublishedAt string                     `json:"published_at,omitempty"`
	Results     []SMMChannelResultResponse `json:"results"`
	CreatedAt   string                     `json:"created_at"`
}

type SMMPostCreateRequest struct {
	Title     string   `json:"title"`
	Body      string   `json:"body"`
	URL       string   `json:"url"`
	PublishAt string   `json:"publish_at"`
	Channels  []string `json:"channels"`
}

func ToSMMPostResponse(post domain.SMMPost) SMMPostResponse {
	results := make([]SMMChannelResultResponse, 0, len(post.Results))
	for _, item := range post.Results {
		results = append(results, SMMChannelResultResponse{
			Channel:     item.Channel,
			OK:          item.OK,
			Error:       item.Error,
			AttemptedAt: formatTime(item.AttemptedAt),
		})
	}
	published := ""
	if post.PublishedAt != nil {
		published = formatTime(*post.PublishedAt)
	}
	return SMMPostResponse{
		ID:          post.ID.String(),
		Title:       post.Title,
		Body:        post.Body,
		URL:         post.URL,
		PublishAt:   formatTime(post.PublishAt),
		Channels:    append([]string(nil), post.Channels...),
		PublishedAt: published,
		Results:     results,
		CreatedAt:   formatTime(post.CreatedAt),
	}
}

func formatTime(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.UTC().Format(time.RFC3339)
}
