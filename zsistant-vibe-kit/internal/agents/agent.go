package agents

import "time"

type Agent struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Slug          string    `json:"slug"`
	Role          string    `json:"role"`
	WorkspaceRoot string    `json:"workspace_root"`
	PersonaPath   string    `json:"persona_path"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Status        string    `json:"status"`
	EnabledChannels []string `json:"enabled_channels,omitempty"`
	ToolPermissions []string `json:"tool_permissions,omitempty"`
	ModelPolicy     string   `json:"model_policy,omitempty"`
	MemoryPolicy    string   `json:"memory_policy,omitempty"`
}
