# 12 — CLI Map

CLI command: `zazi`

## Command groups

```text
zazi doctor
zazi init
zazi serve
zazi agent
zazi chat
zazi job
zazi channel
zazi peer
zazi skill
zazi persona
zazi release
```

## Command concepts

### doctor

Checks config, data directory, provider config, writable paths, web port availability.

### init

Creates default local config and data directories.

### serve

Starts daemon + web UI.

### agent

Create/list/show/update/pause/delete agents.

### chat

Send a local message to an agent.

### job

List/status/resume/cancel jobs.

### channel

Setup/bind/test Telegram, Discord, WhatsApp.

### peer

Allow/revoke/list inter-agent communication rules.

### skill

Analyze/scan/import-plan external skills.

### persona

Show/propose/apply persona trainer patches.

### release

Prepare release, show checklist, request human approval.

## UX principle

Every command should have:

- clear success output
- helpful next step
- machine-readable mode eventually
- no secret printing
