# 01 — Product Brief

## Name

**Zsistant**

## CLI

`zazi`

## One-line promise

Run many isolated personal/company agents on one lightweight local machine.

## Primary users

1. Non-technical users who want a personal AI assistant through Telegram, Discord, WhatsApp, web, or CLI.
2. Developers who want to vibe-code and inspect every step.
3. Small companies that want many role-specific agents doing separate tasks.
4. Power users migrating from OpenClaw-like systems who want simpler, safer local control.

## Core user stories

### Personal agent

A user installs Zsistant, creates an agent, connects Telegram, and messages it like a personal assistant.

### Company swarm without chaos

A company creates many agents: sales, support, finance, devops, QA, research. They cannot see each other's files unless explicitly allowed.

### Agent-to-agent request

The user allows Agent A to ask Agent B for a file, summary, or decision. The request is typed, logged, and permissioned.

### LLM fallback

A job starts on one model. If the provider times out, rate-limits, or fails, Zsistant retries or switches providers without losing job state.

### Persona trainer

The agent adapts to the user: short/direct for sales, technical for developers, friendly/street for casual users, formal for business users.

### OpenClaw migration

The user asks the old tool to export a big Markdown summary. Zsistant ingests the summary, extracts agent settings, personas, workflows, and skills, then proposes a safe migration plan.

## Product principles

- Local-first.
- Human-readable state.
- Isolated by default.
- No free background chatter.
- Visible progress.
- Human-approved release.
- External skills are untrusted until reviewed.
- Cheap models for cheap tasks; strong models for hard tasks.
- Restartable jobs.

## MVP definition

MVP is reached when a user can:

1. Install and run `zazi` locally.
2. Open the web UI.
3. Create two isolated agents.
4. Chat with one agent through CLI and web.
5. Connect Telegram and Discord test channels.
6. See job status and retries.
7. Configure at least two LLM providers with fallback.
8. Allow two agents to exchange a typed request.
9. Scan an external skill folder and get a risk report.
10. Run tests and manually approve a GitHub release.
