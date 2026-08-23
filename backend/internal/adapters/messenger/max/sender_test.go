package maxmsg

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"palomnik/internal/ports"
)

func TestSenderNotConfiguredWithoutToken(t *testing.T) {
	s := New("", "", time.Second)
	if s.Configured() {
		t.Fatal("expected unconfigured")
	}
	if err := s.Send(context.Background(), "max", "1", "hi"); err != ports.ErrMessengerNotConfigured {
		t.Fatalf("got %v", err)
	}
}

func TestSenderPostsToUser(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/messages" {
			t.Fatalf("path %s", r.URL.Path)
		}
		if r.URL.Query().Get("user_id") != "42" {
			t.Fatalf("query %s", r.URL.RawQuery)
		}
		if r.Header.Get("Authorization") != "secret-token" {
			t.Fatalf("authorization %q", r.Header.Get("Authorization"))
		}
		body, _ := io.ReadAll(r.Body)
		if !strings.Contains(string(body), `"text":"hello"`) {
			t.Fatalf("body %s", body)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	s := New("secret-token", server.URL, time.Second)
	if err := s.Send(context.Background(), "max", "42", "hello"); err != nil {
		t.Fatalf("send: %v", err)
	}
}

func TestSenderPostsToChatPrefix(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("chat_id") != "99" {
			t.Fatalf("query %s", r.URL.RawQuery)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	s := New("secret-token", server.URL, time.Second)
	if err := s.Send(context.Background(), "max", "chat:99", "hello"); err != nil {
		t.Fatalf("send: %v", err)
	}
}
