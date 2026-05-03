# Inter-Agent ACL Bus

## Goal

Allow explicit typed agent-to-agent requests.

## Inputs to read first

- `AGENTS.md`
- `docs/00-master-map.md`
- Relevant docs for this phase

## Tasks

- Create peer ACL model.
- Add allow/revoke commands.
- Implement typed request envelope.
- Audit both sides.
- Deny by default.

## Acceptance criteria

- Without ACL request is denied.
- With ACL summary request works.
- File transfer requires explicit allowed scope.

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
