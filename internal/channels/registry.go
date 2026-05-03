package channels

// ChannelRegistry holds ported adapters for the supported channels.
type ChannelRegistry struct {
    Telegram *TelegramAdapter
    Discord  *DiscordAdapter
    WhatsApp *WhatsAppAdapter
}

// NewChannelRegistry returns a simple registry with default (dry-run) adapters.
// This provides a convenient entry point for wiring adapters in the root project.
func NewChannelRegistry() *ChannelRegistry {
    return &ChannelRegistry{
        Telegram: NewTelegramAdapter(""),
        Discord:  NewDiscordAdapter(""),
        WhatsApp: NewWhatsAppAdapter("", "", ""),
    }
}
