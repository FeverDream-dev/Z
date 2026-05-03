package agents

import (
	"path/filepath"
	"strings"
)

// IsAllowedPath guards requestedPath against escaping the given agentRoot.
// Returns true if the target is inside the agent's workspace.
func IsAllowedPath(agentRoot, requestedPath string) bool {
	absRoot, err := filepath.Abs(agentRoot)
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
