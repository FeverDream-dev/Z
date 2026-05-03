# Telegram Adapter

## Goal

Add Telegram setup and test-message routing.

## Inputs to read first

- `AGENTS.md`
- `docs/00-master-map.md`
- Relevant docs for this phase

## Tasks

- Implement token config storage without logging token.
- Support long polling for local dev.
- Map chat ID to agent.
- Send response back.
- Add dry-run/test mode.

## Acceptance criteria

- Test mode works without real token.
- Real token use requires explicit user action.
- Inbound message creates job.

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
