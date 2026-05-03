package runtime

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

var stateMu sync.RWMutex

func statePath(basePath, assistantID string) string {
	return filepath.Join(basePath, "assistants", assistantID, "runtime_state.json")
}

func loadState(basePath, assistantID string) RuntimeState {
	stateMu.RLock()
	defer stateMu.RUnlock()
	p := statePath(basePath, assistantID)
	b, err := os.ReadFile(p)
	if err != nil {
		return DefaultRuntimeState()
	}
	var s RuntimeState
	if err := json.Unmarshal(b, &s); err != nil {
		return DefaultRuntimeState()
	}
	return s
}

func saveState(basePath, assistantID string, s RuntimeState) error {
	stateMu.Lock()
	defer stateMu.Unlock()
	p := statePath(basePath, assistantID)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return fmt.Errorf("creating state dir: %w", err)
	}
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p, b, 0o644)
}

func LoadAssistantState(basePath, assistantID string) RuntimeState {
	return loadState(basePath, assistantID)
}

func SaveAssistantState(basePath, assistantID string, s RuntimeState) error {
	return saveState(basePath, assistantID, s)
}
