package config

import (
    "fmt"
)

// Config holds the runtime configuration for Zsistant.
//
// - DataPath is the base directory for agent data, caches, logs, etc.
// - LogLevel controls verbose output for the local daemon.
// - ServerPort is the HTTP server port the local API should listen on.
// - LLMProviders is a placeholder for future provider-specific configs.
// - Secrets holds sensitive data; it is intentionally excluded from String()
//   and any normal logging.
type Config struct {
    DataPath     string                 `yaml:"data_path" json:"data_path"`
    LogLevel     string                 `yaml:"log_level" json:"log_level"`
    ServerPort   int                    `yaml:"server_port" json:"server_port"`
    LLMProviders map[string]interface{} `yaml:"llm_providers,omitempty" json:"llm_providers,omitempty"`
    Secrets      map[string]string      `yaml:"secrets,omitempty" json:"-"`
}

// String implements fmt.Stringer. Secrets are redacted to avoid leaking values
// into logs or user interfaces.
func (c *Config) String() string {
    // Redact secrets, but show other fields for quick verification.
    return fmt.Sprintf("Config{DataPath:%q, LogLevel:%q, ServerPort:%d, LLMProviders:%v, Secrets:[redacted]}",
        c.DataPath, c.LogLevel, c.ServerPort, c.LLMProviders)
}

// (NewDefault defined in defaults.go)
