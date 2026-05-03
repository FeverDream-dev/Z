package channels

import (
	"testing"
)

func TestNewDiscordAdapterDryRun(t *testing.T) {
	d := NewDiscordAdapter("")
	if !d.IsDryRun() {
		t.Fatal("expected dry-run mode when token is empty")
	}
}

func TestNewDiscordAdapterReal(t *testing.T) {
	d := NewDiscordAdapter("MTExMjIyMzMzNDQ0NTU1NjY2Nzc3.ODg5OQ.aaa_bbb_ccc")
	if d.IsDryRun() {
		t.Fatal("expected real mode when token is provided")
	}
}

func TestBindChannel(t *testing.T) {
	d := NewDiscordAdapter("")
	d.BindChannel("123456789", "agent1")
	if d.AgentForChannel("123456789") != "agent1" {
		t.Fatalf("expected agent1 for channel 123456789, got %s", d.AgentForChannel("123456789"))
	}
}

func TestDiscordSendMessageDryRun(t *testing.T) {
	d := NewDiscordAdapter("")
	if err := d.SendMessage("123456789", "hello"); err != nil {
		t.Fatalf("dry-run send should not error: %v", err)
	}
}

func TestTestEvent(t *testing.T) {
	d := NewDiscordAdapter("")
	d.BindChannel("123456789", "agent1")
	msg := d.TestEvent("123456789", "tester", "hello discord")
	if msg.AgentID != "agent1" {
		t.Fatalf("expected agent1, got %s", msg.AgentID)
	}
	if msg.Content != "hello discord" {
		t.Fatalf("expected 'hello discord', got %s", msg.Content)
	}
	if msg.Channel != "discord" {
		t.Fatalf("expected channel discord, got %s", msg.Channel)
	}
}

func TestValidateDiscordToken(t *testing.T) {
	if err := ValidateDiscordToken("MTExMjIyMzMzNDQ0NTU1NjY2Nzc3.ODg5OQ.aaa_bbb_ccc"); err != nil {
		t.Fatalf("valid token failed: %v", err)
	}
	if err := ValidateDiscordToken(""); err == nil {
		t.Fatal("empty token should fail")
	}
	if err := ValidateDiscordToken("short"); err == nil {
		t.Fatal("short token should fail")
	}
}
