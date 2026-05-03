package runtime

// RuntimeState is the per-assistant runtime configuration and current status.
type RuntimeState struct {
	// Control
	Enabled              bool   `json:"enabled"`
	AutonomyLevel        string `json:"autonomy_level"` // full, semi, manual
	Status               string `json:"status"`           // idle, sleeping, running, waiting_for_approval, error, paused

	// Timing
	LastRunAt            *string `json:"last_run_at,omitempty"`
	NextRunAt            *string `json:"next_run_at,omitempty"`
	RuntimeIntervalSec   int     `json:"runtime_interval_sec"`  // how often to check (default 60)
	QuietHoursStart      string  `json:"quiet_hours_start,omitempty"` // e.g. "22:00"
	QuietHoursEnd        string  `json:"quiet_hours_end,omitempty"`   // e.g. "08:00"

	// Budget and limits
	TokenBudgetPerDay        int `json:"token_budget_per_day"`
	TokenBudgetUsedToday     int `json:"token_budget_used_today"`
	ActionBudgetPerDay       int `json:"action_budget_per_day"`
	ActionBudgetUsedToday    int `json:"action_budget_used_today"`
	MaxLLMCallsPerHour       int `json:"max_llm_calls_per_hour"`
	LLMCallsThisHour         int `json:"llm_calls_this_hour"`
	MaxActionsPerRun         int `json:"max_actions_per_run"`
	MaxRuntimeSecPerRun      int `json:"max_runtime_sec_per_run"`

	// Failure tracking
	FailureCount         int       `json:"failure_count"`
	FailureCountToday    int       `json:"failure_count_today"`
	FailureCountResetAt  *string   `json:"failure_count_reset_at,omitempty"`
	ConsecutiveFailures  int       `json:"consecutive_failures"`
	LastError            string    `json:"last_error,omitempty"`

	// Permissions
	AllowedTools         []string `json:"allowed_tools,omitempty"`
	AllowedChannels      []string `json:"allowed_channels,omitempty"`
	ApprovalRequiredFor  []string `json:"approval_required_for,omitempty"`

	// Model preferences
	CheapModel           string `json:"cheap_model,omitempty"`
	ExpensiveModel       string `json:"expensive_model,omitempty"`

	// Current work
	CurrentTaskID        string `json:"current_task_id,omitempty"`
	PendingApprovalCount int    `json:"pending_approval_count"`

	// Metadata
	UpdatedAt            string `json:"updated_at,omitempty"`
}

// DefaultRuntimeState returns a sensible default.
func DefaultRuntimeState() RuntimeState {
	return RuntimeState{
		Enabled:             true,
		AutonomyLevel:       "semi",
		Status:              "idle",
		RuntimeIntervalSec:  60,
		TokenBudgetPerDay:   10000,
		ActionBudgetPerDay:  100,
		MaxLLMCallsPerHour:  10,
		MaxActionsPerRun:    5,
		MaxRuntimeSecPerRun: 30,
		CheapModel:          "gpt-4o-mini",
		ExpensiveModel:      "gpt-4o",
		AllowedTools:        []string{},
		AllowedChannels:     []string{},
		ApprovalRequiredFor: []string{"delete", "send_external", "spend_tokens"},
	}
}
