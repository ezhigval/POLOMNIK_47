package onec

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"polomnik/internal/config"
	"polomnik/internal/domain"
	"polomnik/internal/ports"
)

func TestODataAccountingAdapterExportBooking(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			_ = json.NewEncoder(w).Encode(map[string]any{"value": []any{}})
		case http.MethodPost:
			_ = json.NewEncoder(w).Encode(map[string]string{"Ref_Key": "order-ref-1"})
		default:
			t.Fatalf("unexpected method: %s", r.Method)
		}
	}))
	defer server.Close()

	adapter := NewODataAccountingAdapter(config.Config{
		OneCBaseURL:            server.URL,
		IntegrationHTTPTimeout: 2 * time.Second,
	})

	booking := mustConfirmedBooking(t)
	result, err := adapter.ExportBooking(t.Context(), booking)
	if err != nil {
		t.Fatalf("export booking: %v", err)
	}
	if result.ExternalID != "order-ref-1" {
		t.Fatalf("expected order-ref-1, got %q", result.ExternalID)
	}
}

func TestODataAccountingAdapterSkipsNewBooking(t *testing.T) {
	t.Parallel()

	adapter := NewODataAccountingAdapter(config.Config{OneCBaseURL: "http://example.local"})
	result, err := adapter.ExportBooking(t.Context(), mustBooking(t))
	if err != nil {
		t.Fatalf("export booking: %v", err)
	}
	if result.Status != ports.IntegrationStatusSynced {
		t.Fatalf("expected synced skip, got %s", result.Status)
	}
}

func mustConfirmedBooking(t *testing.T) domain.Booking {
	t.Helper()
	booking := mustBooking(t)
	booking.Status = domain.BookingStatusConfirmed
	return booking
}
