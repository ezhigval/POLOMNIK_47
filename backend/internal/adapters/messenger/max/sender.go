package maxmsg

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"palomnik/internal/ports"
)

const defaultAPIBase = "https://platform-api2.max.ru"

// Sender talks to the official Max Bot API (platform-api2.max.ru).
// Token goes in Authorization (not a query param). Address is user_id, or chat:<chat_id>.
type Sender struct {
	token      string
	apiBase    string
	httpClient *http.Client
}

func New(token, apiBase string, timeout time.Duration) Sender {
	if timeout <= 0 {
		timeout = 15 * time.Second
	}
	base := strings.TrimRight(strings.TrimSpace(apiBase), "/")
	if base == "" {
		base = defaultAPIBase
	}
	return Sender{
		token:   strings.TrimSpace(token),
		apiBase: base,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

var _ ports.MessengerPort = Sender{}

func (s Sender) Configured() bool {
	return s.token != ""
}

func (s Sender) Send(ctx context.Context, channel, address, text string) error {
	if !s.Configured() {
		return ports.ErrMessengerNotConfigured
	}
	if ch := strings.ToLower(strings.TrimSpace(channel)); ch != "" && ch != "max" {
		return ports.ErrMessengerNotConfigured
	}
	address = strings.TrimSpace(address)
	if address == "" {
		return ports.ErrMessengerNotConfigured
	}

	endpoint, err := s.messagesURL(address)
	if err != nil {
		return err
	}

	payload, err := json.Marshal(map[string]string{"text": text})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", s.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return fmt.Errorf("max send message: status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	return nil
}

func (s Sender) messagesURL(address string) (string, error) {
	q := url.Values{}
	if len(address) >= 5 && strings.EqualFold(address[:5], "chat:") {
		id := strings.TrimSpace(address[5:])
		if id == "" {
			return "", ports.ErrMessengerNotConfigured
		}
		q.Set("chat_id", id)
	} else {
		q.Set("user_id", address)
	}
	return s.apiBase + "/messages?" + q.Encode(), nil
}
