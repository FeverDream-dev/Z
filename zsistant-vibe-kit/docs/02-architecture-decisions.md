# 02 — Architecture Decisions

## ADR-001: One daemon, many logical agents

Decision: Start with one daemon process that manages many logical agents.

Reasoning:

- Lowest memory overhead.
- Easier multi-agent scheduling.
- Easier web UI and local API.
- Better for Raspberry Pi and small servers.

Future option: Allow dangerous or expensive jobs to run in isolated worker processes or containers.

## ADR-002: Isolated workspace roots

Each agent has a root path. Tools must operate inside that root unless a human grants explicit external access.

Suggested root:

```text
~/.zazi/agents/<agent-id>/
```

Suggested children:

```text
persona.md
profile.json
memory/
workspace/
jobs/
audit.jsonl
inbox.jsonl
outbox.jsonl
```

## ADR-003: Human-readable persona files

Every agent should have a `persona.md` that can be edited by a human.

The persona trainer may propose patches, but should not silently rewrite the whole personality.

## ADR-004: Provider router, not provider lock-in

LLM calls go through a router.

The router decides:

- Which provider/model to use.
- Whether to retry.
- Whether to switch provider.
- Whether to pause for human approval.
- Whether a job is likely looping.

## ADR-005: Skill analysis before skill execution

External skills/plugins are treated as untrusted content.

MVP behavior:

1. Read skill metadata and instructions.
2. Summarize capabilities.
3. Detect risky commands/instructions.
4. Propose a Zsistant-native implementation plan.
5. Ask the user before installing or executing anything.

## ADR-006: Event log before complex observability

Use append-only JSONL audit/status logs first. Add metrics dashboards later.

Every job should write status events:

- queued
- started
- model_selected
- tool_called
- retry_scheduled
- provider_failed
- fallback_selected
- waiting_for_human
- completed
- failed
- loop_suspected
- resumed

## ADR-007: Embedded web UI

The daemon should serve the web UI. This simplifies install and makes `zazi serve` enough for local use.

## ADR-008: No release without human approval

Z Assistant may prepare commits, changelogs, tags, and release notes. It must ask before pushing or publishing.
