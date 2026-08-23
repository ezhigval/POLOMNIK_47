package whatsapp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"palomnik/internal/domain"
	"palomnik/internal/ports"
)

const defaultGraphBase = "https://graph.facebook.com/v21.0"

// Sender uses official WhatsApp Cloud API (Graph). Unofficial WhatsApp is out of scope.
type Sender struct {
	token         string
	phoneNumberID string
	graphBase     string
	httpClient    *http.Client
}

func New(token, phoneNumberID, graphBase string, timeout time.Duration) Sender {
	if timeout <= 0 {
		timeout = 15 * time.Second
	}
	base := strings.TrimRight(strings.TrimSpace(graphBase), "/")
	if base == "" {
		base = defaultGraphBase
	}
	return Sender{
		token:         strings.TrimSpace(token),
		phoneNumberID: strings.TrimSpace(phoneNumberID),
		graphBase:     base,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

var _ ports.MessengerPort = Sender{}

func (s Sender) Configured() bool {
	return s.token != "" && s.phoneNumberID != ""
}

func (s Sender) Send(ctx context.Context, channel, address, text string) error {
	if !s.Configured() {
		return ports.ErrMessengerNotConfigured
	}
	if ch := strings.ToLower(strings.TrimSpace(channel)); ch != "" && ch != "whatsapp" {
		return ports.ErrMessengerNotConfigured
	}

	to := strings.TrimPrefix(domain.NormalizePhone(address), "+")
	if to == "" {
		to = digitsOnly(address)
	}
	if to == "" || strings.TrimSpace(text) == "" {
		return ports.ErrMessengerNotConfigured
	}

	payload, err := json.Marshal(map[string]any{
		"messaging_product": "whatsapp",
		"to":                to,
		"type":              "text",
		"text":              map[string]string{"body": text},
	})
	if err != nil {
		return err
	}

	endpoint := s.graphBase + "/" + s.phoneNumberID + "/messages"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+s.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return fmt.Errorf("whatsapp send message: status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	return nil
}

func digitsOnly(value string) string {
	var b strings.Builder
	for _, r := range value {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}
