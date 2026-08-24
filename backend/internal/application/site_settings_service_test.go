package application

import (
	"context"
	"testing"

	"palomnik/internal/adapters/repository/memory"
	"palomnik/internal/domain"
)

func TestSiteSettingsPublicUsesMailForwardFallbackWhenEmpty(t *testing.T) {
	store := memory.NewStore()
	svc := NewSiteSettingsService(store, domain.SiteSettings{
		MailForwardTo: []string{"smailikin70@yandex.ru"},
	})

	got, err := svc.Public(context.Background())
	if err != nil {
		t.Fatalf("public: %v", err)
	}
	if len(got.MailForwardTo) != 1 || got.MailForwardTo[0] != "smailikin70@yandex.ru" {
		t.Fatalf("fallback: %v", got.MailForwardTo)
	}

	_, err = svc.Update(context.Background(), domain.SiteSettings{
		MailForwardTo: []string{"tikhvin-palomnik@yandex.ru"},
	})
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	got, err = svc.Public(context.Background())
	if err != nil {
		t.Fatalf("public after update: %v", err)
	}
	if len(got.MailForwardTo) != 1 || got.MailForwardTo[0] != "tikhvin-palomnik@yandex.ru" {
		t.Fatalf("stored list: %v", got.MailForwardTo)
	}

	if _, err := svc.Update(context.Background(), domain.SiteSettings{}); err != nil {
		t.Fatalf("clear: %v", err)
	}
	got, err = svc.Public(context.Background())
	if err != nil {
		t.Fatalf("public after clear: %v", err)
	}
	if len(got.MailForwardTo) != 1 || got.MailForwardTo[0] != "smailikin70@yandex.ru" {
		t.Fatalf("fallback after clear: %v", got.MailForwardTo)
	}
}
