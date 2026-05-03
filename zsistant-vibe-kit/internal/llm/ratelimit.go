package llm

import (
	"fmt"
	"time"
)

// RateLimitProvider simulates a provider that returns rate-limit errors.
type RateLimitProvider struct {
	name string
}

// NewRateLimitProvider creates a provider that simulates rate limiting.
func NewRateLimitProvider(name string) *RateLimitProvider {
	return &RateLimitProvider{name: name}
}

// Complete simulates a rate-limit error.
func (r *RateLimitProvider) Complete(prompt string) (string, error) {
	return "", fmt.Errorf("provider %s: rate limit exceeded (429)", r.name)
}

// Health returns degraded status.
func (r *RateLimitProvider) Health() ProviderHealth {
	return ProviderHealth{
		Name:      r.name,
		Status:    "degraded",
		LastError: "simulated rate limit",
		UpdatedAt: time.Now(),
	}
}
