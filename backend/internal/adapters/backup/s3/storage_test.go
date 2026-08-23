package s3

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"polomnik/internal/ports"
)

func TestStorageNotConfiguredWithoutKeys(t *testing.T) {
	store := New("", "", "", "", "", "", time.Second)
	if store.Configured() {
		t.Fatal("expected unconfigured storage")
	}
	if err := store.Upload(context.Background(), ports.BackupObject{Name: "a.sql.gz", Body: bytes.NewReader(nil)}); err != ports.ErrBackupStorageNotConfigured {
		t.Fatalf("expected not configured, got %v", err)
	}
}

func TestStorageUploadsWithSigV4(t *testing.T) {
	var gotAuth string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("method %s", r.Method)
		}
		if !strings.Contains(r.URL.Path, "/polomnik/nightly/dump.sql.gz") {
			t.Fatalf("path %s", r.URL.Path)
		}
		gotAuth = r.Header.Get("Authorization")
		body, _ := io.ReadAll(r.Body)
		if string(body) != "dump" {
			t.Fatalf("body %q", body)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	store := New(server.URL, "ru-central1", "polomnik", "key", "secret", "nightly", time.Second)
	if !store.Configured() {
		t.Fatal("expected configured storage")
	}
	err := store.Upload(context.Background(), ports.BackupObject{
		Name: "dump.sql.gz",
		Body: bytes.NewReader([]byte("dump")),
	})
	if err != nil {
		t.Fatalf("upload: %v", err)
	}
	if !strings.HasPrefix(gotAuth, "AWS4-HMAC-SHA256 Credential=key/") {
		t.Fatalf("authorization %q", gotAuth)
	}
}
