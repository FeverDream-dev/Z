package channels

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"
)

// WhatsAppAdapter connects WhatsApp Cloud API to Zsistant agents.
type WhatsAppAdapter struct {
	phoneNumberID string
	accessToken   string
	verifyToken   string
	bindings      map[string]string // phone_number -> agent_id
	client        *http.Client
	dryRun        bool
}

// NewWhatsAppAdapter creates a new WhatsApp adapter.
// If accessToken is empty, the adapter operates in dry-run mode.
func NewWhatsAppAdapter(phoneNumberID, accessToken, verifyToken string) *WhatsAppAdapter {
	return &WhatsAppAdapter{
		phoneNumberID: phoneNumberID,
		accessToken:   accessToken,
		verifyToken:   verifyToken,
		bindings:      make(map[string]string),
		client:        &http.Client{Timeout: 30 * time.Second},
		dryRun:        accessToken == "",
	}
}

// IsDryRun returns true if the adapter is in test/dry-run mode.
func (w *WhatsAppAdapter) IsDryRun() bool {
	return w.dryRun
}

// BindPhone maps a WhatsApp phone number to an agent ID.
func (w *WhatsAppAdapter) BindPhone(phoneNumber, agentID string) {
	w.bindings[phoneNumber] = agentID
}

// AgentForPhone returns the agent ID bound to a phone number, or empty string.
func (w *WhatsAppAdapter) AgentForPhone(phoneNumber string) string {
	return w.bindings[phoneNumber]
}

// SendMessage sends a text message to a WhatsApp phone number.
// In dry-run mode, it prints to stdout instead of calling the API.
func (w *WhatsAppAdapter) SendMessage(phoneNumber, text string) error {
	if w.dryRun {
		fmt.Printf("[DRY-RUN] Would send WhatsApp to %s: %s\n", phoneNumber, text)
		return nil
	}
	return fmt.Errorf("WhatsApp Cloud API not yet implemented (use dry-run mode)")
}

// VerifyWebhook checks the hub.verify_token and returns the hub.challenge.
// This is used during Meta webhook subscription setup.
func (w *WhatsAppAdapter) VerifyWebhook(mode, token, challenge string) (string, error) {
	if mode != "subscribe" {
		return "", fmt.Errorf("invalid mode: %s", mode)
	}
	if token != w.verifyToken {
		return "", fmt.Errorf("verify token mismatch")
	}
	return challenge, nil
}

// TestEvent simulates receiving a WhatsApp message in dry-run mode.
func (w *WhatsAppAdapter) TestEvent(phoneNumber, fromName, text string) Message {
    phoneID := phoneNumber
    return Message{
        ID:               fmt.Sprintf("wa-%d", time.Now().Unix()),
        AgentID:          w.AgentForPhone(phoneID),
        Content:          text,
        Channel:          "whatsapp",
        ChannelType:      "whatsapp",
        ChannelAccountID: phoneID,
        ConversationID:   phoneID,
        SenderID:         "",
        Attachments:      nil,
        CreatedAt:        time.Now(),
        RawEventRef:      "test-whatsapp-" + fmt.Sprint(time.Now().Unix()),
    }
}

// ValidateWebhookSignature verifies the X-Hub-Signature-256 header.
func ValidateWebhookSignature(body []byte, signature, appSecret string) bool {
	if signature == "" || appSecret == "" {
		return false
	}
	mac := hmac.New(sha256.New, []byte(appSecret))
	mac.Write(body)
	expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(signature), []byte(expected))
}
