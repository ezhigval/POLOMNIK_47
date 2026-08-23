package application_test

import (
	"context"
	"testing"
	"time"

	"palomnik/internal/adapters/repository/memory"
	"palomnik/internal/application"
	"palomnik/internal/domain"
	"palomnik/internal/ports"
)

func TestUpdateProfileChangesName(t *testing.T) {
	ctx := context.Background()
	store := memory.NewStore()
	svc := application.NewAuthService(store, store, nil, nil, application.SocialAuthConfig{}, "secret", time.Hour, "http://localhost:3000", store)

	created, err := svc.Register(ctx, application.RegisterInput{
		Email:    "keep@example.com",
		Phone:    "+79001230001",
		Name:     "Старое Имя",
		Password: "password1",
	})
	if err != nil {
		t.Fatalf("register: %v", err)
	}

	updated, err := svc.UpdateProfile(ctx, created.User.ID, application.UpdateProfileInput{
		Name:  "Новое Имя",
		Email: "keep@example.com",
		Phone: "+79001230001",
	})
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if updated.Name != "Новое Имя" || updated.Email != "keep@example.com" {
		t.Fatalf("unexpected profile: %+v", updated)
	}
	if updated.PasswordHash != "" {
		t.Fatal("password hash must not leak")
	}
}

func TestUpdateProfileRejectsDuplicateEmail(t *testing.T) {
	ctx := context.Background()
	store := memory.NewStore()
	svc := application.NewAuthService(store, store, nil, nil, application.SocialAuthConfig{}, "secret", time.Hour, "http://localhost:3000", store)

	first, err := svc.Register(ctx, application.RegisterInput{
		Email:    "first@example.com",
		Phone:    "+79001230002",
		Name:     "Первый",
		Password: "password1",
	})
	if err != nil {
		t.Fatalf("register first: %v", err)
	}
	second, err := svc.Register(ctx, application.RegisterInput{
		Email:    "second@example.com",
		Phone:    "+79001230003",
		Name:     "Второй",
		Password: "password1",
	})
	if err != nil {
		t.Fatalf("register second: %v", err)
	}

	_, err = svc.UpdateProfile(ctx, second.User.ID, application.UpdateProfileInput{
		Name:  "Второй",
		Email: first.User.Email,
		Phone: second.User.Phone,
	})
	if err != domain.ErrDuplicateEmail {
		t.Fatalf("got %v", err)
	}
}

func TestUpdateProfileRequiresCallWhenPhoneChanges(t *testing.T) {
	ctx := context.Background()
	store := memory.NewStore()
	phones := &fakePhone{available: true, status: ports.PhoneCheckPending}
	svc := application.NewAuthService(store, store, phones, nil, application.SocialAuthConfig{}, "secret", time.Hour, "http://localhost:3000", store)

	start, err := svc.StartPhoneVerification(ctx, "+79001230004")
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	phones.status = ports.PhoneCheckConfirmed
	created, err := svc.Register(ctx, application.RegisterInput{
		Email:        "phone@example.com",
		Phone:        "+79001230004",
		Name:         "Телефон",
		Password:     "password1",
		PhoneCheckID: start.CheckID,
	})
	if err != nil {
		t.Fatalf("register: %v", err)
	}

	_, err = svc.UpdateProfile(ctx, created.User.ID, application.UpdateProfileInput{
		Name:  "Телефон",
		Email: "phone@example.com",
		Phone: "+79001230005",
	})
	if err != application.ErrPhoneVerificationRequired {
		t.Fatalf("got %v", err)
	}

	phones.status = ports.PhoneCheckPending
	next, err := svc.StartPhoneVerification(ctx, "+79001230005")
	if err != nil {
		t.Fatalf("start new phone: %v", err)
	}
	phones.status = ports.PhoneCheckConfirmed
	updated, err := svc.UpdateProfile(ctx, created.User.ID, application.UpdateProfileInput{
		Name:         "Телефон",
		Email:        "phone@example.com",
		Phone:        "+79001230005",
		PhoneCheckID: next.CheckID,
	})
	if err != nil {
		t.Fatalf("update with check: %v", err)
	}
	if updated.Phone != "+79001230005" {
		t.Fatalf("phone not updated: %q", updated.Phone)
	}
}
