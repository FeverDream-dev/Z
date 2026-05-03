# 11 — Web UI Map

The web UI should be the control tower.

## Primary screens

### Dashboard

- daemon status
- active agents
- active jobs
- provider health
- channel health
- recent approvals needed

### Agent detail

- persona summary
- workspace summary
- connected channels
- active jobs
- memory notes
- peer permissions
- tool permissions

### Chat

- local web chat with selected agent
- job status sidebar
- attachments
- approval prompts

### Channels

- Telegram setup
- Discord setup
- WhatsApp setup
- local web chat links
- channel-to-agent bindings

### Jobs

- chronological jobs
- current step
- retries/fallbacks
- logs
- resume/pause/cancel

### Skills

- scan skill
- risk report
- translation plan
- approve generated implementation

### Releases

- test status
- changelog draft
- diff summary
- human approval button/command

## UX requirements

- Always show whether an agent is active, paused, waiting, or failed.
- Always show when a job is retrying or switching provider.
- Do not expose raw secrets.
- Make destructive actions ask for confirmation.
- Make channel setup beginner-friendly.

## Design direction

Visual identity:

- Dark control tower mode.
- Clear status chips.
- Agent cards.
- Activity timeline.
- Friendly but not childish.

Use assets from `assets/brand/` and mockups from `assets/mockups/` as starting points.
