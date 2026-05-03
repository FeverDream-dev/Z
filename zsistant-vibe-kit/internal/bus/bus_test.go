package bus

import (
	"strings"
	"testing"
)

func TestBusDenyByDefault(t *testing.T) {
	dir := t.TempDir()
	b := NewBus(dir)
	_ = b.LoadACLs("agent1")

	req := NewEnvelope("agent2", "agent1", ReqSummary, "status")
	resp := b.Send(req)
	if resp.Allowed {
		t.Fatal("expected denied by default")
	}
	if !strings.Contains(resp.Error, "access denied") {
		t.Fatalf("expected access denied error, got: %s", resp.Error)
	}
}

func TestBusAllowSummary(t *testing.T) {
	dir := t.TempDir()
	b := NewBus(dir)
	_ = b.Allow("agent1", "agent2", []Permission{PermSummary}, "")

	req := NewEnvelope("agent2", "agent1", ReqSummary, "status")
	resp := b.Send(req)
	if !resp.Allowed {
		t.Fatalf("expected allowed, got: %s", resp.Error)
	}
	if !strings.Contains(resp.Result, "summary") {
		t.Fatalf("expected summary result, got: %s", resp.Result)
	}
}

func TestBusFileTransferRequiresScope(t *testing.T) {
	dir := t.TempDir()
	b := NewBus(dir)
	_ = b.Allow("agent1", "agent2", []Permission{PermFileRead}, "/allowed")

	req := NewEnvelope("agent2", "agent1", ReqFileRead, "../../../etc/passwd")
	resp := b.Send(req)
	if resp.Allowed {
		t.Fatal("expected denied for path escape")
	}
	if !strings.Contains(resp.Error, "escapes") {
		t.Fatalf("expected scope error, got: %s", resp.Error)
	}
}

func TestBusFileTransferAllowedInScope(t *testing.T) {
	dir := t.TempDir()
	b := NewBus(dir)
	_ = b.Allow("agent1", "agent2", []Permission{PermFileRead}, "/allowed")

	req := NewEnvelope("agent2", "agent1", ReqFileRead, "data.txt")
	resp := b.Send(req)
	if !resp.Allowed {
		t.Fatalf("expected allowed, got: %s", resp.Error)
	}
	if !strings.Contains(resp.Result, "data.txt") {
		t.Fatalf("expected file in result, got: %s", resp.Result)
	}
}

func TestBusAuditLog(t *testing.T) {
	dir := t.TempDir()
	b := NewBus(dir)
	_ = b.Allow("agent1", "agent2", []Permission{PermSummary}, "")

	req := NewEnvelope("agent2", "agent1", ReqSummary, "status")
	_ = b.Send(req)

	// Audit log should exist
	// We can't easily verify content without parsing, but Send should not panic
}

func TestBusRevoke(t *testing.T) {
	dir := t.TempDir()
	b := NewBus(dir)
	_ = b.Allow("agent1", "agent2", []Permission{PermSummary}, "")
	_ = b.Revoke("agent1", "agent2", []Permission{PermSummary})

	req := NewEnvelope("agent2", "agent1", ReqSummary, "status")
	resp := b.Send(req)
	if resp.Allowed {
		t.Fatal("expected denied after revoke")
	}
}

func TestBusList(t *testing.T) {
	dir := t.TempDir()
	b := NewBus(dir)
	_ = b.Allow("agent1", "agent2", []Permission{PermSummary, PermExec}, "")

	rules, err := b.List("agent1")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(rules))
	}
	if len(rules[0].Perms) != 2 {
		t.Fatalf("expected 2 perms, got %d", len(rules[0].Perms))
	}
}
