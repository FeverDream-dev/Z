package runtime

import (
	
	
	"testing"
	"time"

	"github.com/FeverDream-dev/zsistant/internal/assistant"
	"github.com/FeverDream-dev/zsistant/internal/jobs"
)

func TestDefaultRuntimeState(t *testing.T) {
	s := DefaultRuntimeState()
	if !s.Enabled {
		t.Error("expected Enabled to be true")
	}
	if s.Status != "idle" {
		t.Errorf("expected status idle, got %s", s.Status)
	}
	if s.TokenBudgetPerDay == 0 {
		t.Error("expected positive token budget")
	}
	if s.RuntimeIntervalSec == 0 {
		t.Error("expected positive interval")
	}
}

func TestLoadSaveRoundTrip(t *testing.T) {
	dir := t.TempDir()
	id := "test-assistant"

	// Save a state
	s1 := DefaultRuntimeState()
	s1.Status = "running"
	s1.TokenBudgetUsedToday = 500
	err := saveState(dir, id, s1)
	if err != nil {
		t.Fatalf("save state: %v", err)
	}

	// Load it back
	s2 := loadState(dir, id)
	if s2.Status != "running" {
		t.Errorf("expected status running, got %s", s2.Status)
	}
	if s2.TokenBudgetUsedToday != 500 {
		t.Errorf("expected 500 used, got %d", s2.TokenBudgetUsedToday)
	}
}

func TestLoadStateNotExists(t *testing.T) {
	dir := t.TempDir()
	s := loadState(dir, "does-not-exist")
	if s.Status != "idle" {
		t.Error("expected default state for missing file")
	}
}

func TestFindDueJobs(t *testing.T) {
	now := time.Now()
	past := now.Add(-time.Hour)
	future := now.Add(time.Hour)

	list := []*jobs.Job{
		{ID: "1", Name: "past", NextRunAt: &past},
		{ID: "2", Name: "future", NextRunAt: &future},
		{ID: "3", Name: "nil"},
	}

	due := findDueJobs(list, now)
	if len(due) != 2 {
		t.Errorf("expected 2 due jobs, got %d", len(due))
	}
}

func TestFindDueJobsNone(t *testing.T) {
	now := time.Now()
	future := now.Add(time.Hour)
	list := []*jobs.Job{
		{ID: "1", Name: "future", NextRunAt: &future},
	}
	due := findDueJobs(list, now)
	if len(due) != 0 {
		t.Errorf("expected 0 due jobs, got %d", len(due))
	}
}

func TestInQuietHours(t *testing.T) {
	state := DefaultRuntimeState()

	// No quiet hours set
	state.QuietHoursStart = ""
	state.QuietHoursEnd = ""
	if inQuietHours(state) {
		t.Error("expected not in quiet hours when unset")
	}

	// Set quiet hours to include right now
	state.QuietHoursStart = "00:00"
	state.QuietHoursEnd = "23:59"
	if !inQuietHours(state) {
		t.Error("expected in quiet hours when range covers full day")
	}

	// Set quiet hours to a future window
	state.QuietHoursStart = "23:00"
	state.QuietHoursEnd = "23:59"
	if inQuietHours(state) {
		t.Error("expected not in quiet hours for future window")
	}
}

func TestNextCheckTime(t *testing.T) {
	state := DefaultRuntimeState()
	state.RuntimeIntervalSec = 60
	now := time.Now()
	next := nextCheckTime(state, now)
	if !next.After(now) {
		t.Error("expected next check after now")
	}
	diff := next.Sub(now)
	if diff < 55*time.Second || diff > 65*time.Second {
		t.Errorf("expected ~60s diff, got %v", diff)
	}
}

func TestResetDailyBudgets(t *testing.T) {
	state := DefaultRuntimeState()
	state.TokenBudgetUsedToday = 999
	state.ActionBudgetUsedToday = 50
	state.ConsecutiveFailures = 3

	now := time.Now()

	// Day hasn't passed
	state.FailureCountResetAt = ptr(now.Format(time.RFC3339))
	resetDailyBudgets(&state, now)
	if state.TokenBudgetUsedToday != 999 {
		t.Error("budgets should not reset when day hasn't passed")
	}

	// Day has passed
	state.FailureCountResetAt = ptr(now.Add(-25 * time.Hour).Format(time.RFC3339))
	resetDailyBudgets(&state, now)
	if state.TokenBudgetUsedToday != 0 {
		t.Error("token budget should reset after a day")
	}
	if state.ActionBudgetUsedToday != 0 {
		t.Error("action budget should reset after a day")
	}
	if state.ConsecutiveFailures != 3 {
		t.Error("consecutive failures should NOT be reset")
	}
}

