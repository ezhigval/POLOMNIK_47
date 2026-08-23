package yookassa

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"palomnik/internal/domain"
	"palomnik/internal/ports"
)

const defaultAPIBase = "https://api.yookassa.ru"

// Gateway calls official YooKassa payments API. Amount is booking.TotalPrice.
// No refunds. Does not change booking status.
type Gateway struct {
	shopID     string
	secretKey  string
	apiBase    string
	returnURL  string
	httpClient *http.Client
}

func New(shopID, secretKey, apiBase, returnURL string, timeout time.Duration) Gateway {
	if timeout <= 0 {
		timeout = 15 * time.Second
	}
	base := strings.TrimRight(strings.TrimSpace(apiBase), "/")
	if base == "" {
		base = defaultAPIBase
	}
	return Gateway{
		shopID:    strings.TrimSpace(shopID),
		secretKey: strings.TrimSpace(secretKey),
		apiBase:   base,
		returnURL: strings.TrimSpace(returnURL),
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

var _ ports.PaymentPort = Gateway{}

func (g Gateway) Configured() bool {
	return g.shopID != "" && g.secretKey != "" && g.returnURL != ""
}

func (g Gateway) CreatePayment(ctx context.Context, booking domain.Booking) (ports.PaymentSession, error) {
	if !g.Configured() {
		return ports.PaymentSession{}, ports.ErrPaymentNotConfigured
	}
	if booking.ID == uuid.Nil {
		return ports.PaymentSession{}, domain.ErrInvalidID
	}
	if booking.TotalPrice <= 0 {
		return ports.PaymentSession{}, domain.ErrInvalidTotalPrice
	}

	payload, err := json.Marshal(map[string]any{
		"amount": map[string]string{
			"value":    fmt.Sprintf("%d.00", booking.TotalPrice),
			"currency": "RUB",
		},
		"capture": true,
		"confirmation": map[string]string{
			"type":       "redirect",
			"return_url": g.returnURL,
		},
		"description": "booking " + booking.ID.String(),
		"metadata": map[string]string{
			"booking_id": booking.ID.String(),
		},
	})
	if err != nil {
		return ports.PaymentSession{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, g.apiBase+"/v3/payments", bytes.NewReader(payload))
	if err != nil {
		return ports.PaymentSession{}, err
	}
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(g.shopID+":"+g.secretKey)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Idempotence-Key", booking.ID.String())

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return ports.PaymentSession{}, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if resp.StatusCode >= 300 {
		return ports.PaymentSession{}, fmt.Errorf("yookassa payments: status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var parsed struct {
		ID           string `json:"id"`
		Confirmation struct {
			ConfirmationURL string `json:"confirmation_url"`
		} `json:"confirmation"`
	}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return ports.PaymentSession{}, fmt.Errorf("yookassa payments: decode: %w", err)
	}
	if strings.TrimSpace(parsed.ID) == "" || strings.TrimSpace(parsed.Confirmation.ConfirmationURL) == "" {
		return ports.PaymentSession{}, fmt.Errorf("yookassa payments: empty confirmation")
	}
	return ports.PaymentSession{
		ID:          strings.TrimSpace(parsed.ID),
		RedirectURL: strings.TrimSpace(parsed.Confirmation.ConfirmationURL),
	}, nil
}
