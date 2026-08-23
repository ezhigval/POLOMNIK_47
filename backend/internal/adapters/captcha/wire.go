package captcha

import (
	"strings"

	"polomnik/internal/adapters/captcha/noop"
	"polomnik/internal/adapters/captcha/smartcaptcha"
	"polomnik/internal/config"
	"polomnik/internal/ports"
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
