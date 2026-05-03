# AGENTS.md — Instructions for OpenCode and Coding Agents

You are building **Zsistant**, a lightweight multi-agent host with the `zazi` CLI.

This folder begins as a **specification and asset pack only**. The application source code must be generated gradually by OpenCode, phase by phase, so the human operator can inspect progress in real time.

## Non-negotiable rules

1. Read `docs/00-master-map.md` before writing any code.
2. Implement one phase at a time from `phases/`.
3. Do not skip tests. Every phase must end with a test command, manual verification steps, and a summary of changed files.
4. Do not push to GitHub, tag, publish, or create releases without explicit human approval.
5. Do not execute third-party skills from ClawHub or any other marketplace by default. Analyze and translate only.
6. Do not use real Telegram, Discord, WhatsApp, email, browser, or GitHub credentials unless the user explicitly provides them for that session.
7. Do not place secrets in source files, logs, screenshots, or docs.
8. Never let one agent read another agent's workspace unless the ACL explicitly allows it.
9. Never let agent-to-agent chatter become free-form background noise. Use typed job/message envelopes and audit logs.
10. Detect loops, stalled jobs, repeated tool failures, and token/provider failures. Recover or pause with a clear status.

## Implementation style

Prefer a small, boring, inspectable architecture over clever abstractions.

Recommended stack unless the human overrides it:

- Core daemon and CLI: Go.
- Local state: SQLite once schemas stabilize; file-backed JSONL is acceptable only for the very first prototype.
- Web UI: embedded static assets; TypeScript is allowed for the UI build, but runtime should remain simple.
- Messaging channels: adapters behind a common inbound/outbound interface.
- LLMs: provider router with retries, fallbacks, budgets, and job-level status.
- Browser testing: use Chrome DevTools MCP when configured by the user.

## Communication style

Keep the human in the loop. After each coding step, report:

- What changed.
- What commands were run.
- What passed or failed.
- What should be tested manually.
- What the next smallest step is.

## Definition of done for every phase

A phase is not done until:

- The code compiles or the project scaffolding command succeeds.
- Unit tests or equivalent checks pass.
- Any web UI route has a basic manual test.
- Security implications are noted.
- The human can inspect the diff and understand the result.
