package yandexgpt

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"palomnik/internal/ports"
)

func TestClientNotConfiguredWithoutKeys(t *testing.T) {
	c := New("key", "", "", "", time.Second)
	if c.Configured() {
		t.Fatal("expected unconfigured")
	}
	if _, err := c.Complete(context.Background(), ports.AICompletionRequest{User: "hi"}); err != ports.ErrAINotConfigured {
		t.Fatalf("got %v", err)
	}
}

func TestClientPostsCompletion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/foundationModels/v1/completion" {
			t.Fatalf("path %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Api-Key secret-key" {
			t.Fatalf("authorization %q", r.Header.Get("Authorization"))
		}
		if r.Header.Get("x-folder-id") != "folder1" {
			t.Fatalf("folder %q", r.Header.Get("x-folder-id"))
		}
		body, _ := io.ReadAll(r.Body)
		var payload map[string]any
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("json: %v", err)
		}
		if payload["modelUri"] != "gpt://folder1/yandexgpt-lite" {
			t.Fatalf("modelUri %v", payload["modelUri"])
		}
		messages, _ := payload["messages"].([]any)
		if len(messages) != 2 {
			t.Fatalf("messages %v", payload["messages"])
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"result":{"alternatives":[{"message":{"role":"assistant","text":"ответ"}}]}}`))
	}))
	defer server.Close()

	c := New("secret-key", "folder1", "", server.URL, time.Second)
	if !c.Configured() {
		t.Fatal("expected configured")
	}
	got, err := c.Complete(context.Background(), ports.AICompletionRequest{
		System: "системный",
		User:   "вопрос",
	})
	if err != nil {
		t.Fatalf("complete: %v", err)
	}
	if got != "ответ" {
		t.Fatalf("got %q", got)
	}
}

func TestClientUsesFullModelURI(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		if !strings.Contains(string(body), `"modelUri":"gpt://folder1/yandexgpt/latest"`) {
			t.Fatalf("body %s", body)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"result":{"alternatives":[{"message":{"text":"ok"}}]}}`))
	}))
	defer server.Close()

	c := New("secret-key", "folder1", "gpt://folder1/yandexgpt/latest", server.URL, time.Second)
	if _, err := c.Complete(context.Background(), ports.AICompletionRequest{User: "hi"}); err != nil {
		t.Fatalf("complete: %v", err)
	}
}

func TestClientRejectsEmptyUser(t *testing.T) {
	c := New("secret-key", "folder1", "", "http://127.0.0.1", time.Second)
	if _, err := c.Complete(context.Background(), ports.AICompletionRequest{System: "sys"}); err != ports.ErrAINotConfigured {
		t.Fatalf("got %v", err)
	}
}
