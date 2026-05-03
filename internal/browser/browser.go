package browser

import (
	"fmt"
	"time"
)

// CapabilityLevel describes the maturity of browser/MCP integration.
type CapabilityLevel string

const (
	LevelNotAvailable  CapabilityLevel = "not_available"
	LevelNeedsSetup    CapabilityLevel = "needs_setup"
	LevelConnected     CapabilityLevel = "connected"
	LevelUsableManual  CapabilityLevel = "usable_manual"
	LevelAutomatable   CapabilityLevel = "automatable"
	LevelObservable    CapabilityLevel = "observable"
)

// MCPServer describes a connected MCP server.
type MCPServer struct {
	Name        string    `json:"name"`
	URL         string    `json:"url"`
	Status      string    `json:"status"`      // connected, disconnected, error
	LastCallAt  time.Time `json:"last_call_at,omitempty"`
	LastError   string    `json:"last_error,omitempty"`
	ToolsCount  int       `json:"tools_count"`
	CreatedAt   time.Time `json:"created_at"`
}

// BrowserSession represents a live browser session for an assistant.
type BrowserSession struct {
	ID          string    `json:"id"`
	AssistantID string    `json:"assistant_id"`
	CurrentURL  string    `json:"current_url,omitempty"`
	LastAction  string    `json:"last_action,omitempty"`
	ScreenshotPath string `json:"screenshot_path,omitempty"`
	ConsoleErrors  []string `json:"console_errors,omitempty"`
	Status      string    `json:"status"`      // active, idle, error
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Config holds platform-level browser/MCP configuration.
// This is honest: if MCP is not connected, Level is "needs_setup" or "not_available".
type Config struct {
	Level           CapabilityLevel `json:"level"`
	MCPServers      []MCPServer     `json:"mcp_servers"`
	ActiveSessions  []BrowserSession `json:"active_sessions"`
	DefaultSetupURL string          `json:"default_setup_url,omitempty"`
	SetupMessage    string          `json:"setup_message,omitempty"`
}

// DefaultConfig returns an honest default for browser/MCP.
func DefaultConfig() Config {
	return Config{
		Level:        LevelNotAvailable,
		MCPServers:   []MCPServer{},
		ActiveSessions: []BrowserSession{},
		SetupMessage:  "Browser MCP is not connected. Configure a Chrome or Playwright MCP server to enable browser actions.",
	}
}

// IsAvailable returns true if browser actions can be performed.
func (c *Config) IsAvailable() bool {
	return c.Level == LevelConnected || c.Level == LevelUsableManual || c.Level == LevelAutomatable || c.Level == LevelObservable
}

// HonestStatus returns a human-readable status string.
func (c *Config) HonestStatus() string {
	switch c.Level {
	case LevelNotAvailable:
		return "Browser MCP is not available in this build."
	case LevelNeedsSetup:
		return "Browser MCP needs setup. " + c.SetupMessage
	case LevelConnected, LevelUsableManual, LevelAutomatable, LevelObservable:
		if len(c.MCPServers) == 0 {
			return "Browser MCP is configured but no servers are connected."
		}
		connected := 0
		for _, s := range c.MCPServers {
			if s.Status == "connected" {
				connected++
			}
		}
		return fmt.Sprintf("Browser MCP: %d/%d servers connected.", connected, len(c.MCPServers))
	default:
		return "Browser MCP status unknown."
	}
}
