package tools

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

// AgentPermissions defines what tools an agent can use.
type AgentPermissions struct {
    AgentID             string   `json:"agent_id"`
    ToolsAllowed        []string `json:"tools_allowed"`
    ToolsDenied         []string `json:"tools_denied"`
    RequiresApprovalFor []string `json:"requires_approval_for"`
    FilesystemScope     string   `json:"filesystem_scope"` // "workspace" or "none"
}

// contains is a helper to check membership (case-sensitive).
func contains(list []string, v string) bool {
    for _, x := range list {
        if x == v {
            return true
        }
    }
    return false
}

// IsAllowed checks if a tool is permitted for this agent.
// Default deny: must be explicitly allowed and not denied.
func (p *AgentPermissions) IsAllowed(toolName string) bool {
    if p == nil {
        return false
    }
    allowed := contains(p.ToolsAllowed, toolName)
    denied := contains(p.ToolsDenied, toolName)
    return allowed && !denied
}

// NeedsApproval checks if a tool requires human approval.
func (p *AgentPermissions) NeedsApproval(toolName string) bool {
    if p == nil {
        return false
    }
    return contains(p.RequiresApprovalFor, toolName)
}

// LoadPermissions reads permissions from agent workspace.
// Path: <workspaceRoot>/permissions.json
func LoadPermissions(workspaceRoot string) (*AgentPermissions, error) {
    if workspaceRoot == "" {
        return DefaultPermissions(""), nil
    }
    ppath := filepath.Join(workspaceRoot, "permissions.json")
    b, err := ioutil.ReadFile(ppath)
    if err != nil {
        if os.IsNotExist(err) {
            // No permissions file yet; return default safe permissions
            return DefaultPermissions(""), nil
        }
        return nil, err
    }
    var p AgentPermissions
    if err := json.Unmarshal(b, &p); err != nil {
        return nil, err
    }
    return &p, nil
}

// SavePermissions writes permissions to agent workspace.
func SavePermissions(workspaceRoot string, p *AgentPermissions) error {
    if workspaceRoot == "" {
        return nil
    }
    if p == nil {
        p = DefaultPermissions("")
    }
    data, err := json.MarshalIndent(p, "", "  ")
    if err != nil {
        return err
    }
    ppath := filepath.Join(workspaceRoot, "permissions.json")
    return ioutil.WriteFile(ppath, data, 0644)
}

// DefaultPermissions returns a safe default (no tools allowed).
func DefaultPermissions(agentID string) *AgentPermissions {
    return &AgentPermissions{
        AgentID:             agentID,
        ToolsAllowed:        []string{},
        ToolsDenied:         []string{},
        RequiresApprovalFor: []string{},
        FilesystemScope:     "workspace",
    }
}
