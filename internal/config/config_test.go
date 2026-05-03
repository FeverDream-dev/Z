package config

import (
    "os"
    "path/filepath"
    "testing"
    "strings"
)

func TestLoadMissingFileReturnsDefaults(t *testing.T) {
    // Use a temp path that does not exist
    tmpDir := t.TempDir()
    path := filepath.Join(tmpDir, "missing.yaml")
    cfg, err := Load(path)
    if err != nil {
        t.Fatalf("Load returned error: %v", err)
    }
    if cfg == nil {
        t.Fatalf("Load returned nil config")
    }
    // defaults
    def := DefaultConfig()
    if cfg.DataPath != def.DataPath || cfg.LogLevel != def.LogLevel || cfg.ServerPort != def.ServerPort {
        t.Fatalf("expected defaults, got %+v (defaults=%+v)", cfg, def)
    }
}

func TestSaveAndLoadRoundTrip(t *testing.T) {
    tmpDir := t.TempDir()
    path := filepath.Join(tmpDir, "config.yaml")
    cfg := DefaultConfig()
    cfg.DataPath = "/tmp/zs-data"
    cfg.LogLevel = "debug"
    cfg.ServerPort = 9090
    cfg.Providers = []Provider{{Name: "test", Config: map[string]string{"a": "b"}}}
    cfg.Secrets = Secrets{ApiKey: "secret-key", Token: "secret-token"}

    if err := Save(path, &cfg); err != nil {
        t.Fatalf("Save failed: %v", err)
    }
    loaded, err := Load(path)
    if err != nil {
        t.Fatalf("Load failed: %v", err)
    }
    if loaded.DataPath != cfg.DataPath || loaded.LogLevel != cfg.LogLevel || loaded.ServerPort != cfg.ServerPort {
        t.Fatalf("mismatch after load: got %+v want %+v", loaded, cfg)
    }
    if len(loaded.Providers) != 1 || loaded.Providers[0].Name != cfg.Providers[0].Name {
        t.Fatalf("provider mismatch: %+v", loaded.Providers)
    }
    if loaded.Secrets.ApiKey != cfg.Secrets.ApiKey || loaded.Secrets.Token != cfg.Secrets.Token {
        t.Fatalf("secrets mismatch: %+v vs %+v", loaded.Secrets, cfg.Secrets)
    }
}

func TestEnsureDirsCreatesSubdirectories(t *testing.T) {
    cfg := DefaultConfig()
    cfg.DataPath = t.TempDir()
    if err := EnsureDirs(&cfg); err != nil {
        t.Fatalf("EnsureDirs failed: %v", err)
    }
    base := cfg.DataPath
    dirs := []string{"agents", "logs", "cache", "backups"}
    for _, d := range dirs {
        p := filepath.Join(base, d)
        fi, err := os.Stat(p)
        if err != nil || !fi.IsDir() {
            t.Fatalf("expected dir %s to exist", p)
        }
    }
}

func TestStringRedactsSecrets(t *testing.T) {
    cfg := DefaultConfig()
    cfg.Secrets = Secrets{ApiKey: "abc123", Token: "tok987"}
    s := cfg.String()
    if strings.Contains(s, "abc123") || strings.Contains(s, "tok987") {
        t.Fatalf("secrets leaked in string: %s", s)
    }
    if !strings.Contains(s, "[redacted]") {
        t.Fatalf("expected redaction in string, got: %s", s)
    }
}
