package onec

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
)

type ODataClient struct {
	baseURL    string
	odataPath  string
	httpClient *http.Client
	username   string
	password   string
}

func NewODataClient(baseURL, odataPath, username, password string, timeout time.Duration) ODataClient {
	if timeout <= 0 {
		timeout = 15 * time.Second
	}
	if odataPath == "" {
		odataPath = "/odata/standard.odata"
	}
	return ODataClient{
		baseURL:    strings.TrimRight(strings.TrimSpace(baseURL), "/"),
		odataPath:  strings.TrimRight(strings.TrimSpace(odataPath), "/"),
		httpClient: &http.Client{Timeout: timeout},
		username:   username,
		password:   password,
	}
}

func (c ODataClient) entityURL(entity string) string {
	return c.baseURL + c.odataPath + "/" + strings.TrimPrefix(entity, "/")
}

func (c ODataClient) Get(ctx context.Context, entity, filter string) (json.RawMessage, error) {
	values := url.Values{}
	if filter != "" {
		values.Set("$filter", filter)
	}
	values.Set("$format", "json")

	requestURL := c.entityURL(entity)
	if encoded := values.Encode(); encoded != "" {
		requestURL += "?" + encoded
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}
	c.applyAuth(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("1c odata get failed: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("1c odata http %d: %s", resp.StatusCode, truncate(string(raw), 512))
	}
	return raw, nil
}

func (c ODataClient) Post(ctx context.Context, entity string, payload any) (string, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.entityURL(entity), bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	c.applyAuth(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("1c odata post failed: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("1c odata http %d: %s", resp.StatusCode, truncate(string(raw), 512))
	}

	var parsed struct {
		RefKey string `json:"Ref_Key"`
	}
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return "", fmt.Errorf("decode 1c odata response: %w", err)
	}
	if parsed.RefKey == "" {
		return "", fmt.Errorf("1c odata response missing Ref_Key")
	}
	return parsed.RefKey, nil
}

func (c ODataClient) applyAuth(req *http.Request) {
	if c.username != "" {
		req.SetBasicAuth(c.username, c.password)
	}
}

func firstODataRefKey(raw json.RawMessage) (string, bool, error) {
	var envelope struct {
		Value []struct {
			RefKey string `json:"Ref_Key"`
		} `json:"value"`
	}
	if err := json.Unmarshal(raw, &envelope); err != nil {
		return "", false, err
	}
	if len(envelope.Value) == 0 {
		return "", false, nil
	}
	return envelope.Value[0].RefKey, true, nil
}

func bookingMarker(bookingID string) string {
	return "palomnik:booking:" + bookingID
}
