package application

import (
	"context"
	"testing"
)

func TestSecretEqual(t *testing.T) {
	if !SecretEqual(" abc ", "abc") {
		t.Fatal("expected equal secrets")
	}
	if SecretEqual("abc", "abd") {
		t.Fatal("expected mismatch")
	}
	if SecretEqual("ab", "abc") {
		t.Fatal("expected length mismatch")
	}
}

func TestWebhookGuardIdempotency(t *testing.T) {
	ctx := context.Background()
	guard := NewWebhookGuard(newMemoryCache())
	if guard.AlreadyProcessed(ctx, "telegram", "42") {
		t.Fatal("fresh key should not be processed")
	}
	guard.Remember(ctx, "telegram", "42")
	if !guard.AlreadyProcessed(ctx, "telegram", "42") {
		t.Fatal("expected remembered update")
	}
	if guard.AlreadyProcessed(ctx, "telegram", "43") {
		t.Fatal("other id must stay fresh")
	}
}
