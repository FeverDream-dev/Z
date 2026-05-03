package jobs

import (
	"fmt"
	"time"
)

// Job represents a durable unit of work for an agent.
type Job struct {
	ID            string     `json:"id"`
	AgentID       string     `json:"agent_id"`
	SourceChannel string     `json:"source_channel"`
	Objective     string     `json:"objective"`
	Status        string     `json:"status"`
	CurrentStep   string     `json:"current_step"`
	Result        string     `json:"result"`
	RetryCount    int        `json:"retry_count"`
	MaxRetries    int        `json:"max_retries"`
	LastError     string     `json:"last_error,omitempty"`
	PausedReason  string     `json:"paused_reason,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	CompletedAt   *time.Time `json:"completed_at,omitempty"`
}

// NewJob creates a new job with a generated ID and initial status.
func NewJob(agentID, sourceChannel, objective string) *Job {
	now := time.Now()
	return &Job{
		ID:            generateID(),
		AgentID:       agentID,
		SourceChannel: sourceChannel,
		Objective:     objective,
		Status:        "queued",
		CurrentStep:   "created",
		RetryCount:    0,
		MaxRetries:    3,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

var idCounter int64

func generateID() string {
	// Timestamp + counter for uniqueness even on rapid calls.
	idCounter++
	return fmt.Sprintf("%s-%06d", time.Now().Format("20060102-150405"), idCounter)
}

// IsTerminal returns true if the job is in a terminal state.
func (j *Job) IsTerminal() bool {
	return j.Status == "completed" || j.Status == "failed" || j.Status == "cancelled"
}

// CanRetry returns true if the job can be retried.
// Failed, paused, and retrying jobs are eligible for retry.
func (j *Job) CanRetry() bool {
	if j.RetryCount >= j.MaxRetries {
		return false
	}
	return j.Status == "failed" || j.Status == "paused" || j.Status == "retrying"
}
