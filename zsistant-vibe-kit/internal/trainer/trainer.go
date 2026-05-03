package trainer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Trainer watches user interactions and proposes persona patches.
type Trainer struct {
	profile *StyleProfile
}

// NewTrainer creates a new trainer.
func NewTrainer() *Trainer {
	return &Trainer{profile: NewStyleProfile()}
}

// Profile returns the trainer's accumulated style profile.
func (t *Trainer) Profile() *StyleProfile {
	return t.profile
}

// Observe processes a user message for style signals.
func (t *Trainer) Observe(message string) {
	for _, s := range Detect(message) {
		t.profile.Apply(s)
	}
}

// ObserveFeedback processes explicit user feedback.
func (t *Trainer) ObserveFeedback(message string) {
	for _, s := range DetectFeedback(message) {
		t.profile.Apply(s)
	}
}

// ProposePatch generates a persona.md patch based on observed signals.
// It only proposes changes when confidence is above threshold.
func (t *Trainer) ProposePatch(threshold int) string {
	var lines []string
	lines = append(lines, "## Learned communication preferences")
	lines = append(lines, "")

	hasProposal := false
	top := t.profile.Top()
	for _, d := range top {
		if d.Score < threshold {
			continue
		}
		hasProposal = true
		switch d.Name {
		case "brevity":
			lines = append(lines, fmt.Sprintf("- User prefers %s responses.", d.Value))
		case "tone":
			lines = append(lines, fmt.Sprintf("- User prefers a %s tone.", d.Value))
		case "explanation_depth":
			lines = append(lines, fmt.Sprintf("- User wants %s explanation depth.", d.Value))
		case "structure":
			lines = append(lines, fmt.Sprintf("- User prefers output in %s.", d.Value))
		case "autonomy":
			lines = append(lines, fmt.Sprintf("- User prefers agent to %s.", strings.ReplaceAll(d.Value, "_", " ")))
		}
	}

	if !hasProposal {
		return ""
	}
	return strings.Join(lines, "\n")
}

// IsMajorChange returns true if a patch involves tone or autonomy changes.
func IsMajorChange(patch string) bool {
	lower := strings.ToLower(patch)
	return strings.Contains(lower, "tone") ||
		strings.Contains(lower, "autonomy") ||
		strings.Contains(lower, "ask first") ||
		strings.Contains(lower, "act with updates") ||
		strings.Contains(lower, "act aggressively")
}

// ReadPersona reads the current persona.md content.
func ReadPersona(workspace string) (string, error) {
	path := filepath.Join(workspace, "persona.md")
	b, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("reading persona: %w", err)
	}
	return string(b), nil
}

// WritePersona writes persona.md content.
func WritePersona(workspace, content string) error {
	path := filepath.Join(workspace, "persona.md")
	return os.WriteFile(path, []byte(content), 0644)
}

// ApplyPatch appends a patch to persona.md if it doesn't already exist.
func ApplyPatch(workspace, patch string) error {
	current, err := ReadPersona(workspace)
	if err != nil {
		return err
	}
	if strings.Contains(current, patch) {
		return fmt.Errorf("patch already applied")
	}
	return WritePersona(workspace, current+"\n\n"+patch)
}
