package captcha

import (
	"testing"

	"polomnik/internal/config"
)

func TestNewFallsBackToNoopWithoutKeys(t *testing.T) {
	port := New(config.Config{CaptchaAdapter: "smartcaptcha"})
	if port.Configured() {
		t.Fatal("expected noop when smartcaptcha keys are empty")
	}
}

func TestNewNoopByDefault(t *testing.T) {
	port := New(config.Config{})
	if port.Configured() || port.Provider() != "noop" {
		t.Fatalf("expected noop, got configured=%v provider=%s", port.Configured(), port.Provider())
	}
}
