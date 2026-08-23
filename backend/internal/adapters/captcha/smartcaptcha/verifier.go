package smartcaptcha

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"polomnik/internal/ports"
)

const defaultValidateURL = "https://smartcaptcha.yandexcloud.net/validate"

type Verifier struct {
	serverKey   string
	clientKey   string
	validateURL string
	httpClient  *http.Client
}

func New(serverKey, clientKey string, timeout time.Duration) *Verifier {
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	return &Verifier{
		serverKey:   strings.TrimSpace(serverKey),
		clientKey:   strings.TrimSpace(clientKey),
		validateURL: defaultValidateURL,
		httpClient:  &http.Client{Timeout: timeout},
	}
}

var _ ports.CaptchaPort = (*Verifier)(nil)

func (v *Verifier) Configured() bool {
	return v != nil && v.serverKey != "" && v.clientKey != ""
}

func (v *Verifier) Provider() string { return "smartcaptcha" }

func (v *Verifier) ClientKey() string {
	if v == nil {
		return ""
	}
	return v.clientKey
}

func (v *Verifier) Verify(ctx context.Context, token string, ip string) error {
	if !v.Configured() {
		return ports.ErrCaptchaNotConfigured
	}
	if strings.TrimSpace(token) == "" {
		return ports.ErrCaptchaFailed
	}

	form := url.Values{}
	form.Set("secret", v.serverKey)
	form.Set("token", strings.TrimSpace(token))
	if ip != "" {
		form.Set("ip", ip)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, v.validateURL, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := v.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var payload struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return fmt.Errorf("smartcaptcha decode: %w", err)
	}
	if !strings.EqualFold(payload.Status, "ok") {
		return ports.ErrCaptchaFailed
	}
	return nil
}
