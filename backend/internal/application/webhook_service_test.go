package application

import (
	"testing"

	"polomnik/internal/config"
)

func TestValidateInboundTokenRequiresTokenInProduction(t *testing.T) {
	service := &WebhookService{
		cfg: config.Config{AppEnv: "production", BitrixInboundToken: ""},
	}

	if err := service.validateInboundToken("anything"); err != ErrUnauthorized {
		t.Fatalf("expected unauthorized in production without token, got %v", err)
	}
}

func TestValidateInboundTokenRequiresTokenForBitrixAdapter(t *testing.T) {
	service := &WebhookService{
		cfg: config.Config{AppEnv: "local", CRMAdapter: "bitrix", BitrixInboundToken: ""},
	}

	if err := service.validateInboundToken("anything"); err != ErrUnauthorized {
		t.Fatalf("expected unauthorized for bitrix without inbound token, got %v", err)
	}
}

func TestValidateInboundTokenAllowsLocalWithoutToken(t *testing.T) {
	service := &WebhookService{
		cfg: config.Config{AppEnv: "local", CRMAdapter: "noop", BitrixInboundToken: ""},
	}

	if err := service.validateInboundToken(""); err != nil {
		t.Fatalf("expected local noop webhook to pass, got %v", err)
	}
}
