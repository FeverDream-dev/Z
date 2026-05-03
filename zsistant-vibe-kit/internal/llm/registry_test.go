package llm

import (
	"testing"
)

func TestTotalProviderCount(t *testing.T) {
	count := TotalProviderCount()
	if count < 52 {
		t.Fatalf("expected at least 52 providers, got %d", count)
	}
}

func TestTotalModelCount(t *testing.T) {
	count := TotalModelCount()
	if count < 100 {
		t.Fatalf("expected at least 100 models, got %d", count)
	}
}

func TestFindProvider(t *testing.T) {
	p, err := FindProvider("openai")
	if err != nil {
		t.Fatalf("expected openai provider: %v", err)
	}
	if p.Name != "openai" {
		t.Fatalf("expected name openai, got %s", p.Name)
	}
	if len(p.Models) == 0 {
		t.Fatal("expected models for openai")
	}
}

func TestFindProviderNotFound(t *testing.T) {
	_, err := FindProvider("nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent provider")
	}
}

func TestFindModel(t *testing.T) {
	p, m, err := FindModel("gpt-4o")
	if err != nil {
		t.Fatalf("expected gpt-4o model: %v", err)
	}
	if p.Name != "openai" {
		t.Fatalf("expected openai provider, got %s", p.Name)
	}
	if m.ID != "gpt-4o" {
		t.Fatalf("expected gpt-4o, got %s", m.ID)
	}
}

func TestProvidersByTier(t *testing.T) {
	cheap := ProvidersByTier("cheap")
	if len(cheap) == 0 {
		t.Fatal("expected some cheap-tier providers")
	}
	strong := ProvidersByTier("strong")
	if len(strong) == 0 {
		t.Fatal("expected some strong-tier providers")
	}
}
