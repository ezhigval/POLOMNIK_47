package domain

import (
	"testing"
	"time"
)

func TestNewSiteSettingsNormalizesContacts(t *testing.T) {
	got, err := NewSiteSettings(SiteSettings{
		SiteName:      "  Тихвинский путь ",
		ContactPhone:  "8 (966) 933-43-21",
		ContactEmail:  " Info@Example.ru ",
		MailForwardTo: []string{"a@b.ru", "a@b.ru", "c@d.ru"},
	}, time.Time{})
	if err != nil {
		t.Fatalf("new: %v", err)
	}
	if got.SiteName != "Тихвинский путь" {
		t.Fatalf("name: %q", got.SiteName)
	}
	if got.ContactPhone != "+79669334321" {
		t.Fatalf("phone: %q", got.ContactPhone)
	}
	if got.ContactEmail != "info@example.ru" {
		t.Fatalf("email: %q", got.ContactEmail)
	}
	if len(got.MailForwardTo) != 2 {
		t.Fatalf("mail forward: %v", got.MailForwardTo)
	}
}

func TestMergeSiteSettingsKeepsBaseWhenOverlayEmpty(t *testing.T) {
	base := SiteSettings{SiteName: "Base", ContactEmail: "a@b.ru"}
	overlay := SiteSettings{Tagline: "New tagline"}
	got := MergeSiteSettings(base, overlay)
	if got.SiteName != "Base" || got.ContactEmail != "a@b.ru" || got.Tagline != "New tagline" {
		t.Fatalf("merge: %+v", got)
	}
}
