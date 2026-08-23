package sber

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"

	"palomnik/internal/domain"
	"palomnik/internal/ports"
)

func TestNotConfiguredWithoutKeys(t *testing.T) {
	g := New("", "pass", "", "https://tikhvin-palomnik.ru/account/trips", time.Second)
	if g.Configured() {
		t.Fatal("expected unconfigured")
	}
	if _, err := g.CreatePayment(context.Background(), domain.Booking{ID: uuid.New(), TotalPrice: 1000}); err != ports.ErrPaymentNotConfigured {
		t.Fatalf("got %v", err)
	}
}

func TestCreatePaymentUsesTotalPriceKopecks(t *testing.T) {
	bookingID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/payment/rest/register.do" {
			t.Fatalf("path %s", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		form := string(body)
		if !strings.Contains(form, "amount=250000") || !strings.Contains(form, "orderNumber="+bookingID.String()) {
			t.Fatalf("form %s", form)
		}
		_, _ = w.Write([]byte(`{"orderId":"sber-order","formUrl":"https://pay.sber/form","errorCode":"0"}`))
	}))
	defer server.Close()

	g := New("user", "secret", server.URL, "https://tikhvin-palomnik.ru/account/trips", time.Second)
	session, err := g.CreatePayment(context.Background(), domain.Booking{ID: bookingID, TotalPrice: 2500})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if session.ID != "sber-order" || session.RedirectURL != "https://pay.sber/form" {
		t.Fatalf("session %+v", session)
	}
}

func TestCreatePaymentRejectsZeroAmount(t *testing.T) {
	g := New("user", "secret", "https://example.invalid", "https://tikhvin-palomnik.ru/account/trips", time.Second)
	if _, err := g.CreatePayment(context.Background(), domain.Booking{ID: uuid.New(), TotalPrice: 0}); err != domain.ErrInvalidTotalPrice {
		t.Fatalf("got %v", err)
	}
}
