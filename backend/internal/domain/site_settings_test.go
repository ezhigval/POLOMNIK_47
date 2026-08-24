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

func TestParseMailForwardListAcceptsNewlinesAndCommas(t *testing.T) {
	got, err := ParseMailForwardList("tikhvin-palomnik@yandex.ru\n smailikin70@yandex.ru, tikhvin-palomnik@yandex.ru")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(got) != 2 || got[0] != "tikhvin-palomnik@yandex.ru" || got[1] != "smailikin70@yandex.ru" {
		t.Fatalf("parsed: %v", got)
	}
	formatted := FormatMailForwardList(got)
	if formatted != "tikhvin-palomnik@yandex.ru\nsmailikin70@yandex.ru" {
		t.Fatalf("format: %q", formatted)
	}
}

func TestMailForwardListOrFallback(t *testing.T) {
	stored := []string{"tikhvin-palomnik@yandex.ru"}
	fallback := []string{"smailikin70@yandex.ru"}
	got := MailForwardListOrFallback(stored, fallback)
	if len(got) != 1 || got[0] != "tikhvin-palomnik@yandex.ru" {
		t.Fatalf("stored wins: %v", got)
	}
	got = MailForwardListOrFallback(nil, fallback)
	if len(got) != 1 || got[0] != "smailikin70@yandex.ru" {
		t.Fatalf("fallback: %v", got)
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
