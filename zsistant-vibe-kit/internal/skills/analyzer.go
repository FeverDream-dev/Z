package skills

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Analyzer reads and analyzes skill packages without executing them.
type Analyzer struct{}

// NewAnalyzer creates a new skill analyzer.
func NewAnalyzer() *Analyzer {
	return &Analyzer{}
}

// AnalyzeFolder reads a skill directory and produces a risk report.
func (a *Analyzer) AnalyzeFolder(path string) (*RiskReport, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("reading skill folder: %w", err)
	}

	report := &RiskReport{
		SkillName: filepath.Base(path),
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := strings.ToLower(entry.Name())
		if !strings.HasSuffix(name, ".md") && !strings.HasSuffix(name, ".txt") && !strings.HasSuffix(name, ".json") {
			continue
		}
		fpath := filepath.Join(path, entry.Name())
		data, err := os.ReadFile(fpath)
		if err != nil {
			continue
		}
		matches := a.Scan(string(data))
		report.Matches = append(report.Matches, matches...)
	}

	report.Score, report.Level = Score(report.Matches)
	report.Summary = a.summarize(report)
	report.Translations = a.planTranslations(report)
	return report, nil
}

// AnalyzeText analyzes a raw skill text (e.g., pasted Markdown).
func (a *Analyzer) AnalyzeText(name, content string) *RiskReport {
	report := &RiskReport{
		SkillName: name,
		Matches:   a.Scan(content),
	}
	report.Score, report.Level = Score(report.Matches)
	report.Summary = a.summarize(report)
	report.Translations = a.planTranslations(report)
	return report
}

// Scan searches content for risk patterns.
func (a *Analyzer) Scan(content string) []Match {
	var matches []Match
	lines := strings.Split(content, "\n")
	for lineNum, line := range lines {
		lower := strings.ToLower(line)
		for _, p := range riskPatterns {
			for _, kw := range p.Keywords {
				if strings.Contains(lower, strings.ToLower(kw)) {
					matches = append(matches, Match{
						Pattern: p,
						Line:    lineNum + 1,
						Snippet: strings.TrimSpace(line),
					})
					break // one match per pattern per line
				}
			}
		}
	}
	return matches
}

func (a *Analyzer) summarize(r *RiskReport) string {
	switch r.Level {
	case RiskLow:
		return fmt.Sprintf("Skill '%s' appears safe. No dangerous patterns detected.", r.SkillName)
	case RiskMedium:
		return fmt.Sprintf("Skill '%s' has moderate risk. Review network/filesystem access before use.", r.SkillName)
	case RiskHigh:
		return fmt.Sprintf("Skill '%s' is HIGH RISK. Contains execution, credential, or obfuscation patterns. Do not run without review.", r.SkillName)
	}
	return ""
}

func (a *Analyzer) planTranslations(r *RiskReport) []string {
	var plans []string
	for _, m := range r.Matches {
		switch m.Pattern.Category {
		case "execution":
			plans = append(plans, fmt.Sprintf("Replace '%s' with Zsistant job template or sandboxed command", m.Pattern.Name))
		case "network":
			plans = append(plans, fmt.Sprintf("Map '%s' to channel adapter or external API integration", m.Pattern.Name))
		case "filesystem":
			plans = append(plans, fmt.Sprintf("Map '%s' to scoped file operations within agent workspace", m.Pattern.Name))
		case "automation":
			plans = append(plans, fmt.Sprintf("Map '%s' to browser workflow or MCP adapter", m.Pattern.Name))
		case "privacy":
			plans = append(plans, fmt.Sprintf("Flag '%s' for manual review - credential access is not auto-translated", m.Pattern.Name))
		case "obfuscation":
			plans = append(plans, fmt.Sprintf("Reject '%s' - obfuscated skills are not translated", m.Pattern.Name))
		}
	}
	// Deduplicate
	seen := make(map[string]bool)
	var unique []string
	for _, p := range plans {
		if !seen[p] {
			seen[p] = true
			unique = append(unique, p)
		}
	}
	return unique
}
