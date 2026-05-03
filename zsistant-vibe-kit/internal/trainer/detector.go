package trainer

import (
	"strings"
)

// Detectors maps keywords to style signals.
var detectors = []struct {
	keywords  []string
	dimension string
	value     string
}{
	// Brevity
	{[]string{"short", "brief", "quick", "tl;dr", "summarize", "concise"}, "brevity", "short"},
	{[]string{"long", "detailed", "explain", "elaborate", "thorough", "in-depth"}, "brevity", "detailed"},

	// Tone
	{[]string{"formal", "professional", "business"}, "tone", "formal"},
	{[]string{"casual", "friendly", "relaxed", "informal"}, "tone", "friendly"},
	{[]string{"technical", "code", "implementation", "architecture"}, "tone", "technical"},

	// Explanation depth
	{[]string{"just do it", "go ahead", "no explanation", "skip details"}, "explanation_depth", "low"},
	{[]string{"explain why", "how does", "why did", "rationale"}, "explanation_depth", "high"},

	// Structure preference
	{[]string{"bullet", "list", "points", "itemize"}, "structure", "bullets"},
	{[]string{"paragraph", "essay", "narrative", "prose"}, "structure", "paragraphs"},
	{[]string{"checklist", "todo", "tasks", "steps"}, "structure", "checklists"},
	{[]string{"table", "grid", "columns", "matrix"}, "structure", "tables"},

	// Autonomy
	{[]string{"ask first", "check with me", "confirm before"}, "autonomy", "ask_first"},
	{[]string{"just go", "act now", "don't wait", "proceed"}, "autonomy", "act_aggressive"},
}

// Detect analyzes a message and returns any matched style signals.
func Detect(message string) []Signal {
	lower := strings.ToLower(message)
	var signals []Signal
	for _, d := range detectors {
		for _, kw := range d.keywords {
			if strings.Contains(lower, kw) {
				signals = append(signals, Signal{
					Dimension: d.dimension,
					Value:     d.value,
					Evidence:  kw,
				})
				break // one signal per detector set
			}
		}
	}
	return signals
}

// DetectFeedback analyzes explicit feedback (corrections, ratings).
func DetectFeedback(message string) []Signal {
	lower := strings.ToLower(message)
	var signals []Signal

	// Corrections indicate the user is refining style
	if strings.Contains(lower, "too long") || strings.Contains(lower, "shorter") {
		signals = append(signals, Signal{Dimension: "brevity", Value: "short", Evidence: "too long/shorter"})
	}
	if strings.Contains(lower, "too short") || strings.Contains(lower, "more detail") {
		signals = append(signals, Signal{Dimension: "brevity", Value: "detailed", Evidence: "too short/more detail"})
	}
	if strings.Contains(lower, "too formal") || strings.Contains(lower, "relax") {
		signals = append(signals, Signal{Dimension: "tone", Value: "friendly", Evidence: "too formal/relax"})
	}
	if strings.Contains(lower, "too casual") || strings.Contains(lower, "more professional") {
		signals = append(signals, Signal{Dimension: "tone", Value: "formal", Evidence: "too casual/professional"})
	}

	return signals
}
