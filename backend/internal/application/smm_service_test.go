package application

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"palomnik/internal/adapters/repository/memory"
	"palomnik/internal/domain"
	"palomnik/internal/ports"
)

type recordingPublisher struct {
	configured bool
	fail       map[string]error
	mu         sync.Mutex
	channels   []string
}

func (p *recordingPublisher) Configured() bool { return p.configured }

func (p *recordingPublisher) Publish(_ context.Context, channel string, _ ports.PublishContent) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.channels = append(p.channels, channel)
	if err, ok := p.fail[channel]; ok {
		return err
	}
	return nil
}

func TestSMMPublishDoesNotRollBackOtherChannels(t *testing.T) {
	store := memory.NewStore()
	pub := &recordingPublisher{
		configured: true,
		fail:       map[string]error{"vk_wall": errors.New("vk down")},
	}
	svc := NewSMMService(store, pub)
	post, err := svc.CreatePost(context.Background(), CreateSMMPostInput{
		Title:     "Паломничество",
		Body:      "Текст из плана",
		PublishAt: time.Now().UTC(),
		Channels:  []string{"site_news", "vk_wall"},
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	saved, err := svc.PublishPost(context.Background(), post.ID)
	if err != nil {
		t.Fatalf("publish: %v", err)
	}
	if saved.PublishedAt == nil {
		t.Fatal("expected published_at")
	}
	if len(saved.Results) != 2 {
		t.Fatalf("results: %+v", saved.Results)
	}
	ok, failed := 0, 0
	for _, item := range saved.Results {
		if item.OK {
			ok++
		} else {
			failed++
		}
	}
	if ok != 1 || failed != 1 {
		t.Fatalf("want 1 ok 1 fail, got %+v", saved.Results)
	}
}

func TestSMMNoopPublisherRecordsFailure(t *testing.T) {
	store := memory.NewStore()
	svc := NewSMMService(store, &recordingPublisher{configured: false})
	post, err := svc.CreatePost(context.Background(), CreateSMMPostInput{
		Title:     "Паломничество",
		Body:      "Текст из плана",
		PublishAt: time.Now().UTC().Add(-time.Minute),
		Channels:  []string{"telegram_channel"},
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	n, err := svc.PublishDue(context.Background(), time.Now().UTC())
	if err != nil || n != 1 {
		t.Fatalf("due: n=%d err=%v", n, err)
	}
	saved, err := svc.GetPost(context.Background(), post.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if len(saved.Results) != 1 || saved.Results[0].OK {
		t.Fatalf("noop must not post: %+v", saved.Results)
	}
}

func TestSMMRejectsUnknownChannel(t *testing.T) {
	store := memory.NewStore()
	svc := NewSMMService(store, &recordingPublisher{configured: true})
	_, err := svc.CreatePost(context.Background(), CreateSMMPostInput{
		Title:     "Паломничество",
		Body:      "Текст",
		PublishAt: time.Now().UTC(),
		Channels:  []string{"instagram"},
	})
	if !errors.Is(err, domain.ErrInvalidPublisherChannel) {
		t.Fatalf("got %v", err)
	}
}
