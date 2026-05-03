package config

import (
    "gopkg.in/yaml.v3"
)

// Secrets holds sensitive configuration values that should never be logged
// or printed in plain text.
// Deprecated: use ProviderKeys for per-provider API keys instead.
type Secrets struct {
    ApiKey string `yaml:"api_key" json:"api_key"`
    Token  string `yaml:"token" json:"token"`
}

// Provider holds provider-specific configuration.
type Provider struct {
    Name   string            `yaml:"name" json:"name"`
    Config map[string]string `yaml:"config" json:"config"`
}

// Config represents runtime configuration for Zsistant.
type Config struct {
	DataPath      string            `yaml:"data_path" json:"data_path"`
	LogLevel      string            `yaml:"log_level" json:"log_level"`
	ServerPort    int               `yaml:"server_port" json:"server_port"`
	Providers     []Provider        `yaml:"providers" json:"providers"`
	Secrets       Secrets           `yaml:"secrets" json:"secrets"`
	ProviderKeys  map[string]string `yaml:"provider_keys,omitempty" json:"-"` // provider_name -> api_key
	DevMode       bool              `yaml:"dev_mode,omitempty" json:"dev_mode"`
	Theme         string            `yaml:"theme,omitempty" json:"theme"`
	DefaultModel  string            `yaml:"default_model,omitempty" json:"default_model"`
}

// GetProviderKey returns a provider API key, or empty string if not configured.
func (c *Config) GetProviderKey(provider string) string {
    if c.ProviderKeys == nil {
        return ""
    }
    return c.ProviderKeys[provider]
}

// String implements a redacted YAML representation of the config.
func (c *Config) String() string {
    if c == nil {
        return "<nil>"
    }
    // Shallow copy to avoid mutating the original
    copy := *c
    if copy.Secrets.ApiKey != "" {
        copy.Secrets.ApiKey = "[redacted]"
    }
    if copy.Secrets.Token != "" {
        copy.Secrets.Token = "[redacted]"
    }
    if copy.ProviderKeys != nil {
        copy.ProviderKeys = map[string]string{"__redacted": "true"}
    }
    out, _ := yaml.Marshal(&copy)
    return string(out)
}
