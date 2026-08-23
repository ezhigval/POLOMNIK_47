package config

import (
	"testing"
)

func TestValidateProductionRejectsDefaults(t *testing.T) {
	cfg := Load()
	cfg.AppEnv = "production"
	cfg.JWTSecret = DefaultJWTSecret
	cfg.AdminToken = DefaultAdminToken
	cfg.InternalAPISecret = DefaultInternalAPISecret
	cfg.DatabaseURL = "postgres://example"

	if err := Validate(cfg); err == nil {
		t.Fatal("expected production validation to fail with default secrets")
	}
}

func TestValidateBitrixRequiresInboundToken(t *testing.T) {
	cfg := Load()
	cfg.CRMAdapter = "bitrix"
	cfg.BitrixInboundToken = ""

	if err := Validate(cfg); err == nil {
		t.Fatal("expected bitrix adapter to require inbound token")
	}
}

func TestValidateProductionTelegramRequiresWorkerURL(t *testing.T) {
	cfg := Load()
	cfg.AppEnv = "production"
	cfg.JWTSecret = "production-jwt-secret-value-32chars"
	cfg.AdminToken = "production-admin-token"
	cfg.InternalAPISecret = "production-internal"
	cfg.DatabaseURL = "postgres://example"
	cfg.NotificationAdapter = "telegram"
	cfg.TelegramAPIBase = "https://api.telegram.org"

	if err := Validate(cfg); err == nil {
		t.Fatal("expected official Telegram API host to fail in production")
	}

	cfg.NotificationAdapter = "noop"
	cfg.MessengerAdapter = "telegram"
	if err := Validate(cfg); err == nil {
		t.Fatal("expected official Telegram API host to fail for messenger adapter")
	}

	cfg.TelegramAPIBase = "https://palomnik-telegram-api.example.workers.dev"
	if err := Validate(cfg); err != nil {
		t.Fatalf("worker url should pass: %v", err)
	}
}

func TestValidateProductionAcceptsAdminTokenLiteral(t *testing.T) {
	cfg := Load()
	cfg.AppEnv = "production"
	cfg.JWTSecret = "production-jwt-secret-value-32chars"
	cfg.AdminToken = "ADMIN_TOKEN"
	cfg.InternalAPISecret = "production-internal"
	cfg.DatabaseURL = "postgres://example"

	if err := Validate(cfg); err != nil {
		t.Fatalf("ADMIN_TOKEN literal should pass production validation: %v", err)
	}
}

func TestValidateLocalAllowsDefaults(t *testing.T) {
	cfg := Load()
	cfg.AppEnv = "local"
	cfg.CRMAdapter = "noop"

	if err := Validate(cfg); err != nil {
		t.Fatalf("expected local config to pass validation, got %v", err)
	}
}
