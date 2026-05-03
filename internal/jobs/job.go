package jobs

import (
	"fmt"
	"time"
)

// Status values for jobs.
const (
	StatusQueued          = "queued"
	StatusRunning         = "running"
	StatusCompleted       = "completed"
	StatusFailed          = "failed"
	StatusPaused          = "paused"
	StatusCancelled       = "cancelled"
	StatusRetrying        = "retrying"
	StatusWaitingApproval = "waiting_for_approval"
)

// ScheduleType describes how a job is triggered.
type ScheduleType string

const (
	ScheduleManual ScheduleType = "manual"
	ScheduleCron   ScheduleType = "cron"
	ScheduleEvent  ScheduleType = "event"
)

// Job represents a durable unit of scheduled or event-triggered work owned by an assistant.
type Job struct {
	ID                   string       `json:"id"`
	AssistantID          string       `json:"assistant_id"`
	Name                 string       `json:"name"`
	Purpose              string       `json:"purpose"`              // what this job does
	Schedule             string       `json:"schedule,omitempty"`  // cron expression or event description
	ScheduleType         ScheduleType `json:"schedule_type,omitempty"`
	RequiredTools        []string     `json:"required_tools,omitempty"`
	RequiredPermissions  []string     `json:"required_permissions,omitempty"`
	Status               string       `json:"status"`
	CurrentStep          string       `json:"current_step"`
	Result               string       `json:"result,omitempty"`
	OutputHistory        []string     `json:"output_history,omitempty"`
	RetryCount           int          `json:"retry_count"`
	MaxRetries           int          `json:"max_retries"`
	LastError            string       `json:"last_error,omitempty"`
	PausedReason         string       `json:"paused_reason,omitempty"`
	NextRunAt            *time.Time   `json:"next_run_at,omitempty"`
	LastRunAt            *time.Time   `json:"last_run_at,omitempty"`
	CreatedAt            time.Time    `json:"created_at"`
	UpdatedAt            time.Time    `json:"updated_at"`
	CompletedAt          *time.Time   `json:"completed_at,omitempty"`
}

// NewJob creates a new job with a generated ID and initial status.
func NewJob(assistantID, name, purpose string) *Job {
	now := time.Now()
	return &Job{
		ID:            generateID(),
		AssistantID:   assistantID,
		Name:          name,
		Purpose:       purpose,
		Status:        StatusQueued,
		CurrentStep:   "created",
		ScheduleType:  ScheduleManual,
		RetryCount:    0,
		MaxRetries:    3,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

// NewScheduledJob creates a recurring job with a cron schedule.
func NewScheduledJob(assistantID, name, purpose, cronExpr string) *Job {
	j := NewJob(assistantID, name, purpose)
	j.Schedule = cronExpr
	j.ScheduleType = ScheduleCron
	return j
}

// NewTriggeredJob creates an event-triggered job.
func NewTriggeredJob(assistantID, name, purpose, eventDesc string) *Job {
	j := NewJob(assistantID, name, purpose)
	j.Schedule = eventDesc
	j.ScheduleType = ScheduleEvent
	return j
}

// IsTerminal returns true if the job is in a terminal state.
func (j *Job) IsTerminal() bool {
	return j.Status == StatusCompleted || j.Status == StatusFailed || j.Status == StatusCancelled
}

// CanRetry returns true if the job can be retried.
func (j *Job) CanRetry() bool {
	if j.RetryCount >= j.MaxRetries {
		return false
	}
	return j.Status == StatusFailed || j.Status == StatusPaused || j.Status == StatusRetrying
}

// HonestSummary returns a human-readable description of the job state.
func (j *Job) HonestSummary() string {
	if j.Status == StatusFailed && j.LastError != "" {
		return fmt.Sprintf("Job '%s' failed: %s (retries: %d/%d)", j.Name, j.LastError, j.RetryCount, j.MaxRetries)
	}
	if j.Status == StatusPaused && j.PausedReason != "" {
		return fmt.Sprintf("Job '%s' is paused: %s", j.Name, j.PausedReason)
	}
	if j.Status == StatusCompleted {
		return fmt.Sprintf("Job '%s' completed.", j.Name)
	}
	if j.Status == StatusQueued && j.NextRunAt != nil {
		return fmt.Sprintf("Job '%s' queued, next run %s.", j.Name, j.NextRunAt.Format("2006-01-02 15:04"))
	}
	return fmt.Sprintf("Job '%s' is %s.", j.Name, j.Status)
}

var idCounter int64

func generateID() string {
	idCounter++
	return fmt.Sprintf("%s-%06d", time.Now().Format("20060102-150405"), idCounter)
}
