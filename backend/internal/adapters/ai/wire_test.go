package ai

import (
	"testing"

	"palomnik/internal/config"
)

func TestNewNoopByDefault(t *testing.T) {
	if port := New(config.Config{}); port.Configured() {
		t.Fatal("expected noop")
	}
}

func TestNewFallsBackWithoutKeys(t *testing.T) {
	if port := New(config.Config{AIAdapter: "yandexgpt"}); port.Configured() {
		t.Fatal("yandexgpt without keys should be noop")
	}
}

func TestNewYandexGPTConfiguredWithKeys(t *testing.T) {
	port := New(config.Config{
		AIAdapter:         "yandexgpt",
		YandexGPTAPIKey:   "key",
		YandexGPTFolderID: "folder",
	})
	if !port.Configured() {
		t.Fatal("expected configured")
	}
}
