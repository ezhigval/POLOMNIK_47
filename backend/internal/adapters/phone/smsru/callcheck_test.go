package smsru

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"polomnik/internal/ports"
)

func TestVerifierUnavailableWithoutAPIID(t *testing.T) {
	v := New("", time.Second)
	if v.Available() {
		t.Fatal("expected unavailable without api_id")
	}
	if _, err := v.Start(context.Background(), "+79001234567"); err != ports.ErrPhoneVerifierNotConfigured {
		t.Fatalf("start: %v", err)
	}
	if _, err := v.Status(context.Background(), "x"); err != ports.ErrPhoneVerifierNotConfigured {
		t.Fatalf("status: %v", err)
	}
}

func TestStartPostsCallcheckAddShape(t *testing.T) {
	var gotPath string
	var gotBody string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		body, _ := io.ReadAll(r.Body)
		gotBody = string(body)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"status":"OK",
			"status_code":100,
			"check_id":"201737-542",
			"call_phone":"78005008275",
			"call_phone_pretty":"+7 (800) 500-8275"
		}`))
	}))
	defer server.Close()

	v := New("test-api-id", time.Second)
	v.baseURL = server.URL

	challenge, err := v.Start(context.Background(), "+7 (900) 123-45-67")
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	if gotPath != "/callcheck/add" {
		t.Fatalf("path: %q", gotPath)
	}
	if !strings.Contains(gotBody, "api_id=test-api-id") {
		t.Fatalf("missing api_id in body: %q", gotBody)
	}
	if !strings.Contains(gotBody, "phone=79001234567") {
		t.Fatalf("phone shape: %q", gotBody)
	}
	if !strings.Contains(gotBody, "json=1") {
		t.Fatalf("missing json=1: %q", gotBody)
	}
	if challenge.CheckID != "201737-542" || challenge.CallPhone != "78005008275" {
		t.Fatalf("challenge: %+v", challenge)
	}
	if challenge.CallPhonePretty != "+7 (800) 500-8275" {
		t.Fatalf("pretty: %q", challenge.CallPhonePretty)
	}
}

func TestStatusMapsCodes(t *testing.T) {
	cases := []struct {
		checkStatus string
		want        ports.PhoneCheckStatus
	}{
		{"400", ports.PhoneCheckPending},
		{"401", ports.PhoneCheckConfirmed},
		{"402", ports.PhoneCheckExpired},
	}
	for _, tc := range cases {
		t.Run(tc.checkStatus, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/callcheck/status" {
					t.Fatalf("path: %s", r.URL.Path)
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = fmt.Fprintf(w, `{"status":"OK","status_code":100,"check_status":%q}`, tc.checkStatus)
			}))
			defer server.Close()

			v := New("test-api-id", time.Second)
			v.baseURL = server.URL
			got, err := v.Status(context.Background(), "201737-542")
			if err != nil {
				t.Fatalf("status: %v", err)
			}
			if got != tc.want {
				t.Fatalf("got %q want %q", got, tc.want)
			}
		})
	}
}
