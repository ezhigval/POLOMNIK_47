package payment

import (
	"testing"

	"palomnik/internal/adapters/payment/noop"
	"palomnik/internal/config"
)

func TestNewDefaultsToNoop(t *testing.T) {
	port := New(config.Config{})
	if _, ok := port.(noop.Gateway); !ok {
		t.Fatalf("expected noop, got %T", port)
	}
	if port.Configured() {
		t.Fatal("noop must not be configured")
	}
}

func TestNewSberWithoutKeysIsNoop(t *testing.T) {
	port := New(config.Config{PaymentAdapter: "sber", PublicSiteURL: "https://tikhvin-palomnik.ru"})
	if port.Configured() {
		t.Fatal("sber without keys must fall back to noop")
	}
}

func TestNewYooKassaWithoutKeysIsNoop(t *testing.T) {
	port := New(config.Config{PaymentAdapter: "yookassa", PublicSiteURL: "https://tikhvin-palomnik.ru"})
	if port.Configured() {
		t.Fatal("yookassa without keys must fall back to noop")
	}
}
