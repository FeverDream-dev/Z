# Zsistant Audit — Current State vs Target

## Current State: Chat Clone

The root project is currently a **chat-first skeleton** with many placeholder packages.

### What's real
| Package | Status |
|---|---|
| `internal/agents` | Basic registry + minimal struct (ID, Name, Role) |
| `internal/config` | YAML loader, secrets redaction, defaults |
| `internal/tools` | Tool broker, file read/write, shell exec, web fetch, permissions |

### What's a placeholder / stub / no-op
| Package | Issue |
|---|---|
| `internal/bus` | `InitBus()` is empty |
| `internal/channels` | `InitChannels()` is empty |
| `internal/sandbox` | `InitSandbox()` is empty |
| `internal/server` | `InitServer()` is empty |
| `internal/skills` | `InitSkills()` is empty |
| `internal/agents/agents.go` | `InitAgents()` is empty |
| `internal/trainer` | `InitTrainer()` is empty |
| `cmd/zazi/main.go` | Only basic agent CRUD (name+role). No persona, channels, tools, jobs. |

### Missing entirely
- `internal/memory` — no memory package
- `internal/knowledge` — no knowledge package
- `internal/jobs` — no job scheduling (vibe-kit has this but not in root)
- `internal/channels/*` — no Telegram, Discord, WhatsApp adapters
- `internal/browser` — no browser / MCP package
- `internal/devmode` — no developer mode infrastructure
- `internal/activity` — no activity timeline / logs package
- `internal/server/server.go` — no assistant-first HTTP routes
- `ui/` — no dashboard, no assistant detail tabs, no settings

### Root `go.mod`
- Only dependency: `gopkg.in/yaml.v3`
- No HTTP server imports
- No JSON streaming
- No channel integration libraries

### Vibe-kit has more code
`zsistant-vibe-kit/` contains a **more mature codebase** with:
- Real job queue + events + scheduling + retry + loop detection
- LLM router with streaming support (OpenAI, Ollama, mock)
- Telegram / Discord / WhatsApp channel adapters
- HTTP server with chat API, SSE streaming, provider health
- Web UI with model picker, settings modal, developer inspector
- Real provider catalog (52 providers, 207 models)
- But … it is still **chat-first**, not assistant-first.

## Target: Assistant Platform

Per `PRODUCT_SCOPE.md` and the briefing folder:

1. **Assistant is the central object** — not Message, not Conversation.
2. **Each assistant owns**: persona, tools, channels, memory, knowledge, jobs, permissions, logs.
3. **Chat is one tab** inside an assistant detail page.
4. **Dashboard** shows: assistants, jobs, channels, health, recent activity, pending approvals.
5. **Settings** has: appearance, providers, API keys, memory policy, channels, tools, Developer Mode.
6. **Developer Mode** exposes: raw requests, tool traces, provider routing, streaming events, MCP/browser diagnostics.
7. **No fake production behavior** — honest unavailable states everywhere.

## Rebuild Plan

1. Port vibe-kit real implementations into root (jobs, channels, LLM, server, UI).
2. Redesign `Agent` → `Assistant` with full profile.
3. Add missing packages: memory, knowledge, browser/MCP, devmode, activity.
4. Refactor server API to be assistant-first.
5. Rebuild Web UI as assistant control center.
6. Refactor CLI to be assistant-first.
7. Remove all placeholder stubs.
8. Build, test, audit for fake data.
