package config

// DefaultConfig returns the baseline configuration for Zsistant.
func DefaultConfig() Config {
    return Config{
        DataPath:   "~/.zazi",
        LogLevel:   "info",
        ServerPort: 8080,
        Providers:  []Provider{},
        Secrets:    Secrets{},
    }
}

// MergeDefaults fills in any missing fields of cfg with defaults.
func MergeDefaults(cfg *Config) {
    if cfg == nil {
        return
    }
    def := DefaultConfig()
    if cfg.DataPath == "" {
        cfg.DataPath = def.DataPath
    }
    if cfg.LogLevel == "" {
        cfg.LogLevel = def.LogLevel
    }
    if cfg.ServerPort == 0 {
        cfg.ServerPort = def.ServerPort
    }
    if cfg.Providers == nil {
        cfg.Providers = def.Providers
    }
}
