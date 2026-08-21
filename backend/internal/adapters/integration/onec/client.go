package onec

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
	username   string
	password   string
	apiKey     string
}

func NewClient(baseURL, username, password, apiKey string, timeout time.Duration) Client {
	if timeout <= 0 {
		timeout = 15 * time.Second
	}
	return Client{
		baseURL:    strings.TrimRight(strings.TrimSpace(baseURL), "/"),
		httpClient: &http.Client{Timeout: timeout},
		username:   username,
		password:   password,
		apiKey:     apiKey,
	}
}

type exchangeResponse struct {
	ExternalID string `json:"external_id"`
	Status     string `json:"status"`
	Message    string `json:"message"`
}

func (c Client) PostJSON(ctx context.Context, path string, payload any) (exchangeResponse, error) {
	if c.baseURL == "" {
		return exchangeResponse{}, fmt.Errorf("1c base url is empty")
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return exchangeResponse{}, fmt.Errorf("marshal 1c request: %w", err)
	}

	url := c.baseURL + "/" + strings.TrimPrefix(path, "/")
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return exchangeResponse{}, fmt.Errorf("create 1c request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	c.applyAuth(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return exchangeResponse{}, fmt.Errorf("1c request failed: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return exchangeResponse{}, fmt.Errorf("read 1c response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return exchangeResponse{}, fmt.Errorf("1c http %d: %s", resp.StatusCode, truncate(string(raw), 512))
	}

	var parsed exchangeResponse
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return exchangeResponse{}, fmt.Errorf("decode 1c response: %w", err)
	}
	if parsed.ExternalID == "" {
		return exchangeResponse{}, fmt.Errorf("1c response missing external_id")
	}
	return parsed, nil
}

func (c Client) applyAuth(req *http.Request) {
	if c.username != "" {
		req.SetBasicAuth(c.username, c.password)
	}
	if c.apiKey != "" {
		req.Header.Set("X-API-Key", c.apiKey)
	}
}

func truncate(value string, max int) string {
	if len(value) <= max {
		return value
	}
	return value[:max] + "..."
}
