# AUDIT_AND_REBUILD_PLAN.md

## 1. Current Architecture Map

### Repository Layout
```
cmd/zazi/                     CLI entrypoint (main.go, 1737 lines)
internal/
  agents/                     Agent registry (JSON file-backed, no DB)
  bus/                        ACL message bus (JSON file-backed)
  channels/                   Telegram, Discord, WhatsApp adapters
  config/                     YAML config loader + secrets
  daemon/                     PLACEHOLDER — Start() is empty
  jobs/                       Job queue + events (JSONL file-backed)
  llm/                        Provider interface + 6 implementations
  sandbox/                    PLACEHOLDER — InitSandbox() is empty
  server/                     HTTP server + 2 embedded HTML templates
  skills/                     Risk analyzer (keyword-based, safe)
  store/                      PLACEHOLDER — InitStore() is empty
  trainer/                    Persona trainer (keyword-based)
assets/                       Static images, SVGs, avatars, mockups
web/                          PLACEHOLDER — StartWebServer() is empty
docs/                         Phase specifications (00–16)
phases/                       Markdown spec files
```

### Tech Stack
- **Language**: Go 1.25.9
- **Module**: `github.com/FeverDream-dev/zsistant`
- **Dependencies**: `gopkg.in/yaml.v3` (only external dependency)
- **Frontend**: No framework. Raw HTML/CSS/JS embedded as Go string templates inside `internal/server/server.go`
- **Build**: No frontend build step. Pure Go `go build`.
- **Persistence**: File-backed JSON/JSONL/YAML. No SQLite, no database.
- **Auth**: None. Unauthenticated.
- **Tests**: Go standard testing. 8–12 tests per package.

### Entrypoints
1. **CLI**: `cmd/zazi/main.go` — 15+ subcommands: version, doctor, init, serve, chat, telegram, discord, whatsapp, job, train, acl, skill, validate, release, provider
2. **HTTP Server**: `internal/server/server.go` — `Server.Run()` starts `http.ListenAndServe`
3. **No other entrypoints** — no separate frontend app, no service worker, no API gateway

---

## 2. Frontend Routes / Pages / Components

### HTTP Routes (from `internal/server/server.go`)
| Route | Method | Handler | Status |
|---|---|---|---|
| `GET /` | GET | `handleDashboard()` | Real — agent grid, dark theme |
| `GET /chat` | GET | `handleChat()` | Real — basic chat with sidebar |
| `GET /health` | GET | `handleHealth` | Real — returns `{"status":"ok"}` |
| `GET /api/agents` | GET | `handleAgents` | Real — lists agents from registry |
| `GET /api/providers` | GET | `handleProviders` | **STUB** — always returns mock health |
| `POST /api/chat` | POST | `handleChatAPI` | **STUB** — always routes to `MockProvider` |
| `GET /api/jobs/:agent_id` | GET | `handleJobsAPI` | Real — reads `audit.jsonl` |
| `GET /assets/brand/*` | GET | FileServer | Real — static assets |

### Pages
1. **Dashboard** (`/`): Dark-themed page with header (logo + "Online" badge), agent grid loaded via `fetch('/api/agents')`, each card links to `/chat?agent=<id>`. No sidebar, no navigation, no search, no settings.
2. **Chat** (`/chat?agent=<id>`): Dark-themed page with header (back link + title), main chat panel (messages list + text input), sidebar (job timeline from `audit.jsonl`). No model picker, no file upload, no settings, no streaming.

### UI Components
- **Header**: Static image + title + status badge
- **Agent Grid**: Cards with name, ID, role, status, "Open Chat" link
- **Message List**: Simple divs with `user`/`bot`/`error` CSS classes. Plain text only.
- **Composer**: Single text input + Send button. No file attach, no model picker, no voice.
- **Sidebar**: Job event timeline (event type, message, timestamp). No settings, no project list.

---

## 3. Backend / API Status

### Provider/Model Integration
- **Interface**: `Provider` — `Complete(prompt) (string, error)` + `Health() ProviderHealth`
- **Real providers**: Ollama Cloud, Generic OpenAI-compatible, Z.AI Coding Plan, OpenCode Zen (4 real)
- **Mock provider**: `MockProvider` — echoes input prefixed with "Echo:"
- **Catalog**: 52 providers, 207 models in `BuiltInProviders` (static data)
- **Critical bug**: Web server `handleChatAPI` and `handleProviders` **ALWAYS** instantiate `NewMockProvider()` — real providers exist but are never wired into the web server
- **CLI**: `runChat` checks for `ollama_api_key` in config and registers Ollama + mock fallback. Other commands (telegram test, discord test, whatsapp test) always use `NewMockProvider()`
- **Streaming**: Not supported. `stream: false` hardcoded in all providers. No SSE, no chunked transfer.
- **Cancellation**: Not supported. Router uses `context.WithTimeout` but no `context.Cancel`.
- **Retries**: Router tries fallback providers on timeout/error. No exponential backoff.
- **Health**: Updated on success/failure. Thread-safe with `sync.Mutex`.

