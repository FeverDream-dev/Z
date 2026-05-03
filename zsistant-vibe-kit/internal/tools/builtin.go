package tools

import (
    "context"
    "io"
    "net/http"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
    "time"
)

// Helper: path containment within workspace
func isWithinWorkspace(workspaceRoot, requestedPath string) bool {
    if workspaceRoot == "" {
        return false
    }
    absWorkspace, err := filepath.Abs(workspaceRoot)
    if err != nil {
        return false
    }
    target := requestedPath
    if !filepath.IsAbs(requestedPath) {
        target = filepath.Join(workspaceRoot, requestedPath)
    }
    absTarget, err := filepath.Abs(target)
    if err != nil {
        return false
    }
    rel, err := filepath.Rel(absWorkspace, absTarget)
    if err != nil {
        return false
    }
    if strings.HasPrefix(rel, "..") {
        return false
    }
    return true
}

// NewFileReadTool creates a file_read tool scoped to workspaceRoot.
func NewFileReadTool(workspaceRoot string) Tool {
    return &fileReadTool{workspaceRoot: workspaceRoot}
}

type fileReadTool struct{ workspaceRoot string }

func (t *fileReadTool) Name() string { return "file_read" }
func (t *fileReadTool) Description() string { return "Read a file from the workspace" }

func (t *fileReadTool) Execute(call ToolCall) ToolResult {
    v, ok := call.Args["path"].(string)
    if !ok || v == "" {
        return ToolResult{CallID: call.ID, ToolName: t.Name(), Success: false, Output: "", Error: "missing path"}
    }
    if !isWithinWorkspace(t.workspaceRoot, v) {
        return ToolResult{CallID: call.ID, ToolName: t.Name(), Success: false, Output: "", Error: "path outside workspace"}
    }
    absPath := filepath.Join(t.workspaceRoot, v)
    data, err := os.ReadFile(absPath)
    if err != nil {
        return ToolResult{CallID: call.ID, ToolName: t.Name(), Success: false, Output: "", Error: err.Error()}
    }
    return ToolResult{CallID: call.ID, ToolName: t.Name(), Success: true, Output: string(data)}
}

// NewFileWriteTool creates a file_write tool scoped to workspaceRoot.
func NewFileWriteTool(workspaceRoot string) Tool {
    return &fileWriteTool{workspaceRoot: workspaceRoot}
}

type fileWriteTool struct{ workspaceRoot string }

func (t *fileWriteTool) Name() string { return "file_write" }
func (t *fileWriteTool) Description() string { return "Write a file in the workspace" }

func (t *fileWriteTool) Execute(call ToolCall) ToolResult {
    pathVal, ok := call.Args["path"].(string)
    if !ok || pathVal == "" {
        return ToolResult{CallID: call.ID, ToolName: t.Name(), Success: false, Output: "", Error: "missing path"}
    }
    if !isWithinWorkspace(t.workspaceRoot, pathVal) {
        return ToolResult{CallID: call.ID, ToolName: t.Name(), Success: false, Output: "", Error: "path outside workspace"}
    }
    content, _ := call.Args["content"].(string)
    absPath := filepath.Join(t.workspaceRoot, pathVal)
    if err := os.MkdirAll(filepath.Dir(absPath), 0755); err != nil {
        return ToolResult{CallID: call.ID, ToolName: t.Name(), Success: false, Output: "", Error: err.Error()}
    }
    if err := os.WriteFile(absPath, []byte(content), 0644); err != nil {
        return ToolResult{CallID: call.ID, ToolName: t.Name(), Success: false, Output: "", Error: err.Error()}
    }
    return ToolResult{CallID: call.ID, ToolName: t.Name(), Success: true, Output: ""}
}

// NewShellExecTool creates a shell_exec tool scoped to workspaceRoot.
func NewShellExecTool(workspaceRoot string) Tool {
    return &shellExecTool{workspaceRoot: workspaceRoot}
}

type shellExecTool struct{ workspaceRoot string }

func (t *shellExecTool) Name() string { return "shell_exec" }
func (t *shellExecTool) Description() string { return "Execute a shell command in the workspace" }

func (t *shellExecTool) Execute(call ToolCall) ToolResult {
    cmdRaw, ok := call.Args["command"].(string)
    if !ok || cmdRaw == "" {
        return ToolResult{CallID: call.ID, ToolName: t.Name(), Success: false, Output: "", Error: "missing command"}
    }
    // Optional timeout in seconds
    timeout := 30
    if v, ok := call.Args["timeout"]; ok {
        switch tt := v.(type) {
        case int:
            timeout = tt
        case int64:
            timeout = int(tt)
        case float64:
            timeout = int(tt)
        }
    }
    ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
    defer cancel()
    // Run via bash -lc to support complex commands
    cmd := exec.CommandContext(ctx, "bash", "-lc", cmdRaw)
    cmd.Dir = t.workspaceRoot
    out, err := cmd.CombinedOutput()
    res := ToolResult{CallID: call.ID, ToolName: t.Name(), Output: string(out)}
    if ctx.Err() != nil {
        res.Success = false
        res.Error = "timeout"
        return res
    }
    if err != nil {
        res.Success = false
        res.Error = err.Error()
        return res
    }
    res.Success = true
    return res
}

// NewWebFetchTool creates a web_fetch tool.
func NewWebFetchTool() Tool {
    return &webFetchTool{}
}

type webFetchTool struct{}

func (t *webFetchTool) Name() string { return "web_fetch" }
func (t *webFetchTool) Description() string { return "Fetch a URL and return body" }

func (t *webFetchTool) Execute(call ToolCall) ToolResult {
    urlVal, ok := call.Args["url"].(string)
    if !ok || urlVal == "" {
        return ToolResult{CallID: call.ID, ToolName: t.Name(), Success: false, Output: "", Error: "missing url"}
    }
    resp, err := http.Get(urlVal)
    if err != nil {
        return ToolResult{CallID: call.ID, ToolName: t.Name(), Success: false, Output: "", Error: err.Error()}
    }
    defer resp.Body.Close()
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return ToolResult{CallID: call.ID, ToolName: t.Name(), Success: false, Output: "", Error: err.Error()}
    }
    return ToolResult{CallID: call.ID, ToolName: t.Name(), Success: true, Output: string(body)}
}
