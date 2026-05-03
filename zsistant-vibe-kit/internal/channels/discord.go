package channels

import (
	"fmt"
	"net/http"
	"time"
)

// DiscordAdapter connects a Discord bot to Zsistant agents.
type DiscordAdapter struct {
	token     string
	bindings  map[string]string // channel_id -> agent_id
	client    *http.Client
	dryRun    bool
}

// NewDiscordAdapter creates a new Discord adapter.
// If token is empty, the adapter operates in dry-run mode.
func NewDiscordAdapter(token string) *DiscordAdapter {
	return &DiscordAdapter{
		token:     token,
		bindings:  make(map[string]string),
		client:    &http.Client{Timeout: 30 * time.Second},
		dryRun:    token == "",
	}
}

// IsDryRun returns true if the adapter is in test/dry-run mode.
func (d *DiscordAdapter) IsDryRun() bool {
	return d.dryRun
}

// BindChannel maps a Discord channel ID to an agent ID.
func (d *DiscordAdapter) BindChannel(channelID, agentID string) {
	d.bindings[channelID] = agentID
}

// AgentForChannel returns the agent ID bound to a channel, or empty string.
func (d *DiscordAdapter) AgentForChannel(channelID string) string {
	return d.bindings[channelID]
}

// SendMessage sends a text message to a Discord channel.
// In dry-run mode, it prints to stdout instead of calling the API.
func (d *DiscordAdapter) SendMessage(channelID, text string) error {
	if d.dryRun {
		fmt.Printf("[DRY-RUN] Would send to Discord channel %s: %s\n", channelID, text)
		return nil
	}
	// Real Discord API call would go here
	return fmt.Errorf("Discord API not yet implemented (use dry-run mode)")
}

// TestEvent simulates receiving a Discord message in dry-run mode.
func (d *DiscordAdapter) TestEvent(channelID, fromName, text string) Message {
    channelIDStr := channelID
    return Message{
        ID:               fmt.Sprintf("discord-%d", time.Now().Unix()),
        AgentID:          d.AgentForChannel(channelID),
        Content:          text,
        Channel:          "discord",
        ChannelType:      "discord",
        ChannelAccountID: channelIDStr,
        ConversationID:   channelIDStr,
        SenderID:         "",
        Attachments:      nil,
        CreatedAt:        time.Now(),
        RawEventRef:      "test-discord-" + fmt.Sprint(time.Now().Unix()),
    }
}

// ValidateToken checks if a token looks like a valid Discord bot token format.
func ValidateDiscordToken(token string) error {
	if token == "" {
		return fmt.Errorf("token is empty")
	}
	// Discord bot tokens are typically long alphanumeric strings
	if len(token) < 20 {
		return fmt.Errorf("token too short for a Discord bot token")
	}
	return nil
}
