package publisher

import (
	"strings"

	"palomnik/internal/adapters/publisher/maxfeed"
	"palomnik/internal/adapters/publisher/noop"
	"palomnik/internal/adapters/publisher/sitenews"
	"palomnik/internal/adapters/publisher/telegramchannel"
	"palomnik/internal/adapters/publisher/vkwall"
	"palomnik/internal/application"
	"palomnik/internal/config"
	"palomnik/internal/ports"
)

func New(cfg config.Config, news *application.NewsService) ports.PublisherPort {
	switch strings.ToLower(strings.TrimSpace(cfg.PublisherAdapter)) {
	case "live":
		mux := liveMux(cfg, news)
		if !mux.Configured() {
			return noop.New()
		}
		return mux
	case "site_news":
		port := sitenews.New(news)
		if !port.Configured() {
			return noop.New()
		}
		return port
	case "telegram_channel":
		port := telegramchannel.New(cfg)
		if !port.Configured() {
			return noop.New()
		}
		return port
	case "vk_wall":
		port := vkwall.New(cfg.VKWallToken, cfg.VKWallOwnerID, cfg.VKAPIBase, cfg.VKAPIVersion, cfg.IntegrationHTTPTimeout)
		if !port.Configured() {
			return noop.New()
		}
		return port
	case "max_feed":
		port := maxfeed.New(cfg.MaxBotToken, cfg.MaxAPIBase, cfg.MaxFeedChatID, cfg.IntegrationHTTPTimeout)
		if !port.Configured() {
			return noop.New()
		}
		return port
	default:
		return noop.New()
	}
}

func liveMux(cfg config.Config, news *application.NewsService) Mux {
	return newMux(map[string]ports.PublisherPort{
		"site_news":        sitenews.New(news),
		"telegram_channel": telegramchannel.New(cfg),
		"vk_wall":          vkwall.New(cfg.VKWallToken, cfg.VKWallOwnerID, cfg.VKAPIBase, cfg.VKAPIVersion, cfg.IntegrationHTTPTimeout),
		"max_feed":         maxfeed.New(cfg.MaxBotToken, cfg.MaxAPIBase, cfg.MaxFeedChatID, cfg.IntegrationHTTPTimeout),
	})
}

func PublisherUsesTelegram(cfg config.Config) bool {
	adapter := strings.ToLower(strings.TrimSpace(cfg.PublisherAdapter))
	if adapter != "live" && adapter != "telegram_channel" {
		return false
	}
	return strings.TrimSpace(cfg.TelegramBotToken) != "" && strings.TrimSpace(cfg.TelegramChannelID) != ""
}
