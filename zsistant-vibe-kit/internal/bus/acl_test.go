package bus

import (
	"testing"
)

func TestACLAllowAndCheck(t *testing.T) {
	dir := t.TempDir()
	s := NewACLStore(dir)
	_ = s.Load("agent1")

	_ = s.Allow("agent1", "agent2", []Permission{PermSummary}, "")
	allowed, _ := s.Check("agent1", "agent2", PermSummary)
	if !allowed {
		t.Fatal("expected summary permission")
	}
	allowed, _ = s.Check("agent1", "agent2", PermExec)
	if allowed {
		t.Fatal("expected no exec permission")
	}
}

func TestACLDenyByDefault(t *testing.T) {
	dir := t.TempDir()
	s := NewACLStore(dir)
	_ = s.Load("agent1")

	allowed, _ := s.Check("agent1", "agent3", PermSummary)
	if allowed {
		t.Fatal("expected denied by default")
	}
}

func TestACLRevoke(t *testing.T) {
	dir := t.TempDir()
	s := NewACLStore(dir)
	_ = s.Load("agent1")

	_ = s.Allow("agent1", "agent2", []Permission{PermSummary, PermFileRead}, "")
	_ = s.Revoke("agent1", "agent2", []Permission{PermSummary})

	allowed, _ := s.Check("agent1", "agent2", PermSummary)
	if allowed {
		t.Fatal("expected summary revoked")
	}
	allowed, _ = s.Check("agent1", "agent2", PermFileRead)
	if !allowed {
		t.Fatal("expected file_read still allowed")
	}
}

func TestACLRevokeAll(t *testing.T) {
	dir := t.TempDir()
	s := NewACLStore(dir)
	_ = s.Load("agent1")

	_ = s.Allow("agent1", "agent2", []Permission{PermSummary}, "")
	_ = s.Revoke("agent1", "agent2", nil)

	allowed, _ := s.Check("agent1", "agent2", PermSummary)
	if allowed {
		t.Fatal("expected all revoked")
	}
}

func TestACLMergePerms(t *testing.T) {
	dir := t.TempDir()
	s := NewACLStore(dir)
	_ = s.Load("agent1")

	_ = s.Allow("agent1", "agent2", []Permission{PermSummary}, "")
	_ = s.Allow("agent1", "agent2", []Permission{PermExec}, "")

	allowed, _ := s.Check("agent1", "agent2", PermSummary)
	if !allowed {
		t.Fatal("expected summary")
	}
	allowed, _ = s.Check("agent1", "agent2", PermExec)
	if !allowed {
		t.Fatal("expected exec")
	}
}

func TestACLSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	s := NewACLStore(dir)
	_ = s.Load("agent1")
	_ = s.Allow("agent1", "agent2", []Permission{PermSummary}, "/tmp")
	_ = s.Save("agent1")

	s2 := NewACLStore(dir)
	_ = s2.Load("agent1")
	allowed, scope := s2.Check("agent1", "agent2", PermSummary)
	if !allowed {
		t.Fatal("expected permission after reload")
	}
	if scope != "/tmp" {
		t.Fatalf("expected scope /tmp, got %s", scope)
	}
}

func TestACLHasAny(t *testing.T) {
	dir := t.TempDir()
	s := NewACLStore(dir)
	_ = s.Load("agent1")

	if s.HasAny("agent1", "agent2") {
		t.Fatal("expected no permissions")
	}
	_ = s.Allow("agent1", "agent2", []Permission{PermSummary}, "")
	if !s.HasAny("agent1", "agent2") {
		t.Fatal("expected has permissions")
	}
}
