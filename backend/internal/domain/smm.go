package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	PublisherChannelSiteNews  = "site_news"
	PublisherChannelTelegram  = "telegram_channel"
	PublisherChannelVKWall    = "vk_wall"
	PublisherChannelMaxFeed   = "max_feed"
)

func AllPublisherChannels() []string {
	return []string{
		PublisherChannelSiteNews,
		PublisherChannelTelegram,
		PublisherChannelVKWall,
		PublisherChannelMaxFeed,
	}
}

func ValidPublisherChannel(channel string) bool {
	switch strings.ToLower(strings.TrimSpace(channel)) {
	case PublisherChannelSiteNews, PublisherChannelTelegram, PublisherChannelVKWall, PublisherChannelMaxFeed:
		return true
	default:
		return false
	}
}

type SMMChannelResult struct {
	Channel     string    `json:"channel"`
	OK          bool      `json:"ok"`
	Error       string    `json:"error"`
	AttemptedAt time.Time `json:"attempted_at"`
}

type SMMPost struct {
	ID          uuid.UUID
	Title       string
	Body        string
	URL         string
	PublishAt   time.Time
	Channels    []string
	PublishedAt *time.Time
	Results     []SMMChannelResult
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type NewSMMPostInput struct {
	ID        uuid.UUID
	Title     string
	Body      string
	URL       string
	PublishAt time.Time
	Channels  []string
	Now       time.Time
}

func NewSMMPost(input NewSMMPostInput) (SMMPost, error) {
	if input.ID == uuid.Nil {
		return SMMPost{}, ErrInvalidID
	}
	title := strings.TrimSpace(input.Title)
	body := strings.TrimSpace(input.Body)
	if title == "" {
		return SMMPost{}, ErrInvalidTitle
	}
	if body == "" {
		return SMMPost{}, ErrInvalidArticleBody
	}
	if input.PublishAt.IsZero() {
		return SMMPost{}, ErrInvalidPublishedAt
	}
	channels, err := NormalizePublisherChannels(input.Channels)
	if err != nil {
		return SMMPost{}, err
	}

	now := input.Now
	if now.IsZero() {
		now = time.Now().UTC()
	}
	return SMMPost{
		ID:        input.ID,
		Title:     title,
		Body:      body,
		URL:       strings.TrimSpace(input.URL),
		PublishAt: input.PublishAt.UTC(),
		Channels:  channels,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func NormalizePublisherChannels(raw []string) ([]string, error) {
	seen := make(map[string]struct{})
	out := make([]string, 0, len(raw))
	for _, item := range raw {
		ch := strings.ToLower(strings.TrimSpace(item))
		if ch == "" {
			continue
		}
		if !ValidPublisherChannel(ch) {
			return nil, ErrInvalidPublisherChannel
		}
		if _, ok := seen[ch]; ok {
			continue
		}
		seen[ch] = struct{}{}
		out = append(out, ch)
	}
	if len(out) == 0 {
		return nil, ErrInvalidPublisherChannel
	}
	return out, nil
}
