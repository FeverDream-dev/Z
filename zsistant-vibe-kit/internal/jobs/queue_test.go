package jobs

import (
	"testing"
	"time"
)

func TestQueueEnqueueAndGet(t *testing.T) {
	dir := t.TempDir()
	q := NewQueue(dir)
	job := NewJob("agent1", "test", "do something")
	if err := q.Enqueue(job); err != nil {
		t.Fatalf("enqueue: %v", err)
	}
	got, err := q.Get(job.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.ID != job.ID {
		t.Fatalf("expected job ID %s, got %s", job.ID, got.ID)
	}
}

func TestQueueList(t *testing.T) {
	dir := t.TempDir()
	q := NewQueue(dir)
	j1 := NewJob("agent1", "test", "task 1")
	j2 := NewJob("agent1", "test", "task 2")
	_ = q.Enqueue(j1)
	_ = q.Enqueue(j2)

	list, err := q.List()
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2 jobs, got %d", len(list))
	}
}

func TestQueueRetry(t *testing.T) {
	dir := t.TempDir()
	q := NewQueue(dir)
	job := NewJob("agent1", "test", "do something")
	job.Status = "failed"
	_ = q.Enqueue(job)

	retried, err := q.Retry(job.ID)
	if err != nil {
		t.Fatalf("retry: %v", err)
	}
	if retried.Status != "retrying" {
		t.Fatalf("expected status retrying, got %s", retried.Status)
	}
	if retried.RetryCount != 1 {
		t.Fatalf("expected retry count 1, got %d", retried.RetryCount)
	}
}

func TestQueueRetryExhausted(t *testing.T) {
	dir := t.TempDir()
	q := NewQueue(dir)
	job := NewJob("agent1", "test", "do something")
	job.RetryCount = 3
	job.MaxRetries = 3
	_ = q.Enqueue(job)

	_, err := q.Retry(job.ID)
	if err == nil {
		t.Fatal("expected error when retries exhausted")
	}
}

func TestQueuePauseAndResume(t *testing.T) {
	dir := t.TempDir()
	q := NewQueue(dir)
	job := NewJob("agent1", "test", "do something")
	_ = q.Enqueue(job)

	paused, err := q.Pause(job.ID, "loop suspected")
	if err != nil {
		t.Fatalf("pause: %v", err)
	}
	if paused.Status != "paused" {
		t.Fatalf("expected status paused, got %s", paused.Status)
	}
	if paused.PausedReason != "loop suspected" {
		t.Fatalf("expected reason 'loop suspected', got %s", paused.PausedReason)
	}

	resumed, err := q.Resume(job.ID)
	if err != nil {
		t.Fatalf("resume: %v", err)
	}
	if resumed.Status != "queued" {
		t.Fatalf("expected status queued, got %s", resumed.Status)
	}
	if resumed.PausedReason != "" {
		t.Fatalf("expected empty reason, got %s", resumed.PausedReason)
	}
}

func TestQueueCancel(t *testing.T) {
	dir := t.TempDir()
	q := NewQueue(dir)
	job := NewJob("agent1", "test", "do something")
	_ = q.Enqueue(job)

	cancelled, err := q.Cancel(job.ID)
	if err != nil {
		t.Fatalf("cancel: %v", err)
	}
	if cancelled.Status != "cancelled" {
		t.Fatalf("expected status cancelled, got %s", cancelled.Status)
	}
}

func TestBackoff(t *testing.T) {
	cases := []struct {
		count    int
		expected time.Duration
	}{
		{0, 1 * time.Second},
		{1, 2 * time.Second},
		{2, 4 * time.Second},
		{3, 8 * time.Second},
		{10, 30 * time.Second}, // capped
	}
	for _, c := range cases {
		got := Backoff(c.count)
		if got != c.expected {
			t.Fatalf("backoff(%d) = %v, want %v", c.count, got, c.expected)
		}
	}
}

func TestLoopDetectorCheck(t *testing.T) {
	ld := NewLoopDetector(3)
	events := []JobEvent{
		{EventType: "job.failed", Message: "error 1"},
		{EventType: "job.failed", Message: "error 2"},
		{EventType: "job.failed", Message: "error 3"},
	}
	isLoop, reason := ld.Check(events)
	if !isLoop {
		t.Fatal("expected loop detected")
	}
	if reason == "" {
		t.Fatal("expected non-empty reason")
	}
}

func TestLoopDetectorReset(t *testing.T) {
	ld := NewLoopDetector(3)
	events := []JobEvent{
		{EventType: "job.failed", Message: "error 1"},
		{EventType: "job.completed", Message: "success"},
		{EventType: "job.failed", Message: "error 2"},
	}
	isLoop, _ := ld.Check(events)
	if isLoop {
		t.Fatal("expected no loop after completion reset")
	}
}

func TestJobIsTerminal(t *testing.T) {
	cases := []struct {
		status string
		want   bool
	}{
		{"completed", true},
		{"failed", true},
		{"cancelled", true},
		{"queued", false},
		{"paused", false},
		{"retrying", false},
	}
	for _, c := range cases {
		j := &Job{Status: c.status}
		if got := j.IsTerminal(); got != c.want {
			t.Fatalf("IsTerminal(%s) = %v, want %v", c.status, got, c.want)
		}
	}
}

func TestJobCanRetry(t *testing.T) {
	j := &Job{Status: "failed", RetryCount: 1, MaxRetries: 3}
	if !j.CanRetry() {
		t.Fatal("expected CanRetry true")
	}
	j.RetryCount = 3
	if j.CanRetry() {
		t.Fatal("expected CanRetry false when max reached")
	}
	j.Status = "completed"
	j.RetryCount = 0
	if j.CanRetry() {
		t.Fatal("expected CanRetry false when terminal")
	}
}
