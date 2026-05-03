# Assistant Profile and Capabilities

An assistant profile is the durable definition of an assistant.

The coding AI should treat this as a product concept, not just a UI form.

## Required profile sections

### Identity

- Name
- Avatar/icon
- Short description
- Long purpose
- Owner
- Status

### Persona

- Tone
- Communication style
- Decision-making style
- Boundaries
- Preferred format
- User relationship

Example:

```text
Name: Executive Assistant
Persona: concise, proactive, careful, permission-seeking for external communication
Style: brief summaries first, details on request
```

### Responsibilities

What this assistant owns.

Examples:

- Manage calendar conflicts
- Prepare daily briefings
- Watch website changes
- Answer Telegram community questions
- Run browser QA
- Route coding tasks to agents

### Model and provider policy

The assistant can have a model strategy.

Examples:

- Default model
- Fallback model
- High-reasoning model for hard tasks
- Cheap model for routine tasks
- Local model for private tasks
- Provider health checks

Do not fake providers. Missing credentials should show setup state.

### Tools

Capabilities the assistant may use.

Examples:

- Browser
- Search
- Files
- Gmail
- Calendar
- Telegram
- Slack
- Discord
- GitHub
- Coding agent
- MCP server
- Shell/sandbox

### Permissions

Every powerful action needs a permission model.

Permission examples:

- Can read calendar
- Can create calendar draft but not send invite
- Can draft email but needs approval to send
- Can browse public web
- Can access selected project files
- Can run local browser QA
- Cannot run shell commands
- Can message Telegram group only when mentioned

### Knowledge

The assistant should know which knowledge sources it can use.

Examples:

- Project documents
- Uploaded files
- Personal notes
- User preferences
- Team documentation
- Previous conversations

### Memory policy

Memory is not the same as conversation history.

Define:

- Whether memory is enabled
- What can be remembered
- What must not be remembered
- Whether memories are global or assistant-specific
- Whether user approval is needed before storing memory

### Channels

Where the assistant can communicate.

Examples:

- Web UI
- CLI
- Telegram
- WhatsApp
- Slack
- Discord
- Email

### Jobs

Scheduled tasks.

Examples:

- Daily briefing at 8 AM
- Weekly project summary every Friday
- Check website uptime every 15 minutes
- Summarize unread emails every morning
- Run browser regression check nightly

### Automations

Event-triggered tasks.

Examples:

- If new email matches rule, summarize
- If user sends Telegram voice note, transcribe and respond
- If GitHub issue opened, triage
- If browser test fails, capture screenshot and report

### Logs

Every assistant should have an inspectable activity timeline.

### Health

The assistant can be:

- Active
- Paused
- Needs setup
- Degraded
- Error
- Waiting for approval

## Capability maturity levels

Use this language in UI and planning.

### Level 0: Not available

Feature is not built. Do not pretend.

### Level 1: Configurable but inactive

Setup UI exists, but user has not configured it.

### Level 2: Connected

Credentials/channel/tool is connected and can be tested.

### Level 3: Usable manually

User can invoke it from chat/UI.

### Level 4: Automatable

Assistant can use it in jobs/automations.

### Level 5: Observable and reliable

Logs, retries, error states, permissions, and audits exist.

Zsistant should clearly show maturity levels instead of hiding incomplete features.
