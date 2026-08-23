package captcha

import (
	"strings"

	"palomnik/internal/adapters/captcha/noop"
	"palomnik/internal/adapters/captcha/smartcaptcha"
	"palomnik/internal/config"
	"palomnik/internal/ports"
)

func New(cfg config.Config) ports.CaptchaPort {
	adapter := strings.ToLower(strings.TrimSpace(cfg.CaptchaAdapter))
	switch adapter {
	case "smartcaptcha":
		verifier := smartcaptcha.New(cfg.SmartCaptchaServerKey, cfg.SmartCaptchaClientKey, cfg.IntegrationHTTPTimeout)
		if !verifier.Configured() {
			return noop.New()
		}
		return verifier
	default:
		return noop.New()
	}
}
