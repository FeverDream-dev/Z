package channels

import (
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "os"
    "strings"
    "time"
)

// TelegramAdapter connects a Telegram bot to Zsistant agents.
type TelegramAdapter struct {
    token    string
    apiBase  string
    bindings map[string]string // chat_id -> agent_id
    client   *http.Client
    dryRun   bool
}

// NewTelegramAdapter creates a new Telegram adapter.
// If token is empty, the adapter operates in dry-run mode.
func NewTelegramAdapter(token string) *TelegramAdapter {
    dry := token == ""
    base := "https://api.telegram.org/bot" + token
    if dry {
        base = ""
    }
    return &TelegramAdapter{
        token:    token,
        apiBase:  base,
        bindings: make(map[string]string),
        client:   &http.Client{Timeout: 30 * time.Second},
        dryRun:   dry,
    }
}

// IsDryRun returns true if the adapter is in test/dry-run mode.
func (t *TelegramAdapter) IsDryRun() bool {
    return t.dryRun
}

// BindChat maps a Telegram chat ID to an agent ID.
func (t *TelegramAdapter) BindChat(chatID, agentID string) {
    t.bindings[chatID] = agentID
}

// AgentForChat returns the agent ID bound to a chat, or empty string.
func (t *TelegramAdapter) AgentForChat(chatID string) string {
    return t.bindings[chatID]
}

// SendMessage sends a text message to a Telegram chat.
// In dry-run mode, it prints to stdout instead of calling the API.
func (t *TelegramAdapter) SendMessage(chatID, text string) error {
    if t.dryRun {
        fmt.Printf("[DRY-RUN] Would send to chat %s: %s\n", chatID, text)
        return nil
    }
    u := fmt.Sprintf("%s/sendMessage", t.apiBase)
    data := url.Values{}
    data.Set("chat_id", chatID)
    data.Set("text", text)
    resp, err := t.client.PostForm(u, data)
    if err != nil {
        return fmt.Errorf("telegram API error: %w", err)
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return fmt.Errorf("telegram API returned %d: %s", resp.StatusCode, string(body))
    }
    return nil
}

// GetUpdates fetches pending updates from Telegram (long polling).
// In dry-run mode, it returns an error.
func (t *TelegramAdapter) GetUpdates(offset int) ([]Update, error) {
    if t.dryRun {
        return nil, fmt.Errorf("cannot poll in dry-run mode")
    }
    u := fmt.Sprintf("%s/getUpdates?offset=%d&limit=10&timeout=30", t.apiBase, offset)
    resp, err := t.client.Get(u)
    if err != nil {
        return nil, fmt.Errorf("telegram API error: %w", err)
    }
    defer resp.Body.Close()
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("reading response body: %w", err)
    }
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("telegram API returned %d: %s", resp.StatusCode, string(body))
    }
    var apiResp telegramAPIResponse
    if err := json.Unmarshal(body, &apiResp); err != nil {
        return nil, fmt.Errorf("parsing telegram response: %w", err)
    }
    if !apiResp.OK {
        return nil, fmt.Errorf("telegram API returned ok=false: %s", string(body))
    }
    return apiResp.Result, nil
}

// Update represents a Telegram update from the Bot API.
type Update struct {
    ID      int            `json:"update_id"`
    Message *UpdateMessage `json:"message,omitempty"`
}

// UpdateMessage represents a message inside a Telegram update.
type UpdateMessage struct {
    ID   int        `json:"message_id"`
    From *UpdateFrom `json:"from,omitempty"`
    Chat *UpdateChat  `json:"chat,omitempty"`
    Date int64      `json:"date"`
    Text string     `json:"text"`
}

// UpdateFrom represents the sender of a message.
type UpdateFrom struct {
    ID        int    `json:"id"`
    IsBot     bool   `json:"is_bot"`
    FirstName string `json:"first_name"`
}

// UpdateChat represents the chat a message was sent in.
type UpdateChat struct {
    ID   int    `json:"id"`
    Type string `json:"type"`
}

// telegramAPIResponse is the response from getUpdates.
type telegramAPIResponse struct {
    OK     bool     `json:"ok"`
    Result []Update `json:"result"`
}

// MessageHandler is a callback for processing inbound Telegram messages.
type MessageHandler func(chatID int, fromName, text string) string

// Listen starts a long-poll loop for Telegram updates.
// It calls the handler for each text message and sends the response back.
// Blocks until the context is cancelled.
func (t *TelegramAdapter) Listen(ctx context.Context, handler MessageHandler) error {
    if t.dryRun {
        return fmt.Errorf("cannot listen in dry-run mode")
    }
    offset := 0
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
        }
        updates, err := t.GetUpdates(offset)
        if err != nil {
            // Log error but continue polling
            fmt.Fprintf(os.Stderr, "[telegram] poll error: %v\n", err)
            time.Sleep(2 * time.Second)
            continue
        }
        for _, u := range updates {
            offset = u.ID + 1
            if u.Message == nil || u.Message.Text == "" {
                continue
            }
            chatID := u.Message.Chat.ID
            fromName := ""
            if u.Message.From != nil {
                fromName = u.Message.From.FirstName
            }
            text := u.Message.Text

            response := handler(chatID, fromName, text)
            if response != "" {
                if err := t.SendMessage(fmt.Sprintf("%d", chatID), response); err != nil {
                    fmt.Fprintf(os.Stderr, "[telegram] send error: %v\n", err)
                }
            }
        }
    }
}

// TestMessage simulates receiving a message in dry-run mode.
func (t *TelegramAdapter) TestMessage(chatID, fromName, text string) Message {
    // Populate the new normalized inbound fields where possible
    // chatID is already a string in this TestMessage signature
    chatIDStr := chatID
    return Message{
        ID:               "test-telegram-" + fmt.Sprint(time.Now().Unix()),
        AgentID:          t.AgentForChat(chatID),
        Content:          text,
        Channel:          "telegram",
        ChannelType:      "telegram",
        ChannelAccountID: chatIDStr,
        ConversationID:   chatIDStr,
        SenderID:         "",
        Attachments:      nil,
        CreatedAt:        time.Now(),
        RawEventRef:      "test-telegram-" + fmt.Sprint(time.Now().Unix()),
    }
}

// RedactToken returns a redacted version of the token for logging.
func RedactToken(token string) string {
    if len(token) <= 8 {
        return "***"
    }
    return token[:4] + "..." + token[len(token)-4:]
}

// ValidateToken checks if a token looks like a valid Telegram bot token format.
func ValidateToken(token string) error {
    if token == "" {
        return fmt.Errorf("token is empty")
    }
    parts := strings.Split(token, ":")
    if len(parts) != 2 {
        return fmt.Errorf("invalid token format: expected numeric_id:alphanumeric")
    }
    if parts[0] == "" || parts[1] == "" {
        return fmt.Errorf("invalid token format: both parts must be non-empty")
    }
    return nil
}
