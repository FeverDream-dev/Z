# 06 — LLM Router and Resilience

## Purpose

The LLM router prevents one failed provider from killing a job.

It also prevents expensive models from being used for cheap tasks by default.

## Provider policy

Each agent can have a model policy:

```text
primary: provider/model
fallbacks:
  - provider/model
  - provider/model
cheap_model: provider/model
strong_model: provider/model
coding_model: provider/model
vision_model: provider/model, optional
max_retries_per_provider: N
timeout_seconds: N
budget:
  daily_tokens
  daily_cost
  job_tokens
```

## Provider types to support early

- OpenAI-compatible HTTP APIs.
- Ollama local and cloud-compatible endpoints.
- Z.AI Coding Plan endpoint for coding workflows.
- Echo/mock provider for local tests.

## Routing strategy

### Routine tasks

Use cheap/fast model.

Examples:

- classify message
- choose persona style
- summarize short status
- detect whether human approval is needed

### Planning and coding

Use coding/strong model.

Examples:

- code generation
- migration plans
- security review
- complex browser failures

### Fallback

Fallback on:

- provider timeout
- rate limit
- quota exhaustion
- transient network failure
- model unavailable
- configured cost cap

Do not fallback silently forever. Record status events and stop if loop detector fires.

## User-facing status language

Bad:

```text
HTTP 504 timeout
```

Better:

```text
Still working. The current model timed out, so I am retrying and will switch providers if needed.
```

## Anti-loop policy

Retries should be finite per checkpoint. A job may keep trying toward the objective, but only with visible state changes, backoff, and loop detection.
