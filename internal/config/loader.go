package config

import (
    "errors"
    "io/ioutil"
    "os"
    "path/filepath"
    "strings"
    "gopkg.in/yaml.v3"
)

// DefaultPath resolves the default config path at ~/.zazi/config.yaml
func DefaultPath() (string, error) {
    home, err := os.UserHomeDir()
    if err != nil {
        return "", err
    }
    cfgPath := filepath.Join(home, ".zazi", "config.yaml")
    return cfgPath, nil
}

// Load reads YAML config from disk and applies defaults for missing fields.
func Load(path string) (*Config, error) {
    cfg := DefaultConfig()
    if path == "" {
        // resolve default path if not provided
        p, err := DefaultPath()
        if err != nil {
            return nil, err
        }
        path = p
    }
    b, err := ioutil.ReadFile(path)
    if err != nil {
        if errors.Is(err, os.ErrNotExist) {
            // Return defaults if file doesn't exist
            return &cfg, nil
        }
        return nil, err
    }
    if err := yaml.Unmarshal(b, &cfg); err != nil {
        return nil, err
    }
    // Apply defaults for missing fields
    MergeDefaults(&cfg)
	// Ensure data path normalization
	cfg.DataPath = ExpandPath(cfg.DataPath)
	// Normalization for nested slices/maps is left simple for MVP
	return &cfg, nil
}

// ExpandPath expands a path that may start with '~' to the user's home directory.
func ExpandPath(p string) string {
	if strings.HasPrefix(p, "~") {
		home, err := os.UserHomeDir()
		if err == nil {
			return filepath.Join(home, p[1:])
		}
	}
	return p
}

// Save writes YAML config to disk with a header comment and ensures dirs exist.
func Save(path string, cfg *Config) error {
    if cfg == nil {
        tmp := DefaultConfig()
        cfg = &tmp
    }
    if path == "" {
        p, err := DefaultPath()
        if err != nil {
            return err
        }
        path = p
    }
    // Ensure parent directories exist
    if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
        return err
    }
    // Marshal with indentation for readability
    out, err := yaml.Marshal(cfg)
    if err != nil {
        return err
    }
    // Add a small header comment
    header := []byte("# Zsistant configuration file. Do not edit unless you know what you are doing.\n")
    data := append(header, out...)
    return os.WriteFile(path, data, 0o644)
}

// EnsureDirs creates subdirectories under the configured data path.
func EnsureDirs(cfg *Config) error {
    if cfg == nil {
        tmp := DefaultConfig()
        cfg = &tmp
    }
    base := expandPath(cfg.DataPath)
    dirs := []string{"agents", "logs", "cache", "backups"}
    for _, d := range dirs {
        p := filepath.Join(base, d)
        if err := os.MkdirAll(p, 0o755); err != nil {
            return err
        }
    }
    return nil
}

// private helpers
func expandPath(p string) string {
    // Expand ~ to user home
    if strings.HasPrefix(p, "~") {
        home, err := os.UserHomeDir()
        if err == nil {
            return filepath.Join(home, p[1:])
        }
    }
    return p
}
