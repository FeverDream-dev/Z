# UI Information Architecture

This is a conceptual IA for the Web UI.

## Global navigation

Suggested top-level areas:

```text
Home
Assistants
Projects
Jobs
Channels
Tools
Browser
Knowledge
Logs
Settings
```

A simpler first version can use:

```text
Assistants
Projects
Automations
Tools
Settings
```

But it must still expose the assistant-first model.

## Home

Purpose: command center.

Sections:

- Active assistants
- Needs attention
- Upcoming jobs
- Recent activity
- Connected channels
- Provider/tool health
- Pending approvals

## Assistants

Purpose: create and manage assistants.

Views:

- Grid/list of assistants
- Assistant templates
- Assistant status filters
- Create assistant

## Assistant detail

Tabs:

- Overview
- Chat
- Channels
- Tools
- Knowledge
- Memory
- Jobs
- Automations
- Browser
- Logs
- Settings
- Developer

## Projects

Purpose: group context, files, chats, and assistants.

Project contains:

- Name
- Description
- Instructions
- Files
- Assistants
- Chats
- Jobs
- Knowledge

## Jobs

Purpose: global schedule and automation visibility.

Views:

- Upcoming
- Running
- Failed
- Completed
- Paused
- By assistant

## Channels

Purpose: manage communication surfaces.

Views:

- Web UI
- Telegram
- WhatsApp
- Discord
- Slack
- CLI
- Email
- Voice/future

Each channel shows connection and routing.

## Tools

Purpose: manage available capabilities.

Views:

- Browser
- Search
- Files
- Email
- Calendar
- MCP servers
- Code agents
- APIs
- Messaging

## Browser

Purpose: manage Chrome/MCP/browser sessions.

Views:

- Connection status
- Active sessions
- Screenshots
- Recent actions
- Browser QA jobs

## Knowledge

Purpose: manage files, docs, and sources.

Views:

- Global knowledge
- Project knowledge
- Assistant knowledge
- Upload/import status

## Logs

Purpose: activity across assistants.

Filter by:

- Assistant
- Tool
- Job
- Channel
- Error
- Approval

## Settings

Sections:

- Account/profile
- Appearance
- Providers/models
- Tools/MCP
- Channels
- Memory
- Privacy/security
- Developer Mode
- Keyboard shortcuts
- About

## Developer Mode overlays

When Developer Mode is enabled:

- Add trace drawer to chat
- Show provider/model details
- Show tool/MCP call details
- Show IDs and timestamps
- Show export options
- Show diagnostics badges

## Important note

This IA is conceptual. The exact layout can change, but the product must remain assistant-first.
