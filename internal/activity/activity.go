package activity

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ActivityEvent is a chronological record of something that happened in the assistant platform.
type ActivityEvent struct {
	ID          string            `json:"id"`
	AssistantID string            `json:"assistant_id"`   // "" = global / workspace level
	EventType   string            `json:"event_type"`     // e.g. "assistant.created", "tool.called", "job.scheduled", "job.failed", "channel.connected", "approval.requested", "memory.saved", "message.sent"
	Message     string            `json:"message"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	Severity    string            `json:"severity"`       // info, warning, error, critical
	CreatedAt   time.Time         `json:"created_at"`
}

// Log writes an event to the global activity log.
func Log(workspaceRoot string, evt ActivityEvent) error {
	if evt.ID == "" {
		evt.ID = generateID()
	}
	if evt.CreatedAt.IsZero() {
		evt.CreatedAt = time.Now()
	}
	if evt.Severity == "" {
		evt.Severity = "info"
	}
	path := filepath.Join(workspaceRoot, "activity.jsonl")
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("opening activity log: %w", err)
	}
	defer f.Close()
	b, err := json.Marshal(evt)
	if err != nil {
		return fmt.Errorf("marshaling activity event: %w", err)
	}
	if _, err := f.Write(append(b, '\n')); err != nil {
		return fmt.Errorf("writing activity event: %w", err)
	}
	return nil
}

// ReadGlobal returns global activity events (workspace level).
func ReadGlobal(workspaceRoot string, limit int) ([]ActivityEvent, error) {
	path := filepath.Join(workspaceRoot, "activity.jsonl")
	return readEvents(path, "", limit)
}

// ReadAssistant returns activity events for a specific assistant.
func ReadAssistant(workspaceRoot, assistantID string, limit int) ([]ActivityEvent, error) {
	// First try assistant-specific log
	path := filepath.Join(workspaceRoot, "assistants", assistantID, "audit.jsonl")
	evts, err := readEvents(path, assistantID, limit)
	if err == nil && len(evts) > 0 {
		// Also merge global events that mention this assistant
		globalPath := filepath.Join(workspaceRoot, "activity.jsonl")
		globalEvts, _ := readEvents(globalPath, assistantID, limit)
		// Deduplicate by ID
		seen := make(map[string]bool)
		for _, e := range evts {
			seen[e.ID] = true
		}
		for _, e := range globalEvts {
			if !seen[e.ID] {
				evts = append(evts, e)
				seen[e.ID] = true
			}
		}
		// reverse to get newest first
		for i, j := 0, len(evts)-1; i < j; i, j = i+1, j-1 {
			evts[i], evts[j] = evts[j], evts[i]
		}
		if len(evts) > limit && limit > 0 {
			evts = evts[:limit]
		}
		return evts, nil
	}
	// Fallback to global log filtering by assistantID
	globalPath := filepath.Join(workspaceRoot, "activity.jsonl")
	return readEvents(globalPath, assistantID, limit)
}

func readEvents(path, filterAssistantID string, limit int) ([]ActivityEvent, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []ActivityEvent{}, nil
		}
		return nil, err
	}
	var out []ActivityEvent
	for _, line := range strings.Split(string(b), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var e ActivityEvent
		if err := json.Unmarshal([]byte(line), &e); err != nil {
			continue
		}
		if filterAssistantID != "" && e.AssistantID != "" && e.AssistantID != filterAssistantID {
			continue
		}
		out = append(out, e)
	}
	// reverse to get newest first (JSONL is append-only)
	for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
		out[i], out[j] = out[j], out[i]
	}
	if limit > 0 && len(out) > limit {
		out = out[:limit]
	}
	return out, nil
}

func generateID() string {
	return fmt.Sprintf("evt-%d", time.Now().UnixNano())
}
