package jobs

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// Queue manages durable jobs for an assistant.
type Queue struct {
	workspace string
	jobDir    string
}

// NewQueue creates a queue for the given agent workspace.
func NewQueue(workspace string) *Queue {
	return &Queue{
		workspace: workspace,
		jobDir:    filepath.Join(workspace, "jobs"),
	}
}

// Enqueue persists a new job to the queue.
func (q *Queue) Enqueue(job *Job) error {
	path := q.jobPath(job.ID)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("creating job dir: %w", err)
	}
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("creating job file: %w", err)
	}
	defer f.Close()
	if err := json.NewEncoder(f).Encode(job); err != nil {
		return fmt.Errorf("encoding job: %w", err)
	}
	return nil
}

// Get retrieves a job by ID.
func (q *Queue) Get(jobID string) (*Job, error) {
	path := q.jobPath(jobID)
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening job file: %w", err)
	}
	defer f.Close()
	var job Job
	if err := json.NewDecoder(f).Decode(&job); err != nil {
		return nil, fmt.Errorf("decoding job: %w", err)
	}
	return &job, nil
}

// Update persists changes to a job.
func (q *Queue) Update(job *Job) error {
	return q.Enqueue(job)
}

// List returns all jobs in the queue.
func (q *Queue) List() ([]*Job, error) {
	entries, err := os.ReadDir(q.jobDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("reading jobs dir: %w", err)
	}
	var jobs []*Job
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		job, err := q.Get(entry.Name())
		if err != nil {
			continue
		}
		jobs = append(jobs, job)
	}
	sort.Slice(jobs, func(i, j int) bool {
		return jobs[i].CreatedAt.After(jobs[j].CreatedAt)
	})
	return jobs, nil
}

// Retry marks a job for retry with exponential backoff.
func (q *Queue) Retry(jobID string) (*Job, error) {
	job, err := q.Get(jobID)
	if err != nil {
		return nil, err
	}
	if !job.CanRetry() {
		return nil, fmt.Errorf("job %s cannot be retried (retries: %d/%d)", jobID, job.RetryCount, job.MaxRetries)
	}
	job.RetryCount++
	job.Status = StatusRetrying
	job.LastError = ""
	job.UpdatedAt = time.Now()
	if err := q.Update(job); err != nil {
		return nil, err
	}
	return job, nil
}

// Pause marks a job as paused with a reason.
func (q *Queue) Pause(jobID, reason string) (*Job, error) {
	job, err := q.Get(jobID)
	if err != nil {
		return nil, err
	}
	job.Status = StatusPaused
	job.PausedReason = reason
	job.UpdatedAt = time.Now()
	if err := q.Update(job); err != nil {
		return nil, err
	}
	return job, nil
}

// Resume marks a paused job as queued.
func (q *Queue) Resume(jobID string) (*Job, error) {
	job, err := q.Get(jobID)
	if err != nil {
		return nil, err
	}
	if job.Status != StatusPaused {
		return nil, fmt.Errorf("job %s is not paused (status: %s)", jobID, job.Status)
	}
	job.Status = StatusQueued
	job.PausedReason = ""
	job.UpdatedAt = time.Now()
	if err := q.Update(job); err != nil {
		return nil, err
	}
	return job, nil
}

// Cancel marks a job as cancelled.
func (q *Queue) Cancel(jobID string) (*Job, error) {
	job, err := q.Get(jobID)
	if err != nil {
		return nil, err
	}
	job.Status = StatusCancelled
	job.UpdatedAt = time.Now()
	if err := q.Update(job); err != nil {
		return nil, err
	}
	return job, nil
}

func (q *Queue) jobPath(jobID string) string {
	return filepath.Join(q.jobDir, jobID)
}

// Backoff computes the delay before the next retry using exponential backoff.
func Backoff(retryCount int) time.Duration {
	base := time.Second
	max := 30 * time.Second
	d := base * (1 << retryCount)
	if d > max {
		return max
	}
	return d
}

// LoopDetector detects repeated failures that may indicate a loop.
type LoopDetector struct {
	threshold int
}

// NewLoopDetector creates a loop detector with the given failure threshold.
func NewLoopDetector(threshold int) *LoopDetector {
	return &LoopDetector{threshold: threshold}
}

// Check analyzes job events and returns true if a loop is suspected.
func (ld *LoopDetector) Check(events []JobEvent) (bool, string) {
	failureCount := 0
	var lastError string
	for _, evt := range events {
		switch evt.EventType {
		case "job.failed", "job.error":
			failureCount++
			lastError = evt.Message
		case "job.completed", "job.cancelled":
			failureCount = 0
		}
	}
	if failureCount >= ld.threshold {
		return true, fmt.Sprintf("loop suspected: %d consecutive failures (last: %s)", failureCount, lastError)
	}
	return false, ""
}