### Persistence
- **Agents**: `~/.zazi/agents/<id>/profile.json` — JSON encoded `Agent` struct
- **Jobs**: `~/.zazi/agents/<id>/jobs/` directory + `audit.jsonl` (line-delimited JSON events)
- **Config**: `~/.zazi/config.yaml` — YAML with secrets map
- **No database**: No SQLite, no PostgreSQL, no in-memory store. Pure filesystem.
- **Durable**: Yes, files are written with `os.WriteFile`.

### Auth
- **Status**: None. Zero authentication.
- **Implications**: Anyone with network access can call any API.

### Settings
- **Status**: None. No settings page, no settings API endpoint, no configuration UI.
- **Config is CLI-only**: Edited via `zazi init` or manual YAML editing.

### MCP/Tools
- **Status**: No MCP integration in web UI.
- **Skill analyzer**: Keyword-based risk scoring. Safe, never executes code.
- **No tool calling**: No function calling in chat. No MCP server registry.
- **Chrome MCP**: Not integrated.

---

## 4. Placeholder / Mock / Stub Inventory

### Critical Production Stubs (MUST FIX)
| File | Line | Issue | Impact |
|---|---|---|---|
| `internal/server/server.go:155-157` | 155-157 | `handleChatAPI` creates router with only `NewMockProvider()` | **Web chat always returns echo responses** |
| `internal/server/server.go:82-85` | 82-85 | `handleProviders` creates router with only `NewMockProvider()` | **Health API always reports mock** |
| `cmd/zazi/main.go:1116` | 1116 | `runWhatsAppTest` uses `llm.NewMockProvider()` | WhatsApp test is fake |
| `cmd/zazi/main.go:1259` | 1259 | `runDiscordTest` uses `llm.NewMockProvider()` | Discord test is fake |
| `cmd/zazi/main.go:1379-1380` | 1379-1380 | `runTelegramTest` uses mock provider (comment says "Use mock provider") | Telegram test is fake |
| `cmd/zazi/main.go:1478` | 1478 | `runTelegramListen` registers only `NewMockProvider()` when no Ollama key | Telegram listen falls back to echo |

### Placeholder Packages (no-op functions)
| File | Line | Issue |
|---|---|---|
| `internal/channels/channels.go:3` | 3 | `InitChannels` is a no-op comment-only function |
| `internal/daemon/daemon.go:3` | 3 | `Start` is a no-op comment-only function |
| `internal/jobs/jobs.go:3` | 3 | `InitJobs` is a no-op comment-only function |
| `internal/llm/llm.go:3` | 3 | `InitLLM` is a no-op comment-only function |
| `internal/sandbox/sandbox.go:3` | 3 | `InitSandbox` is a no-op comment-only function |
| `internal/skills/skills.go:3` | 3 | `InitSkills` is a no-op comment-only function |
| `internal/store/store.go:3` | 3 | `InitStore` is a no-op comment-only function |
| `web/server.go:3` | 3 | `StartWebServer` is a no-op comment-only function |

### Fake/Placeholder UI Elements
| File | Line | Issue |
|---|---|---|
| `internal/server/server.go:296` | 296 | `<span class="status">Online</span>` — hardcoded, not from health check |
| `internal/server/server.go:300` | 300 | `Loading agents...` — no skeleton, no loading state management |
| `internal/server/server.go:380` | 380 | Static welcome message: "Hello! Send a message to start chatting." |
| `assets/brand/channels/*-placeholder.svg` | — | Channel icons are labeled "placeholder" |
| `assets/mockups/*.svg` | — | Wireframe mockups are not real UI |

### Console logs in production JS
| File | Line | Issue |
|---|---|---|
| `internal/server/server.go:461` | 461 | `console.error('Failed to load events', e)` in chat page |

### Config placeholders
| File | Line | Issue |
|---|---|---|
| `internal/config/config.go:12` | 12 | `LLMProviders` is described as "placeholder for future provider-specific configs" |

---

## 5. Missing Features vs Modern AI Apps

### ChatGPT-like Missing
- [ ] **Projects/Workspaces** — No project concept. Agents are flat.
- [ ] **Custom Instructions** — No per-agent or per-project instructions UI.
- [ ] **Memory** — No memory/personalization controls.
- [ ] **Tools** — No tool calling, no plugin architecture.
- [ ] **Canvas/Workspace** — No split-pane artifact editing.
- [ ] **Model Picker** — No UI to select model. Hardcoded or config-only.
- [ ] **Voice** — No voice input/output architecture.
- [ ] **Settings** — No settings screen at all.

