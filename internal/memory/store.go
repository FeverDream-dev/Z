package memory

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// MemoryEntry is a durable learned fact about the user, project, or assistant.
type MemoryEntry struct {
	ID           string    `json:"id"`
	AssistantID  string    `json:"assistant_id"` // "" = global
	Scope        string    `json:"scope"`        // "global" or "assistant"
	Category     string    `json:"category"`     // e.g. "preference", "fact", "project", "learned"
	Content      string    `json:"content"`
	Source       string    `json:"source,omitempty"` // how it was learned
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Approved     bool      `json:"approved"`         // user-approved?
}

// Store manages memory persistence per workspace.
type Store struct {
	basePath string
}

// NewStore creates a memory store at basePath.
func NewStore(basePath string) *Store {
	return &Store{basePath: basePath}
}

func (s *Store) memoryPath() string {
	return filepath.Join(s.basePath, "memory.jsonl")
}

// Add appends a new memory entry. If assistantID is empty, scope is global.
func (s *Store) Add(assistantID, category, content, source string, approved bool) (*MemoryEntry, error) {
	scope := "assistant"
	if assistantID == "" {
		scope = "global"
	}
	now := time.Now()
	m := MemoryEntry{
		ID:          generateID(),
		AssistantID: assistantID,
		Scope:       scope,
		Category:    category,
		Content:     content,
		Source:      source,
		CreatedAt:   now,
		UpdatedAt:   now,
		Approved:    approved,
	}
	if err := s.appendEntry(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

// List returns memories, optionally filtered by assistantID ("" = global).
func (s *Store) List(assistantID string) ([]MemoryEntry, error) {
	entries, err := s.readAll()
	if err != nil {
		return nil, err
	}
	var out []MemoryEntry
	for _, e := range entries {
		if assistantID == "" && e.Scope == "global" {
			out = append(out, e)
		} else if assistantID != "" && (e.Scope == "global" || e.AssistantID == assistantID) {
			out = append(out, e)
		}
	}
	// newest first
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// Get retrieves a memory by ID.
func (s *Store) Get(id string) (*MemoryEntry, error) {
	entries, err := s.readAll()
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		if e.ID == id {
			return &e, nil
		}
	}
	return nil, fmt.Errorf("memory not found: %s", id)
}

// Update modifies an existing memory entry.
func (s *Store) Update(id string, fn func(m *MemoryEntry)) error {
	entries, err := s.readAll()
	if err != nil {
		return err
	}
	found := false
	for i := range entries {
		if entries[i].ID == id {
			fn(&entries[i])
			entries[i].UpdatedAt = time.Now()
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("memory not found: %s", id)
	}
	return s.writeAll(entries)
}

// Delete removes a memory entry by ID.
func (s *Store) Delete(id string) error {
	entries, err := s.readAll()
	if err != nil {
		return err
	}
	var filtered []MemoryEntry
	for _, e := range entries {
		if e.ID != id {
			filtered = append(filtered, e)
		}
	}
	if len(filtered) == len(entries) {
		return fmt.Errorf("memory not found: %s", id)
	}
	return s.writeAll(filtered)
}

func (s *Store) appendEntry(m *MemoryEntry) error {
	path := s.memoryPath()
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	_, err = f.Write(append(b, '\n'))
	return err
}

func (s *Store) readAll() ([]MemoryEntry, error) {
	path := s.memoryPath()
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var out []MemoryEntry
	for _, line := range strings.Split(string(b), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var m MemoryEntry
		if err := json.Unmarshal([]byte(line), &m); err != nil {
			continue
		}
		out = append(out, m)
	}
	return out, nil
}

func (s *Store) writeAll(entries []MemoryEntry) error {
	path := s.memoryPath()
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	for _, m := range entries {
		b, err := json.Marshal(m)
		if err != nil {
			return err
		}
		if _, err := f.Write(append(b, '\n')); err != nil {
			return err
		}
	}
	return nil
}

func generateID() string {
	return fmt.Sprintf("mem-%d", time.Now().UnixNano())
}
