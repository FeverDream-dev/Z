package config

import (
    "os"
    "path/filepath"
    "testing"
)

// TestLoadMissingFileReturnsDefaults verifies that loading a non-existent file
// yields default configuration values.
func TestLoadMissingFileReturnsDefaults(t *testing.T) {
    // Create a temp dir and point to a non-existent file inside it
    dir := t.TempDir()
    path := filepath.Join(dir, "config.yaml")

    // Ensure the file does not exist
    os.Remove(path)

    cfg, err := Load(path)
    if err != nil {
        t.Fatalf("Load should not fail for missing file: %v", err)
    }
    if cfg.DataPath == "" || cfg.LogLevel == "" || cfg.ServerPort == 0 {
        t.Fatalf("Defaults not applied: %#v", cfg)
    }
}

// TestSaveAndLoadRoundTrip ensures a config can be saved and then loaded back
// with the same values (modulo non-exported fields).
func TestSaveAndLoadRoundTrip(t *testing.T) {
    dir := t.TempDir()
    path := filepath.Join(dir, "config.yaml")

    cfg := NewDefault()
    cfg.DataPath = "/tmp/zsistant"
    cfg.LogLevel = "debug"
    cfg.ServerPort = 9090
    cfg.Secrets = map[string]string{"token": "s3cr3t"}

    if err := Save(path, cfg); err != nil {
        t.Fatalf("Save failed: %v", err)
    }

    loaded, err := Load(path)
    if err != nil {
        t.Fatalf("Load failed: %v", err)
    }
    if loaded.DataPath != cfg.DataPath || loaded.LogLevel != cfg.LogLevel || loaded.ServerPort != cfg.ServerPort {
        t.Fatalf("Loaded config does not match saved: %#v vs %#v", loaded, cfg)
    }
    // Secrets should be loaded as well but not printed; ensure key exists
    if loaded.Secrets == nil {
        t.Fatalf("Secrets should be preserved through round-trip (nil found)")
    }
}

// TestEnsureDirs creates the expected subdirectories under the provided DataPath.
func TestEnsureDirsCreatesSubdirs(t *testing.T) {
    dir := t.TempDir()
    cfg := NewDefault()
    cfg.DataPath = dir // use temp dir as base
    if err := EnsureDirs(cfg); err != nil {
        t.Fatalf("EnsureDirs failed: %v", err)
    }
    subdirs := []string{"agents", "logs", "cache", "backups"}
    for _, s := range subdirs {
        p := filepath.Join(dir, s)
        info, err := os.Stat(p)
        if err != nil {
            t.Fatalf("Expected subdir %s to exist, error: %v", p, err)
        }
        if !info.IsDir() {
            t.Fatalf("Expected %s to be a directory", p)
        }
    }
}

// TestStringRedactsSecrets ensures that the String() method does not leak secrets.
func TestStringRedactsSecrets(t *testing.T) {
    cfg := NewDefault()
    cfg.Secrets = map[string]string{"password": "hunter2"}
    s := cfg.String()
    if !contains(s, "Secrets:[redacted]") {
        t.Fatalf("Secrets were not redacted in String(): %s", s)
    }
}

// contains is a tiny helper for substring checks.
func contains(s, sub string) bool {
    return len(s) >= len(sub) && (stringContains(s, sub))
}

func stringContains(s, substr string) bool {
    return len(substr) == 0 || (bytesIndex(s, substr) >= 0)
}

func bytesIndex(s, substr string) int {
    return indexOf([]byte(s), []byte(substr))
}

func indexOf(a, b []byte) int {
    // simple implementation to avoid importing bytes
    la, lb := len(a), len(b)
    if lb == 0 {
        return 0
    }
    for i := 0; i+lb <= la; i++ {
        if string(a[i:i+lb]) == string(b) {
            return i
        }
    }
    return -1
}
