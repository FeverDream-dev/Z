# Web UI: Assistant Manager

The Web UI should not open directly into a lonely chat screen.

It should make clear that Zsistant manages many assistants.

## Home screen

The home screen should answer:

- What assistants do I have?
- Which are active?
- Which need attention?
- Which jobs are due soon?
- Which integrations are connected?
- What happened recently?

## Assistant cards

Each assistant card should show:

- Name
- Icon/avatar
- Purpose
- Status: active, paused, needs setup, error
- Main channels: web, Telegram, WhatsApp, Slack, Discord, CLI
- Tool count
- Knowledge count
- Job count
- Last activity
- Current model/provider
- Attention badge if something failed or needs approval

Example assistant card:

```text
Personal Chief of Staff
Status: Active
Channels: Web, Telegram
Tools: Calendar, Gmail, Browser, Notes
Jobs: 4 scheduled
Last activity: Prepared daily briefing 12 minutes ago
Needs approval: 2 email drafts
```

## Assistant detail page

Clicking an assistant should open a full control center.

Suggested tabs:

1. Overview
2. Chat
3. Channels
4. Tools
5. Knowledge
6. Memory
7. Jobs
8. Automations
9. Browser
10. Logs
11. Settings
12. Developer

## Overview tab

The overview should show:

- What this assistant does
- Current health
- Connected channels
- Enabled tools
- Running jobs
- Recent actions
- Pending approvals
- Recent errors
- Quick actions

## Chat tab

The chat tab is for direct conversation with that assistant.

It should show the assistant's identity and active context.

## Channels tab

The channels tab manages where the assistant can be reached.

Examples:

- Web UI
- Telegram bot
- WhatsApp gateway
- Discord server/channel
- Slack workspace/channel
- CLI
- Email

Each channel should have:

- Connection status
- Setup instructions
- Routing rules
- Allowed users/groups
- Mention rules
- Safety controls
- Last inbound/outbound message

## Tools tab

The tools tab shows what the assistant can do.

Examples:

- Browser/Chrome
- Web search
- Gmail
- Calendar
- Files
- Shell/sandbox
- MCP servers
- HTTP/API calls
- Other LLM providers
- Coding agents
- Notion/Drive/GitHub

Each tool should have:

- Enabled/disabled state
- Permission level
- Setup status
- Last used
- Failure count
- Audit log

## Knowledge tab

Knowledge is what the assistant can reference.

Examples:

- Uploaded PDFs
- Notes
- Project docs
- Personal preferences
- Company docs
- Saved webpages
- Conversation summaries
- Manual instructions

## Memory tab

Memory is durable learned context.

The user should be able to:

- View memories
- Add memories
- Edit memories
- Delete memories
- Set memory policy
- Choose what is global vs assistant-specific

## Jobs tab

Jobs are scheduled responsibilities.

Show:

- Job name
- Schedule
- Next run
- Last run
- Status
- Output
- Failures
- Owner assistant
- Retry policy
- Approval requirement

## Automations tab

Automations are event-triggered workflows.

Examples:

- When new email from boss arrives, summarize it.
- When website changes, notify me.
- When Sentry error appears, ask Coding Assistant to investigate.
- When Telegram keyword appears, route to Community Assistant.

## Browser tab

Shows browser or Chrome/MCP capability.

Should include:

- Connection status
- Active sessions
- Last screenshot
- Current URL
- Recent actions
- Permissions
- Open/capture/inspect controls if real

## Logs tab

Human-readable timeline:

- Message received
- Tool called
- Job started
- Browser screenshot captured
- File read
- Email draft created
- Approval requested
- Error occurred

## Developer tab

Visible only in Developer Mode or only for advanced users.

Shows traces and raw diagnostics.

## Creation flow

Creating an assistant should ask:

- What is its purpose?
- Which persona should it use?
- Which channels should it connect to?
- Which tools should it have?
- What knowledge should it use?
- What can it do without approval?
- What must always require approval?
- Should it run scheduled jobs?

## Assistant templates

Templates are allowed if they create real configurable assistants.

Examples:

- Personal Chief of Staff
- Browser QA Assistant
- Coding Agent Supervisor
- Research Assistant
- Telegram Community Assistant
- Content Pipeline Assistant
- Calendar/Inbox Assistant
- Memory Curator
- Automation Supervisor

Templates should not be fake cards. They should create actual assistant records/configuration.
