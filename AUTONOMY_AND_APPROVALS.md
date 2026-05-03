# Autonomy and Approvals

Permission layer for autonomous assistant actions.

## Why This Exists

Zsistant is NOT an always-thinking agent. It runs on a schedule, but when it *does* need to act (especially with an LLM or external tool), it should ask for permission unless the user trusts it completely.

## Autonomy Levels

Each assistant has a runtime state field `autonomy_level`:

- **`none`** — Every LLM call requires approval. Nothing runs without explicit human permission.
- **`semi`** — Default. Heuristic rules decide when to ask.
- **`full`** — No approval needed. The assistant runs freely within budget.

### Heuristic Rules (`semi`)

Approval required when ANY of:
- Consecutive failures >= 2
- Remaining token budget < 1000
- Job is event-triggered (`schedule_type == "event"`)

## Approval Requests

### Model (`internal/approvals/Request`)

- `ID`, `AssistantID`, `TaskID`
- `ActionSummary` — human-readable description
- `RiskLevel` — `low`, `medium`, `high`, `critical`
- `Status` — `pending`, `approved`, `denied`, `expired`
- `RequestedAt`, `ResolvedAt`, `ExpiresAt` (default 24h)

### Store (`internal/approvals/store.go`)

- JSONL persistence at `{basePath}/approvals.jsonl`
- `Create(req)` — auto-assigns UUID
- `List()` — all requests
- `Get(id)` / `Resolve(id, status, approver)`
- `PendingCount()` — for dashboard badge

## Engine Integration

`internal/runtime/engine.go`:

- Before calling `executeWithLLM`, checks `requiresApproval()`
- If approval needed, creates request with `e.approvals.Create()` and returns error
- Job is marked `paused` with `LastError = "approval required"`
- When user resolves via `/api/approvals/{id}`, a later tick may retry

## Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/approvals` | List all requests (filter `?assistant=`) |
| POST | `/api/approvals` | Create request manually |
| POST | `/api/approvals/{id}` | Resolve: body `{status: "approved|denied"}` |
| GET | `/api/approvals/pending` | Count pending only |

## UI

- **Runtime tab** — shows current autonomy level, budget, controls
- **Approvals tab** — list of requests with Approve/Deny buttons
- Dashboard badge for pending approvals (planned)
