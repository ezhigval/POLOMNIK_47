package config

import (
	"fmt"
	"strings"
)

const (
	DefaultJWTSecret         = "dev-jwt-secret-change-me"
	DefaultAdminToken        = "dev-admin-token"
	DefaultInternalAPISecret = "dev-internal-api-secret"
)

func (c Config) IsProduction() bool {
	return strings.EqualFold(c.AppEnv, "production")
}

func Validate(cfg Config) error {
	if cfg.IsProduction() {
		if err := requireStrongSecret("JWT_SECRET", cfg.JWTSecret, DefaultJWTSecret, 32); err != nil {
			return err
		}
		if err := requireStrongSecret("ADMIN_TOKEN", cfg.AdminToken, DefaultAdminToken, 11); err != nil {
			return err
		}
		if err := requireStrongSecret("INTERNAL_API_SECRET", cfg.InternalAPISecret, DefaultInternalAPISecret, 16); err != nil {
			return err
		}
		if cfg.DatabaseURL == "" {
			return fmt.Errorf("DATABASE_URL is required when APP_ENV=production")
		}
	}

	if strings.EqualFold(cfg.CRMAdapter, "bitrix") && cfg.BitrixInboundToken == "" {
		return fmt.Errorf("BITRIX_INBOUND_TOKEN is required when CRM_ADAPTER=bitrix")
	}

	telegramLive := strings.EqualFold(cfg.NotificationAdapter, "telegram") || strings.EqualFold(cfg.MessengerAdapter, "telegram") || publisherUsesTelegram(cfg)
	if cfg.IsProduction() && telegramLive && cfg.TelegramAPIHostIsOfficial() {
		return fmt.Errorf("TELEGRAM_API_BASE must be the outbound Worker URL when Telegram notifications, messenger, or telegram_channel publisher are live in production")
	}

	return nil
}

func publisherUsesTelegram(cfg Config) bool {
	adapter := strings.ToLower(strings.TrimSpace(cfg.PublisherAdapter))
	if adapter != "live" && adapter != "telegram_channel" {
		return false
	}
	return strings.TrimSpace(cfg.TelegramBotToken) != "" && strings.TrimSpace(cfg.TelegramChannelID) != ""
}

func requireStrongSecret(name, value, forbiddenDefault string, minLen int) error {
	if value == "" {
		return fmt.Errorf("%s must be set when APP_ENV=production", name)
	}
	if value == forbiddenDefault {
		return fmt.Errorf("%s must not use the default development value in production", name)
	}
	if len(value) < minLen {
		return fmt.Errorf("%s must be at least %d characters in production", name, minLen)
	}
	return nil
}
