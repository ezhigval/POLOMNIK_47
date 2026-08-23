package sber

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"

	"palomnik/internal/domain"
	"palomnik/internal/ports"
)

const defaultAPIBase = "https://securepayments.sberbank.ru"

// Gateway calls official Sber register.do. Amount is booking.TotalPrice in rubles → kopecks.
// No refunds. Does not change booking status.
type Gateway struct {
	userName   string
	password   string
	apiBase    string
	returnURL  string
	httpClient *http.Client
}

func New(userName, password, apiBase, returnURL string, timeout time.Duration) Gateway {
	if timeout <= 0 {
		timeout = 15 * time.Second
	}
	base := strings.TrimRight(strings.TrimSpace(apiBase), "/")
	if base == "" {
		base = defaultAPIBase
	}
	return Gateway{
		userName:  strings.TrimSpace(userName),
		password:  strings.TrimSpace(password),
		apiBase:   base,
		returnURL: strings.TrimSpace(returnURL),
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

var _ ports.PaymentPort = Gateway{}

func (g Gateway) Configured() bool {
	return g.userName != "" && g.password != "" && g.returnURL != ""
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

	form := url.Values{}
	form.Set("userName", g.userName)
	form.Set("password", g.password)
	form.Set("amount", fmt.Sprintf("%d", booking.TotalPrice*100))
	form.Set("orderNumber", booking.ID.String())
	form.Set("returnUrl", g.returnURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, g.apiBase+"/payment/rest/register.do", strings.NewReader(form.Encode()))
	if err != nil {
		return ports.PaymentSession{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return ports.PaymentSession{}, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if resp.StatusCode >= 300 {
		return ports.PaymentSession{}, fmt.Errorf("sber register: status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var parsed struct {
		OrderID   string `json:"orderId"`
		FormURL   string `json:"formUrl"`
		ErrorCode string `json:"errorCode"`
		ErrorMsg  string `json:"errorMessage"`
	}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return ports.PaymentSession{}, fmt.Errorf("sber register: decode: %w", err)
	}
	if parsed.ErrorCode != "" && parsed.ErrorCode != "0" {
		return ports.PaymentSession{}, fmt.Errorf("sber register: %s %s", parsed.ErrorCode, parsed.ErrorMsg)
	}
	if strings.TrimSpace(parsed.FormURL) == "" {
		return ports.PaymentSession{}, fmt.Errorf("sber register: empty formUrl")
	}
	return ports.PaymentSession{
		ID:          strings.TrimSpace(parsed.OrderID),
		RedirectURL: strings.TrimSpace(parsed.FormURL),
	}, nil
}
