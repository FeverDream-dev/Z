package trainer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDetectBrevityShort(t *testing.T) {
	sigs := Detect("Give me a short summary")
	if len(sigs) == 0 {
		t.Fatal("expected signals")
	}
	found := false
	for _, s := range sigs {
		if s.Dimension == "brevity" && s.Value == "short" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected brevity:short signal, got %v", sigs)
	}
}

func TestDetectBrevityDetailed(t *testing.T) {
	sigs := Detect("I need a detailed explanation")
	found := false
	for _, s := range sigs {
		if s.Dimension == "brevity" && s.Value == "detailed" {
			found = true
		}
	}
	if !found {
		t.Fatal("expected brevity:detailed signal")
	}
}

func TestDetectTone(t *testing.T) {
	sigs := Detect("Use a formal tone please")
	found := false
	for _, s := range sigs {
		if s.Dimension == "tone" && s.Value == "formal" {
			found = true
		}
	}
	if !found {
		t.Fatal("expected tone:formal signal")
	}
}

func TestDetectStructure(t *testing.T) {
	sigs := Detect("Show me a bullet list")
	found := false
	for _, s := range sigs {
		if s.Dimension == "structure" && s.Value == "bullets" {
			found = true
		}
	}
	if !found {
		t.Fatal("expected structure:bullets signal")
	}
}

func TestDetectFeedback(t *testing.T) {
	sigs := DetectFeedback("That was too long, make it shorter")
	found := false
	for _, s := range sigs {
		if s.Dimension == "brevity" && s.Value == "short" {
			found = true
		}
	}
	if !found {
		t.Fatal("expected brevity:short from feedback")
	}
}

func TestStyleProfileAccumulate(t *testing.T) {
	sp := NewStyleProfile()
	sp.Apply(Signal{Dimension: "brevity", Value: "short"})
	sp.Apply(Signal{Dimension: "brevity", Value: "short"})
	sp.Apply(Signal{Dimension: "brevity", Value: "detailed"})

	top := sp.Top()
	b, ok := top["brevity"]
	if !ok {
		t.Fatal("expected brevity dimension")
	}
	if b.Value != "short" {
		t.Fatalf("expected short to win, got %s", b.Value)
	}
	if b.Score != 2 {
		t.Fatalf("expected score 2, got %d", b.Score)
	}
}

func TestTrainerProposePatch(t *testing.T) {
	tr := NewTrainer()
	// Not enough evidence
	tr.Observe("short")
	patch := tr.ProposePatch(3)
	if patch != "" {
		t.Fatalf("expected no proposal with threshold 3, got: %s", patch)
	}

	// Accumulate enough evidence
	tr.Observe("short summary")
	tr.Observe("keep it brief")
	patch = tr.ProposePatch(3)
	if patch == "" {
		t.Fatal("expected proposal with threshold 3")
	}
	if !strings.Contains(patch, "short") {
		t.Fatalf("expected 'short' in patch, got: %s", patch)
	}
}

func TestTrainerDoesNotInferSensitiveTraits(t *testing.T) {
	tr := NewTrainer()
	// Messages that might suggest personal info
	tr.Observe("I am feeling depressed today")
	tr.Observe("My political views are conservative")
	tr.Observe("I earn $100k per year")

	patch := tr.ProposePatch(1)
	// Should not propose anything since these are not style signals
	if strings.Contains(patch, "depressed") || strings.Contains(patch, "political") || strings.Contains(patch, "earn") {
		t.Fatal("trainer should not infer sensitive traits")
	}
}

func TestIsMajorChange(t *testing.T) {
	if !IsMajorChange("User prefers a formal tone.") {
		t.Fatal("expected tone to be major change")
	}
	if !IsMajorChange("User prefers agent to ask first.") {
		t.Fatal("expected autonomy to be major change")
	}
	if IsMajorChange("User prefers short responses.") {
		t.Fatal("expected brevity to not be major change")
	}
}

func TestReadWritePersona(t *testing.T) {
	dir := t.TempDir()
	content := "# Test Agent\n\nRole: helper\n"
	if err := os.WriteFile(filepath.Join(dir, "persona.md"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	read, err := ReadPersona(dir)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if read != content {
		t.Fatalf("expected %q, got %q", content, read)
	}

	newContent := "# Updated"
	if err := WritePersona(dir, newContent); err != nil {
		t.Fatalf("write: %v", err)
	}
	read, _ = ReadPersona(dir)
	if read != newContent {
		t.Fatalf("expected %q, got %q", newContent, read)
	}
}

func TestApplyPatch(t *testing.T) {
	dir := t.TempDir()
	base := "# Agent\n"
	_ = os.WriteFile(filepath.Join(dir, "persona.md"), []byte(base), 0644)

	patch := "## Learned\n- short"
	if err := ApplyPatch(dir, patch); err != nil {
		t.Fatalf("apply: %v", err)
	}

	read, _ := ReadPersona(dir)
	if !strings.Contains(read, patch) {
		t.Fatalf("expected patch in persona, got: %s", read)
	}

	// Re-applying should fail
	if err := ApplyPatch(dir, patch); err == nil {
		t.Fatal("expected error re-applying same patch")
	}
}
