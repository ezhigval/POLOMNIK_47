package maxfeed

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"palomnik/internal/ports"
)

func TestPublisherNotConfiguredWithoutChatID(t *testing.T) {
	p := New("token", "", "", time.Second)
	if p.Configured() {
		t.Fatal("expected unconfigured")
	}
	if err := p.Publish(context.Background(), "max_feed", ports.PublishContent{Title: "A"}); err != ports.ErrPublisherNotConfigured {
		t.Fatalf("got %v", err)
	}
}

func TestPublisherPostsToChat(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/messages" {
			t.Fatalf("path %s", r.URL.Path)
		}
		if r.URL.Query().Get("chat_id") != "99" {
			t.Fatalf("query %s", r.URL.RawQuery)
		}
		if r.Header.Get("Authorization") != "max-token" {
			t.Fatalf("authorization %q", r.Header.Get("Authorization"))
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	p := New("max-token", server.URL, "99", time.Second)
	if !p.Configured() {
		t.Fatal("expected configured")
	}
	if err := p.Publish(context.Background(), "max_feed", ports.PublishContent{Title: "Заголовок"}); err != nil {
		t.Fatalf("publish: %v", err)
	}
}