func TestRequiresApproval(t *testing.T) {
	state := DefaultRuntimeState()
	a := assistant.Assistant{ID: "test", Name: "Test"}
	job := &jobs.Job{ID: "j1", Name: "job", ScheduleType: jobs.ScheduleManual}

	// full autonomy → no approval
	state.AutonomyLevel = "full"
	if requiresApproval(&state, a, job) {
		t.Error("full autonomy should not require approval")
	}

	// none autonomy → always approval
	state.AutonomyLevel = "none"
	if !requiresApproval(&state, a, job) {
		t.Error("none autonomy should always require approval")
	}

	// semi autonomy + event job → approval
	state.AutonomyLevel = "semi"
	job.ScheduleType = jobs.ScheduleEvent
	if !requiresApproval(&state, a, job) {
		t.Error("event job should require approval under semi")
	}

	// semi autonomy + high failures → approval
	job.ScheduleType = jobs.ScheduleManual
	state.ConsecutiveFailures = 3
	if !requiresApproval(&state, a, job) {
		t.Error("high failures should require approval under semi")
	}

	// semi autonomy + low budget → approval
	state.ConsecutiveFailures = 0
	state.TokenBudgetPerDay = 500
	state.TokenBudgetUsedToday = 450
	if !requiresApproval(&state, a, job) {
		t.Error("low remaining budget should require approval under semi")
	}
}

func TestEngineStartStop(t *testing.T) {
	dir := t.TempDir()
	reg := assistant.NewRegistry(dir)
	eng := NewEngine(dir, reg, nil, nil, nil)

	if eng.IsRunning() {
		t.Error("expected not running before start")
	}

	eng.Start()
	if !eng.IsRunning() {
		t.Error("expected running after start")
	}

	eng.Stop()
	if eng.IsRunning() {
		t.Error("expected not running after stop")
	}
}

func TestTickAssistantDisabled(t *testing.T) {
	dir := t.TempDir()
	reg := assistant.NewRegistry(dir)
	_, _ = reg.Create("disabled-test", "Disabled Test")

	// Disable it
	s := loadState(dir, "disabled-test")
	s.Enabled = false
	_ = saveState(dir, "disabled-test", s)

	eng := NewEngine(dir, reg, nil, nil, nil)
	eng.Start()
	defer eng.Stop()

	a, _ := reg.Get("disabled-test")
	eng.tickAssistant(*a)

	// No crash; no state change
	s2 := loadState(dir, "disabled-test")
	if s2.Enabled {
		t.Error("expected disabled to stay disabled")
	}
}

func TestTickAssistantNoJobs(t *testing.T) {
	dir := t.TempDir()
	reg := assistant.NewRegistry(dir)
	_, _ = reg.Create("idle-test", "Idle Test")

	eng := NewEngine(dir, reg, nil, nil, nil)
	eng.Start()
	defer eng.Stop()

	a, _ := reg.Get("idle-test")
	eng.tickAssistant(*a)

	s := loadState(dir, "idle-test")
	if s.Status != "idle" {
		t.Errorf("expected idle, got %s", s.Status)
	}
	if s.NextRunAt == nil || *s.NextRunAt == "" {
		t.Error("expected next check time to be set")
	}
}

func TestTickAssistantBudgetExceeded(t *testing.T) {
	dir := t.TempDir()
	reg := assistant.NewRegistry(dir)
	_, _ = reg.Create("broke-test", "Broke Test")

	// Use up all tokens
	s := loadState(dir, "broke-test")
	s.TokenBudgetUsedToday = s.TokenBudgetPerDay + 1
	_ = saveState(dir, "broke-test", s)

	eng := NewEngine(dir, reg, nil, nil, nil)
	eng.Start()
	defer eng.Stop()

	a, _ := reg.Get("broke-test")
	eng.tickAssistant(*a)

	s2 := loadState(dir, "broke-test")
	if s2.Status == "running" {
		t.Error("expected not running when budget exceeded")
	}
}
