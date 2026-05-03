package tools

import (
    "testing"
)

func TestPermissions_IsAllowed_DefaultDeny(t *testing.T) {
    p := DefaultPermissions("agent1")
    if p.IsAllowed("foo") {
        t.Fatalf("expected default deny for foo")
    }
}

func TestPermissions_IsAllowed_ExplicitAllow(t *testing.T) {
    p := &AgentPermissions{AgentID: "a1", ToolsAllowed: []string{"foo"}}
    if !p.IsAllowed("foo") {
        t.Fatalf("expected foo to be allowed")
    }
}

func TestPermissions_IsAllowed_DenyOverrides(t *testing.T) {
    p := &AgentPermissions{AgentID: "a1", ToolsAllowed: []string{"foo"}, ToolsDenied: []string{"foo"}}
    if p.IsAllowed("foo") {
        t.Fatalf("deny override not respected")
    }
}

func TestPermissions_NeedsApproval(t *testing.T) {
    p := &AgentPermissions{AgentID: "a1", RequiresApprovalFor: []string{"bar"}}
    if !p.NeedsApproval("bar") {
        t.Fatalf("expected NeedsApproval to be true for bar")
    }
}

func TestPermissions_SaveAndLoad(t *testing.T) {
    dir := t.TempDir()
    p := &AgentPermissions{AgentID: "a1", ToolsAllowed: []string{"foo"}, ToolsDenied: []string{}, RequiresApprovalFor: []string{}}
    if err := SavePermissions(dir, p); err != nil {
        t.Fatalf("save perms: %v", err)
    }
    loaded, err := LoadPermissions(dir)
    if err != nil {
        t.Fatalf("load perms: %v", err)
    }
    if loaded.AgentID != p.AgentID {
        t.Fatalf("agentID mismatch: %s vs %s", loaded.AgentID, p.AgentID)
    }
    if len(loaded.ToolsAllowed) != len(p.ToolsAllowed) || loaded.ToolsAllowed[0] != p.ToolsAllowed[0] {
        t.Fatalf("permissions not preserved on load")
    }
    // Cleanup is automatic by t.TempDir
}
