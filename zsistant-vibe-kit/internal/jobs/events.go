package jobs

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// JobEvent is an append-only audit/status record.
type JobEvent struct {
	ID        string            `json:"id"`
	JobID     string            `json:"job_id"`
	AgentID   string            `json:"agent_id"`
	EventType string            `json:"event_type"`
	Message   string            `json:"message"`
	Metadata  map[string]string `json:"metadata,omitempty"`
	CreatedAt time.Time         `json:"created_at"`
}

// AppendEvent writes a single event as a JSON line to the agent's audit log.
func AppendEvent(agentWorkspace string, evt JobEvent) error {
	path := filepath.Join(agentWorkspace, "audit.jsonl")
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("opening audit log: %w", err)
	}
	defer f.Close()

	b, err := json.Marshal(evt)
	if err != nil {
		return fmt.Errorf("marshaling event: %w", err)
	}
	if _, err := f.Write(append(b, '\n')); err != nil {
		return fmt.Errorf("writing event: %w", err)
	}
	return nil
}
