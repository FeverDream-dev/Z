package runtime

import (
	"fmt"
	"sync"
	"time"

	"github.com/FeverDream-dev/zsistant/internal/activity"
	"github.com/FeverDream-dev/zsistant/internal/assistant"
	"github.com/FeverDream-dev/zsistant/internal/approvals"
	"github.com/FeverDream-dev/zsistant/internal/jobs"
	"github.com/FeverDream-dev/zsistant/internal/llm"
	"github.com/FeverDream-dev/zsistant/internal/tools"
)

// Engine is the autonomous assistant runtime.
type Engine struct {
	mu          sync.RWMutex
	basePath    string
	running     bool
	stopCh      chan struct{}
	tick        time.Duration
	registry    *assistant.Registry
	approvals   *approvals.Store
	factory     *llm.Factory
	broker      *tools.Broker
}

// NewEngine creates the runtime engine.
func NewEngine(basePath string, reg *assistant.Registry, appStore *approvals.Store, factory *llm.Factory, broker *tools.Broker) *Engine {
	return &Engine{
		basePath:  basePath,
		stopCh:    make(chan struct{}),
		tick:      30 * time.Second,
		registry:  reg,
		approvals: appStore,
		factory:   factory,
		broker:    broker,
	}
}

// Start begins the background worker loop.
func (e *Engine) Start() {
	e.mu.Lock()
	if e.running {
		e.mu.Unlock()
		return
	}
	e.running = true
	e.mu.Unlock()
	go e.loop()
}

// Stop gracefully shuts down the runtime.
func (e *Engine) Stop() {
	e.mu.Lock()
	if !e.running {
		e.mu.Unlock()
		return
	}
	e.running = false
	e.mu.Unlock()
	close(e.stopCh)
}

// IsRunning reports whether the engine is active.
func (e *Engine) IsRunning() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.running
}

func (e *Engine) loop() {
	ticker := time.NewTicker(e.tick)
	defer ticker.Stop()
	for {
		select {
		case <-e.stopCh:
			return
		case <-ticker.C:
			e.tickAll()
		}
	}
}

// tickAll checks every assistant for due work.
func (e *Engine) tickAll() {
	list, err := e.registry.List()
	if err != nil {
		return
	}
	for _, a := range list {
		e.tickAssistant(a)
	}
}

// RunNow immediately runs a single assistant tick regardless of schedule.
func (e *Engine) RunNow(a assistant.Assistant) {
	e.tickAssistant(a)
}

// tickAssistant runs one assistant tick: deterministic checks only, no LLM unless needed.
func (e *Engine) tickAssistant(a assistant.Assistant) {
	state := loadState(e.basePath, a.ID)
	if !state.Enabled || state.Status == "paused" {
		return
	}
	if inQuietHours(state) {
		return
	}
	if state.ConsecutiveFailures >= 5 {
		state.Status = "error"
		state.LastError = "too many consecutive failures"
		saveState(e.basePath, a.ID, state)
		return
	}

	now := time.Now()
	resetDailyBudgets(&state, now)
	if state.TokenBudgetUsedToday >= state.TokenBudgetPerDay {
		return
	}

	q := jobs.NewQueue(assistantWorkspace(e.basePath, a.ID))
	jlist, _ := q.List()
	due := findDueJobs(jlist, now)
	if len(due) == 0 {
		state.Status = "idle"
		state.NextRunAt = ptr(nextCheckTime(state, now).Format(time.RFC3339))
		saveState(e.basePath, a.ID, state)
		return
	}

	state.Status = "running"
	saveState(e.basePath, a.ID, state)

	for _, job := range due {
		if state.ActionBudgetUsedToday >= state.ActionBudgetPerDay {
			break
		}
		if state.MaxActionsPerRun > 0 && state.ActionBudgetUsedToday >= state.MaxActionsPerRun {
			break
		}
		e.runJob(a, state, q, job)
	}

	state.Status = "idle"
	state.LastRunAt = ptr(now.Format(time.RFC3339))
	state.NextRunAt = ptr(nextCheckTime(state, now).Format(time.RFC3339))
	saveState(e.basePath, a.ID, state)
}

// runJob executes a single job deterministically. LLM only called for complex decisions.
func (e *Engine) runJob(a assistant.Assistant, state RuntimeState, q *jobs.Queue, job *jobs.Job) {
	start := time.Now()
	job.Status = jobs.StatusRunning
	job.LastRunAt = &start
	q.Update(job)

	activity.Log(e.basePath, activity.ActivityEvent{
		AssistantID: a.ID,
		EventType:   "job.started",
		Message:     fmt.Sprintf("Job %s started", job.Name),
		Severity:    "info",
		Metadata:    map[string]string{"job_id": job.ID, "name": job.Name},
	})

	var result string
	var err error

	switch job.ScheduleType {
	case jobs.ScheduleManual:
		result, err = e.executeManual(a, state, job)
	default:
		result, err = e.executeScheduled(a, state, job)
	}

	if err != nil {
		job.Status = jobs.StatusFailed
		job.LastError = err.Error()
		job.RetryCount++
		state.FailureCount++
		state.FailureCountToday++
		state.ConsecutiveFailures++
		state.LastError = err.Error()
		if job.RetryCount < job.MaxRetries {
			backoff := jobs.Backoff(job.RetryCount)
			next := time.Now().Add(backoff)
			job.NextRunAt = &next
			job.Status = jobs.StatusRetrying
		}
	} else {
		job.Status = jobs.StatusCompleted
		job.Result = result
		state.ConsecutiveFailures = 0
		state.ActionBudgetUsedToday++
	}
	job.UpdatedAt = time.Now()
	q.Update(job)

	activity.Log(e.basePath, activity.ActivityEvent{
		AssistantID: a.ID,
		EventType:   "job.completed",
		Message:     fmt.Sprintf("Job %s: %s", job.Name, job.Status),
		Severity:    map[bool]string{true: "error", false: "info"}[err != nil],
		Metadata:    map[string]string{"job_id": job.ID, "status": job.Status},
	})
}

