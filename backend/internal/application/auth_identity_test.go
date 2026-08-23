package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"palomnik/internal/adapters/repository/memory"
	"palomnik/internal/domain"
	"palomnik/internal/ports"
)

func testAuthService(store *memory.Store) *AuthService {
	return NewAuthService(store, store, nil, nil, SocialAuthConfig{}, "secret", time.Hour, "http://localhost:3000", store)
}

func TestOAuthLoginCreatesIdentity(t *testing.T) {
	ctx := context.Background()
	store := memory.NewStore()
	svc := testAuthService(store)

	result, err := svc.OAuthLogin(ctx, OAuthLoginInput{
		Provider:            "Yandex",
		Subject:             "ya-100",
		Email:               "oauth@example.com",
		Name:                "Яндекс Пользователь",
		ConsentPersonalData: true,
		ConsentTerms:        true,
	})
	if err != nil {
		t.Fatalf("oauth login: %v", err)
	}
	if result.Linked || result.Merged {
		t.Fatalf("fresh login should not be a link/merge: %+v", result)
	}
	if !result.Created {
		t.Fatalf("expected Created=true for new oauth user")
	}

	identities, err := svc.ListIdentities(ctx, result.User.ID)
	if err != nil {
		t.Fatalf("list identities: %v", err)
	}
	if len(identities) != 1 || identities[0].Provider != "yandex" || identities[0].Subject != "ya-100" {
		t.Fatalf("expected yandex identity, got %+v", identities)
	}

	again, err := svc.OAuthLogin(ctx, OAuthLoginInput{
		Provider: "yandex",
		Subject:  "ya-100",
		Name:     "Яндекс Пользователь",
	})
	if err != nil {
		t.Fatalf("second oauth login: %v", err)
	}
	if again.User.ID != result.User.ID {
		t.Fatalf("expected same user, got %s vs %s", again.User.ID, result.User.ID)
	}
	if again.Created {
		t.Fatalf("existing oauth login should not set Created")
	}
}

func TestOAuthLoginRequiresConsentForNewUser(t *testing.T) {
	ctx := context.Background()
	svc := testAuthService(memory.NewStore())
	_, err := svc.OAuthLogin(ctx, OAuthLoginInput{
		Provider: "yandex",
		Subject:  "ya-no-consent",
		Name:     "Без согласия",
	})
	if !errors.Is(err, domain.ErrConsentRequired) {
		t.Fatalf("expected ErrConsentRequired, got %v", err)
	}
}

func TestOAuthLinkUnusedIdentity(t *testing.T) {
	ctx := context.Background()
	store := memory.NewStore()
	svc := testAuthService(store)

	registered, err := svc.Register(ctx, RegisterInput{
		Email:    "cabinet@example.com",
		Phone:    "+79001230001",
		Name:     "Кабинет",
		Password: "password1",
	})
	if err != nil {
		t.Fatalf("register: %v", err)
	}

	linked, err := svc.OAuthLogin(ctx, OAuthLoginInput{
		Provider:     "vk",
		Subject:      "vk-77",
		Name:         "VK User",
		SessionToken: registered.Token,
	})
	if err != nil {
		t.Fatalf("link: %v", err)
	}
	if !linked.Linked || linked.Merged {
		t.Fatalf("expected link without merge, got %+v", linked)
	}
	if linked.User.ID != registered.User.ID {
		t.Fatalf("link changed current user")
	}

	identities, err := svc.ListIdentities(ctx, registered.User.ID)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(identities) != 1 || identities[0].Provider != "vk" {
		t.Fatalf("expected vk identity, got %+v", identities)
	}
}

func TestOAuthMergeMovesBookingAndKeepsEmail(t *testing.T) {
	ctx := context.Background()
	store := memory.NewStore()
	svc := testAuthService(store)

	current, err := svc.Register(ctx, RegisterInput{
		Email:    "keep@example.com",
		Phone:    "+79001230002",
		Name:     "Текущий",
		Password: "password1",
	})
	if err != nil {
		t.Fatalf("register: %v", err)
	}

	other, err := svc.OAuthLogin(ctx, OAuthLoginInput{
		Provider:            "yandex",
		Subject:             "ya-merge",
		Email:               "other@example.com",
		Name:                "Другой",
		ConsentPersonalData: true,
		ConsentTerms:        true,
	})
	if err != nil {
		t.Fatalf("oauth other: %v", err)
	}

	tour := testTour()
	if _, err := store.CreateTour(ctx, tour); err != nil {
		t.Fatalf("create tour: %v", err)
	}
	userID := other.User.ID
	booking, err := domain.NewBooking(domain.NewBookingInput{
		ID:          testUUID("22222222-2222-2222-2222-222222222222"),
		Tour:        tour,
		UserID:      &userID,
		Name:        "Другой",
		Phone:       "+79001230099",
		PeopleCount: 1,
	})
	if err != nil {
		t.Fatalf("new booking: %v", err)
	}
	if _, err := store.CreateBooking(ctx, booking); err != nil {
		t.Fatalf("store booking: %v", err)
	}

	passenger, err := domain.NewPassenger(domain.NewPassengerInput{
		ID:        testUUID("33333333-3333-3333-3333-333333333333"),
		UserID:    other.User.ID,
		Name:      "Спутник",
		Phone:     "+79001230088",
		BirthDate: time.Date(1991, 2, 3, 0, 0, 0, 0, time.UTC),
		Passport:  "4010 111111",
	})
	if err != nil {
		t.Fatalf("new passenger: %v", err)
	}
	if _, err := store.CreatePassenger(ctx, passenger); err != nil {
		t.Fatalf("store passenger: %v", err)
	}

	merged, err := svc.OAuthLogin(ctx, OAuthLoginInput{
		Provider:     "yandex",
		Subject:      "ya-merge",
		Name:         "Другой",
		SessionToken: current.Token,
	})
	if err != nil {
		t.Fatalf("merge: %v", err)
	}
	if !merged.Linked || !merged.Merged {
		t.Fatalf("expected linked merge, got %+v", merged)
	}
	if merged.User.Email != "keep@example.com" {
		t.Fatalf("email should stay on current user, got %q", merged.User.Email)
	}
	foundEmail := false
	for _, conflict := range merged.Conflicts {
		if conflict.Field == "email" {
			foundEmail = true
		}
	}
	if !foundEmail {
		t.Fatalf("expected email conflict, got %+v", merged.Conflicts)
	}

	list, err := svc.ListMyBookings(ctx, current.User.ID, ports.Pagination{Page: 1, Limit: 20})
	if err != nil {
		t.Fatalf("list bookings: %v", err)
	}
	if len(list.Items) != 1 || list.Items[0].ID != booking.ID {
		t.Fatalf("booking was not moved, got %+v", list.Items)
	}

	moved, err := store.ListPassengers(ctx, current.User.ID)
	if err != nil {
		t.Fatalf("list passengers: %v", err)
	}
	if len(moved) != 1 || moved[0].ID != passenger.ID {
		t.Fatalf("passenger was not moved, got %+v", moved)
	}

	if _, err := store.GetUserByID(ctx, other.User.ID); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("source user should be deleted, got %v", err)
	}
}
