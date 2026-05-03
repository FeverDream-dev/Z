package llm

import (
	"strings"
	"testing"
)

func TestMockProviderComplete(t *testing.T) {
	p := NewMockProvider()
	resp, err := p.Complete("hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(resp, "Echo: ") {
		t.Fatalf("expected echo prefix, got: %s", resp)
	}
	if !strings.Contains(resp, "hello") {
		t.Fatalf("expected response to contain prompt, got: %s", resp)
	}
}
