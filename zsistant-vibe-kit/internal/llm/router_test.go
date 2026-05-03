package llm

import (
	"testing"
	"time"
)

func TestRouterSuccess(t *testing.T) {
	r := NewRouter(5 * time.Second)
	r.Register(TaskCheap, NewMockProvider())

	resp, provider, err := r.Route(TaskCheap, "hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != "Echo: hello" {
		t.Fatalf("unexpected response: %s", resp)
	}
	if provider != "mock" {
		t.Fatalf("unexpected provider: %s", provider)
	}
}

func TestRouterFallbackOnTimeout(t *testing.T) {
	r := NewRouter(100 * time.Millisecond)
	r.Register(TaskCheap, NewTimeoutProvider("slow"))
	r.Register(TaskCheap, NewMockProvider())

	resp, provider, err := r.Route(TaskCheap, "hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != "Echo: hello" {
		t.Fatalf("unexpected response: %s", resp)
	}
	if provider != "mock" {
		t.Fatalf("expected fallback to mock, got: %s", provider)
	}
}

func TestRouterFallbackOnError(t *testing.T) {
	r := NewRouter(5 * time.Second)
	r.Register(TaskCheap, NewRateLimitProvider("limited"))
	r.Register(TaskCheap, NewMockProvider())

	resp, provider, err := r.Route(TaskCheap, "hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != "Echo: hello" {
		t.Fatalf("unexpected response: %s", resp)
	}
	if provider != "mock" {
		t.Fatalf("expected fallback to mock, got: %s", provider)
	}
}

func TestRouterAllFail(t *testing.T) {
	r := NewRouter(100 * time.Millisecond)
	r.Register(TaskCheap, NewTimeoutProvider("slow"))
	r.Register(TaskCheap, NewRateLimitProvider("limited"))

	_, _, err := r.Route(TaskCheap, "hello")
	if err == nil {
		t.Fatal("expected error when all providers fail")
	}
}

func TestRouterNoProviders(t *testing.T) {
	r := NewRouter(5 * time.Second)
	_, _, err := r.Route(TaskStrong, "hello")
	if err == nil {
		t.Fatal("expected error when no providers registered")
	}
}

func TestRouterHealth(t *testing.T) {
	r := NewRouter(5 * time.Second)
	r.Register(TaskCheap, NewMockProvider())
	r.Register(TaskStrong, NewTimeoutProvider("slow"))

	health := r.Health()
	if len(health) != 2 {
		t.Fatalf("expected 2 health entries, got %d", len(health))
	}

	hasMock := false
	hasSlow := false
	for _, h := range health {
		if h.Name == "mock" && h.Status == "healthy" {
			hasMock = true
		}
		if h.Name == "slow" && h.Status == "degraded" {
			hasSlow = true
		}
	}
	if !hasMock {
		t.Fatal("expected healthy mock provider")
	}
	if !hasSlow {
		t.Fatal("expected degraded slow provider")
	}
}

func TestRouterTaskTypeRouting(t *testing.T) {
	r := NewRouter(5 * time.Second)
	r.Register(TaskCoding, NewNamedMockProvider("coder", 0))

	resp, provider, err := r.Route(TaskCoding, "write a function")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if provider != "coder" {
		t.Fatalf("expected coder provider, got: %s", provider)
	}
	if resp != "Echo: write a function" {
		t.Fatalf("unexpected response: %s", resp)
	}

	// TaskCheap has no providers
	_, _, err = r.Route(TaskCheap, "hello")
	if err == nil {
		t.Fatal("expected error for unregistered task type")
	}
}
