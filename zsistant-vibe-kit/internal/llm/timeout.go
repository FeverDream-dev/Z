package llm

import (
	"fmt"
	"time"
)

// TimeoutProvider simulates a provider that always times out.
type TimeoutProvider struct {
	name string
}

// NewTimeoutProvider creates a provider that simulates timeouts.
func NewTimeoutProvider(name string) *TimeoutProvider {
	return &TimeoutProvider{name: name}
}

// Complete simulates a timeout by sleeping longer than any reasonable timeout.
func (t *TimeoutProvider) Complete(prompt string) (string, error) {
	time.Sleep(5 * time.Second)
	return "", fmt.Errorf("provider %s: request timed out", t.name)
}

// Health returns degraded status.
func (t *TimeoutProvider) Health() ProviderHealth {
	return ProviderHealth{
		Name:      t.name,
		Status:    "degraded",
		LastError: "simulated timeout",
		UpdatedAt: time.Now(),
	}
}
