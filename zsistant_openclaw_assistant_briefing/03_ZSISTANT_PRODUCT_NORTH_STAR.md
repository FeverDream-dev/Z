# Zsistant Product North Star

## Product statement

Zsistant / zazi is an assistant operating workspace.

It lets a user create, configure, use, monitor, and improve multiple AI assistants that can act through channels, tools, browsers, models, jobs, files, and memory.

## North star experience

A user opens Zsistant and sees their AI workforce.

They can see:

- Which assistants exist
- Which assistants are active
- Which channels are connected
- Which jobs are scheduled
- What each assistant has done recently
- What needs approval
- What failed
- What tools are available
- What knowledge each assistant uses
- Which model/provider powers each assistant
- Which browser sessions or MCP servers are connected

The user can click an assistant and configure it deeply.

## The core object is Assistant

Everything should orbit around Assistants.

A conversation belongs to an assistant.
A job belongs to an assistant.
A tool permission belongs to an assistant.
A Telegram channel routes to an assistant.
A memory namespace belongs to an assistant.
A browser session can be assigned to an assistant.

## Product surfaces

Zsistant should support or plan for these surfaces:

- Web UI
- CLI through `zazi`
- Telegram
- WhatsApp
- Discord
- Slack
- Browser/Chrome MCP
- Local desktop/browser automation
- Future mobile/voice

## Product areas

### Assistant Manager

The home of all assistants.

### Assistant Detail

The control center for one assistant.

### Conversations

Chat history, but attached to an assistant and context.

### Tools

Registry of actions the assistant can use.

### Channels

Communication surfaces connected to assistants.

### Knowledge

Files, documents, notes, memories, and project context.

### Jobs

Scheduled and recurring tasks.

### Automations

Event-based triggers.

### Browser

Chrome/MCP/browser sessions and visual inspection.

### Developer Mode

Advanced traces, logs, payloads, tool calls, and debugging.

### Settings

Product-level preferences, provider credentials, permissions, appearance, data controls.

## Product promise

Zsistant should feel like:

- A personal command center
- A team of assistants
- A browser/computer-using operator
- A multi-channel AI hub
- A developer-observable agent platform
- A beautiful consumer product

It should not feel like:

- A demo chat UI
- A template dashboard
- A toy chatbot
- A model wrapper
- A fake integration showcase

## Decision rule

When choosing between two features, pick the one that makes assistants more capable, observable, persistent, and useful across time.
