package integration

import (
	"polomnik/internal/adapters/integration/bitrix"
	"polomnik/internal/adapters/integration/noop"
	"polomnik/internal/adapters/integration/onec"
	"polomnik/internal/adapters/integration/recording"
	"polomnik/internal/config"
	"polomnik/internal/domain"
	"polomnik/internal/ports"
)

func NewCRM(cfg config.Config, refs ports.IntegrationReferenceRepository, outbox ports.OutboxRepository) ports.CRMPort {
	switch cfg.CRMAdapter {
	case "bitrix":
		return recording.NewCRM(bitrix.NewCRMAdapter(cfg), refs, outbox, recording.CRMOptions{
			ExternalSystem: domain.ExternalSystemBitrix24,
		})
	case "noop", "":
		return noop.NewCRMAdapter()
	default:
		return noop.NewCRMAdapter()
	}
}

func NewAccounting(cfg config.Config, refs ports.IntegrationReferenceRepository, outbox ports.OutboxRepository) ports.AccountingPort {
	switch cfg.AccountingAdapter {
	case "onec":
		return recording.NewAccounting(onec.NewAccountingAdapter(cfg), refs, outbox, recording.AccountingOptions{
			ExternalSystem: domain.ExternalSystemOneC,
		})
	case "noop", "":
		return noop.NewAccountingAdapter()
	default:
		return noop.NewAccountingAdapter()
	}
}

func CRMEnabled(cfg config.Config) bool {
	return cfg.CRMAdapter == "bitrix"
}

func AccountingEnabled(cfg config.Config) bool {
	return cfg.AccountingAdapter == "onec"
}

func NewCRMInbound(cfg config.Config) ports.CRMInboundPort {
	switch cfg.CRMAdapter {
	case "bitrix":
		return bitrix.NewCRMAdapter(cfg)
	default:
		return noop.NewCRMInboundAdapter()
	}
}
