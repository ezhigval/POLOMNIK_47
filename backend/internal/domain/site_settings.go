package domain

import (
	"strings"
	"time"
	"unicode/utf8"
)

const (
	siteSettingsRowID        int16 = 1
	maxSiteFieldLen                = 500
	maxSiteDescriptionLen          = 2000
	maxMailForwardRecipients       = 20
)

type SiteSettings struct {
	SiteName            string
	FullName            string
	Tagline             string
	Description         string
	Region              string
	DepartureCity       string
	ParentOrgName       string
	ParentOrgURL        string
	ContactPhone        string
	ContactPhoneDisplay string
	ContactEmail        string
	MailForwardTo       []string
	UpdatedAt           time.Time
}

func SiteSettingsID() int16 {
	return siteSettingsRowID
}

func NewSiteSettings(input SiteSettings, now time.Time) (SiteSettings, error) {
	settings := SiteSettings{
		SiteName:            trimSiteField(input.SiteName, maxSiteFieldLen),
		FullName:            trimSiteField(input.FullName, maxSiteFieldLen),
		Tagline:             trimSiteField(input.Tagline, maxSiteFieldLen),
		Description:         trimSiteField(input.Description, maxSiteDescriptionLen),
		Region:              trimSiteField(input.Region, 32),
		DepartureCity:       trimSiteField(input.DepartureCity, maxSiteFieldLen),
		ParentOrgName:       trimSiteField(input.ParentOrgName, maxSiteFieldLen),
		ParentOrgURL:        trimSiteField(input.ParentOrgURL, maxSiteFieldLen),
		ContactPhoneDisplay: trimSiteField(input.ContactPhoneDisplay, 64),
	}

	phone := strings.TrimSpace(input.ContactPhone)
	if phone != "" {
		normalized := NormalizePhone(phone)
		if normalized == "" {
			return SiteSettings{}, ErrInvalidPhone
		}
		settings.ContactPhone = normalized
	}

	email := strings.TrimSpace(strings.ToLower(input.ContactEmail))
	if email != "" {
		if !strings.Contains(email, "@") || strings.Contains(email, " ") {
			return SiteSettings{}, ErrInvalidEmail
		}
		settings.ContactEmail = email
	}

	forwards, err := NormalizeMailForwardList(input.MailForwardTo)
	if err != nil {
		return SiteSettings{}, err
	}
	settings.MailForwardTo = forwards

	if now.IsZero() {
		now = time.Now().UTC()
	}
	settings.UpdatedAt = now
	return settings, nil
}

func NormalizeMailForwardList(raw []string) ([]string, error) {
	seen := make(map[string]struct{})
	out := make([]string, 0, len(raw))
	for _, item := range raw {
		email := strings.TrimSpace(strings.ToLower(item))
		if email == "" {
			continue
		}
		if !strings.Contains(email, "@") || strings.Contains(email, " ") {
			return nil, ErrInvalidEmail
		}
		if _, ok := seen[email]; ok {
			continue
		}
		if len(out) >= maxMailForwardRecipients {
			return nil, ErrInvalidEmail
		}
		seen[email] = struct{}{}
		out = append(out, email)
	}
	return out, nil
}

func FormatMailForwardList(emails []string) string {
	return strings.Join(emails, "\n")
}

func ParseMailForwardList(raw string) ([]string, error) {
	return NormalizeMailForwardList(ParseLooseAddressList(raw))
}

func MailForwardListOrFallback(stored, fallback []string) []string {
	if len(stored) > 0 {
		return append([]string(nil), stored...)
	}
	if len(fallback) == 0 {
		return nil
	}
	return append([]string(nil), fallback...)
}

func MergeSiteSettings(base, overlay SiteSettings) SiteSettings {
	out := base
	if overlay.SiteName != "" {
		out.SiteName = overlay.SiteName
	}
	if overlay.FullName != "" {
		out.FullName = overlay.FullName
	}
	if overlay.Tagline != "" {
		out.Tagline = overlay.Tagline
	}
	if overlay.Description != "" {
		out.Description = overlay.Description
	}
	if overlay.Region != "" {
		out.Region = overlay.Region
	}
	if overlay.DepartureCity != "" {
		out.DepartureCity = overlay.DepartureCity
	}
	if overlay.ParentOrgName != "" {
		out.ParentOrgName = overlay.ParentOrgName
	}
	if overlay.ParentOrgURL != "" {
		out.ParentOrgURL = overlay.ParentOrgURL
	}
	if overlay.ContactPhone != "" {
		out.ContactPhone = overlay.ContactPhone
	}
	if overlay.ContactPhoneDisplay != "" {
		out.ContactPhoneDisplay = overlay.ContactPhoneDisplay
	}
	if overlay.ContactEmail != "" {
		out.ContactEmail = overlay.ContactEmail
	}
	if overlay.MailForwardTo != nil {
		out.MailForwardTo = append([]string(nil), overlay.MailForwardTo...)
	}
	if !overlay.UpdatedAt.IsZero() {
		out.UpdatedAt = overlay.UpdatedAt
	}
	return out
}

func trimSiteField(value string, maxLen int) string {
	value = strings.TrimSpace(value)
	if maxLen > 0 && utf8.RuneCountInString(value) > maxLen {
		runes := []rune(value)
		value = string(runes[:maxLen])
	}
	return value
}
