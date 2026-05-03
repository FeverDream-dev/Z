package channels

import "time"

// Attachment represents a simple inbound attachment for a message.
type Attachment struct {
    Type string `json:"type"`
    URL  string `json:"url"`
    Name string `json:"name"`
}

// Message is a normalized inbound message from any channel.
type Message struct {
    ID               string       `json:"id"`
    AgentID          string       `json:"agent_id"`
    Content          string       `json:"text"`
    Channel          string       `json:"channel"`
    ChannelType      string       `json:"channel_type"`
    ChannelAccountID string       `json:"channel_account_id"`
    ConversationID   string       `json:"conversation_id"`
    SenderID         string       `json:"sender_id"`
    Attachments      []Attachment `json:"attachments"`
    CreatedAt        time.Time    `json:"timestamp"`
    RawEventRef      string       `json:"raw_event_ref"`
}

// NewMessage creates a new inbound message.
// The new signature accepts all possible inbound fields; callers that don't
// have certain data can pass empty strings or nil attachments.
func NewMessage(
    agentID string,
    content string,
    channel string,
    channelType string,
    channelAccountID string,
    conversationID string,
    senderID string,
    rawEventRef string,
    attachments []Attachment,
) Message {
    return Message{
        ID:               generateMessageID(),
        AgentID:          agentID,
        Content:          content,
        Channel:          channel,
        ChannelType:      channelType,
        ChannelAccountID: channelAccountID,
        ConversationID:   conversationID,
        SenderID:         senderID,
        Attachments:      attachments,
        CreatedAt:        time.Now(),
        RawEventRef:      rawEventRef,
    }
}

func generateMessageID() string {
    return time.Now().Format("20060102-150405-000000000")
}
