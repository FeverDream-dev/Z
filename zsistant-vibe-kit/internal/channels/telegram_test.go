package channels

import (
    "encoding/json"
    "strings"
    "testing"
)

func TestNewTelegramAdapterDryRun(t *testing.T) {
	a := NewTelegramAdapter("")
	if !a.IsDryRun() {
		t.Fatal("expected dry-run mode when token is empty")
	}
}

func TestNewTelegramAdapterReal(t *testing.T) {
	a := NewTelegramAdapter("123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11")
	if a.IsDryRun() {
		t.Fatal("expected real mode when token is provided")
	}
}

func TestBindChat(t *testing.T) {
	a := NewTelegramAdapter("")
	a.BindChat("12345", "agent1")
	if a.AgentForChat("12345") != "agent1" {
		t.Fatalf("expected agent1 for chat 12345, got %s", a.AgentForChat("12345"))
	}
}

func TestSendMessageDryRun(t *testing.T) {
	a := NewTelegramAdapter("")
	// Should not error in dry-run mode
	if err := a.SendMessage("12345", "hello"); err != nil {
		t.Fatalf("dry-run send should not error: %v", err)
	}
}

func TestTestMessage(t *testing.T) {
	a := NewTelegramAdapter("")
	a.BindChat("12345", "agent1")
	msg := a.TestMessage("12345", "tester", "hello world")
	if msg.AgentID != "agent1" {
		t.Fatalf("expected agent1, got %s", msg.AgentID)
	}
	if msg.Content != "hello world" {
		t.Fatalf("expected 'hello world', got %s", msg.Content)
	}
	if msg.Channel != "telegram" {
		t.Fatalf("expected channel telegram, got %s", msg.Channel)
	}
}

func TestRedactToken(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"", "***"},
		{"short", "***"},
		{"123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11", "1234...ew11"},
	}
	for _, c := range cases {
		got := RedactToken(c.input)
		if got != c.expected {
			t.Fatalf("RedactToken(%q) = %q, want %q", c.input, got, c.expected)
		}
	}
}

func TestValidateToken(t *testing.T) {
	if err := ValidateToken("123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"); err != nil {
		t.Fatalf("valid token failed: %v", err)
	}
	if err := ValidateToken(""); err == nil {
		t.Fatal("empty token should fail")
	}
	if err := ValidateToken("notoken"); err == nil {
		t.Fatal("missing colon should fail")
	}
	if err := ValidateToken("123:"); err == nil {
		t.Fatal("empty second part should fail")
	}
}

func TestGetUpdatesDryRun(t *testing.T) {
	a := NewTelegramAdapter("")
	_, err := a.GetUpdates(0)
	if err == nil {
		t.Fatal("expected error in dry-run mode")
	}
	if !strings.Contains(err.Error(), "dry-run") {
		t.Fatalf("expected dry-run error, got: %v", err)
	}
}

func TestUpdateJSONParsing(t *testing.T) {
    jsonStr := `{
        "ok": true,
        "result": [
            {
                "update_id": 12345,
                "message": {
                    "message_id": 1,
                    "from": {"id": 111, "is_bot": false, "first_name": "John"},
                    "chat": {"id": 111, "type": "private"},
                    "date": 1234567890,
                    "text": "Hello"
                }
            }
        ]
    }`
    var apiResp telegramAPIResponse
    if err := json.Unmarshal([]byte(jsonStr), &apiResp); err != nil {
        t.Fatalf("unexpected error unmarshalling: %v", err)
    }
    if len(apiResp.Result) != 1 {
        t.Fatalf("expected 1 update, got %d", len(apiResp.Result))
    }
    upd := apiResp.Result[0]
    if upd.ID != 12345 {
        t.Fatalf("expected update_id 12345, got %d", upd.ID)
    }
    if upd.Message == nil {
        t.Fatalf("expected message in update")
    }
    if upd.Message.Text != "Hello" {
        t.Fatalf("expected text 'Hello', got %q", upd.Message.Text)
    }
    if upd.Message.Chat == nil || upd.Message.Chat.ID != 111 {
        t.Fatalf("unexpected chat id: %#v", upd.Message.Chat)
    }
    if upd.Message.From == nil || upd.Message.From.FirstName != "John" {
        t.Fatalf("unexpected from name: %#v", upd.Message.From)
    }
}

func TestUpdateNilMessage(t *testing.T) {
    jsonStr := `{
        "ok": true,
        "result": [
            {"update_id": 999, "message": null}
        ]
    }`
    var apiResp telegramAPIResponse
    if err := json.Unmarshal([]byte(jsonStr), &apiResp); err != nil {
        t.Fatalf("unexpected error unmarshalling: %v", err)
    }
    if len(apiResp.Result) != 1 {
        t.Fatalf("expected 1 update, got %d", len(apiResp.Result))
    }
    if apiResp.Result[0].Message != nil {
        t.Fatalf("expected nil message, got %#v", apiResp.Result[0].Message)
    }
}
