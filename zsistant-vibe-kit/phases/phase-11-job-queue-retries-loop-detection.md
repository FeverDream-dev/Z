# Job Queue, Retries, and Loop Detection

## Goal

Make jobs durable, resumable, and loop-aware.

## Inputs to read first

- `AGENTS.md`
- `docs/00-master-map.md`
- Relevant docs for this phase

## Tasks

- Persist job checkpoints.
- Add retry backoff.
- Detect repeated failed actions.
- Pause loop-suspected jobs.
- Expose resume/cancel.

## Acceptance criteria

- Repeated failure fixture pauses job.
- Resume command exists.
- UI shows stalled/paused status.

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
