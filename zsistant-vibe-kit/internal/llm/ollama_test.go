package llm

import (
    "testing"
)

func TestNewOllamaProvider(t *testing.T) {
    p := NewOllamaProvider("test-key", "")
    if p == nil {
        t.Fatalf("expected provider, got nil")
    }
    if p.baseURL != "https://ollama.com" {
        t.Fatalf("unexpected baseURL: %s", p.baseURL)
    }
    if p.model != "gemma3:4b" {
        t.Fatalf("expected default model gemma3:4b, got %s", p.model)
    }
    if p.apiKey != "test-key" {
        t.Fatalf("unexpected apiKey: %s", p.apiKey)
    }
    h := p.Health()
    if h.Name != "ollama" {
        t.Fatalf("expected health name 'ollama', got %s", h.Name)
    }
    if h.Status != "healthy" {
        t.Fatalf("expected healthy status, got %s", h.Status)
    }
    if h.LastError != "" {
        t.Fatalf("expected empty last_error, got %s", h.LastError)
    }
}

func TestNewOllamaProviderWithKey(t *testing.T) {
    p := NewOllamaProvider("abcd1234", "gemma3:4b")
    if p == nil {
        t.Fatalf("expected provider, got nil")
    }
    h := p.Health()
    if h.Status != "healthy" {
        t.Fatalf("expected healthy status with key, got %s", h.Status)
    }
    if h.LastError != "" {
        t.Fatalf("expected empty last_error with key, got %s", h.LastError)
    }
}

func TestNewOllamaProviderNoKey(t *testing.T) {
    p := NewOllamaProvider("", "gemma3:4b")
    if p == nil {
        t.Fatalf("expected provider, got nil")
    }
    h := p.Health()
    if h.Status != "unhealthy" {
        t.Fatalf("expected unhealthy status without key, got %s", h.Status)
    }
    if h.LastError != "no API key" {
        t.Fatalf("expected LastError 'no API key', got %q", h.LastError)
    }
}

func TestOllamaProviderHealth(t *testing.T) {
    p := NewOllamaProvider("key", "gemma3:4b")
    // Health should reflect current state and be accessible
    h := p.Health()
    if h.Name != "ollama" {
        t.Fatalf("health name mismatch: %s", h.Name)
    }
    // After creation with key, should be healthy
    if h.Status != "healthy" {
        t.Fatalf("expected healthy, got %s", h.Status)
    }
}
