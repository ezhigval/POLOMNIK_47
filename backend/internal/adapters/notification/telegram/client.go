package telegram

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"polomnik/internal/config"
	"polomnik/internal/ports"
)

type Client struct {
	token      string
	httpClient *http.Client
	apiBase    string
}

func NewClient(cfg config.Config) Client {
	apiBase := strings.TrimRight(cfg.TelegramAPIBase, "/")
	if apiBase == "" {
		apiBase = "https://api.telegram.org"
	}

	return Client{
		token:   cfg.TelegramBotToken,
		apiBase: apiBase,
		httpClient: &http.Client{
			Timeout: cfg.IntegrationHTTPTimeout,
		},
	}
}

func (c Client) SendMessage(ctx context.Context, chatID string, text string) error {
	if c.token == "" || chatID == "" {
		return ports.ErrNotificationNotConfigured
	}

	endpoint := fmt.Sprintf("%s/bot%s/sendMessage", c.apiBase, c.token)
	form := url.Values{}
	form.Set("chat_id", chatID)
	form.Set("text", text)
	form.Set("parse_mode", "HTML")
	form.Set("disable_web_page_preview", "true")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return fmt.Errorf("telegram sendMessage: status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	return nil
}

func (c Client) Configured() bool {
	return c.token != ""
}
