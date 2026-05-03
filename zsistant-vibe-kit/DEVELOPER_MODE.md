# Developer Mode Guide

This document covers building, testing, debugging, and extending Zsistant.

## Prerequisites

- Go 1.25.9+
- No external services required (mock provider works offline)
- Optional: Ollama Cloud API key for live LLM testing

## Build & Run

```bash
# Build the binary
go build -o zazi ./cmd/zazi

# Run unit tests
go test ./...

# Run with verbose test output
go test -v ./...

# Run a specific package's tests
go test -v ./internal/server/
go test -v ./internal/llm/
go test -v ./internal/config/

# Initialize and start
./zazi init
./zazi serve
# → http://0.0.0.0:8080

# Build with version info
go build -ldflags "-X main.zaziVersion=v0.1.0 -X main.zaziCommit=$(git rev-parse HEAD)" -o zazi ./cmd/zazi
./zazi version
# → zazi version v0.1.0 (commit abc123...)
```

## Project Structure

```
cmd/zazi/main.go           CLI hub — all subcommands live here
internal/
  server/server.go          HTTP server, routes, SSE streaming (~560 lines)
  server/server_test.go     Server tests
  llm/
    provider.go             Provider + ProviderHealth interfaces
    router.go               Router: ordered fallback, ChainFor()
    registry.go             52 providers, 207 models (BuiltInProviders)
    ollama.go               Ollama Cloud — Complete() + Stream()
    openai.go               OpenAI-compatible — Complete()
    zai_coding.go           Z.AI Coding Plan — Complete()
    opencode.go             OpenCode Zen — Complete()
    mock.go                 Echo fallback — always registered last
    stream.go               Streamer interface
  channels/
    telegram.go             Real adapter: Listen(), SendMessage(), JSON parsing
    discord.go              Dry-run: validate token, simulate events
    whatsapp.go             Dry-run: webhook verify, simulate events
  agents/registry.go        Agent CRUD (JSON-backed)
  jobs/queue.go             Job queue with pause/resume/cancel
  jobs/loop.go              Loop detector (stall/retry detection)
  bus/bus.go                Inter-agent ACL bus
  skills/analyzer.go        External skill risk analysis
  trainer/trainer.go        Persona trainer (observe, propose, apply)
  config/
    config.go               Config struct (Secrets map[string]string)
    defaults.go             Default values
    loader.go               Load/Save/DefaultPath/EnsureDirs
ui/
  index.html                Premium UI (sidebar, chat, settings, dev inspector)
  app.css                   Dark/light/system themes, accent colors, responsive
  app.js                    SSE streaming chat, settings management, dev mode
```

## Data Directory Layout

All data lives under `~/.zazi/`:

```
~/.zazi/
  config.yaml               Main config + secrets
  conversations.jsonl       Chat history (JSONL, one JSON object per line)
  agents/
    <agent-id>.json         Agent definitions
  logs/                     Daemon logs
  cache/                    Provider response cache
  backups/                  Config backups
```

## Server Architecture

The server (`internal/server/server.go`) wires routes in `buildRouter()`:

1. Reads `~/.zazi/config.yaml` for secrets
2. Registers providers based on available keys:
   - `ollama_api_key` → Ollama Cloud (with streaming)
   - `openai_api_key` → OpenAI-compatible
   - `zai_api_key` → Z.AI Coding Plan
   - `opencode_api_key` → OpenCode Zen
   - Mock provider always registered as final fallback
3. Routes:
   - `/`, `/chat` → serve `ui/index.html` from disk
   - `/api/chat` → non-streaming response
   - `/api/chat/stream` → SSE streaming (`event: chunk`, `event: done`, `event: error`)
   - `/api/providers` → active provider health status
   - `/api/models` → flat model catalog from BuiltInProviders
   - `/api/settings` → GET/POST secrets management
   - `/api/conversations` → GET/POST JSONL persistence

## SSE Streaming Protocol

The frontend (`ui/app.js`) connects to `POST /api/chat/stream`:

