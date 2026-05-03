# Jobs and Scheduler

Background task system for Zsistant assistants.

## What It Does

- Jobs are durable, owned by one assistant
- The engine wakes, finds jobs whose `NextRunAt` is in the past, and executes them
- No cron daemon: Go goroutine loop + time comparison
- Honest states: if nothing is due, nothing runs

## Job Model (`internal/jobs/job.go`)

- `ID`, `AssistantID`, `Name`, `Purpose`
- `ScheduleType`: `manual`, `cron`, `event`
- `Schedule`: cron expression or event name
- `Status`: `pending`, `running`, `completed`, `failed`, `paused`
- `NextRunAt`, `LastRunAt`, `CreatedAt`
- `RetryCount`, `MaxRetries`, `LastError`

## Queue (`internal/jobs/queue.go`)

- `Enqueue(job)` — writes job JSON to `{basePath}/assistants/{id}/jobs/{job-id}.json`
- `List()` — returns all jobs for an assistant
- `Get(id)` / `Update(job)` / `Delete(id)`

## Scheduler

Inside `internal/runtime/engine.go`:

- `findDueJobs(list, now)` filters jobs where `NextRunAt < now`
- `runJob()` picks `executeManual`, `executeScheduled`, or `executeWithLLM`
- Budget enforcement happens before any execution
- Retry logic uses exponential backoff via `job.Backoff()`

## Execution Strategy

| Job Type | Execution Path |
|----------|---------------|
| `manual` | `executeManual()` — deterministic, no LLM |
| `scheduled` (no purpose) | Deterministic, no LLM |
| `scheduled` (needs reasoning) | `executeWithLLM()` — cheap model, approval gate |

## Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/assistants/{id}/jobs` | List assistant jobs |
| POST | `/api/assistants/{id}/jobs` | Create a job |
| GET | `/api/jobs` | Global jobs (stub) |

## UI

- Jobs tab shows list with status badges
- Click **Create Job** in Jobs tab to make a manual job
- Runtime tab shows next check time and budget impact
