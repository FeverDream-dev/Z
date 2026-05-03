package llm

import (
    "time"
)

// MockProvider is a local echo provider for testing without external API keys.
type MockProvider struct {
    name   string
    delay  time.Duration
    health ProviderHealth
}

// NewMockProvider creates a new mock provider.
func NewMockProvider() *MockProvider {
    return &MockProvider{
        name:  "mock",
        delay: 0,
        health: ProviderHealth{
            Name:      "mock",
            Status:    "healthy",
            UpdatedAt: time.Now(),
        },
    }
}

// NewNamedMockProvider creates a mock with a custom name and optional delay.
func NewNamedMockProvider(name string, delay time.Duration) *MockProvider {
    return &MockProvider{
        name:  name,
        delay: delay,
        health: ProviderHealth{
            Name:      name,
            Status:    "healthy",
            UpdatedAt: time.Now(),
        },
    }
}

// Complete returns an echo response for the given prompt.
func (m *MockProvider) Complete(prompt string) (string, error) {
    if m.delay > 0 {
        time.Sleep(m.delay)
    }
    return "Echo: " + prompt, nil
}

// Health returns the provider's health status.
func (m *MockProvider) Health() ProviderHealth {
    return m.health
}
