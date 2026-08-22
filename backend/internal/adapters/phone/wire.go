package phone

import (
	"strings"

	"polomnik/internal/adapters/phone/noop"
	"polomnik/internal/adapters/phone/smsru"
	"polomnik/internal/config"
	"polomnik/internal/ports"
)

// New builds a PhoneVerifier. PHONE_ADAPTER=smsru uses sms.ru callcheck (not SMS).
// Empty adapter or missing SMSRU_API_ID → noop (Available=false).
func New(cfg config.Config) ports.PhoneVerifier {
	adapter := strings.ToLower(strings.TrimSpace(cfg.PhoneAdapter))
	apiID := strings.TrimSpace(cfg.SMSRUAPIID)

	switch adapter {
	case "smsru", "call", "callcheck":
		if apiID == "" {
			return noop.New()
		}
		return smsru.New(apiID, cfg.IntegrationHTTPTimeout)
	default:
		return noop.New()
	}
}
