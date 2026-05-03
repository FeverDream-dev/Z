# Local Chat Loop

## Goal

Send a CLI message to an agent through the runtime using a mock/echo provider.

## Inputs to read first

- `AGENTS.md`
- `docs/00-master-map.md`
- Relevant docs for this phase

## Tasks

- Create normalized inbound message type.
- Create job record.
- Route to echo/mock LLM provider.
- Append audit/job events.
- Return response in CLI.

## Acceptance criteria

- `zazi chat --agent ...` returns a response.
- Job events are written.
- No external LLM key required.

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
