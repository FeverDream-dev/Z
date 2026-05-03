package approvals

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	StatusPending  = "pending"
	StatusApproved = "approved"
	StatusDenied   = "denied"
	StatusExpired  = "expired"
)

// Request represents a human approval request for an autonomous action.
type Request struct {
	ID              string    `json:"id"`
	AssistantID     string    `json:"assistant_id"`
	TaskID          string    `json:"task_id,omitempty"`
	ActionSummary   string    `json:"action_summary"`
	RiskLevel       string    `json:"risk_level"` // low, medium, high, critical
	RequestedAt     time.Time `json:"requested_at"`
	Status          string    `json:"status"`
	ApprovedBy      string    `json:"approved_by,omitempty"`
	ResolvedAt      *time.Time `json:"resolved_at,omitempty"`
	ExpiresAt       *time.Time `json:"expires_at,omitempty"`
}

// Store persists approval requests per workspace.
type Store struct {
	basePath string
	mu       sync.RWMutex
}

func NewStore(basePath string) *Store {
	return &Store{basePath: basePath}
}

func (s *Store) path() string {
	return filepath.Join(s.basePath, "approvals.jsonl")
}

func (s *Store) Create(req *Request) error {
	if req.ID == "" {
		req.ID = fmt.Sprintf("apr-%d", time.Now().UnixNano())
	}
	if req.RequestedAt.IsZero() {
		req.RequestedAt = time.Now()
	}
	if req.Status == "" {
		req.Status = StatusPending
	}
	if req.ExpiresAt == nil {
		t := req.RequestedAt.Add(24 * time.Hour)
		req.ExpiresAt = &t
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.append(req)
}

func (s *Store) List(assistantID string, status string, limit int) ([]Request, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	entries, err := s.readAll()
	if err != nil {
		return nil, err
	}
	var out []Request
	for _, r := range entries {
		if assistantID != "" && r.AssistantID != assistantID {
			continue
		}
		if status != "" && r.Status != status {
			continue
		}
		if r.ExpiresAt != nil && time.Now().After(*r.ExpiresAt) && r.Status == StatusPending {
			r.Status = StatusExpired
		}
		out = append(out, r)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].RequestedAt.After(out[j].RequestedAt)
	})
	if limit > 0 && len(out) > limit {
		out = out[:limit]
	}
	return out, nil
}

func (s *Store) Get(id string) (*Request, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	entries, err := s.readAll()
	if err != nil {
		return nil, err
	}
	for i := range entries {
		if entries[i].ID == id {
			return &entries[i], nil
		}
	}
	return nil, fmt.Errorf("approval request not found: %s", id)
}

func (s *Store) Resolve(id, by, decision string) (*Request, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	entries, err := s.readAll()
	if err != nil {
		return nil, err
	}
	for i := range entries {
		if entries[i].ID == id {
			entries[i].Status = decision
			entries[i].ApprovedBy = by
			now := time.Now()
			entries[i].ResolvedAt = &now
			if err := s.writeAll(entries); err != nil {
				return nil, err
			}
			return &entries[i], nil
		}
	}
	return nil, fmt.Errorf("approval request not found: %s", id)
}

func (s *Store) PendingCount(assistantID string) int {
	list, _ := s.List(assistantID, StatusPending, 0)
	return len(list)
}

func (s *Store) append(req *Request) error {
	p := s.path()
	os.MkdirAll(filepath.Dir(p), 0o755)
	b, _ := json.Marshal(req)
	f, err := os.OpenFile(p, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(append(b, '\n'))
	return err
}

func (s *Store) readAll() ([]Request, error) {
	p := s.path()
	b, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var out []Request
	for _, line := range strings.Split(string(b), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var r Request
		if err := json.Unmarshal([]byte(line), &r); err == nil {
			out = append(out, r)
		}
	}
	return out, nil
}

func (s *Store) writeAll(list []Request) error {
	p := s.path()
	f, err := os.Create(p)
	if err != nil {
		return err
	}
	defer f.Close()
	for _, r := range list {
		b, _ := json.Marshal(r)
		if _, err := f.Write(append(b, '\n')); err != nil {
			return err
		}
	}
	return nil
}
