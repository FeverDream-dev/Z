# LLM Router and Fallback

## Goal

Replace echo-only runtime with provider router, timeout policy, and fallback events.

## Inputs to read first

- `AGENTS.md`
- `docs/00-master-map.md`
- Relevant docs for this phase

## Tasks

- Define provider config.
- Implement mock providers for success/timeout/rate-limit.
- Route cheap/strong/coding tasks.
- Record fallback events.
- Expose provider health in UI.

## Acceptance criteria

- Timeout fixture triggers fallback.
- Job completes with fallback.
- Status is user-friendly.

## Required end-of-step report

```text
Changed:
- ...

Validated:
- command: ...
- result: ...

Manual check:
- ...

Risks / notes:
- ...

Next step:
- ...
```

## Human approval needed before

- Using real credentials.
- Installing global tools.
- Pushing to GitHub.
- Running external skill/plugin code.
