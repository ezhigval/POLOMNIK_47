package vkwall

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"palomnik/internal/adapters/publisher/postformat"
	"palomnik/internal/ports"
)

const (
	defaultAPIBase    = "https://api.vk.com"
	defaultAPIVersion = "5.199"
)

// Publisher posts to a VK community wall via official wall.post.
// owner_id for a community must be negative (e.g. -123456).
type Publisher struct {
	token      string
	ownerID    string
	apiBase    string
	apiVersion string
	httpClient *http.Client
}

func New(token, ownerID, apiBase, apiVersion string, timeout time.Duration) Publisher {
	if timeout <= 0 {
		timeout = 15 * time.Second
	}
	base := strings.TrimRight(strings.TrimSpace(apiBase), "/")
	if base == "" {
		base = defaultAPIBase
	}
	version := strings.TrimSpace(apiVersion)
	if version == "" {
		version = defaultAPIVersion
	}
	return Publisher{
		token:      strings.TrimSpace(token),
		ownerID:    strings.TrimSpace(ownerID),
		apiBase:    base,
		apiVersion: version,
		httpClient: &http.Client{Timeout: timeout},
	}
}

var _ ports.PublisherPort = Publisher{}

func (p Publisher) Configured() bool {
	return p.token != "" && p.ownerID != ""
}

func (p Publisher) Publish(ctx context.Context, channel string, content ports.PublishContent) error {
	if !p.Configured() {
		return ports.ErrPublisherNotConfigured
	}
	if ch := strings.ToLower(strings.TrimSpace(channel)); ch != "" && ch != "vk_wall" {
		return ports.ErrPublisherNotConfigured
	}

	message := postformat.Post(content)
	if strings.TrimSpace(message) == "" {
		return ports.ErrPublisherNotConfigured
	}

	form := url.Values{}
	form.Set("access_token", p.token)
	form.Set("v", p.apiVersion)
	form.Set("owner_id", p.ownerID)
	form.Set("from_group", "1")
	form.Set("message", message)

	endpoint := p.apiBase + "/method/wall.post"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
	if resp.StatusCode >= 300 {
		return fmt.Errorf("vk wall.post: status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var parsed struct {
		Error *struct {
			Code    int    `json:"error_code"`
			Message string `json:"error_msg"`
		} `json:"error"`
	}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return fmt.Errorf("vk wall.post: decode: %w", err)
	}
	if parsed.Error != nil {
		return fmt.Errorf("vk wall.post: %d %s", parsed.Error.Code, parsed.Error.Message)
	}
	return nil
}
