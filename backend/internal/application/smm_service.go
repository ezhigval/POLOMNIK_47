package application

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"palomnik/internal/domain"
	"palomnik/internal/ports"
)

type SMMService struct {
	posts     ports.SMMPostRepository
	publisher ports.PublisherPort
}

func NewSMMService(posts ports.SMMPostRepository, publisher ports.PublisherPort) *SMMService {
	return &SMMService{posts: posts, publisher: publisher}
}

type CreateSMMPostInput struct {
	Title     string
	Body      string
	URL       string
	PublishAt time.Time
	Channels  []string
}

func (s *SMMService) CreatePost(ctx context.Context, input CreateSMMPostInput) (domain.SMMPost, error) {
	if s == nil || s.posts == nil {
		return domain.SMMPost{}, domain.ErrNotFound
	}
	post, err := domain.NewSMMPost(domain.NewSMMPostInput{
		ID:        uuid.New(),
		Title:     input.Title,
		Body:      input.Body,
		URL:       input.URL,
		PublishAt: input.PublishAt,
		Channels:  input.Channels,
	})
	if err != nil {
		return domain.SMMPost{}, err
	}
	return s.posts.CreateSMMPost(ctx, post)
}

func (s *SMMService) GetPost(ctx context.Context, id uuid.UUID) (domain.SMMPost, error) {
	if s == nil || s.posts == nil {
		return domain.SMMPost{}, domain.ErrNotFound
	}
	return s.posts.GetSMMPost(ctx, id)
}

func (s *SMMService) ListPosts(ctx context.Context, pagination ports.Pagination) (ports.SMMPostList, error) {
	if s == nil || s.posts == nil {
		return ports.SMMPostList{}, nil
	}
	return s.posts.ListSMMPosts(ctx, pagination)
}

func (s *SMMService) DeletePost(ctx context.Context, id uuid.UUID) error {
	if s == nil || s.posts == nil {
		return domain.ErrNotFound
	}
	return s.posts.DeleteSMMPost(ctx, id)
}

func (s *SMMService) PublishPost(ctx context.Context, id uuid.UUID) (domain.SMMPost, error) {
	if s == nil || s.posts == nil {
		return domain.SMMPost{}, domain.ErrNotFound
	}
	post, err := s.posts.GetSMMPost(ctx, id)
	if err != nil {
		return domain.SMMPost{}, err
	}
	return s.publish(ctx, post)
}

func (s *SMMService) PublishDue(ctx context.Context, now time.Time) (int, error) {
	if s == nil || s.posts == nil {
		return 0, nil
	}
	due, err := s.posts.ListDueSMMPosts(ctx, now)
	if err != nil {
		return 0, err
	}
	count := 0
	for _, post := range due {
		if _, err := s.publish(ctx, post); err != nil {
			return count, err
		}
		count++
	}
	return count, nil
}

func (s *SMMService) publish(ctx context.Context, post domain.SMMPost) (domain.SMMPost, error) {
	now := time.Now().UTC()
	results := make([]domain.SMMChannelResult, 0, len(post.Channels))
	content := ports.PublishContent{Title: post.Title, Body: post.Body, URL: post.URL}
	for _, channel := range post.Channels {
		result := domain.SMMChannelResult{Channel: channel, AttemptedAt: now}
		if s.publisher == nil || !s.publisher.Configured() {
			result.Error = ports.ErrPublisherNotConfigured.Error()
		} else if err := s.publisher.Publish(ctx, channel, content); err != nil {
			if errors.Is(err, ports.ErrPublisherNotConfigured) {
				result.Error = ports.ErrPublisherNotConfigured.Error()
			} else {
				result.Error = err.Error()
			}
		} else {
			result.OK = true
		}
		results = append(results, result)
	}
	post.Results = results
	post.PublishedAt = &now
	post.UpdatedAt = now
	return s.posts.SaveSMMPost(ctx, post)
}
