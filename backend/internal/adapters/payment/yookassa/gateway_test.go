package yookassa

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"

	"palomnik/internal/domain"
	"palomnik/internal/ports"
)

func TestNotConfiguredWithoutKeys(t *testing.T) {
	g := New("shop", "", "", "https://tikhvin-palomnik.ru/account/trips", time.Second)
	if g.Configured() {
		t.Fatal("expected unconfigured")
	}
	if _, err := g.CreatePayment(context.Background(), domain.Booking{ID: uuid.New(), TotalPrice: 1000}); err != ports.ErrPaymentNotConfigured {
		t.Fatalf("got %v", err)
	}
}

func TestCreatePaymentUsesTotalPrice(t *testing.T) {
	bookingID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v3/payments" {
			t.Fatalf("path %s", r.URL.Path)
		}
		if r.Header.Get("Idempotence-Key") != bookingID.String() {
			t.Fatalf("idempotence %q", r.Header.Get("Idempotence-Key"))
		}
		body, _ := io.ReadAll(r.Body)
		var payload map[string]any
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("json: %v", err)
		}
		amount, _ := payload["amount"].(map[string]any)
		if amount["value"] != "2500.00" || amount["currency"] != "RUB" {
			t.Fatalf("amount %+v", amount)
		}
		_, _ = w.Write([]byte(`{"id":"yk-1","confirmation":{"confirmation_url":"https://yookassa.ru/pay"}}`))
	}))
	defer server.Close()

	g := New("shop", "secret", server.URL, "https://tikhvin-palomnik.ru/account/trips", time.Second)
	session, err := g.CreatePayment(context.Background(), domain.Booking{ID: bookingID, TotalPrice: 2500})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if session.ID != "yk-1" || session.RedirectURL != "https://yookassa.ru/pay" {
		t.Fatalf("session %+v", session)
	}
}
