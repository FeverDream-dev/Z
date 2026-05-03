# Zsistant

A lightweight, local-first multi-agent AI workspace. Run agents on your machine, connect to 50+ LLM providers, chat via web UI or messaging channels, and keep everything under your control.

**Motto:** Many agents. One light host. Your machine stays in control.

## What it is

Zsistant is a Go binary (`zazi`) that serves a local web UI and exposes an HTTP API for agent management, streaming chat, and LLM provider routing. It stores all data locally in `~/.zazi/` — no database, no cloud dependency, no telemetry.

## Quick start

```bash
# Build
go build -o zazi ./cmd/zazi

# Initialize config and data directories
./zazi init

# Start the server
./zazi serve
# → Listening on http://0.0.0.0:8080

# Open in browser
open http://localhost:8080
```

## Configuration

All config lives in `~/.zazi/config.yaml`. Secrets (API keys, tokens) are stored in the `secrets` map and never appear in logs.

You can manage secrets via the web UI (Settings modal) or by editing the file directly:

```yaml
# ~/.zazi/config.yaml
data_path: ~/.zazi
log_level: info
server_port: 8080
secrets:
  ollama_api_key: "your-ollama-key"
  openai_api_key: "sk-..."
  openai_base_url: "https://api.openai.com/v1"
  telegram_my-agent_token: "123456:ABC-DEF"
```

See `.env.example` for the full list of supported secret keys.

## CLI commands

```
zazi serve                      Start web server (default :8080)
zazi init                       Initialize config and data dirs
zazi doctor                     Run diagnostics
zazi version                    Print version info

zazi chat --agent=<id> --message=<msg>       Chat from CLI
zazi provider list                            List 52 providers / 207 models
zazi provider show --name=ollama              Show provider details

zazi agent create <id> --name=<n> --role=<r>  Create an agent
zazi agent list                                List agents
zazi agent show <id>                           Show agent details
zazi agent delete <id>                         Delete an agent

zazi telegram setup --agent=<id> --token=<t>  Configure Telegram
zazi telegram listen --agent=<id>             Start listening
zazi telegram send --agent=<id> --chat=<c> --message=<m>

zazi job list --agent=<id>                    List jobs
zazi job pause --agent=<id> --job=<jid>
zazi job resume --agent=<id> --job=<jid>
zazi job cancel --agent=<id> --job=<jid>

zazi train observe --agent=<id> --message=<msg>
zazi train propose --agent=<id>
zazi train apply --agent=<id> --patch=<text>

zazi acl allow --agent=<id> --peer=<pid> --perm=<p>
zazi acl revoke --agent=<id> --peer=<pid>
zazi acl list --agent=<id>

zazi skill analyze --path=<folder>
zazi validate ui
zazi release prepare --version=vX.Y.Z
```

## Web UI

The server serves a premium dark/light themed UI at `/` and `/chat`:

- **Sidebar** with conversation list, new-chat button, and delete
- **Streaming chat** via SSE (`POST /api/chat/stream`)
- **Model picker** with 207 models across 52 providers
- **Settings modal** to manage API keys (saved to `~/.zazi/config.yaml`)
- **Developer Mode** inspector to view raw SSE events and provider metadata

## API endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/` | Web UI |
| GET | `/chat` | Web UI (alias) |
| POST | `/api/chat` | Send message, get response |
| POST | `/api/chat/stream` | SSE streaming chat |
| GET | `/api/providers` | Active provider health |
| GET | `/api/models` | All models (id, name, provider) |
| GET | `/api/settings` | Get secrets (redacted values) |
| POST | `/api/settings` | Update secrets |
| GET | `/api/conversations` | List conversations (JSONL) |
| POST | `/api/conversations` | Save conversation |

## LLM providers

4 providers are wired with real HTTP calls:

| Provider | Secret key | Streaming |
|----------|-----------|-----------|
| Ollama Cloud | `ollama_api_key` | Yes |
| OpenAI-compatible | `openai_api_key` + `openai_base_url` | No (planned) |
| Z.AI Coding Plan | `zai_api_key` | No (planned) |
| OpenCode Zen | `opencode_api_key` | No (planned) |

48 additional providers are cataloged in `internal/llm/registry.go` (207 models total) and available in the model picker. They become active once their API key is configured.

## Channel adapters

| Channel | Status | Command |
|---------|--------|---------|
| Telegram | Real long-poll, JSON parsing, graceful shutdown | `zazi telegram listen` |
| Discord | Dry-run mode (validate token, simulate events) | `zazi discord test` |
| WhatsApp | Dry-run mode (webhook verify, simulate events) | `zazi whatsapp test` |

## Architecture

```
cmd/zazi/main.go           CLI entry point (15+ subcommands)
internal/
  server/server.go          HTTP server, SSE streaming, API routes
  llm/
    provider.go             Provider + Streamer interfaces
    router.go               Router with fallback chains
    registry.go             52 providers, 207 models
    ollama.go               Ollama Cloud (real)
    openai.go               OpenAI-compatible (real)
    zai_coding.go           Z.AI Coding Plan (real)
    opencode.go             OpenCode Zen (real)
    mock.go                 Echo fallback (always last)
    stream.go               Streamer interface + Ollama SSE
  channels/
    telegram.go             Real Telegram adapter
    discord.go              Discord adapter (dry-run)
    whatsapp.go             WhatsApp adapter (dry-run)
  agents/                   Agent registry (JSON-backed)
  jobs/                     Job queue (JSONL-backed)
  bus/                      Inter-agent ACL bus
  skills/                   Skill risk analyzer
  trainer/                  Persona trainer
  config/                   Config load/save (~/.zazi/config.yaml)
ui/
  index.html                Premium UI template
  app.css                   Dark/light themes, responsive
  app.js                    SSE chat, settings, dev mode
```

All persistence is file-backed (JSON/JSONL/YAML) under `~/.zazi/`. No database required.

## Development

```bash
go mod tidy
go build ./...
go test ./...
```

For detailed development instructions, testing, and debugging, see [DEVELOPER_MODE.md](DEVELOPER_MODE.md).

## Stack

- Go 1.25.9
- Single dependency: `gopkg.in/yaml.v3`
- Vanilla HTML/CSS/JS frontend (no framework, no build step)
- CI: `.github/workflows/ci.yml` runs `go test ./...` on push/PR

## Core identity

- Product name: **Zsistant**
- CLI command: **`zazi`**
- Coding agent name: **Z Assistant**
- Repository: `github.com/FeverDream-dev/zsistant`

## Philosophy

Zsistant should feel like a lightweight local control tower for many personal and company agents. Simple enough for a non-technical user, but structured enough for a developer to inspect every decision, endpoint, prompt, and permission.

The main differentiator is not "agent features" — it is **speed, loop detection, proactivity, stability, isolation, and running many agents on modest hardware**.
