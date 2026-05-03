package assistant

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Status values for an Assistant lifecycle.
const (
	StatusCreated    = "created"
	StatusActive     = "active"
	StatusPaused     = "paused"
	StatusNeedsSetup = "needs_setup"
	StatusDegraded   = "degraded"
	StatusError      = "error"
	StatusWaiting    = "waiting_for_approval"
)

// Persona defines how the assistant behaves.
type Persona struct {
	Tone           string `json:"tone"`
	Style          string `json:"style"`               // concise, detailed, bullet-first...
	RoleDescription string `json:"role_description"`    // e.g. "concise senior engineer"
	Boundaries     string `json:"boundaries"`          // what it should not do
	DecisionPolicy string `json:"decision_policy"`     // when to ask for approval
	FormattingPref string `json:"formatting_pref"`     // markdown, plain, summary-first...
}

// MemoryPolicy defines what the assistant can remember.
type MemoryPolicy struct {
	Enabled            bool     `json:"enabled"`
	AutoSave           bool     `json:"auto_save"`           // auto-save memories or require approval
	Scope              string   `json:"scope"`               // "global", "assistant-only", "both"
	RetentionDays      int      `json:"retention_days"`      // 0 = forever
	SensitiveBlocklist []string `json:"sensitive_blocklist"` // topics to never remember
}

// ChannelConfig defines a communication surface attached to an assistant.
type ChannelConfig struct {
	ChannelType    string            `json:"channel_type"`     // web_ui, telegram, whatsapp, discord, slack, cli, email
	ChannelID      string            `json:"channel_id"`       // e.g. bot token, webhook URL, channel name
	Status         string            `json:"status"`           // connected, needs_setup, paused, error
	AllowedUsers   []string          `json:"allowed_users,omitempty"`
	RequireMention bool              `json:"require_mention"`  // in groups
	Settings       map[string]string `json:"settings,omitempty"`
	LastError      string            `json:"last_error,omitempty"`
	CreatedAt      time.Time         `json:"created_at"`
}

// Assistant is the central object in Zsistant.
type Assistant struct {
	ID                string            `json:"id"`
	Name              string            `json:"name"`
	Slug              string            `json:"slug"`
	Description       string            `json:"description"`       // short purpose
	Purpose           string            `json:"purpose"`           // long-form responsibility statement
	AvatarPath        string            `json:"avatar_path,omitempty"`
	Owner             string            `json:"owner,omitempty"`

	// Identity / Persona
	Persona          Persona           `json:"persona"`
	Responsibilities []string          `json:"responsibilities,omitempty"`

	// Provider / model policy
	DefaultModel    string            `json:"default_model,omitempty"`
	FallbackModel   string            `json:"fallback_model,omitempty"`
	ProviderName    string            `json:"provider_name,omitempty"`

	// Permissions and safety
	ToolPermissions  []string          `json:"tool_permissions,omitempty"`
	ApprovalRequired []string          `json:"approval_required,omitempty"`

	// Memory and knowledge (simple string arrays for MVP compatibility with server)
	Memory    []string          `json:"memory"`
	Knowledge []string          `json:"knowledge"`
	MemoryPolicy MemoryPolicy    `json:"memory_policy"`
	KnowledgeIDs []string        `json:"knowledge_ids,omitempty"`

	// Channels
	Channels []ChannelConfig    `json:"channels,omitempty"`
	ChannelNames []string        `json:"channel_names,omitempty"`   // simple string list for backwards compat

	// Jobs
	JobsEnabled bool             `json:"jobs_enabled"`

	// Lifecycle
	Status           string        `json:"status"`
	CreatedAt        time.Time     `json:"created_at"`
	UpdatedAt        time.Time     `json:"updated_at"`
	LastActivityAt   *time.Time    `json:"last_activity_at,omitempty"`
}

// IsHealthy returns true if the assistant is in an operational state.
func (a *Assistant) IsHealthy() bool {
	return a.Status == StatusActive || a.Status == StatusCreated
}