### Claude-like Missing
- [ ] **Artifacts** — No artifact extraction from chat.
- [ ] **Computer Use** — No browser/computer-use orientation.
- [ ] **Long Context** — No explicit long-context handling UI.

### Perplexity-like Missing
- [ ] **Web Search** — No web search integration.
- [ ] **Citations** — No source/citation system.
- [ ] **Spaces** — No search/pin organization.

### Apple-like UX Missing
- [ ] **Polished Empty States** — "Loading agents..." is minimal.
- [ ] **Smooth Animations** — No transitions, no microinteractions.
- [ ] **Typography** — System font only, no font loading.
- [ ] **Responsive** — Desktop-only layout, no mobile consideration.
- [ ] **Accessibility** — No ARIA labels, no focus management, no keyboard shortcuts.

---

## 6. Prioritized Rebuild Plan

### Phase A: Critical Infrastructure (Week 1)
1. **Wire real providers into web server** — Replace `NewMockProvider()` in `handleChatAPI` and `handleProviders` with dynamic provider selection from config. This is the #1 blocker.
2. **Add streaming support to provider layer** — SSE endpoint, chunked response, frontend EventSource.
3. **Add settings API + settings page** — `GET/POST /api/settings`, minimal HTML settings form.
4. **Add real provider configuration flow** — API key input, provider selection, model selection, test connection.

### Phase B: Chat Architecture (Week 1-2)
5. **Conversation persistence** — Conversations as first-class objects (not just jobs). JSON file-backed.
6. **Message model enrichment** — Roles, timestamps, provider/model metadata, token usage (when available), latency.
7. **Chat actions** — Stop generation, regenerate, copy, delete message, edit user message.
8. **Error states** — Human-readable errors, retry buttons, provider health indicators.

### Phase C: UI/UX Rebuild (Week 2-3)
9. **New dashboard** — Sidebar with conversations, projects, search. Clean empty states.
10. **New chat page** — Model picker in composer, file upload (or hidden if not implemented), streaming indicator, markdown rendering, code blocks.
11. **Settings panel** — Sections: Appearance, Models, API Keys, Memory, Tools, Developer Mode, Privacy, Shortcuts, About.
12. **Developer Mode toggle** — Real toggle that exposes: raw request/response inspector, token/latency display, provider routing details, streaming event log, provider health panel.

### Phase D: Advanced Features (Week 3-4)
13. **Projects/Workspaces** — Create project, add chats to project, project-level instructions, project sidebar.
14. **Canvas/Artifacts** — Split-pane workspace, text/markdown artifacts, create/revise/save from chat.
15. **MCP/Tools panel** — MCP server configuration UI, status display, tool call trace.
16. **Research mode** — Honest disabled state with setup instructions if web search not available.

### Phase E: Polish & QA (Week 4)
17. **Remove all placeholders** — Delete or implement the 8 no-op init functions.
18. **Browser QA** — Screenshots of every page, console error checks, interaction tests.
19. **Tests** — API route tests, provider tests, E2E smoke test.
20. **Docs** — README, DEVELOPER_MODE.md, MCP_INTEGRATION.md, QA_SCREENSHOTS.md.

---

## 7. Acceptance Criteria

The rebuild is complete when:

- [ ] `handleChatAPI` routes to real providers (not mock) when configured
- [ ] `handleProviders` reports real provider health (not mock)
- [ ] Settings page exists and is reachable from UI
- [ ] Developer Mode toggle exists and meaningfully changes the UI
- [ ] Streaming works end-to-end (provider → API → frontend)
- [ ] All 8 no-op placeholder init functions are either implemented or removed
- [ ] No `console.log` or `console.error` in production frontend code
- [ ] No hardcoded "Online" status — uses real health check
- [ ] Chat has model picker, stop button, and copy action
- [ ] Dashboard has sidebar navigation
- [ ] Build passes (`go build ./...`)
- [ ] Tests pass (`go test ./...`)
- [ ] Browser QA completed with screenshots
- [ ] README updated with setup instructions
- [ ] `.env.example` updated
- [ ] UI feels like a serious AI product (not a starter template)

---

## 8. Immediate Next Steps

1. **Fix web server provider wiring** — This is the single most impactful change. One edit to `internal/server/server.go` lines 155-157 switches the web chat from fake to real.
2. **Add `/api/settings` endpoint** — Minimal settings read/write.
3. **Add settings cog to chat page** — Link to settings, or inline settings panel.
4. **Add model picker to chat composer** — Dropdown populated from `BuiltInProviders`.

---

*Audit completed. Begin implementation.*
