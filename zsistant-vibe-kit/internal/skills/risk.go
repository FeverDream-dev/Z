package skills

import (
	"fmt"
	"strings"
)

// RiskLevel categorizes the danger of a skill.
type RiskLevel string

const (
	RiskLow    RiskLevel = "low"
	RiskMedium RiskLevel = "medium"
	RiskHigh   RiskLevel = "high"
)

// RiskPattern describes a detectable dangerous pattern.
type RiskPattern struct {
	Name        string
	Keywords    []string
	Weight      int
	Category    string
	Description string
}

// riskPatterns is the ordered list of patterns to scan for.
var riskPatterns = []RiskPattern{
	{
		Name:        "shell_execution",
		Keywords:    []string{"bash", "sh", "cmd", "powershell", "exec", "system(", "subprocess", "spawn", "shell"},
		Weight:      3,
		Category:    "execution",
		Description: "Executes shell commands or system calls",
	},
	{
		Name:        "dependency_install",
		Keywords:    []string{"npm install", "pip install", "apt-get", "brew install", "go get", "cargo install", "gem install"},
		Weight:      2,
		Category:    "execution",
		Description: "Installs external dependencies",
	},
	{
		Name:        "network_request",
		Keywords:    []string{"curl", "wget", "fetch(", "http.request", "axios", "requests.get", "urllib", "api call"},
		Weight:      2,
		Category:    "network",
		Description: "Makes external network requests",
	},
	{
		Name:        "file_deletion",
		Keywords:    []string{"rm -rf", "rmdir", "del ", "remove(", "unlink(", "delete(", "shred", "wipe"},
		Weight:      3,
		Category:    "filesystem",
		Description: "Deletes or destroys files",
	},
	{
		Name:        "credential_access",
		Keywords:    []string{"password", "secret", "token", "private key", "wallet", "ssh key", "credential", "api_key", "auth"},
		Weight:      3,
		Category:    "privacy",
		Description: "Accesses sensitive credentials or keys",
	},
	{
		Name:        "code_evaluation",
		Keywords:    []string{"eval(", "exec(", "compile(", "dynamic", "runtime.exec", "Function(", "new Function"},
		Weight:      3,
		Category:    "execution",
		Description: "Evaluates or compiles code dynamically",
	},
	{
		Name:        "browser_automation",
		Keywords:    []string{"puppeteer", "playwright", "selenium", "chromedp", "browser.launch", "chrome"},
		Weight:      2,
		Category:    "automation",
		Description: "Controls browser or UI automation",
	},
	{
		Name:        "home_access",
		Keywords:    []string{"$HOME", "~/.", "%USERPROFILE%", "/home/", "os.homedir", "user.home"},
		Weight:      1,
		Category:    "filesystem",
		Description: "Accesses user's home directory",
	},
	{
		Name:        "elevated_privilege",
		Keywords:    []string{"sudo", "admin", "root", "elevated", "runas", "setuid"},
		Weight:      3,
		Category:    "execution",
		Description: "Requires elevated privileges",
	},
	{
		Name:        "obfuscation",
		Keywords:    []string{"base64", "encoded", "obfuscated", "hex decode", "atob", "btoa", "rot13"},
		Weight:      2,
		Category:    "obfuscation",
		Description: "Contains obfuscated or encoded content",
	},
	{
		Name:        "download_script",
		Keywords:    []string{"curl | bash", "wget | sh", "download and run", "pipe to shell", "invoke-expression"},
		Weight:      4,
		Category:    "execution",
		Description: "Downloads and executes remote scripts",
	},
}

// Match represents a single detected risk pattern instance.
type Match struct {
	Pattern RiskPattern
	Line    int
	Snippet string
}

// Score computes the total risk score and level from matches.
func Score(matches []Match) (int, RiskLevel) {
	total := 0
	for _, m := range matches {
		total += m.Pattern.Weight
	}
	switch {
	case total < 2:
		return total, RiskLow
	case total <= 5:
		return total, RiskMedium
	default:
		return total, RiskHigh
	}
}

// RiskReport summarizes the analysis of a skill.
type RiskReport struct {
	SkillName    string
	Score        int
	Level        RiskLevel
	Matches      []Match
	Summary      string
	Translations []string
}

// String returns a human-readable risk report.
func (r *RiskReport) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "Skill: %s\n", r.SkillName)
	fmt.Fprintf(&b, "Risk Score: %d (%s)\n", r.Score, r.Level)
	fmt.Fprintf(&b, "Matches: %d\n", len(r.Matches))
	for _, m := range r.Matches {
		fmt.Fprintf(&b, "  - %s (%s, +%d): %s\n", m.Pattern.Name, m.Pattern.Category, m.Pattern.Weight, m.Pattern.Description)
	}
	if r.Summary != "" {
		fmt.Fprintf(&b, "Summary: %s\n", r.Summary)
	}
	if len(r.Translations) > 0 {
		fmt.Fprintf(&b, "Translation Plan:\n")
		for _, tr := range r.Translations {
			fmt.Fprintf(&b, "  - %s\n", tr)
		}
	}
	return b.String()
}
