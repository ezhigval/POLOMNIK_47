package config

import (
	"testing"
	"time"
)

func TestLoadUsesDefaults(t *testing.T) {
	t.Setenv("APP_ENV", "")
	t.Setenv("HTTP_ADDR", "")
	t.Setenv("LOG_LEVEL", "")
	t.Setenv("CRM_ADAPTER", "")
	t.Setenv("ACCOUNTING_ADAPTER", "")
	t.Setenv("SHUTDOWN_TIMEOUT", "")

	cfg := Load()

	if cfg.AppEnv != "local" {
		t.Fatalf("expected default app env local, got %q", cfg.AppEnv)
	}
	if cfg.HTTPAddr != ":8080" {
		t.Fatalf("expected default http addr :8080, got %q", cfg.HTTPAddr)
	}
	if cfg.LogLevel != "info" {
		t.Fatalf("expected default log level info, got %q", cfg.LogLevel)
	}
	if cfg.CRMAdapter != "noop" {
		t.Fatalf("expected default crm adapter noop, got %q", cfg.CRMAdapter)
	}
	if cfg.AccountingAdapter != "noop" {
		t.Fatalf("expected default accounting adapter noop, got %q", cfg.AccountingAdapter)
	}
	if cfg.UploadDir != "./data/uploads" {
		t.Fatalf("expected default upload dir ./data/uploads, got %q", cfg.UploadDir)
	}
	if cfg.UploadPublicBaseURL != "http://localhost:8080" {
		t.Fatalf("expected default upload public base url, got %q", cfg.UploadPublicBaseURL)
	}
	if cfg.UploadMaxBytes != 5*1024*1024 {
		t.Fatalf("expected default upload max bytes 5MiB, got %d", cfg.UploadMaxBytes)
	}
	if cfg.ShutdownTimeout != 10*time.Second {
		t.Fatalf("expected default shutdown timeout 10s, got %s", cfg.ShutdownTimeout)
	}
}

func TestLoadUsesEnvironment(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("HTTP_ADDR", ":9090")
	t.Setenv("LOG_LEVEL", "debug")
	t.Setenv("CRM_ADAPTER", "fake")
	t.Setenv("ACCOUNTING_ADAPTER", "fake")
	t.Setenv("SHUTDOWN_TIMEOUT", "2s")

	cfg := Load()

	if cfg.AppEnv != "test" {
		t.Fatalf("expected app env from env, got %q", cfg.AppEnv)
	}
	if cfg.HTTPAddr != ":9090" {
		t.Fatalf("expected http addr from env, got %q", cfg.HTTPAddr)
	}
	if cfg.LogLevel != "debug" {
		t.Fatalf("expected log level from env, got %q", cfg.LogLevel)
	}
	if cfg.CRMAdapter != "fake" {
		t.Fatalf("expected crm adapter from env, got %q", cfg.CRMAdapter)
	}
	if cfg.AccountingAdapter != "fake" {
		t.Fatalf("expected accounting adapter from env, got %q", cfg.AccountingAdapter)
	}
	if cfg.ShutdownTimeout != 2*time.Second {
		t.Fatalf("expected shutdown timeout from env, got %s", cfg.ShutdownTimeout)
	}
}

func TestEffectiveTelegramLoginFallsBackToBotToken(t *testing.T) {
	t.Setenv("TELEGRAM_BOT_TOKEN", "main-token")
	t.Setenv("TELEGRAM_BOT_USERNAME", "PalomnikBot")
	t.Setenv("TELEGRAM_LOGIN_BOT_TOKEN", "")
	t.Setenv("TELEGRAM_LOGIN_BOT_USERNAME", "")

	cfg := Load()
	if cfg.EffectiveTelegramLoginBotToken() != "main-token" {
		t.Fatalf("token fallback: %q", cfg.EffectiveTelegramLoginBotToken())
	}
	if cfg.EffectiveTelegramLoginBotUsername() != "PalomnikBot" {
		t.Fatalf("username fallback: %q", cfg.EffectiveTelegramLoginBotUsername())
	}
}