// Job represents a background task associated with an assistant.
type Job struct {
	ID          string    `json:"id"`
	AssistantID string    `json:"assistant_id"`
	Type        string    `json:"type"`
	Message     string    `json:"message"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Status      string    `json:"status"`
	Result      string    `json:"result"`
}

// Registry holds a persistent collection of assistants.
type Registry struct {
	BasePath string
	mu       sync.RWMutex
}

// NewRegistry creates a new Registry bound to a storage path.
func NewRegistry(basePath string) *Registry {
	return &Registry{BasePath: basePath}
}

// listUnsafe reads assistants without lock — caller must hold r.mu.
func (r *Registry) listUnsafe() ([]Assistant, error) {
	path := filepath.Join(r.BasePath, "assistants.json")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return []Assistant{}, nil
	}
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var out []Assistant
	if err := json.Unmarshal(b, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// List fetches all assistants from storage.
func (r *Registry) List() ([]Assistant, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.listUnsafe()
}

// Get returns a single assistant by id.
func (r *Registry) Get(id string) (*Assistant, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	as, err := r.listUnsafe()
	if err != nil {
		return nil, err
	}
	for _, a := range as {
		if a.ID == id {
			aa := a
			return &aa, nil
		}
	}
	return nil, errors.New("not found")
}

// Create adds a new assistant.
func (r *Registry) Create(id, name string) (*Assistant, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	as, err := r.listUnsafe()
	if err != nil {
		return nil, err
	}
	for _, a := range as {
		if a.ID == id {
			return nil, errors.New("already exists")
		}
	}
	now := time.Now()
	a := Assistant{
		ID: id, Name: name,
		Memory: []string{}, Knowledge: []string{}, Channels: []ChannelConfig{},
		ChannelNames: []string{}, Status: StatusCreated,
		CreatedAt: now, UpdatedAt: now,
		MemoryPolicy: MemoryPolicy{Scope: "assistant-only"},
	}
	as = append(as, a)
	if err := r.saveAll(as); err != nil {
		return nil, err
	}
	return &a, nil
}

// Update replaces an existing assistant.
func (r *Registry) Update(a *Assistant) error {
	if a == nil {
		return errors.New("nil assistant")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	list, err := r.listUnsafe()
	if err != nil {
		return err
	}
	found := false
	for i := range list {
		if list[i].ID == a.ID {
			a.UpdatedAt = time.Now()
			list[i] = *a
			found = true
			break
		}
	}
	if !found {
		return errors.New("not found")
	}
	return r.saveAll(list)
}

// Delete removes an assistant.
func (r *Registry) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	list, err := r.listUnsafe()
	if err != nil {
		return err
	}
	idx := -1
	for i := range list {
		if list[i].ID == id {
			idx = i
			break
		}
	}
	if idx < 0 {
		return errors.New("not found")
	}
	list = append(list[:idx], list[idx+1:]...)
	return r.saveAll(list)
}

// saveAll writes the current list to disk.
func (r *Registry) saveAll(list []Assistant) error {
	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}
	path := filepath.Join(r.BasePath, "assistants.json")
	if err := os.MkdirAll(r.BasePath, 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// JobStore persists jobs for assistants.
type JobStore struct {
	BasePath string
	mu       sync.RWMutex
}

func NewJobStore(basePath string) *JobStore {
	return &JobStore{BasePath: basePath}
}

// List reads all jobs for an assistant.
func (js *JobStore) List(assistantID string) ([]Job, error) {
	js.mu.RLock()
	defer js.mu.RUnlock()
	path := filepath.Join(js.BasePath, "jobs_"+assistantID+".jsonl")
	var out []Job
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []Job{}, nil
		}
		return nil, err
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	for _, line := range splitLines(string(data)) {
		if line == "" {
			continue
		}
		var j Job
		if err := json.Unmarshal([]byte(line), &j); err == nil {
			out = append(out, j)
		}
	}
	return out, nil
}

func (js *JobStore) Append(assistantID string, job Job) error {
	js.mu.Lock()
	defer js.mu.Unlock()
	path := filepath.Join(js.BasePath, "jobs_"+assistantID+".jsonl")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	b, _ := json.Marshal(job)
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.Write(b); err != nil {
		return err
	}
	f.Write([]byte("\n"))
	return nil
}

// Update updates an existing job in the jsonl store by appending updated version.
func (js *JobStore) Update(assistantID string, job Job) error {
	job.UpdatedAt = time.Now()
	return js.Append(assistantID, job)
}

func splitLines(s string) []string {
	var lines []string
	cur := make([]byte, 0, 256)
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '\n' {
			lines = append(lines, string(cur))
			cur = cur[:0]
			continue
		}
		cur = append(cur, c)
	}
	if len(cur) > 0 {
		lines = append(lines, string(cur))
	}
	return lines
}
