package messenger

import (
	"strings"

	maxmsg "palomnik/internal/adapters/messenger/max"
	"palomnik/internal/adapters/messenger/noop"
	msgtelegram "palomnik/internal/adapters/messenger/telegram"
	"palomnik/internal/adapters/messenger/whatsapp"
	"palomnik/internal/config"
	"palomnik/internal/ports"
)

func New(cfg config.Config) ports.MessengerPort {
	switch strings.ToLower(strings.TrimSpace(cfg.MessengerAdapter)) {
	case "telegram":
		sender := msgtelegram.New(cfg)
		if !sender.Configured() {
			return noop.New()
		}
		return sender
	case "max":
		sender := maxmsg.New(cfg.MaxBotToken, cfg.MaxAPIBase, cfg.IntegrationHTTPTimeout)
		if !sender.Configured() {
			return noop.New()
		}
		return sender
	case "whatsapp":
		sender := whatsapp.New(cfg.WhatsAppToken, cfg.WhatsAppPhoneNumberID, cfg.WhatsAppGraphBase, cfg.IntegrationHTTPTimeout)
		if !sender.Configured() {
			return noop.New()
		}
		return sender
	default:
		return noop.New()
	}
}
