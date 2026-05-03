package main

import (
  "bytes"
  "io"
  "os"
  "path/filepath"
  "strings"
  "testing"
)

// captureOutput captures stdout produced by a function call.
func captureOutput(f func()) string {
  old := os.Stdout
  r, w, _ := os.Pipe()
  os.Stdout = w
  f()
  w.Close()
  os.Stdout = old
  var buf bytes.Buffer
  io.Copy(&buf, r)
  return buf.String()
}

func TestVersionOutput(t *testing.T) {
  out := captureOutput(func() {
    run([]string{"version"})
  })
  if !strings.Contains(out, "v0.0.1-dev") {
    t.Fatalf("expected version output to contain 'v0.0.1-dev', got: %q", out)
  }
}

func TestDoctorOutputContainsStatuses(t *testing.T) {
  // Use a temp HOME to avoid touching real home directory
  tmpdir, err := os.MkdirTemp("", "zazi-doctor-test-")
  if err != nil {
    t.Fatal(err)
  }
  defer os.RemoveAll(tmpdir)
  oldHome := os.Getenv("HOME")
  os.Setenv("HOME", tmpdir)
  defer os.Setenv("HOME", oldHome)

  out := captureOutput(func() {
    run([]string{"doctor"})
  })
  if !strings.Contains(out, "Config file:") {
    t.Fatalf("expected doctor output to include Config file, got: %q", out)
  }
  if !strings.Contains(out, "Data path:") {
    t.Fatalf("expected doctor output to include Data path, got: %q", out)
  }
  // Should mention statuses like OK or MISSING
  if !strings.Contains(out, "OK") && !strings.Contains(out, "MISSING") && !strings.Contains(out, "ERROR") {
    t.Fatalf("expected doctor output to include a status label (OK/MISSING/ERROR), got: %q", out)
  }
}

func TestHelpOutputContainsUsage(t *testing.T) {
  out := captureOutput(func() {
    run([]string{"--help"})
  })
  if !strings.Contains(out, "Usage:") {
    t.Fatalf("expected help usage, got: %q", out)
  }
  if !strings.Contains(out, "zazi version") {
    t.Fatalf("expected help to include version command, got: %q", out)
  }
}

func TestInitCreatesConfigAndDirs(t *testing.T) {
  tmpdir, err := os.MkdirTemp("", "zazi-init-test-")
  if err != nil {
    t.Fatal(err)
  }
  defer os.RemoveAll(tmpdir)
  oldHome := os.Getenv("HOME")
  os.Setenv("HOME", tmpdir)
  defer os.Setenv("HOME", oldHome)

  if err := runInit(); err != nil {
    t.Fatalf("init failed: %v", err)
  }
  base := filepath.Join(tmpdir, ".zazi")
  // Check base and subdirectories exist
  for _, d := range []string{base, filepath.Join(base, "agents"), filepath.Join(base, "logs"), filepath.Join(base, "cache"), filepath.Join(base, "backups")} {
    if stat, err := os.Stat(d); err != nil || !stat.IsDir() {
      t.Fatalf("expected directory %s to exist, err=%v", d, err)
    }
  }
  cfg := filepath.Join(base, "config.yaml")
  if stat, err := os.Stat(cfg); err != nil || stat.IsDir() {
    t.Fatalf("expected config.yaml to exist at %s", cfg)
  }
  // Validate its contents contain the new skeleton header (header written by Save)
  b, err := os.ReadFile(cfg)
  if err != nil {
    t.Fatalf("reading config.yaml: %v", err)
  }
  if !strings.Contains(string(b), "zsistant config") {
    t.Fatalf("expected config.yaml to contain skeleton header, got: %q", string(b))
  }
}

func TestReleasePrepareRequiresVersion(t *testing.T) {
  code := runReleasePrepare([]string{})
  if code == 0 {
    t.Fatal("expected non-zero exit without version")
  }
}

func TestReleasePublishRequiresApproval(t *testing.T) {
  code := runReleasePublish([]string{"--version=v1.0.0"})
  if code == 0 {
    t.Fatal("expected non-zero exit without approval")
  }
}

func TestReleasePublishRequiresLongApproval(t *testing.T) {
  code := runReleasePublish([]string{"--version=v1.0.0", "--approve=short"})
  if code == 0 {
    t.Fatal("expected non-zero exit with short approval")
  }
}

func TestReleasePrepareOutputsChecklist(t *testing.T) {
  out := captureOutput(func() {
    runReleasePrepare([]string{"--version=v1.0.0"})
  })
  if !strings.Contains(out, "Release Checklist:") {
    t.Fatalf("expected checklist in output, got: %q", out)
  }
  if !strings.Contains(out, "v1.0.0") {
    t.Fatalf("expected version in output, got: %q", out)
  }
}
