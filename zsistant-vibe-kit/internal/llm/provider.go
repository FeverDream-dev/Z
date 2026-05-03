package llm

import "time"

// Provider is the interface for LLM backends.
type Provider interface {
	Complete(prompt string) (string, error)
	Health() ProviderHealth
}

// ProviderHealth describes the current state of a provider.
type ProviderHealth struct {
	Name      string    `json:"name"`
	Status    string    `json:"status"` // "healthy", "degraded", "unhealthy"
	Latency   time.Duration `json:"latency"`
	LastError string    `json:"last_error,omitempty"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TaskType categorizes workloads for routing decisions.
type TaskType string

const (
	TaskCheap  TaskType = "cheap"  // Simple, low-cost tasks
	TaskStrong TaskType = "strong" // Complex reasoning tasks
	TaskCoding TaskType = "coding" // Code generation tasks
)