func (e *Engine) executeManual(a assistant.Assistant, state RuntimeState, job *jobs.Job) (string, error) {
	// Deterministic: just record the job and return a summary
	return fmt.Sprintf("Manual job '%s' recorded. No LLM call needed for deterministic execution.", job.Name), nil
}

func (e *Engine) executeScheduled(a assistant.Assistant, state RuntimeState, job *jobs.Job) (string, error) {
	// Scheduled tasks: simple deterministic execution first
	// Only call LLM if job description implies reasoning is needed
	needsLLM := len(job.Purpose) > 0 && state.LLMCallsThisHour < state.MaxLLMCallsPerHour
	if !needsLLM {
		return fmt.Sprintf("Scheduled job '%s' executed without LLM (deterministic).", job.Name), nil
	}
	return e.executeWithLLM(a, state, job)
}

func (e *Engine) executeWithLLM(a assistant.Assistant, state RuntimeState, job *jobs.Job) (string, error) {
	modelID := state.CheapModel
	if modelID == "" {
		modelID = a.DefaultModel
	}
	if modelID == "" {
		modelID = "gpt-4o-mini"
	}

	if e.approvals != nil && requiresApproval(&state, a, job) {
		expires := time.Now().Add(24 * time.Hour)
		req := &approvals.Request{
			AssistantID:    a.ID,
			TaskID:         job.ID,
			ActionSummary:  fmt.Sprintf("LLM call for job '%s': %s", job.Name, job.Purpose),
			RiskLevel:      "medium",
			RequestedAt:    time.Now(),
			Status:         "pending",
			ExpiresAt:      &expires,
		}
		if err := e.approvals.Create(req); err != nil {
			return "", fmt.Errorf("failed to create approval: %w", err)
		}
		return "", fmt.Errorf("approval required for LLM call on job '%s'", job.Name)
	}

	prov, _, _, err := e.factory.Create(modelID)
	if err != nil {
		return "", fmt.Errorf("provider not available: %w", err)
	}
	prompt := fmt.Sprintf("You are %s. Your purpose: %s.\nTask: %s\nComplete this concisely.", a.Name, a.Purpose, job.Purpose)
	resp, err := prov.Complete(prompt)
	if err != nil {
		return "", err
	}
	state.LLMCallsThisHour++
	state.TokenBudgetUsedToday += len(prompt) + len(resp)
	return resp, nil
}

func requiresApproval(state *RuntimeState, a assistant.Assistant, job *jobs.Job) bool {
	if state.AutonomyLevel == "full" {
		return false
	}
	if state.AutonomyLevel == "none" {
		return true
	}
	if job.ScheduleType == "event" {
		return true
	}
	if state.ConsecutiveFailures >= 2 {
		return true
	}
	if state.TokenBudgetPerDay-state.TokenBudgetUsedToday < 1000 {
		return true
	}
	return false
}

// findDueJobs returns jobs whose NextRunAt is past or nil.
func findDueJobs(list []*jobs.Job, now time.Time) []*jobs.Job {
	var out []*jobs.Job
	for _, j := range list {
		if j.Status == jobs.StatusPaused || j.Status == jobs.StatusCancelled {
			continue
		}
		if j.NextRunAt == nil || j.NextRunAt.Before(now) || j.NextRunAt.Equal(now) {
			out = append(out, j)
		}
	}
	return out
}

func inQuietHours(state RuntimeState) bool {
	if state.QuietHoursStart == "" || state.QuietHoursEnd == "" {
		return false
	}
	now := time.Now()
	start, _ := time.Parse("15:04", state.QuietHoursStart)
	end, _ := time.Parse("15:04", state.QuietHoursEnd)
	if start.IsZero() || end.IsZero() {
		return false
	}
	current := now.Hour()*60 + now.Minute()
	s := start.Hour()*60 + start.Minute()
	e := end.Hour()*60 + end.Minute()
	if s < e {
		return current >= s && current < e
	}
	return current >= s || current < e
}

func resetDailyBudgets(state *RuntimeState, now time.Time) {
	if state.FailureCountResetAt == nil {
		state.FailureCountToday = 0
		state.TokenBudgetUsedToday = 0
		state.ActionBudgetUsedToday = 0
		state.LLMCallsThisHour = 0
		state.FailureCountResetAt = ptr(now.Format(time.RFC3339))
		return
	}
	resetTime, _ := time.Parse(time.RFC3339, *state.FailureCountResetAt)
	if resetTime.IsZero() || now.Day() != resetTime.Day() {
		state.FailureCountToday = 0
		state.TokenBudgetUsedToday = 0
		state.ActionBudgetUsedToday = 0
		state.LLMCallsThisHour = 0
		state.FailureCountResetAt = ptr(now.Format(time.RFC3339))
	}
}

func nextCheckTime(state RuntimeState, now time.Time) time.Time {
	interval := time.Duration(state.RuntimeIntervalSec) * time.Second
	if interval <= 0 {
		interval = 30 * time.Second
	}
	return now.Add(interval)
}

func assistantWorkspace(basePath, assistantID string) string {
	return fmt.Sprintf("%s/assistants/%s", basePath, assistantID)
}

func ptr(s string) *string { return &s }
