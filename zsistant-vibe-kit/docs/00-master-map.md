# 00 — Zsistant Master Map

Zsistant is a lightweight local-first multi-agent host.

The target user can install it, open a web UI, create multiple isolated agents, connect channels such as Telegram/Discord/WhatsApp/local chat/CLI, and let each agent run jobs with clear status, memory boundaries, LLM fallback, and human-approved release automation.

## Product thesis

OpenClaw-like capability can be made lighter by separating the ideas:

- One small local supervisor daemon.
- Many isolated logical agents.
- Per-agent workspaces and permissions.
- Channel adapters that feed messages into the same job bus.
- LLM provider router with fallback and budgets.
- Persona trainer that adapts tone and workflow to the user.
- Skill importer that analyzes external skills but does not run them by default.
- Human-visible status for long loops.

## Target architecture

```text
User / Channels
  ├─ CLI: zazi
  ├─ Local Web Chat
  ├─ Telegram Bot Adapter
  ├─ Discord Adapter
  └─ WhatsApp Cloud API Adapter
        │
        ▼
Zsistant Daemon
  ├─ HTTP/Web UI Server
  ├─ Agent Registry
  ├─ Job Queue
  ├─ LLM Router
  ├─ Tool Broker
  ├─ Persona Trainer
  ├─ Skill Analyzer
  ├─ Inter-Agent ACL Bus
  └─ Audit/Status Log
        │
        ▼
Agent Workspaces
  ├─ agents/<id>/persona.md
  ├─ agents/<id>/memory/
  ├─ agents/<id>/workspace/
  ├─ agents/<id>/jobs/
  └─ agents/<id>/audit.jsonl
```

## Default architectural choice

Use **one daemon managing many logical agents**.

Why:

- Cheaper than one daemon per agent.
- Easier to run 10–100 agents on modest hardware.
- Centralized job scheduling and status.
- Centralized provider fallback.
- Easier web UI.

Isolation is achieved through:

- Per-agent workspace roots.
- Path-scope checks.
- Per-agent tool permissions.
- Per-agent channel permissions.
- Inter-agent ACLs.
- Optional future worker process/container isolation for dangerous tools.

## What must feel different

Zsistant should not feel like a black-box agent that vanishes for minutes.

Every job should expose:

- Current step.
- Last useful action.
- Current provider/model.
- Retry count.
- Next retry time.
- Loop/stall suspicion.
- Whether human input is needed.

## MVP channels

Mandatory:

- CLI.
- Local web chat.
- Telegram.
- Discord.
- WhatsApp.

Recommended implementation order:

1. CLI and local web chat.
2. Telegram.
3. Discord.
4. WhatsApp.

## What is not MVP

- Full email automation.
- Full browser control for end users.
- Third-party skill execution.
- Multi-machine clustering.
- Public marketplace.
- Payment/subscription backend.

These can be designed now but implemented later.
