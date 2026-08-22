package domain

import (
	"testing"
	"time"
)

func TestNotificationRoutingFromTelegramRecipients(t *testing.T) {
	legacy := TelegramRecipients{
		BookingUsernames: []string{"ezhigval", "manager_one"},
		SupportUsernames: []string{"support_bot"},
	}
	got := NotificationRoutingFromTelegramRecipients(legacy, time.Time{})

	booking := RecipientsForEvent(got, NotificationEventBookingCreated)
	if len(booking) != 2 || booking[0].Channel != NotificationChannelTelegram || booking[0].Address != "ezhigval" {
		t.Fatalf("booking_created: %+v", booking)
	}
	status := RecipientsForEvent(got, NotificationEventBookingStatusChanged)
	if len(status) != 2 || status[1].Address != "manager_one" {
		t.Fatalf("booking_status_changed must copy booking list, got %+v", status)
	}
	support := RecipientsForEvent(got, NotificationEventSupportMessage)
	if len(support) != 1 || support[0].Address != "support_bot" {
		t.Fatalf("support: %+v", support)
	}
}

func TestNewNotificationRecipientNormalizesChannels(t *testing.T) {
	tg, err := NewNotificationRecipient("telegram", "@EzhigVal")
	if err != nil {
		t.Fatalf("telegram: %v", err)
	}
	if tg.Address != "ezhigval" {
		t.Fatalf("telegram address: %q", tg.Address)
	}

	maxRecipient, err := NewNotificationRecipient("max", "+7 (999) 999-99-99")
	if err != nil {
		t.Fatalf("max: %v", err)
	}
	if maxRecipient.Address != "+79999999999" {
		t.Fatalf("max address: %q", maxRecipient.Address)
	}

	if _, err := NewNotificationRecipient("whatsapp", "x"); err != ErrInvalidNotificationChannel {
		t.Fatalf("expected invalid channel, got %v", err)
	}
	if _, err := NewNotificationRecipient("max", "12"); err != ErrInvalidNotificationAddress {
		t.Fatalf("expected invalid max address, got %v", err)
	}
}

func TestNewNotificationRoutingDedupes(t *testing.T) {
	routing, err := NewNotificationRouting(map[NotificationEventKind][]NotificationRecipient{
		NotificationEventBookingCreated: {
			{Channel: NotificationChannelTelegram, Address: "@EzhigVal"},
			{Channel: NotificationChannelTelegram, Address: "ezhigval"},
			{Channel: NotificationChannelMax, Address: "+79991112233"},
		},
	}, time.Time{})
	if err != nil {
		t.Fatalf("new: %v", err)
	}
	list := RecipientsForEvent(routing, NotificationEventBookingCreated)
	if len(list) != 2 {
		t.Fatalf("expected 2 recipients after dedupe, got %+v", list)
	}
	if len(RecipientsForEvent(routing, NotificationEventSupportMessage)) != 0 {
		t.Fatal("missing events must stay empty")
	}
}
