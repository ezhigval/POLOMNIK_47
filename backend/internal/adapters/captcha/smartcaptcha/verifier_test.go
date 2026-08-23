package smartcaptcha

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"palomnik/internal/ports"
)

func TestVerifierNotConfiguredWithoutKeys(t *testing.T) {
	v := New("", "", time.Second)
	if v.Configured() {
		t.Fatal("expected unconfigured verifier")
	}
	if err := v.Verify(context.Background(), "token", "127.0.0.1"); err != ports.ErrCaptchaNotConfigured {
		t.Fatalf("expected not configured, got %v", err)
	}
}

func TestVerifierAcceptsOkStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		if !containsAll(string(body), "secret=server", "token=good-token") {
			t.Fatalf("unexpected form: %s", body)
		}
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}))
	defer server.Close()

	v := New("server", "client", time.Second)
	v.validateURL = server.URL
	if !v.Configured() {
		t.Fatal("expected configured verifier")
	}
	if err := v.Verify(context.Background(), "good-token", "127.0.0.1"); err != nil {
		t.Fatalf("verify: %v", err)
	}
}

func TestVerifierRejectsFailedStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"status":"failed"}`))
	}))
	defer server.Close()

	v := New("server", "client", time.Second)
	v.validateURL = server.URL
	if err := v.Verify(context.Background(), "bad", ""); err != ports.ErrCaptchaFailed {
		t.Fatalf("expected captcha failed, got %v", err)
	}
}

func containsAll(haystack string, parts ...string) bool {
	for _, part := range parts {
		if !contains(haystack, part) {
			return false
		}
	}
	return true
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 || (func() bool {
		for i := 0; i+len(sub) <= len(s); i++ {
			if s[i:i+len(sub)] == sub {
				return true
			}
		}
		return false
	})())
}
