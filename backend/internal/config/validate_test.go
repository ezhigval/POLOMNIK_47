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

func TestValidateLocalAllowsDefaults(t *testing.T) {
	cfg := Load()
	cfg.AppEnv = "local"
	cfg.CRMAdapter = "noop"

	if err := Validate(cfg); err != nil {
		t.Fatalf("expected local config to pass validation, got %v", err)
	}
}
