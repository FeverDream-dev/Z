package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/FeverDream-dev/zsistant/internal/agents"
)

func TestHealthEndpoint(t *testing.T) {
	srv := New("127.0.0.1:0", t.TempDir())
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()
	srv.mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var body map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if body["status"] != "ok" {
		t.Fatalf("expected status ok, got %s", body["status"])
	}
}

func TestDashboardPage(t *testing.T) {
	srv := New("127.0.0.1:0", t.TempDir())
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	srv.mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "<title>Zsistant</title>") {
		t.Fatalf("expected Zsistant title in response")
	}
	if !strings.Contains(rr.Body.String(), "conversationList") {
		t.Fatalf("expected conversationList in response")
	}
}

func TestChatPage(t *testing.T) {
	srv := New("127.0.0.1:0", t.TempDir())
	req := httptest.NewRequest(http.MethodGet, "/chat", nil)
	rr := httptest.NewRecorder()
	srv.mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "<title>Zsistant</title>") {
		t.Fatalf("expected Zsistant title in response")
	}
	if !strings.Contains(rr.Body.String(), "messages") {
		t.Fatalf("expected messages container in response")
	}
}

func TestAgentsEndpoint(t *testing.T) {
	tmp := t.TempDir()
	base := filepath.Join(tmp, ".zazi")
	reg := agents.NewRegistry(base)
	if _, err := reg.Create("agent1", "Agent One", "tester"); err != nil {
		t.Fatalf("create agent: %v", err)
	}

	srv := New("127.0.0.1:0", base)
	req := httptest.NewRequest(http.MethodGet, "/api/agents", nil)
	rr := httptest.NewRecorder()
	srv.mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var list []agents.Agent
	if err := json.Unmarshal(rr.Body.Bytes(), &list); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 agent, got %d", len(list))
	}
	if list[0].ID != "agent1" {
		t.Fatalf("expected agent1, got %s", list[0].ID)
	}
}

func TestAgentsEndpointCreate(t *testing.T) {
	tmp := t.TempDir()
	base := filepath.Join(tmp, ".zazi")
	srv := New("127.0.0.1:0", base)

	payload := map[string]string{"id": "new-agent", "name": "New Agent", "role": "tester"}
	b, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/agents", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	srv.mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestChatAPI(t *testing.T) {
	tmp := t.TempDir()
	base := filepath.Join(tmp, ".zazi")
	reg := agents.NewRegistry(base)
	if _, err := reg.Create("chat-agent", "Chat Agent", "tester"); err != nil {
		t.Fatalf("create agent: %v", err)
	}

	srv := New("127.0.0.1:0", base)
	payload := chatRequest{AgentID: "chat-agent", Message: "hello"}
	b, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/chat", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	srv.mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}
	var resp chatResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if resp.Status != "completed" {
		t.Fatalf("expected completed status, got %s", resp.Status)
	}
	if resp.Response == "" {
		t.Fatalf("expected non-empty response")
	}
	if resp.JobID == "" {
		t.Fatal("expected job_id to be set")
	}
}

func TestChatAPIMissingMessage(t *testing.T) {
	srv := New("127.0.0.1:0", t.TempDir())
	payload := map[string]string{"message": ""}
	b, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/chat", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	srv.mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestChatAPIAutoCreateAgent(t *testing.T) {
	srv := New("127.0.0.1:0", t.TempDir())
	payload := chatRequest{AgentID: "auto-agent", Message: "hello"}
	b, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/chat", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	srv.mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 (auto-create agent), got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestJobsAPI(t *testing.T) {
	tmp := t.TempDir()
	base := filepath.Join(tmp, ".zazi")
	reg := agents.NewRegistry(base)
	a, err := reg.Create("job-agent", "Job Agent", "tester")
	if err != nil {
		t.Fatalf("create agent: %v", err)
	}

	// Write a test event to audit.jsonl
	auditPath := filepath.Join(a.WorkspaceRoot, "audit.jsonl")
	evt := map[string]interface{}{
		"job_id":     "job-1",
		"agent_id":   "job-agent",
		"event_type": "job.created",
		"message":    "test",
	}
	f, _ := os.Create(auditPath)
	b, _ := json.Marshal(evt)
	f.Write(append(b, '\n'))
	f.Close()

	srv := New("127.0.0.1:0", base)
	req := httptest.NewRequest(http.MethodGet, "/api/jobs/job-agent", nil)
	rr := httptest.NewRecorder()
	srv.mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var events []map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &events); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0]["event_type"] != "job.created" {
		t.Fatalf("expected job.created event, got %v", events[0]["event_type"])
	}
}

func TestJobsAPIMethodNotAllowed(t *testing.T) {
	srv := New("127.0.0.1:0", t.TempDir())
	req := httptest.NewRequest(http.MethodPost, "/api/jobs/agent1", nil)
	rr := httptest.NewRecorder()
	srv.mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rr.Code)
	}
}

func TestServeBrandAsset(t *testing.T) {
	assetDir := t.TempDir()
	dummyFile := filepath.Join(assetDir, "test-icon.png")
	if err := os.WriteFile(dummyFile, []byte("PNG"), 0644); err != nil {
		t.Fatalf("write dummy: %v", err)
	}

	srv := NewWithAssets("127.0.0.1:0", t.TempDir(), assetDir)
	req := httptest.NewRequest(http.MethodGet, "/assets/brand/test-icon.png", nil)
	rr := httptest.NewRecorder()
	srv.mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 for asset, got %d", rr.Code)
	}
}
