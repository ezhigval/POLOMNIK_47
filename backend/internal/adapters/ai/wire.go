package ai

import (
	"strings"

	"palomnik/internal/adapters/ai/noop"
	"palomnik/internal/adapters/ai/yandexgpt"
	"palomnik/internal/config"
	"palomnik/internal/ports"
)

func New(cfg config.Config) ports.AIPort {
	switch strings.ToLower(strings.TrimSpace(cfg.AIAdapter)) {
	case "yandexgpt":
		client := yandexgpt.New(cfg.YandexGPTAPIKey, cfg.YandexGPTFolderID, cfg.YandexGPTModel, cfg.YandexGPTAPIBase, cfg.IntegrationHTTPTimeout)
		if !client.Configured() {
			return noop.New()
		}
		return client
	default:
		return noop.New()
	}
}
