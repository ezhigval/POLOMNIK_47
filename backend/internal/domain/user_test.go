package domain

import "testing"

func TestFillEmptyProfileCopiesOnlyEmptyFields(t *testing.T) {
	into := User{Name: "Анна", Email: "anna@example.com", Phone: "+79001111111", PasswordHash: "keep"}
	from := User{Name: "Борис", Email: "boris@example.com", Phone: "+79002222222", PasswordHash: "other"}

	got, conflicts := FillEmptyProfile(into, from)
	if got.Name != "Анна" || got.Email != "anna@example.com" || got.Phone != "+79001111111" {
		t.Fatalf("kept fields changed: %+v", got)
	}
	if got.PasswordHash != "keep" {
		t.Fatalf("password should stay on target, got %q", got.PasswordHash)
	}
	if len(conflicts) != 3 {
		t.Fatalf("expected name/email/phone conflicts, got %+v", conflicts)
	}
}

func TestFillEmptyProfileFillsBlanksAndCopiesPassword(t *testing.T) {
	into := User{Name: "Анна"}
	from := User{Name: "Анна", Email: "anna@example.com", Phone: "+79001111111", PasswordHash: "hash"}

	got, conflicts := FillEmptyProfile(into, from)
	if len(conflicts) != 0 {
		t.Fatalf("unexpected conflicts: %+v", conflicts)
	}
	if got.Email != "anna@example.com" || got.Phone != "+79001111111" || got.PasswordHash != "hash" {
		t.Fatalf("empty fields were not copied: %+v", got)
	}
}

func TestFillEmptyProfileIgnoresCaseOnEmail(t *testing.T) {
	into := User{Email: "Anna@example.com"}
	from := User{Email: "anna@example.com"}
	got, conflicts := FillEmptyProfile(into, from)
	if len(conflicts) != 0 {
		t.Fatalf("same email in different case should not conflict, got %+v", conflicts)
	}
	if got.Email != "Anna@example.com" {
		t.Fatalf("kept email changed: %q", got.Email)
	}
}

func TestWithProfileAllowsEmptyPhone(t *testing.T) {
	user := User{Name: "Старое", Phone: "+79001111111"}
	got, err := user.WithProfile("Новое Имя", "user@example.com", "")
	if err != nil {
		t.Fatalf("with profile: %v", err)
	}
	if got.Name != "Новое Имя" || got.Email != "user@example.com" || got.Phone != "" {
		t.Fatalf("unexpected profile: %+v", got)
	}
}

func TestWithProfileRejectsEmptyName(t *testing.T) {
	user := User{Name: "Анна"}
	if _, err := user.WithProfile("  ", "", ""); err != ErrInvalidContactName {
		t.Fatalf("got %v", err)
	}
}
