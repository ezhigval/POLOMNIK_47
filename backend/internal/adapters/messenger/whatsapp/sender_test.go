package whatsapp

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

func TestSenderNotConfiguredWithoutKeys(t *testing.T) {
	s := New("", "", "", time.Second)
	if s.Configured() {
		t.Fatal("expected unconfigured")
	}
	if err := s.Send(context.Background(), "whatsapp", "+79001234567", "hi"); err != ports.ErrMessengerNotConfigured {
		t.Fatalf("got %v", err)
	}
}

func TestSenderPostsCloudAPI(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/PHONEID/messages" {
			t.Fatalf("path %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer waba-token" {
			t.Fatalf("authorization %q", r.Header.Get("Authorization"))
		}
		body, _ := io.ReadAll(r.Body)
		raw := string(body)
		if !strings.Contains(raw, `"messaging_product":"whatsapp"`) || !strings.Contains(raw, `"to":"79001234567"`) {
			t.Fatalf("body %s", raw)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"messages":[{"id":"wamid.1"}]}`))
	}))
	defer server.Close()

	s := New("waba-token", "PHONEID", server.URL, time.Second)
	if !s.Configured() {
		t.Fatal("expected configured")
	}
	if err := s.Send(context.Background(), "whatsapp", "+79001234567", "hello"); err != nil {
		t.Fatalf("send: %v", err)
	}
}
