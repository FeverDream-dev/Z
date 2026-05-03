# Config and Local State

## Goal

Create local config loading and data directory management.

## Inputs to read first

- `AGENTS.md`
- `docs/00-master-map.md`
- Relevant docs for this phase

## Tasks

- Define config file format.
- Implement default path resolution.
- Create data directories.
- Protect secrets from logs.
- Document config examples.

## Acceptance criteria

- `zazi init` creates config/data folders.
- Repeated init is idempotent.
- Secrets are never printed.

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
