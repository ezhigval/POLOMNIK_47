package domain

import (
	"testing"
	"time"
)

func TestParseTelegramUsernameList(t *testing.T) {
	got, err := ParseTelegramUsernameList("@EzhigVal, other_user\n@ezhigval")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	want := []string{"ezhigval", "other_user"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %v, want %v", got, want)
		}
	}

	if _, err := ParseTelegramUsernameList("ab"); err != ErrInvalidTelegramUsername {
		t.Fatalf("expected ErrInvalidTelegramUsername, got %v", err)
	}
	if _, err := ParseTelegramUsernameList("1badname"); err != ErrInvalidTelegramUsername {
		t.Fatalf("expected leading letter, got %v", err)
	}
}

func TestNewTelegramRecipientsStoresBothLists(t *testing.T) {
	got, err := NewTelegramRecipients("@manager_one", "support_bot", time.Time{})
	if err != nil {
		t.Fatalf("new: %v", err)
	}
	if len(got.BookingUsernames) != 1 || got.BookingUsernames[0] != "manager_one" {
		t.Fatalf("booking: %v", got.BookingUsernames)
	}
	if len(got.SupportUsernames) != 1 || got.SupportUsernames[0] != "support_bot" {
		t.Fatalf("support: %v", got.SupportUsernames)
	}
}

func TestResolveTelegramTargetsSkipsUnknownAndUsesFallbackOnlyWhenEmpty(t *testing.T) {
	bindings := map[string]string{"ezhigval": "111"}

	withList := ResolveTelegramTargets([]string{"ezhigval", "unknown_user"}, bindings, "999")
	if len(withList.ChatIDs) != 1 || withList.ChatIDs[0] != "111" {
		t.Fatalf("resolved chat ids: %v", withList.ChatIDs)
	}
	if len(withList.Waiting) != 1 || withList.Waiting[0] != "unknown_user" {
		t.Fatalf("waiting: %v", withList.Waiting)
	}

	fallback := ResolveTelegramTargets(nil, bindings, "999")
	if len(fallback.ChatIDs) != 1 || fallback.ChatIDs[0] != "999" {
		t.Fatalf("fallback: %v", fallback.ChatIDs)
	}

	empty := ResolveTelegramTargets([]string{"unknown_user"}, bindings, "999")
	if len(empty.ChatIDs) != 0 {
		t.Fatalf("unknown list must not use fallback, got %v", empty.ChatIDs)
	}
}
