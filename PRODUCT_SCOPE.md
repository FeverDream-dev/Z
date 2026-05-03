# Zsistant / zazi — Product Scope

> This document is derived from `zsistant_openclaw_assistant_briefing/`.  
> Read that folder before editing code.

## One-sentence summary

Zsistant is a personal / team AI operating layer where the user **creates, configures, monitors, and commands real assistants** that can communicate through multiple channels, use tools, browse through Chrome / MCP, remember context, run scheduled jobs, coordinate with other agents, and expose advanced developer controls when needed.

## What this is NOT

- **Not** a ChatGPT clone.  
- **Not** a single generic chatbot.  
- **Not** a demo UI with fake integrations.

## Central object: Assistant

Everything orbits around `Assistant`.

- A conversation belongs to an assistant.
- A job belongs to an assistant.
- A tool permission belongs to an assistant.
- A Telegram channel routes to an assistant.
- A memory namespace belongs to an assistant.
- A browser session can be assigned to an assistant.

## Product hierarchy (correct)

1. Workspace
2. Assistants
3. Channels
4. Tools
5. Memory & knowledge
6. Jobs & automations
7. Conversations
8. Messages

## Required product areas

| Area | Must exist | Notes |
|---|---|---|
| Assistant manager | Yes | Real assistants, not fake cards |
| Assistant detail / control center | Yes | Overview + tabs |
| Chat tab | Yes | One surface of an assistant |
| Assistant profile / persona | Yes | Identity, tone, role, boundaries |
| Tools registry | Yes | Per-assistant tool enablement + permissions |
| MCP / browser / Chrome area | Yes | Honest unavailable state if not connected |
| Channels area | Yes | Web UI, CLI, Telegram, WhatsApp, Discord, Slack, Email — honest setup states |
| Knowledge area | Yes | Files, docs, notes, project context |
| Memory area | Yes | Inspectable / editable durable context |
| Jobs & automations | Yes | Scheduled + event-triggered work |
| Logs / activity timeline | Yes | Human-readable chronology |
| Settings cog & settings page | Yes | Real sections |
| Developer Mode | Yes | Toggle + meaningful advanced UI |
| Provider / model configuration | Yes | Honest missing-credential states |

## Honesty rules

- No fake success.
- No fake job history.
- No fake screenshots.
- No fake browser actions.
- No fake Telegram / Discord / WhatsApp integrations.
- If a feature is not connected, show a setup state.
- If a feature is not built, say it is unavailable in this build.
- Every visible element must correspond to real state.

## Maturity levels

| Level | Meaning |
|---|---|
| 0 - Not available | Feature is not built. Do not pretend. |
| 1 - Configurable but inactive | Setup UI exists, but not configured. |
| 2 - Connected | Credentials / channel / tool connected and can be tested. |
| 3 - Usable manually | User can invoke it from chat / UI. |
| 4 - Automatable | Assistant can use it in jobs / automations. |
| 5 - Observable and reliable | Logs, retries, error states, permissions, audits exist. |

## UI principles

- **Assistant-centered navigation**.  Normal users see: Assistants, Projects, Jobs, Channels, Tools, Settings.
- **Simple first, powerful on demand**.  Developer details hidden unless Developer Mode is on.
- **Status is visible**.  Connected / disabled / needs-setup / failed / running / approval-needed.
- **Honest empty states**.  "No assistants yet. Create your first assistant: a persistent AI worker with tools, memory, channels, and scheduled jobs."
- **Polished but functional**.  Every visible element is real.  No fake graphs, cards, or activity.

## Technology constraints

- Go 1.20+
- YAML config, JSON / JSONL file-backed persistence (no database for MVP)
- No frontend build step — raw HTML / CSS / JS served by Go
- Honest states for all integrations

## Decision rule

When choosing between two features, pick the one that makes assistants more capable, observable, persistent, and useful across time.
