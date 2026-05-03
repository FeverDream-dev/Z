package server

import (
	"fmt"
	"net/http"
	"time"
)

// ValidationResult records the outcome of a UI validation check.
type ValidationResult struct {
	Name      string    `json:"name"`
	Passed    bool      `json:"passed"`
	URL       string    `json:"url"`
	StatusCode int      `json:"status_code,omitempty"`
	Error     string    `json:"error,omitempty"`
	Duration  time.Duration `json:"duration"`
	CheckedAt time.Time `json:"checked_at"`
}

// Validate runs a series of health checks against the local server.
func (s *Server) Validate(baseURL string) []ValidationResult {
	if baseURL == "" {
		baseURL = "http://" + s.addr
	}

	checks := []struct {
		name string
		path string
	}{
		{"health", "/health"},
		{"dashboard", "/"},
		{"chat_page", "/chat"},
		{"agents_api", "/api/agents"},
		{"providers_api", "/api/providers"},
	}

	client := &http.Client{Timeout: 5 * time.Second}
	var results []ValidationResult

	for _, c := range checks {
		start := time.Now()
		url := baseURL + c.path
		resp, err := client.Get(url)
		duration := time.Since(start)

		res := ValidationResult{
			Name:      c.name,
			URL:       url,
			Duration:  duration,
			CheckedAt: time.Now(),
		}

		if err != nil {
			res.Passed = false
			res.Error = err.Error()
		} else {
			res.StatusCode = resp.StatusCode
			res.Passed = resp.StatusCode == http.StatusOK
			resp.Body.Close()
		}

		results = append(results, res)
	}

	return results
}

// FormatReport produces a human-readable validation report.
func FormatReport(results []ValidationResult) string {
	var passed, failed int
	for _, r := range results {
		if r.Passed {
			passed++
		} else {
			failed++
		}
	}

	var out string
	out += fmt.Sprintf("UI Validation Report — %s\n", time.Now().Format("2006-01-02 15:04:05"))
	out += fmt.Sprintf("Passed: %d | Failed: %d\n\n", passed, failed)

	for _, r := range results {
		status := "PASS"
		if !r.Passed {
			status = "FAIL"
		}
		out += fmt.Sprintf("[%s] %s (%s) — %d in %v\n", status, r.Name, r.URL, r.StatusCode, r.Duration)
		if r.Error != "" {
			out += fmt.Sprintf("  Error: %s\n", r.Error)
		}
	}
	return out
}
