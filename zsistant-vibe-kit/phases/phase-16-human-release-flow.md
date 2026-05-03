# Human Release Flow

## Goal

Prepare release automation gated by explicit human approval.

## Inputs to read first

- `AGENTS.md`
- `docs/00-master-map.md`
- Relevant docs for this phase

## Tasks

- Add release checklist.
- Generate changelog draft.
- Run tests before release.
- Ask explicit approval.
- Only then push/tag/publish.

## Acceptance criteria

- Prepare command does not publish.
- Approval wording is required.
- Release notes mention tests and risks.

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
