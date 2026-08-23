// Package operator — единый источник реквизитов оператора персональных данных.
// После получения настоящих реквизитов заменить placeholders только здесь
// и в frontend/src/lib/operator-config.ts (зеркало для UI).
package operator

import (
	"os"
	"strings"
)

// Config holds operator details for legal documents.
type Config struct {
	Name            string
	INN             string
	OGRN            string
	LegalAddress    string
	ActualAddress   string
	Email           string
	Phone           string
	Website         string
	Regions         []string
	PublicSiteName  string
	PublicSiteFull  string
}

// Default returns operator config with placeholders from environment or defaults.
func Default() Config {
	siteURL := strings.TrimRight(env("PUBLIC_SITE_URL", "https://tikhvin-palomnik.ru"), "/")
	return Config{
		Name:           env("OPERATOR_NAME", "название"),
		INN:            env("OPERATOR_INN", "—"),
		OGRN:           env("OPERATOR_OGRN", "—"),
		LegalAddress:   env("OPERATOR_LEGAL_ADDRESS", "—"),
		ActualAddress:  env("OPERATOR_ACTUAL_ADDRESS", "—"),
		Email:          env("OPERATOR_EMAIL", env("NEXT_PUBLIC_CONTACT_EMAIL", "info@tikhvin-palomnik.ru")),
		Phone:          env("OPERATOR_PHONE", env("NEXT_PUBLIC_CONTACT_PHONE_DISPLAY", "+7 966 933-43-21")),
		Website:        siteURL,
		Regions:        []string{"Санкт-Петербург", "Ленинградская область", "иные регионы РФ по мере необходимости"},
		PublicSiteName: env("NEXT_PUBLIC_SITE_NAME", "Тихвинский путь"),
		PublicSiteFull: env("NEXT_PUBLIC_SITE_FULL_NAME", "РПЦ Тихвинская Епархия. Паломническая служба «Под покровом Божией Матери «Тихвинская»»"),
	}
}

func env(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}
