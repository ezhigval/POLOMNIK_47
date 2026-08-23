package yandexgpt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"palomnik/internal/ports"
)

const (
	defaultAPIBase     = "https://llm.api.cloud.yandex.net"
	completionPath     = "/foundationModels/v1/completion"
	defaultModel       = "yandexgpt-lite"
	defaultMaxTokens   = "2000"
	defaultTemperature = 0.3
)

// Client calls official Yandex Cloud Foundation Models completion.
// Auth is a service-account API key (Authorization: Api-Key), not a short-lived IAM token.
type Client struct {
	apiKey     string
	folderID   string
	model      string
	apiBase    string
	httpClient *http.Client
}

func New(apiKey, folderID, model, apiBase string, timeout time.Duration) Client {
	if timeout <= 0 {
		timeout = 15 * time.Second
	}
	base := strings.TrimRight(strings.TrimSpace(apiBase), "/")
	if base == "" {
		base = defaultAPIBase
	}
	name := strings.TrimSpace(model)
	if name == "" {
		name = defaultModel
	}
	return Client{
		apiKey:   strings.TrimSpace(apiKey),
		folderID: strings.TrimSpace(folderID),
		model:    name,
		apiBase:  base,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

var _ ports.AIPort = Client{}

func (c Client) Configured() bool {
	return c.apiKey != "" && c.folderID != ""
}

func (c Client) Complete(ctx context.Context, req ports.AICompletionRequest) (string, error) {
	if !c.Configured() {
		return "", ports.ErrAINotConfigured
	}
	user := strings.TrimSpace(req.User)
	if user == "" {
		return "", ports.ErrAINotConfigured
	}

	messages := make([]map[string]string, 0, 2)
	if system := strings.TrimSpace(req.System); system != "" {
		messages = append(messages, map[string]string{"role": "system", "text": system})
	}
	messages = append(messages, map[string]string{"role": "user", "text": user})

	payload, err := json.Marshal(map[string]any{
		"modelUri": c.modelURI(),
		"completionOptions": map[string]any{
			"stream":      false,
			"temperature": defaultTemperature,
			"maxTokens":   defaultMaxTokens,
		},
		"messages": messages,
	})
	if err != nil {
		return "", err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.apiBase+completionPath, bytes.NewReader(payload))
	if err != nil {
		return "", err
	}
	httpReq.Header.Set("Authorization", "Api-Key "+c.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-folder-id", c.folderID)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("yandexgpt completion: status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var parsed struct {
		Result struct {
			Alternatives []struct {
				Message struct {
					Text string `json:"text"`
				} `json:"message"`
			} `json:"alternatives"`
		} `json:"result"`
		Error *struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return "", fmt.Errorf("yandexgpt completion: decode: %w", err)
	}
	if parsed.Error != nil && strings.TrimSpace(parsed.Error.Message) != "" {
		return "", fmt.Errorf("yandexgpt completion: %s", parsed.Error.Message)
	}
	if len(parsed.Result.Alternatives) == 0 {
		return "", fmt.Errorf("yandexgpt completion: empty alternatives")
	}
	text := strings.TrimSpace(parsed.Result.Alternatives[0].Message.Text)
	if text == "" {
		return "", fmt.Errorf("yandexgpt completion: empty text")
	}
	return text, nil
}

func (c Client) modelURI() string {
	if strings.Contains(c.model, "://") {
		return c.model
	}
	return "gpt://" + c.folderID + "/" + strings.TrimPrefix(c.model, "/")
}
