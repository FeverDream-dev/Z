package config

// defaultConfig returns a Config populated with MVP defaults. Used by tests and
// initialization flow to ensure a predictable baseline.
func defaultConfig() Config {
    return Config{
        DataPath:     "~/.zazi",
        LogLevel:     "info",
        ServerPort:   8080,
        LLMProviders: map[string]interface{}{},
        Secrets:      map[string]string{},
    }
}

// NewDefault is a small wrapper returning a pointer to a default Config instance.
func NewDefault() *Config {
    c := defaultConfig()
    return &c
}
