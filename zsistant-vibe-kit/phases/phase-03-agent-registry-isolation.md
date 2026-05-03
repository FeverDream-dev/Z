# Agent Registry and Isolation

## Goal

Create agents with isolated workspaces and path guard.

## Inputs to read first

- `AGENTS.md`
- `docs/00-master-map.md`
- Relevant docs for this phase

## Tasks

- Create/list/show agent profile.
- Create persona.md per agent.
- Create workspace directory.
- Implement path-scope guard.
- Test denied cross-agent path access.

## Acceptance criteria

- Can create two agents.
- Each has its own workspace.
- Agent A cannot read Agent B workspace by default.

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
