package agents

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Registry struct {
	basePath string
}

// New creates a new agent Registry at the given base path.
func New(basePath string) *Registry {
	return &Registry{basePath: basePath}
}

// NewRegistry is a compatibility wrapper returning the same registry type.
func NewRegistry(basePath string) *Registry {
	return New(basePath)
}

func (r *Registry) WorkspacePath(id string) string {
	return filepath.Join(r.basePath, "agents", id)
}

func (r *Registry) Create(id, name, role string) (*Agent, error) {
	root := r.WorkspacePath(id)
	dirs := []string{root, filepath.Join(root, "memory"), filepath.Join(root, "workspace"), filepath.Join(root, "jobs")}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			return nil, err
		}
	}
	personaPath := filepath.Join(root, "persona.md")
	if err := os.WriteFile(personaPath, []byte("# Persona for "+name+"\n\nRole: "+role+"\n"), 0644); err != nil {
		return nil, err
	}
	a := Agent{
		ID:              id,
		Name:            name,
		Slug:            slugify(id),
		Role:            role,
		WorkspaceRoot:   root,
		PersonaPath:     personaPath,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		Status:          "created",
		EnabledChannels: []string{"web"},
		ToolPermissions: []string{},
		ModelPolicy:     "default",
		MemoryPolicy:    "session",
	}
	f, err := os.Create(filepath.Join(root, "profile.json"))
	if err != nil {
		return nil, err
	}
	if err := json.NewEncoder(f).Encode(a); err != nil {
		f.Close()
		return nil, err
	}
	f.Close()
	for _, fname := range []string{"audit.jsonl", "inbox.jsonl", "outbox.jsonl"} {
		if _, err := os.Create(filepath.Join(root, fname)); err != nil {
			return nil, err
		}
	}
	return &a, nil
}

func (r *Registry) List() ([]Agent, error) {
	base := filepath.Join(r.basePath, "agents")
	f, err := os.Open(base)
	if err != nil {
		if os.IsNotExist(err) {
			return []Agent{}, nil
		}
		return nil, err
	}
	defer f.Close()
	entries, err := f.Readdir(-1)
	if err != nil {
		return nil, err
	}
	var out []Agent
	for _, e := range entries {
		if e.IsDir() {
			path := filepath.Join(base, e.Name(), "profile.json")
			b, err := os.ReadFile(path)
			if err != nil {
				continue
			}
			var a Agent
			if err := json.Unmarshal(b, &a); err != nil {
				continue
			}
			out = append(out, a)
		}
	}
	return out, nil
}

func (r *Registry) Get(id string) (*Agent, error) {
	b, err := os.ReadFile(filepath.Join(r.basePath, "agents", id, "profile.json"))
	if err != nil {
		return nil, err
	}
	var a Agent
	if err := json.Unmarshal(b, &a); err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *Registry) Delete(id string) error {
	return os.RemoveAll(filepath.Join(r.basePath, "agents", id))
}

func (r *Registry) IsAllowedPath(agentID, requestedPath string) bool {
	root := filepath.Join(r.basePath, "agents", agentID)
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return false
	}
	absTarget, err := filepath.Abs(requestedPath)
	if err != nil {
		return false
	}
	rel, err := filepath.Rel(absRoot, absTarget)
	if err != nil {
		return false
	}
	if strings.HasPrefix(rel, "..") {
		return false
	}
	return true
}

func slugify(s string) string {
	s = strings.ToLower(s)
	out := make([]rune, 0, len(s))
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			out = append(out, r)
		}
	}
	if len(out) == 0 {
		return "agent"
	}
	return string(out)
}
