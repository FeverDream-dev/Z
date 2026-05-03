package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestValidateHealth(t *testing.T) {
	dir := t.TempDir()
	s := NewWithAssets("localhost:0", dir, "assets/brand")

	ts := httptest.NewServer(s.mux)
	defer ts.Close()

	results := s.Validate(ts.URL)

	healthFound := false
	for _, r := range results {
		if r.Name == "health" {
			healthFound = true
			if !r.Passed {
				t.Fatalf("health check failed: %s", r.Error)
			}
			if r.StatusCode != http.StatusOK {
				t.Fatalf("expected 200, got %d", r.StatusCode)
			}
		}
	}
	if !healthFound {
		t.Fatal("health check not found in results")
	}
}

func TestValidateDashboard(t *testing.T) {
	dir := t.TempDir()
	s := NewWithAssets("localhost:0", dir, "assets/brand")

	ts := httptest.NewServer(s.mux)
	defer ts.Close()

	results := s.Validate(ts.URL)

	var dashboard *ValidationResult
	for i := range results {
		if results[i].Name == "dashboard" {
			dashboard = &results[i]
			break
		}
	}
	if dashboard == nil {
		t.Fatal("dashboard check not found")
	}
	if !dashboard.Passed {
		t.Fatalf("dashboard check failed: %s", dashboard.Error)
	}
}

func TestFormatReport(t *testing.T) {
	results := []ValidationResult{
		{Name: "health", Passed: true, URL: "http://localhost/health", StatusCode: 200},
		{Name: "dashboard", Passed: false, URL: "http://localhost/", StatusCode: 0, Error: "connection refused"},
	}
	report := FormatReport(results)
	if !strings.Contains(report, "Passed: 1") {
		t.Fatal("expected passed count")
	}
	if !strings.Contains(report, "Failed: 1") {
		t.Fatal("expected failed count")
	}
	if !strings.Contains(report, "FAIL") {
		t.Fatal("expected FAIL in report")
	}
}
