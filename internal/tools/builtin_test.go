package tools

import (
    "net/http/httptest"
    "net/http"
    "os"
    "path/filepath"
    "testing"
)

func TestFileReadTool_WithinWorkspace(t *testing.T) {
    workspace, cleanup := t.TempDir(), func() {}
    _ = cleanup
    // write a file inside workspace
    if err := os.WriteFile(filepath.Join(workspace, "data.txt"), []byte("content"), 0644); err != nil {
        t.Fatalf("write file: %v", err)
    }
    tool := NewFileReadTool(workspace)
    call := ToolCall{ID: "cr1", ToolName: tool.Name(), Args: map[string]interface{}{"path": "data.txt"}}
    res := tool.Execute(call)
    if !res.Success {
        t.Fatalf("expected success, got error: %s", res.Error)
    }
    if res.Output != "content" {
        t.Fatalf("unexpected content: %q", res.Output)
    }
}

func TestFileReadTool_OutsideWorkspace(t *testing.T) {
    workspace := t.TempDir()
    // Create a path that attempts to escape
    tool := NewFileReadTool(workspace)
    call := ToolCall{ID: "cr2", ToolName: tool.Name(), Args: map[string]interface{}{"path": "../outside.txt"}}
    res := tool.Execute(call)
    if res.Success {
        t.Fatalf("expected failure for outside path")
    }
}

func TestFileWriteTool_WithinWorkspace(t *testing.T) {
    workspace := t.TempDir()
    tool := NewFileWriteTool(workspace)
    call := ToolCall{ID: "cw1", ToolName: tool.Name(), Args: map[string]interface{}{"path": "out.txt", "content": "hello"}}
    res := tool.Execute(call)
    if !res.Success {
        t.Fatalf("expected success writing inside workspace: %s", res.Error)
    }
    data, err := os.ReadFile(filepath.Join(workspace, "out.txt"))
    if err != nil {
        t.Fatalf("read back: %v", err)
    }
    if string(data) != "hello" {
        t.Fatalf("unexpected file content: %q", string(data))
    }
}

func TestFileWriteTool_OutsideWorkspace(t *testing.T) {
    workspace := t.TempDir()
    tool := NewFileWriteTool(workspace)
    // Attempt to write outside; use relative path that escapes
    call := ToolCall{ID: "cw2", ToolName: tool.Name(), Args: map[string]interface{}{"path": "../outside2.txt", "content": "x"}}
    res := tool.Execute(call)
    if res.Success {
        t.Fatalf("expected failure for outside path write")
    }
}

func TestShellExecTool_BasicCommand(t *testing.T) {
    workspace := t.TempDir()
    tool := NewShellExecTool(workspace)
    call := ToolCall{ID: "se1", ToolName: tool.Name(), Args: map[string]interface{}{"command": "echo hello"}}
    res := tool.Execute(call)
    if !res.Success {
        t.Fatalf("expected success for echo, got: %s", res.Error)
    }
    if res.Output != "hello\n" {
        t.Fatalf("unexpected output: %q", res.Output)
    }
}

func TestShellExecTool_Timeout(t *testing.T) {
    workspace := t.TempDir()
    tool := NewShellExecTool(workspace)
    call := ToolCall{ID: "se2", ToolName: tool.Name(), Args: map[string]interface{}{"command": "sleep 2", "timeout": 1}}
    res := tool.Execute(call)
    if res.Success {
        t.Fatalf("expected timeout, got success")
    }
}

func TestWebFetchTool_Basic(t *testing.T) {
    // Create a test HTTP server
    srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("OK"))
    }))
    defer srv.Close()

    tool := NewWebFetchTool()
    call := ToolCall{ID: "wf1", ToolName: tool.Name(), Args: map[string]interface{}{"url": srv.URL}}
    res := tool.Execute(call)
    if !res.Success {
        t.Fatalf("web_fetch failed: %s", res.Error)
    }
    if res.Output != "OK" {
        t.Fatalf("unexpected body: %q", res.Output)
    }
}
