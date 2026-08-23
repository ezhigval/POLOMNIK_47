package phone_test

import (
	"testing"
	"time"

	"palomnik/internal/adapters/phone"
	"palomnik/internal/config"
)

func TestNewNoopWithoutCredentials(t *testing.T) {
	cfg := config.Config{
		PhoneAdapter:           "smsru",
		SMSRUAPIID:             "",
		IntegrationHTTPTimeout: time.Second,
	}
	v := phone.New(cfg)
	if v.Available() {
		t.Fatal("expected unavailable without SMSRU_API_ID")
	}
}

func TestNewSMSRUWithAPIID(t *testing.T) {
	cfg := config.Config{
		PhoneAdapter:           "smsru",
		SMSRUAPIID:             "test-id",
		IntegrationHTTPTimeout: time.Second,
	}
	v := phone.New(cfg)
	if !v.Available() {
		t.Fatal("expected available with api_id")
	}
}

func TestNewCallAlias(t *testing.T) {
	cfg := config.Config{
		PhoneAdapter: "call",
		SMSRUAPIID:   "test-id",
	}
	if !phone.New(cfg).Available() {
		t.Fatal("call alias should enable smsru adapter")
	}
}
