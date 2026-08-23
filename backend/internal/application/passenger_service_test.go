package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"palomnik/internal/adapters/repository/memory"
	"palomnik/internal/domain"
)

func TestPassengerServiceCRUD(t *testing.T) {
	ctx := context.Background()
	store := memory.NewStore()
	svc := NewPassengerService(store)

	user, err := domain.NewUser(domain.RegisterUserInput{
		ID:           uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		Email:        "pax@example.com",
		Phone:        "+79001230001",
		Name:         "Хозяин",
		PasswordHash: "hashhash",
	})
	if err != nil {
		t.Fatalf("user: %v", err)
	}
	if _, err := store.CreateUser(ctx, user); err != nil {
		t.Fatalf("store user: %v", err)
	}

	created, err := svc.Create(ctx, user.ID, PassengerInput{
		Name:      "Анна Паломница",
		Phone:     "89001112233",
		BirthDate: time.Date(1985, 3, 1, 0, 0, 0, 0, time.UTC),
		Passport:  "4010 654321",
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.Phone != "+79001112233" {
		t.Fatalf("phone: %q", created.Phone)
	}

	list, err := svc.List(ctx, user.ID)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(list) != 1 || list[0].ID != created.ID {
		t.Fatalf("list: %+v", list)
	}

	updated, err := svc.Update(ctx, user.ID, created.ID, PassengerInput{
		Name:      "Анна Иванова",
		Phone:     "+79001112233",
		BirthDate: time.Date(1985, 3, 1, 0, 0, 0, 0, time.UTC),
		Passport:  "4010 654321",
	})
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if updated.Name != "Анна Иванова" {
		t.Fatalf("name: %q", updated.Name)
	}

	other := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	if _, err := svc.Update(ctx, other, created.ID, PassengerInput{
		Name:      "Чужой",
		Phone:     "+79001112233",
		BirthDate: time.Date(1985, 3, 1, 0, 0, 0, 0, time.UTC),
		Passport:  "4010 654321",
	}); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("foreign update: %v", err)
	}

	if err := svc.Delete(ctx, user.ID, created.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}
	list, err = svc.List(ctx, user.ID)
	if err != nil {
		t.Fatalf("list after delete: %v", err)
	}
	if len(list) != 0 {
		t.Fatalf("expected empty list, got %+v", list)
	}
}
