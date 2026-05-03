package bus

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Permission levels for inter-agent access.
type Permission string

const (
	PermSummary   Permission = "summary"    // status/summary requests
	PermFileRead  Permission = "file_read"  // read files in scope
	PermFileWrite Permission = "file_write" // write files in scope
	PermExec      Permission = "exec"       // execute commands
)

// ACLRule defines what a peer agent is allowed to do.
type ACLRule struct {
	PeerID    string       `json:"peer_id"`
	AgentID   string       `json:"agent_id"`
	Perms     []Permission `json:"perms"`
	Scope     string       `json:"scope,omitempty"` // path scope for file ops
	CreatedAt time.Time    `json:"created_at"`
}

// ACLStore manages access rules per agent.
type ACLStore struct {
	mu     sync.RWMutex
	base   string
	rules  map[string][]ACLRule // agent_id -> rules
}

// NewACLStore creates an ACL store at the given base path.
func NewACLStore(base string) *ACLStore {
	return &ACLStore{
		base:  base,
		rules: make(map[string][]ACLRule),
	}
}

// loadPath returns the file path for an agent's ACL rules.
func (s *ACLStore) loadPath(agentID string) string {
	return filepath.Join(s.base, "agents", agentID, "acl.json")
}

// Load reads ACL rules for an agent from disk.
func (s *ACLStore) Load(agentID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	path := s.loadPath(agentID)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			s.rules[agentID] = nil
			return nil
		}
		return fmt.Errorf("reading acl: %w", err)
	}
	var rules []ACLRule
	if err := json.Unmarshal(data, &rules); err != nil {
		return fmt.Errorf("parsing acl: %w", err)
	}
	s.rules[agentID] = rules
	return nil
}

// Save persists ACL rules for an agent.
func (s *ACLStore) Save(agentID string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	path := s.loadPath(agentID)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("creating acl dir: %w", err)
	}
	data, err := json.MarshalIndent(s.rules[agentID], "", "  ")
	if err != nil {
		return fmt.Errorf("encoding acl: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

// Allow grants a permission to a peer.
func (s *ACLStore) Allow(agentID, peerID string, perms []Permission, scope string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	rules := s.rules[agentID]
	var updated []ACLRule
	found := false
	for _, r := range rules {
		if r.PeerID == peerID {
			// Merge permissions
			permSet := make(map[Permission]bool)
			for _, p := range r.Perms {
				permSet[p] = true
			}
			for _, p := range perms {
				permSet[p] = true
			}
			var newPerms []Permission
			for p := range permSet {
				newPerms = append(newPerms, p)
			}
			r.Perms = newPerms
			if scope != "" {
				r.Scope = scope
			}
			found = true
		}
		updated = append(updated, r)
	}
	if !found {
		updated = append(updated, ACLRule{
			PeerID:    peerID,
			AgentID:   agentID,
			Perms:     perms,
			Scope:     scope,
			CreatedAt: time.Now(),
		})
	}
	s.rules[agentID] = updated
	return nil
}

// Revoke removes permissions from a peer.
func (s *ACLStore) Revoke(agentID, peerID string, perms []Permission) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	rules := s.rules[agentID]
	var updated []ACLRule
	for _, r := range rules {
		if r.PeerID == peerID {
			if len(perms) == 0 {
				continue // revoke all
			}
			permSet := make(map[Permission]bool)
			for _, p := range r.Perms {
				permSet[p] = true
			}
			for _, p := range perms {
				delete(permSet, p)
			}
			if len(permSet) == 0 {
				continue // no perms left
			}
			var newPerms []Permission
			for p := range permSet {
				newPerms = append(newPerms, p)
			}
			r.Perms = newPerms
		}
		updated = append(updated, r)
	}
	s.rules[agentID] = updated
	return nil
}

// List returns all rules for an agent.
func (s *ACLStore) List(agentID string) []ACLRule {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.rules[agentID]
}

// Check returns true if peer has the requested permission.
func (s *ACLStore) Check(agentID, peerID string, perm Permission) (bool, string) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rules := s.rules[agentID]
	for _, r := range rules {
		if r.PeerID == peerID {
			for _, p := range r.Perms {
				if p == perm {
					return true, r.Scope
				}
			}
		}
	}
	return false, ""
}

// HasAny returns true if peer has any permission.
func (s *ACLStore) HasAny(agentID, peerID string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rules := s.rules[agentID]
	for _, r := range rules {
		if r.PeerID == peerID && len(r.Perms) > 0 {
			return true
		}
	}
	return false
}
