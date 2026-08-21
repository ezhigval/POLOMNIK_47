package bitrix

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
	webhookURL string
	httpClient *http.Client
}

func NewClient(webhookURL string, timeout time.Duration) Client {
	if timeout <= 0 {
		timeout = 15 * time.Second
	}
	return Client{
		webhookURL: strings.TrimRight(strings.TrimSpace(webhookURL), "/"),
		httpClient: &http.Client{Timeout: timeout},
	}
}

type apiResponse struct {
	Result          json.RawMessage `json:"result"`
	Error           string          `json:"error"`
	ErrorDesc       string          `json:"error_description"`
	ErrorDescLegacy string          `json:"errorDescription"`
}

func (c Client) Call(ctx context.Context, method string, params map[string]any) (json.RawMessage, error) {
	if c.webhookURL == "" {
		return nil, fmt.Errorf("bitrix webhook url is empty")
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("marshal bitrix request: %w", err)
	}

	url := c.webhookURL + "/" + strings.TrimPrefix(method, "/")
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create bitrix request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("bitrix request failed: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read bitrix response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("bitrix http %d: %s", resp.StatusCode, truncate(string(raw), 512))
	}

	var parsed apiResponse
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return nil, fmt.Errorf("decode bitrix response: %w", err)
	}

	errMsg := parsed.ErrorDesc
	if errMsg == "" {
		errMsg = parsed.ErrorDescLegacy
	}
	if parsed.Error != "" {
		if errMsg == "" {
			errMsg = parsed.Error
		}
		return nil, fmt.Errorf("bitrix api %s: %s", parsed.Error, errMsg)
	}

	return parsed.Result, nil
}

func decodeID(raw json.RawMessage) (string, error) {
	if len(raw) == 0 || string(raw) == "null" {
		return "", fmt.Errorf("empty bitrix id result")
	}

	var asString string
	if err := json.Unmarshal(raw, &asString); err == nil && asString != "" {
		return asString, nil
	}

	var asNumber float64
	if err := json.Unmarshal(raw, &asNumber); err == nil && asNumber > 0 {
		return fmt.Sprintf("%.0f", asNumber), nil
	}

	return "", fmt.Errorf("unexpected bitrix id result: %s", truncate(string(raw), 128))
}

func decodeFirstListID(raw json.RawMessage) (string, bool, error) {
	if len(raw) == 0 || string(raw) == "null" {
		return "", false, nil
	}

	var items []map[string]any
	if err := json.Unmarshal(raw, &items); err != nil {
		return "", false, fmt.Errorf("decode bitrix list: %w", err)
	}
	if len(items) == 0 {
		return "", false, nil
	}

	id, ok := items[0]["ID"]
	if !ok {
		return "", false, fmt.Errorf("bitrix list item has no ID")
	}
	switch v := id.(type) {
	case string:
		return v, true, nil
	case float64:
		return fmt.Sprintf("%.0f", v), true, nil
	default:
		return fmt.Sprintf("%v", v), true, nil
	}
}

func truncate(value string, max int) string {
	if len(value) <= max {
		return value
	}
	return value[:max] + "..."
}
