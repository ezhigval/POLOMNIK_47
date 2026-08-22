package application

import (
	"testing"
)

func TestValidateInboundTokenRequiresTokenWhenRequired(t *testing.T) {
	service := &WebhookService{requireToken: true}

	if err := service.validateInboundToken("anything"); err != ErrUnauthorized {
		t.Fatalf("expected unauthorized without token, got %v", err)
	}
}

func TestValidateInboundTokenAllowsLocalWithoutToken(t *testing.T) {
	service := &WebhookService{requireToken: false}

	if err := service.validateInboundToken(""); err != nil {
		t.Fatalf("expected local webhook to pass, got %v", err)
	}
}

func TestValidateInboundTokenRejectsMismatch(t *testing.T) {
	service := &WebhookService{inboundToken: "secret", requireToken: true}

	if err := service.validateInboundToken("other"); err != ErrUnauthorized {
		t.Fatalf("expected unauthorized for mismatch, got %v", err)
	}
}
