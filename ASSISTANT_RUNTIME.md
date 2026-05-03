# Assistant Runtime

The autonomous execution layer of Zsistant.

## What It Does

- Wakes every N seconds (default 30)
- Checks every enabled assistant for due work
- Runs jobs **without calling LLMs** unless truly needed
- Enforces budgets: tokens/day, actions/day, LLM calls/hour
- Requires human approval for risky actions when autonomy is not `full`
- Logs every event to JSONL activity feed
- Goes back to sleep

## Architecture

### Engine (`internal/runtime/engine.go`)

- `Start()` / `Stop()` / `IsRunning()`
- `tickAll()` loops through every assistant
- `tickAssistant()` runs the 5-layer execution model

### 5-Layer Execution Model

1. **Deterministic checks** — quiet hours, budget, enabled status
2. **Lightweight planning** — find due jobs, check retries
3. **LLM call** (only if needed) — `executeWithLLM()` uses the cheapest available model
4. **Action execution** — update job status, record outputs
5. **Persistence + sleep** — write state, log activity, done

### Approval Layer

- `requiresApproval()` gates LLM calls based on:
  - `AutonomyLevel`: `none` (always), `semi` (heuristic), `full` (never)
  - Heuristics: >2 failures, <1000 tokens left, event-triggered jobs
- Creates an `approvals.Request` in `internal/approvals/store.go`
- Returns an error, engine marks job as `paused` and moves on

### State Persistence

`internal/runtime/state.go` + `persistence.go`

Stored per-assistant at `{basePath}/assistants/{id}/runtime_state.json`.

Fields include:
- `enabled`, `status`, `autonomy_level`
- `token_budget_used_today`, `token_budget_per_day`
- `actions_used_today`, `action_budget_per_day`
- `llm_calls_this_hour`, `max_llm_calls_per_hour`
- `consecutive_failures`, `last_failure_at`
- `last_check_at`, `next_check_at`
- `cheap_model`, `expensive_model`

## Endpoints

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/assistants/{id}/run` | Run assistant immediately |
| POST | `/api/assistants/{id}/pause` | Disable assistant |
| POST | `/api/assistants/{id}/resume` | Re-enable assistant |
| GET | `/api/assistants/{id}/state` | Get runtime state |
| PUT | `/api/assistants/{id}/state` | Update runtime state |
| GET | `/api/runtime/status` | Engine uptime, tick, running flag |
| GET | `/api/runtime/activity` | Global activity feed |

## UI

Runtime tab in assistant detail shows:
- Enabled/Status/Autonomy/Interval/Last check/Next check
- Token and action budget usage
- Failure count
- Pause/Resume/Run Now buttons
- Link to Approvals tab
