package jobs

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewJob(t *testing.T) {
	j := NewJob("agent1", "cli", "test objective")
	if j.AgentID != "agent1" {
		t.Fatalf("expected agent_id agent1, got %s", j.AgentID)
	}
	if j.SourceChannel != "cli" {
		t.Fatalf("expected source_channel cli, got %s", j.SourceChannel)
	}
	if j.Objective != "test objective" {
		t.Fatalf("expected objective test objective, got %s", j.Objective)
	}
	if j.Status != "queued" {
		t.Fatalf("expected status queued, got %s", j.Status)
	}
	if j.ID == "" {
		t.Fatal("expected job id to be set")
	}
}

func TestAppendEvent(t *testing.T) {
	tmp := t.TempDir()
	evt := JobEvent{
		JobID:     "job-1",
		AgentID:   "agent1",
		EventType: "job.created",
		Message:   "test event",
	}
	if err := AppendEvent(tmp, evt); err != nil {
		t.Fatalf("append event: %v", err)
	}
	b, err := os.ReadFile(filepath.Join(tmp, "audit.jsonl"))
	if err != nil {
		t.Fatalf("read audit: %v", err)
	}
	if !strings.Contains(string(b), "job.created") {
		t.Fatalf("expected audit to contain event type, got: %s", string(b))
	}
	if !strings.Contains(string(b), "test event") {
		t.Fatalf("expected audit to contain message, got: %s", string(b))
	}
}
