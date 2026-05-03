package tools

import (
    "os"
    "path/filepath"
    "testing"
)

// Helper to create a temp workspace with given permissions.json content
func createWorkspaceWithPerms(t *testing.T, permJSON string) (string, func()) {
    t.Helper()
    dir, err := os.MkdirTemp("", "broker-workspace-")
    if err != nil { t.Fatalf("temp dir: %v", err) }
    // Write permissions.json if provided
    if permJSON != "" {
        if err := os.WriteFile(filepath.Join(dir, "permissions.json"), []byte(permJSON), 0644); err != nil {
            t.Fatalf("write perms: %v", err)
        }
    }
    // Return cleanup function
    return dir, func() { os.RemoveAll(dir) }
}

func TestBroker_Call_AllowedTool(t *testing.T) {
    workspace, cleanup := createWorkspaceWithPerms(t, `{
  "agent_id": "agent1",
  "tools_allowed": ["file_read"],
  "tools_denied": [],
  "requires_approval_for": [],
  "filesystem_scope": "workspace"
}`)
    defer cleanup()

    // Setup broker and register a single tool
    broker := NewBroker(workspace)
    broker.Register(NewFileReadTool(workspace))

    // Create a file to read
    if err := os.MkdirAll(workspace, 0755); err != nil {
        t.Fatalf("mkdir: %v", err)
    }
    msg := []byte("hello world")
    if err := os.WriteFile(filepath.Join(workspace, "data.txt"), msg, 0644); err != nil {
        t.Fatalf("write data: %v", err)
    }

    call := ToolCall{ID: "c1", ToolName: "file_read", Args: map[string]interface{}{"path": "data.txt"}}
    res := broker.Call(call)
    if !res.Success {
        t.Fatalf("expected success, got error: %s", res.Error)
    }
    if res.Output != string(msg) {
        t.Fatalf("unexpected output: %q", res.Output)
    }
}

func TestBroker_Call_DeniedTool(t *testing.T) {
    workspace, cleanup := createWorkspaceWithPerms(t, `{
  "agent_id": "agent1",
  "tools_allowed": [],
  "tools_denied": [],
  "requires_approval_for": [],
  "filesystem_scope": "workspace"
}`)
    defer cleanup()

    broker := NewBroker(workspace)
    broker.Register(NewFileReadTool(workspace))

    call := ToolCall{ID: "c2", ToolName: "file_read", Args: map[string]interface{}{"path": "data.txt"}}
    res := broker.Call(call)
    if res.Success {
        t.Fatalf("expected denial, got success")
    }
    if res.Error != "permission denied" {
        t.Fatalf("expected permission denied, got %q", res.Error)
    }
}

func TestBroker_Call_NeedsApproval(t *testing.T) {
    workspace, cleanup := createWorkspaceWithPerms(t, `{
  "agent_id": "agent1",
  "tools_allowed": ["shell_exec"],
  "tools_denied": [],
  "requires_approval_for": ["shell_exec"],
  "filesystem_scope": "workspace"
}`)
    defer cleanup()

    broker := NewBroker(workspace)
    broker.Register(NewShellExecTool(workspace))
    call := ToolCall{ID: "c3", ToolName: "shell_exec", Args: map[string]interface{}{"command": "echo hi"}}
    res := broker.Call(call)
    if res.Success {
        t.Fatalf("expected approval required error, got success")
    }
    if res.Error != "requires approval" {
        t.Fatalf("expected requires approval, got %q", res.Error)
    }
}

func TestBroker_ListTools(t *testing.T) {
    workspace, cleanup := createWorkspaceWithPerms(t, `{"agent_id":"a","tools_allowed":["file_read","shell_exec"],"tools_denied":[],"filesystem_scope":"workspace"}`)
    defer cleanup()
    broker := NewBroker(workspace)
    broker.Register(NewFileReadTool(workspace))
    broker.Register(NewShellExecTool(workspace))
    infos := broker.ListTools()
    if len(infos) != 2 {
        t.Fatalf("expected 2 tools, got %d", len(infos))
    }
}

func TestBroker_Call_Logging(t *testing.T) {
    workspace, cleanup := createWorkspaceWithPerms(t, `{"agent_id":"a","tools_allowed":["file_read"],"tools_denied":[],"filesystem_scope":"workspace"}`)
    defer cleanup()
    broker := NewBroker(workspace)
    broker.Register(NewFileReadTool(workspace))
    // Prepare a file
    if err := os.WriteFile(filepath.Join(workspace, "logdata.txt"), []byte("log"), 0644); err != nil {
        t.Fatalf("write logdata: %v", err)
    }
    call := ToolCall{ID: "c4", ToolName: "file_read", Args: map[string]interface{}{"path": "logdata.txt"}}
    _ = broker.Call(call)
    // Try to read logs; we assume audit log file is audit.jsonl
    auditPath := filepath.Join(workspace, "audit.jsonl")
    if _, err := os.Stat(auditPath); err != nil {
        // If no audit file, still allow test to pass if AppendEvent is in-memory; fail gracefully
        // Return early to avoid flake; but for determinism, fail the test if file not present
        t.Fatalf("audit log not found at %s", auditPath)
    } else {
        data, _ := os.ReadFile(auditPath)
        if len(data) == 0 {
            t.Fatalf("audit log empty; expected entries")
        }
    }
}
