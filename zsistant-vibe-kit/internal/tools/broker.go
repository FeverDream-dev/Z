package tools

import (
    "fmt"

    "github.com/FeverDream-dev/zsistant/internal/jobs"
)

// Broker mediates tool calls with permission checking and audit logging.
type Broker struct {
    tools         map[string]Tool
    workspaceRoot string // agent workspace root for permission loading
}

// NewBroker creates a broker with the given tools registered.
func NewBroker(workspaceRoot string) *Broker {
    return &Broker{
        tools:         make(map[string]Tool),
        workspaceRoot: workspaceRoot,
    }
}

// Register adds a tool to the broker.
func (b *Broker) Register(t Tool) {
    if b.tools == nil {
        b.tools = make(map[string]Tool)
    }
    b.tools[t.Name()] = t
}

// Call executes a tool call after checking permissions.
// Returns ToolResult. Logs tool.called, tool.completed/tool.failed events.
func (b *Broker) Call(call ToolCall) ToolResult {
    // Load permissions for this workspace
    perm, err := LoadPermissions(b.workspaceRoot)
    if err != nil || perm == nil {
        perm = DefaultPermissions("")
    }

    // Permission check
    if !perm.IsAllowed(call.ToolName) {
        return ToolResult{CallID: call.ID, ToolName: call.ToolName, Success: false, Output: "", Error: "permission denied"}
    }

    // Approval check
    if perm.NeedsApproval(call.ToolName) {
        return ToolResult{CallID: call.ID, ToolName: call.ToolName, Success: false, Output: "", Error: "requires approval"}
    }

    // Audit log: tool.called
    jobs.AppendEvent(b.workspaceRoot, jobs.JobEvent{
        EventType: "tool.called",
        Message:   fmt.Sprintf("Calling tool %s", call.ToolName),
        Metadata: map[string]string{
            "tool":     call.ToolName,
            "agent_id": call.AgentID,
            "job_id":   call.JobID,
        },
    })

    // Execute the tool if registered
    tool, ok := b.tools[call.ToolName]
    if !ok {
        return ToolResult{CallID: call.ID, ToolName: call.ToolName, Success: false, Output: "", Error: "unknown tool"}
    }
    result := tool.Execute(call)

    // Audit log: tool.completed or tool.failed
    if result.Success {
        jobs.AppendEvent(b.workspaceRoot, jobs.JobEvent{
            EventType: "tool.completed",
            Message:   fmt.Sprintf("Tool %s completed", call.ToolName),
            Metadata: map[string]string{
                "tool": call.ToolName,
            },
        })
    } else {
        jobs.AppendEvent(b.workspaceRoot, jobs.JobEvent{
            EventType: "tool.failed",
            Message:   fmt.Sprintf("Tool %s failed", call.ToolName),
            Metadata: map[string]string{
                "tool":  call.ToolName,
                "error": result.Error,
            },
        })
    }

    // Return the raw result
    result.CallID = call.ID
    result.ToolName = call.ToolName
    return result
}

// ListTools returns info about all registered tools.
func (b *Broker) ListTools() []ToolInfo {
    infos := []ToolInfo{}
    for _, t := range b.tools {
        infos = append(infos, ToolInfo{Name: t.Name(), Description: t.Description()})
    }
    return infos
}

type ToolInfo struct {
    Name        string `json:"name"`
    Description string `json:"description"`
}
