package channels

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func TestNewWhatsAppAdapterDryRun(t *testing.T) {
	w := NewWhatsAppAdapter("1234567890", "", "verify123")
	if !w.IsDryRun() {
		t.Fatal("expected dry-run mode when accessToken is empty")
	}
}

func TestNewWhatsAppAdapterReal(t *testing.T) {
	w := NewWhatsAppAdapter("1234567890", "EAAB real_token", "verify123")
	if w.IsDryRun() {
		t.Fatal("expected real mode when accessToken is provided")
	}
}

func TestBindPhone(t *testing.T) {
	w := NewWhatsAppAdapter("", "", "")
	w.BindPhone("+1234567890", "agent1")
	if w.AgentForPhone("+1234567890") != "agent1" {
		t.Fatalf("expected agent1 for phone +1234567890, got %s", w.AgentForPhone("+1234567890"))
	}
}

func TestWhatsAppSendMessageDryRun(t *testing.T) {
	w := NewWhatsAppAdapter("", "", "")
	if err := w.SendMessage("+1234567890", "hello"); err != nil {
		t.Fatalf("dry-run send should not error: %v", err)
	}
}

func TestVerifyWebhook(t *testing.T) {
	w := NewWhatsAppAdapter("", "", "my_verify_token")
	challenge, err := w.VerifyWebhook("subscribe", "my_verify_token", "abc123")
	if err != nil {
		t.Fatalf("verify webhook: %v", err)
	}
	if challenge != "abc123" {
		t.Fatalf("expected challenge abc123, got %s", challenge)
	}
}

func TestVerifyWebhookBadMode(t *testing.T) {
	w := NewWhatsAppAdapter("", "", "my_verify_token")
	_, err := w.VerifyWebhook("deny", "my_verify_token", "abc123")
	if err == nil {
		t.Fatal("expected error for invalid mode")
	}
}

func TestVerifyWebhookBadToken(t *testing.T) {
	w := NewWhatsAppAdapter("", "", "my_verify_token")
	_, err := w.VerifyWebhook("subscribe", "wrong_token", "abc123")
	if err == nil {
		t.Fatal("expected error for mismatched token")
	}
}

func TestWhatsAppTestEvent(t *testing.T) {
	w := NewWhatsAppAdapter("", "", "")
	w.BindPhone("+1234567890", "agent1")
	msg := w.TestEvent("+1234567890", "tester", "hello whatsapp")
	if msg.AgentID != "agent1" {
		t.Fatalf("expected agent1, got %s", msg.AgentID)
	}
	if msg.Content != "hello whatsapp" {
		t.Fatalf("expected 'hello whatsapp', got %s", msg.Content)
	}
	if msg.Channel != "whatsapp" {
		t.Fatalf("expected channel whatsapp, got %s", msg.Channel)
	}
}

func TestValidateWebhookSignature(t *testing.T) {
	body := []byte(`{"test":"data"}`)
	secret := "my_app_secret"
	// Generate a valid signature
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	validSig := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	if !ValidateWebhookSignature(body, validSig, secret) {
		t.Fatal("expected valid signature to pass")
	}
	if ValidateWebhookSignature(body, "sha256=invalid", secret) {
		t.Fatal("expected invalid signature to fail")
	}
	if ValidateWebhookSignature(body, "", secret) {
		t.Fatal("expected empty signature to fail")
	}
	if ValidateWebhookSignature(body, validSig, "") {
		t.Fatal("expected empty secret to fail")
	}
}
