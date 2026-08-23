package vkwall

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"palomnik/internal/ports"
)

func TestPublisherNotConfiguredWithoutKeys(t *testing.T) {
	p := New("token", "", "", "", time.Second)
	if p.Configured() {
		t.Fatal("expected unconfigured")
	}
	if err := p.Publish(context.Background(), "vk_wall", ports.PublishContent{Title: "A"}); err != ports.ErrPublisherNotConfigured {
		t.Fatalf("got %v", err)
	}
}

func TestPublisherPostsWall(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/method/wall.post" {
			t.Fatalf("path %s", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		values, err := url.ParseQuery(string(body))
		if err != nil {
			t.Fatalf("form: %v", err)
		}
		if values.Get("access_token") != "vk-token" {
			t.Fatalf("token %q", values.Get("access_token"))
		}
		if values.Get("owner_id") != "-123456" {
			t.Fatalf("owner_id %q", values.Get("owner_id"))
		}
		if values.Get("from_group") != "1" {
			t.Fatalf("from_group %q", values.Get("from_group"))
		}
		if values.Get("v") != "5.199" {
			t.Fatalf("v %q", values.Get("v"))
		}
		if values.Get("message") == "" {
			t.Fatal("empty message")
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":{"post_id":17}}`))
	}))
	defer server.Close()

	p := New("vk-token", "-123456", server.URL, "", time.Second)
	if err := p.Publish(context.Background(), "vk_wall", ports.PublishContent{Title: "Заголовок", Body: "Текст"}); err != nil {
		t.Fatalf("publish: %v", err)
	}
}

func TestPublisherSurfacesVKError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"error":{"error_code":5,"error_msg":"User authorization failed"}}`))
	}))
	defer server.Close()

	p := New("vk-token", "-1", server.URL, "5.199", time.Second)
	if err := p.Publish(context.Background(), "vk_wall", ports.PublishContent{Title: "A"}); err == nil {
		t.Fatal("expected vk error")
	}
}
