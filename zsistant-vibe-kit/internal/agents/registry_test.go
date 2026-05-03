package agents

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCreateTwoAgentsIsolation(t *testing.T) {
	tmp := t.TempDir()
	reg := NewRegistry(filepath.Join(tmp, ".zazi"))

	a1, err := reg.Create("agent1", "Agent One", "tester")
	if err != nil {
		t.Fatalf("create agent1: %v", err)
	}
	a2, err := reg.Create("agent2", "Agent Two", "tester")
	if err != nil {
		t.Fatalf("create agent2: %v", err)
	}

	if _, err := os.Stat(reg.WorkspacePath(a1.ID)); err != nil {
		t.Fatalf("agent1 workspace missing: %v", err)
	}
	if _, err := os.Stat(reg.WorkspacePath(a2.ID)); err != nil {
		t.Fatalf("agent2 workspace missing: %v", err)
	}

	if reg.IsAllowedPath(a1.ID, reg.WorkspacePath(a2.ID)) {
		t.Fatalf("expected agent1 not to access agent2 workspace")
	}
	if !reg.IsAllowedPath(a1.ID, reg.WorkspacePath(a1.ID)) {
		t.Fatalf("expected allowed path for agent1 self workspace")
	}

	if err := reg.Delete(a1.ID); err != nil {
		t.Fatalf("delete agent1: %v", err)
	}
	if err := reg.Delete(a2.ID); err != nil {
		t.Fatalf("delete agent2: %v", err)
	}

	if _, err := os.Stat(reg.WorkspacePath(a1.ID)); !os.IsNotExist(err) {
		t.Fatalf("agent1 workspace not removed")
	}
	if _, err := os.Stat(reg.WorkspacePath(a2.ID)); !os.IsNotExist(err) {
		t.Fatalf("agent2 workspace not removed")
	}
	if a1.CreatedAt.IsZero() || a2.CreatedAt.IsZero() {
		t.Fatal("expected CreatedAt to be set")
	}
	if a1.UpdatedAt.IsZero() || a2.UpdatedAt.IsZero() {
		t.Fatal("expected UpdatedAt to be set")
	}
	time.Sleep(1 * time.Millisecond)
}

func TestAgentFilesCreated(t *testing.T) {
	tmp := t.TempDir()
	reg := NewRegistry(filepath.Join(tmp, ".zazi"))

	a, err := reg.Create("test-agent", "Test Agent", "assistant")
	if err != nil {
		t.Fatalf("create agent: %v", err)
	}

	root := reg.WorkspacePath(a.ID)
	for _, d := range []string{"memory", "workspace", "jobs"} {
		p := filepath.Join(root, d)
		if stat, err := os.Stat(p); err != nil || !stat.IsDir() {
			t.Fatalf("expected directory %s to exist", p)
		}
	}

	for _, f := range []string{"persona.md", "profile.json", "audit.jsonl", "inbox.jsonl", "outbox.jsonl"} {
		p := filepath.Join(root, f)
		if _, err := os.Stat(p); err != nil {
			t.Fatalf("expected file %s to exist: %v", p, err)
		}
	}

	content, err := os.ReadFile(filepath.Join(root, "persona.md"))
	if err != nil {
		t.Fatalf("read persona.md: %v", err)
	}
	if string(content) != "# Persona for Test Agent\n\nRole: assistant\n" {
		t.Fatalf("unexpected persona.md content: %s", string(content))
	}

	loaded, err := reg.Get(a.ID)
	if err != nil {
		t.Fatalf("get agent: %v", err)
	}
	if loaded.ID != a.ID || loaded.Name != a.Name {
		t.Fatalf("loaded agent mismatch")
	}

	list, err := reg.List()
	if err != nil {
		t.Fatalf("list agents: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 agent, got %d", len(list))
	}
}
