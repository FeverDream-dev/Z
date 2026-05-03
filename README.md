# Zsistant / zazi — Assistant-First AI Operating Platform

> **Important:** This is NOT a ChatGPT clone. It is an assistant platform where the user creates, manages, and commands multiple AI assistants.

## What Zsistant Is

Zsistant is a personal/team AI operating layer. The user creates **assistants** — persistent digital workers with their own identity, tools, channels, memory, knowledge, jobs, and permissions.

## Philosophy

- **Central object is Assistant.** Not Message. Not Conversation.
- **Chat is one surface.** One tab inside an assistant detail page, alongside Overview, Jobs, Channels, Tools, Knowledge, Memory, Browser, Logs, and Settings.
- **Honest capability states.** If a feature is not configured, the UI shows a setup state. If it is not built, it says so. No fake integrations, no fake job history, no fake browser screenshots.
- **Multi-channel.** Assistants can be reached via Web UI (today), and designed for Telegram, Discord, Slack, WhatsApp, CLI, and Email.
- **Developer Mode.** A real toggle that exposes raw traces, provider diagnostics, tool calls, and streaming events to power users.

## Project Structure

```
cmd/zazi/                    CLI entrypoint
internal/
  assistant/                 Core assistant model + registry + jobs
  config/                    YAML config loader with provider key support
  llm/                       LLM provider interface, router, streaming
  channels/                  Telegram, Discord, WhatsApp adapters
  tools/                     Tool broker, permissions, built-in tools
  jobs/                      Job scheduling, queue, events, retry logic
  memory/                    Memory store (edit, global vs assistant-scoped)
  knowledge/                 Knowledge sources (uploads, indexing status)
  browser/                   Browser/MCP session with honest unavailable state
  activity/                  Activity timeline / logs package
  devmode/                   Developer traces, diagnostics, toggle
  server/                    Assistant-first HTTP server with SSE streaming
ui/                          Assistant control center HTML/CSS/JS
PRODUCT_SCOPE.md             Product scope document
AUDIT.md                      Audit of current state vs target
```

## Quick Start

```bash
# Build everything
go build ./...

# Initialize config
go run ./cmd/zazi init

# Start the web server
go run ./cmd/zazi serve --addr=:8080

# Or use the CLI
go run ./cmd/zazi assistant create my-assistant --name="My Assistant"
go run ./cmd/zazi assistant list
```

## API Endpoints (Assistant-First)

| Method | Path | Description |
|--------|------|-------------|
| GET | `/` | Dashboard UI |
| GET | `/api/assistants` | List assistants |
| POST | `/api/assistants` | Create assistant |
| GET | `/api/assistants/:id` | Get assistant detail |
| PUT | `/api/assistants/:id` | Update assistant |
| DELETE | `/api/assistants/:id` | Delete assistant |
| GET | `/api/assistants/:id/channels` | Assistant channels |
| GET | `/api/assistants/:id/memory` | Assistant memory |
| GET | `/api/assistants/:id/knowledge` | Assistant knowledge |
| GET | `/api/assistants/:id/jobs` | Assistant jobs |
| GET | `/api/assistants/:id/logs` | Assistant activity logs |
| GET | `/api/assistants/:id/browser` | Browser/MCP status |
| POST | `/api/assistants/:id/chat` | Chat with assistant |
| POST | `/api/assistants/:id/chat/stream` | SSE streaming chat |
| GET | `/api/tools` | Tool registry |
| GET | `/api/providers` | Provider health |
| GET | `/api/models` | Model catalog |
| GET/POST | `/api/settings` | Settings |
| GET | `/api/jobs` | Global jobs |
| GET | `/api/activity` | Global activity |

## Honest States

- **Browser/MCP:** Shows "not connected" with setup instructions.
- **Channels:** Show "needs setup" when tokens are missing.
- **Providers:** Return real health. If no API key, return honest "unconfigured" with guidance.
- **Jobs:** Show "no jobs scheduled yet" if none exist. No fake cron history.
- **Memory/Knowledge:** Show empty but educative states.

## UI Views

1. **Home Dashboard** — Active assistants, channels, jobs, recent activity.
2. **Assistants Grid** — Create, manage, delete assistants.
3. **Assistant Detail** — Tabs: Overview, Chat, Channels, Tools, Knowledge, Memory, Jobs, Browser, Logs, Settings, Developer.
4. **Settings** — Theme, Developer Mode, API keys.
5. **Jobs / Channels / Tools** — Global views.

## Tech Stack

- **Go 1.20+**
- **File-backed persistence** (JSON/JSONL/YAML) — no database for MVP
- **No frontend build step** — raw HTML/CSS/JS served by Go
- **Honest UI** — every visible element corresponds to real state

## License

MIT
