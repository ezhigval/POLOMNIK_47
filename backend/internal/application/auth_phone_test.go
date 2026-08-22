package application_test

import (
	"context"
	"testing"
	"time"

	"polomnik/internal/adapters/repository/memory"
	"polomnik/internal/application"
	"polomnik/internal/ports"
)

type fakePhone struct {
	available bool
	status    ports.PhoneCheckStatus
	started   bool
}

func (f *fakePhone) Available() bool { return f.available }

func (f *fakePhone) Start(_ context.Context, phone string) (ports.PhoneChallenge, error) {
	if !f.available {
		return ports.PhoneChallenge{}, ports.ErrPhoneVerifierNotConfigured
	}
	f.started = true
	return ports.PhoneChallenge{
		CheckID:         "chk-1",
		CallPhone:       "78005008275",
		CallPhonePretty: "+7 (800) 500-8275",
		ExpiresIn:       300,
	}, nil
}

func (f *fakePhone) Status(_ context.Context, _ string) (ports.PhoneCheckStatus, error) {
	if !f.available {
		return "", ports.ErrPhoneVerifierNotConfigured
	}
	return f.status, nil
}

func TestAuthMethodsPhoneUnavailable(t *testing.T) {
	store := memory.NewStore()
	svc := application.NewAuthService(store, store, &fakePhone{available: false}, "secret", time.Hour)
	methods := svc.AuthMethods()
	if methods.PhoneCall.Available {
		t.Fatal("expected phone unavailable")
	}
	if methods.PhoneCall.Message == "" {
		t.Fatal("expected unavailable message")
	}
}

func TestStartPhoneVerificationUnavailable(t *testing.T) {
	store := memory.NewStore()
	svc := application.NewAuthService(store, store, &fakePhone{available: false}, "secret", time.Hour)
	_, err := svc.StartPhoneVerification(context.Background(), "+79001234567")
	if err != application.ErrPhoneVerificationUnavailable {
		t.Fatalf("got %v", err)
	}
}

func TestRegisterRequiresConfirmedCallWhenPhoneAvailable(t *testing.T) {
	store := memory.NewStore()
	phones := &fakePhone{available: true, status: ports.PhoneCheckPending}
	svc := application.NewAuthService(store, store, phones, "secret", time.Hour)

	_, err := svc.Register(context.Background(), application.RegisterInput{
		Name:     "Иван Тест",
		Phone:    "+79001234567",
		Password: "password1",
	})
	if err != application.ErrPhoneVerificationRequired {
		t.Fatalf("got %v", err)
	}

	start, err := svc.StartPhoneVerification(context.Background(), "+79001234567")
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	phones.status = ports.PhoneCheckConfirmed

	_, err = svc.Register(context.Background(), application.RegisterInput{
		Name:         "Иван Тест",
		Phone:        "+79001234567",
		Password:     "password1",
		PhoneCheckID: start.CheckID,
	})
	if err != nil {
		t.Fatalf("register: %v", err)
	}
}

func TestRegisterWithoutPhoneAdapter(t *testing.T) {
	store := memory.NewStore()
	svc := application.NewAuthService(store, store, nil, "secret", time.Hour)
	_, err := svc.Register(context.Background(), application.RegisterInput{
		Name:     "Иван Тест",
		Phone:    "+79001234567",
		Password: "password1",
	})
	if err != nil {
		t.Fatalf("register without verifier: %v", err)
	}
}