```
Request:  {"message": "Hello", "model": "llama3"}
Response: Content-Type: text/event-stream

event: chunk
data: "Hello"

event: chunk
data: " world"

event: done
data: {"job_id":"...","provider":"ollama","status":"completed"}
```

Error case:

```
event: error
data: {"error":"provider timeout"}
```

The SSE parser in `app.js` tracks `event:` and `data:` lines across chunks for correct parsing.

## Testing

### Unit Tests

```bash
go test ./...                    # All tests
go test -v ./internal/server/    # Server tests with output
go test -race ./...              # Race detector
```

### Server Test Details

Server tests (`internal/server/server_test.go`) verify:
- `TestDashboardPage` — expects `<title>Zsistant</title>` in response
- `TestChatPage` — same as dashboard (both serve `ui/index.html`)
- `TestChatAPI` — accepts any non-empty response from provider chain
- `TestHealthEndpoint` — returns provider health array

The `handleApp()` function searches for `ui/index.html` at two paths (relative and `../../ui/`) for test compatibility.

### Manual Testing

```bash
# Start server
./zazi serve

# Health check
curl http://localhost:8080/api/providers

# Non-streaming chat
curl -X POST http://localhost:8080/api/chat \
  -H "Content-Type: application/json" \
  -d '{"message":"Hello"}'

# Streaming chat
curl -N -X POST http://localhost:8080/api/chat/stream \
  -H "Content-Type: application/json" \
  -d '{"message":"Hello","model":"llama3"}'

# Model catalog
curl http://localhost:8080/api/models | jq '. | length'
# → 207

# Settings
curl http://localhost:8080/api/settings
```

### Browser Testing

Open `http://localhost:8080` in a browser. The UI should:
- Show sidebar with conversation list
- Show chat input at bottom
- Show model picker with 207 models
- Show Settings gear icon (opens modal for API keys)
- Show Developer toggle (enables raw event inspector)

## Adding a New LLM Provider

1. Create `internal/llm/myprovider.go`
2. Implement the `Provider` interface:

```go
type Provider interface {
    Complete(prompt string) (string, error)
    Health() ProviderHealth
}
```

3. Optionally implement `Streamer` for SSE:

```go
type Streamer interface {
    Stream(prompt string, chunkCh chan<- string, doneCh chan<- StreamResult)
}
```

4. Register in `buildRouter()` (`internal/server/server.go`) when the relevant secret key is present
5. Add provider info to `BuiltInProviders` in `internal/llm/registry.go`
6. Run `go test ./...` to verify

## Adding a New Channel Adapter

1. Create `internal/channels/mychannel.go`
2. Implement adapter with `Listen(ctx, handler)` and `SendMessage()` methods
3. Add CLI subcommands in `cmd/zazi/main.go`
4. Start with dry-run mode (simulate events), upgrade to real API when ready

## Debugging

### Verbose Logging

Set `log_level: debug` in `~/.zazi/config.yaml`:

```yaml
log_level: debug
```

### Doctor Check

```bash
./zazi doctor
# Checks: config file exists, data path writable, prints config summary
```

### Provider Health

```bash
curl http://localhost:8080/api/providers
# Returns array of {name, status, latency_ms}
```

### Developer Mode (UI)

Toggle Developer Mode in the web UI to see:
- Raw SSE events in real-time
- Provider name and response metadata
- Job IDs for each response

## Known Limitations

- **Streaming** only implemented for Ollama. OpenAI, Z.AI, and OpenCode providers fall back to `Complete()` then stream the full response.
- **No authentication** — the server binds `0.0.0.0:8080` with zero auth middleware. Do not expose to public networks without a reverse proxy.
- **Model picker UX** — 207 models in a flat dropdown is unwieldy. Grouping by provider or adding search is planned.
- **Conversations** are saved to localStorage in the frontend. Server-side persistence via `/api/conversations` is wired but the frontend still uses localStorage as the primary store.
- **Discord/WhatsApp** adapters are dry-run only. They validate tokens and simulate events but do not make real API calls.

## CI

`.github/workflows/ci.yml` runs `go test ./...` on every push and pull request.
