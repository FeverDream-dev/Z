package tools

// Tool is the interface every tool must implement.
type Tool interface {
    // Name returns the tool's identifier (e.g. "file_read").
    Name() string
    // Description returns a human-readable description.
    Description() string
    // Execute runs the tool with the given call and returns a result.
    Execute(call ToolCall) ToolResult
}

// ToolCall represents a request to execute a tool.
type ToolCall struct {
    ID       string                 `json:"id"`
    ToolName string                 `json:"tool_name"`
    AgentID  string                 `json:"agent_id"`
    JobID    string                 `json:"job_id,omitempty"`
    Args     map[string]interface{} `json:"args"`
}

// ToolResult represents the outcome of a tool execution.
type ToolResult struct {
    CallID   string `json:"call_id"`
    ToolName string `json:"tool_name"`
    Success  bool   `json:"success"`
    Output   string `json:"output"`
    Error    string `json:"error,omitempty"`
}
