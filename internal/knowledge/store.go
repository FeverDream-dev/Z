package knowledge

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// ProcessingStatus describes where the knowledge source is in its lifecycle.
type ProcessingStatus string

const (
	StatusPending    ProcessingStatus = "pending"
	StatusIndexing   ProcessingStatus = "indexing"
	StatusReady      ProcessingStatus = "ready"
	StatusError      ProcessingStatus = "error"
	StatusProcessing ProcessingStatus = "processing"
)

// KnowledgeSource is a document, file, webpage, or note that an assistant can reference.
type KnowledgeSource struct {
	ID            string           `json:"id"`
	Name          string           `json:"name"`
	Type          string           `json:"type"`          // pdf, markdown, webpage, note, docx, txt
	Path          string           `json:"path"`          // local path or URL
	SizeBytes     int64            `json:"size_bytes"`
	Description   string           `json:"description,omitempty"`
	Tags          []string         `json:"tags,omitempty"`
	Status        ProcessingStatus `json:"status"`
	StatusMessage string           `json:"status_message,omitempty"`
	AssistantIDs  []string         `json:"assistant_ids,omitempty"` // attached to which assistants
	ProjectID     string           `json:"project_id,omitempty"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
}

// Store manages knowledge sources.
type Store struct {
	basePath string
}

// NewStore creates a knowledge store.
func NewStore(basePath string) *Store {
	return &Store{basePath: basePath}
}

func (s *Store) knowledgePath() string {
	return filepath.Join(s.basePath, "knowledge.jsonl")
}

// Create registers a new knowledge source.
func (s *Store) Create(src *KnowledgeSource) (*KnowledgeSource, error) {
	if src.ID == "" {
		src.ID = generateID()
	}
	if src.Status == "" {
		src.Status = StatusPending
	}
	now := time.Now()
	src.CreatedAt = now
	src.UpdatedAt = now
	if err := s.appendEntry(src); err != nil {
		return nil, err
	}
	return src, nil
}

// List returns all knowledge sources, optionally filtered by assistantID.
func (s *Store) List(assistantID string) ([]KnowledgeSource, error) {
	entries, err := s.readAll()
	if err != nil {
		return nil, err
	}
	if assistantID == "" {
		return entries, nil
	}
	var out []KnowledgeSource
	for _, e := range entries {
		for _, aid := range e.AssistantIDs {
			if aid == assistantID {
				out = append(out, e)
				break
			}
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// Get retrieves a knowledge source by ID.
func (s *Store) Get(id string) (*KnowledgeSource, error) {
	entries, err := s.readAll()
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		if e.ID == id {
			return &e, nil
		}
	}
	return nil, fmt.Errorf("knowledge source not found: %s", id)
}

// Update modifies a knowledge source.
func (s *Store) Update(id string, fn func(src *KnowledgeSource)) error {
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
		return fmt.Errorf("knowledge source not found: %s", id)
	}
	return s.writeAll(entries)
}

// Delete removes a knowledge source.
func (s *Store) Delete(id string) error {
	entries, err := s.readAll()
	if err != nil {
		return err
	}
	var filtered []KnowledgeSource
	for _, e := range entries {
		if e.ID != id {
			filtered = append(filtered, e)
		}
	}
	if len(filtered) == len(entries) {
		return fmt.Errorf("knowledge source not found: %s", id)
	}
	return s.writeAll(filtered)
}

// AttachToAssistant links a knowledge source to an assistant.
func (s *Store) AttachToAssistant(id, assistantID string) error {
	return s.Update(id, func(src *KnowledgeSource) {
		for _, a := range src.AssistantIDs {
			if a == assistantID {
				return
			}
		}
		src.AssistantIDs = append(src.AssistantIDs, assistantID)
	})
}

// DetachFromAssistant removes an assistant link from a knowledge source.
func (s *Store) DetachFromAssistant(id, assistantID string) error {
	return s.Update(id, func(src *KnowledgeSource) {
		var filtered []string
		for _, a := range src.AssistantIDs {
			if a != assistantID {
				filtered = append(filtered, a)
			}
		}
		src.AssistantIDs = filtered
	})
}

func (s *Store) appendEntry(src *KnowledgeSource) error {
	path := s.knowledgePath()
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	b, err := json.Marshal(src)
	if err != nil {
		return err
	}
	_, err = f.Write(append(b, '\n'))
	return err
}

func (s *Store) readAll() ([]KnowledgeSource, error) {
	path := s.knowledgePath()
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var out []KnowledgeSource
	for _, line := range strings.Split(string(b), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var src KnowledgeSource
		if err := json.Unmarshal([]byte(line), &src); err != nil {
			continue
		}
		out = append(out, src)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

func (s *Store) writeAll(entries []KnowledgeSource) error {
	path := s.knowledgePath()
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	for _, src := range entries {
		b, err := json.Marshal(src)
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
	return fmt.Sprintf("know-%d", time.Now().UnixNano())
}
