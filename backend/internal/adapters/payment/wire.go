package payment

import (
	"strings"

	"palomnik/internal/adapters/payment/noop"
	"palomnik/internal/adapters/payment/sber"
	"palomnik/internal/adapters/payment/yookassa"
	"palomnik/internal/config"
	"palomnik/internal/ports"
)

func New(cfg config.Config) ports.PaymentPort {
	returnURL := strings.TrimSpace(cfg.PaymentReturnURL)
	if returnURL == "" {
		base := strings.TrimRight(strings.TrimSpace(cfg.PublicSiteURL), "/")
		if base != "" {
			returnURL = base + "/account/trips"
		}
	}
	switch strings.ToLower(strings.TrimSpace(cfg.PaymentAdapter)) {
	case "sber":
		client := sber.New(cfg.SberUserName, cfg.SberPassword, cfg.SberAPIBase, returnURL, cfg.IntegrationHTTPTimeout)
		if !client.Configured() {
			return noop.New()
		}
		return client
	case "yookassa":
		client := yookassa.New(cfg.YooKassaShopID, cfg.YooKassaSecretKey, cfg.YooKassaAPIBase, returnURL, cfg.IntegrationHTTPTimeout)
		if !client.Configured() {
			return noop.New()
		}
		return client
	default:
		return noop.New()
	}
}
