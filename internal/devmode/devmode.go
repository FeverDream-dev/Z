package devmode

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Settings controls whether Developer Mode is active.
type Settings struct {
	Enabled bool `json:"enabled"`
}

// Trace is a developer-level event recording tool calls, provider routing, etc.
type Trace struct {
	ID            string                 `json:"id"`
	AssistantID   string                 `json:"assistant_id,omitempty"`
	TraceType     string                 `json:"trace_type"`       // "provider_call", "tool_call", "streaming_event", "mcp_call", "browser_action", "fallback", "error"
	Message       string                 `json:"message"`
	Payload       map[string]interface{} `json:"payload,omitempty"`
	DurationMs    int64                  `json:"duration_ms,omitempty"`
	ProviderName  string                 `json:"provider_name,omitempty"`
	ModelID       string                 `json:"model_id,omitempty"`
	TokensIn      int                    `json:"tokens_in,omitempty"`
	TokensOut     int                    `json:"tokens_out,omitempty"`
	CreatedAt     time.Time              `json:"created_at"`
}

// TraceStore persists developer traces.
type TraceStore struct {
	workspaceRoot string
}

// NewTraceStore creates a trace store for the given workspace.
func NewTraceStore(workspaceRoot string) *TraceStore {
	return &TraceStore{workspaceRoot: workspaceRoot}
}

// Append writes a trace event.
func (ts *TraceStore) Append(t Trace) error {
	if t.ID == "" {
		t.ID = generateID()
	}
	if t.CreatedAt.IsZero() {
		t.CreatedAt = time.Now()
	}
	path := filepath.Join(ts.workspaceRoot, "traces.jsonl")
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("opening trace log: %w", err)
	}
	defer f.Close()
	b, err := json.Marshal(t)
	if err != nil {
		return fmt.Errorf("marshaling trace: %w", err)
	}
	if _, err := f.Write(append(b, '\n')); err != nil {
		return fmt.Errorf("writing trace: %w", err)
	}
	return nil
}

// Read returns the latest traces, optionally filtered by assistantID or traceType.
func (ts *TraceStore) Read(limit int, assistantID, traceType string) ([]Trace, error) {
	path := filepath.Join(ts.workspaceRoot, "traces.jsonl")
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []Trace{}, nil
		}
		return nil, err
	}
	var out []Trace
	for _, line := range strings.Split(string(b), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var t Trace
		if err := json.Unmarshal([]byte(line), &t); err != nil {
			continue
		}
		if assistantID != "" && t.AssistantID != "" && t.AssistantID != assistantID {
			continue
		}
		if traceType != "" && t.TraceType != traceType {
			continue
		}
		out = append(out, t)
	}
	// reverse: newest first
	for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
		out[i], out[j] = out[j], out[i]
	}
	if limit > 0 && len(out) > limit {
		out = out[:limit]
	}
	return out, nil
}

// Diagnostics holds runtime diagnostic data exposed in Developer Mode.
type Diagnostics struct {
	WorkspaceRoot     string                 `json:"workspace_root"`
	ConfigPath        string                 `json:"config_path"`
	GoVersion         string                 `json:"go_version"`
	BuildVersion      string                 `json:"build_version"`
	ActiveAssistants  int                    `json:"active_assistants"`
	ConnectedChannels int                    `json:"connected_channels"`
	RegisteredTools   int                    `json:"registered_tools"`
	PendingJobs       int                    `json:"pending_jobs"`
	FailedJobsLast24h int                    `json:"failed_jobs_last_24h"`
	ProviderHealth    []map[string]string    `json:"provider_health"`
	UptimeSeconds     int64                  `json:"uptime_seconds"`
	MemoryUsageMB     int64                  `json:"memory_usage_mb,omitempty"`
}

func generateID() string {
	return fmt.Sprintf("trc-%d", time.Now().UnixNano())
}
