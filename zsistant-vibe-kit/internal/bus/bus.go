package bus

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Bus routes inter-agent requests with ACL enforcement.
type Bus struct {
	store *ACLStore
	base  string
}

// NewBus creates a message bus at the given base path.
func NewBus(base string) *Bus {
	return &Bus{
		store: NewACLStore(base),
		base:  base,
	}
}

// LoadACLs loads ACL rules for an agent.
func (b *Bus) LoadACLs(agentID string) error {
	return b.store.Load(agentID)
}

// Send routes a request envelope and returns the response.
// Denies by default if no ACL rule exists.
func (b *Bus) Send(env *Envelope) *Response {
	// Load recipient ACLs
	_ = b.store.Load(env.To)

	allowed, scope := b.store.Check(env.To, env.From, Permission(env.Type))
	if !allowed {
		// Audit denial
		_ = b.audit(env.To, env.From, env.Type, false, "no acl rule")
		return NewResponse(env, false, "", fmt.Sprintf("access denied: %s has no %s permission from %s", env.To, env.Type, env.From))
	}

	// Audit approval
	_ = b.audit(env.To, env.From, env.Type, true, scope)

	// Process based on type
	var result string
	var err string
	switch env.Type {
	case ReqSummary:
		result = fmt.Sprintf("summary for %s: ok", env.To)
	case ReqFileRead:
		if scope != "" && !filepath.IsLocal(env.Payload) {
			err = "path escapes allowed scope"
		} else {
			result = fmt.Sprintf("would read %s (scope: %s)", env.Payload, scope)
		}
	case ReqFileWrite:
		if scope != "" && !filepath.IsLocal(env.Payload) {
			err = "path escapes allowed scope"
		} else {
			result = fmt.Sprintf("would write %s (scope: %s)", env.Payload, scope)
		}
	case ReqExec:
		result = fmt.Sprintf("would execute: %s", env.Payload)
	default:
		err = "unknown request type"
	}

	if err != "" {
		return NewResponse(env, false, "", err)
	}
	return NewResponse(env, true, result, "")
}

// Allow grants a permission.
func (b *Bus) Allow(agentID, peerID string, perms []Permission, scope string) error {
	_ = b.store.Load(agentID)
	if err := b.store.Allow(agentID, peerID, perms, scope); err != nil {
		return err
	}
	return b.store.Save(agentID)
}

// Revoke removes permissions.
func (b *Bus) Revoke(agentID, peerID string, perms []Permission) error {
	_ = b.store.Load(agentID)
	if err := b.store.Revoke(agentID, peerID, perms); err != nil {
		return err
	}
	return b.store.Save(agentID)
}

// List returns ACL rules for an agent.
func (b *Bus) List(agentID string) ([]ACLRule, error) {
	if err := b.store.Load(agentID); err != nil {
		return nil, err
	}
	return b.store.List(agentID), nil
}

// audit logs ACL decisions.
func (b *Bus) audit(agentID, peerID string, reqType RequestType, allowed bool, scope string) error {
	path := filepath.Join(b.base, "agents", agentID, "bus_audit.jsonl")
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	status := "denied"
	if allowed {
		status = "allowed"
	}
	line := fmt.Sprintf(`{"time":"%s","peer":"%s","type":"%s","status":"%s","scope":"%s"}`+"\n",
		time.Now().Format(time.RFC3339), peerID, reqType, status, scope)
	_, err = f.WriteString(line)
	return err
}
