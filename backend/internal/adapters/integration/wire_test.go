package integration

import (
	"testing"

	"palomnik/internal/adapters/integration/noop"
	"palomnik/internal/config"
	"palomnik/internal/ports"
)

func TestNewCRMReturnsNoopAdapterByDefault(t *testing.T) {
	cfg := config.Config{CRMAdapter: "noop"}
	adapter := NewCRM(cfg, nil, nil)
	if adapter == nil {
		t.Fatal("expected CRM port implementation, got nil")
	}
	if _, ok := adapter.(noop.CRMAdapter); !ok {
		t.Fatalf("expected noop CRM adapter, got %T", adapter)
	}
}

func TestNewCRMWiresBitrixStub(t *testing.T) {
	cfg := config.Config{CRMAdapter: "bitrix"}
	adapter := NewCRM(cfg, nil, nil)
	if adapter == nil {
		t.Fatal("expected bitrix CRM stub, got nil")
	}
}

func TestNewAccountingReturnsNoopAdapterByDefault(t *testing.T) {
	cfg := config.Config{AccountingAdapter: "noop"}
	adapter := NewAccounting(cfg, nil, nil)
	if adapter == nil {
		t.Fatal("expected accounting port implementation, got nil")
	}
	if _, ok := adapter.(noop.AccountingAdapter); !ok {
		t.Fatalf("expected noop accounting adapter, got %T", adapter)
	}
}

func TestNewAccountingWiresOneCStub(t *testing.T) {
	cfg := config.Config{AccountingAdapter: "onec"}
	adapter := NewAccounting(cfg, nil, nil)
	if adapter == nil {
		t.Fatal("expected 1c accounting stub, got nil")
	}
}

func TestCRMEnabledOnlyForBitrix(t *testing.T) {
	if CRMEnabled(config.Config{CRMAdapter: "noop"}) {
		t.Fatal("noop must not enable CRM integration")
	}
	if !CRMEnabled(config.Config{CRMAdapter: "bitrix"}) {
		t.Fatal("bitrix must enable CRM integration flag")
	}
}

func TestAccountingEnabledOnlyForOneC(t *testing.T) {
	if AccountingEnabled(config.Config{AccountingAdapter: "noop"}) {
		t.Fatal("noop must not enable accounting integration")
	}
	if !AccountingEnabled(config.Config{AccountingAdapter: "onec"}) {
		t.Fatal("onec must enable accounting integration flag")
	}
}

func TestNoopCRMPushBookingReturnsNotConfigured(t *testing.T) {
	t.Parallel()

	adapter := noop.NewCRMAdapter()
	result, err := adapter.PushBooking(t.Context(), mustBooking(t))
	if err != nil {
		t.Fatalf("push booking: %v", err)
	}
	if result.Status != ports.IntegrationStatusNotConfigured {
		t.Fatalf("expected not_configured, got %s", result.Status)
	}
}
