# Channels and Surfaces

A major assistant-platform feature is that the assistant is not trapped inside the web app.

The user should be able to reach assistants from the places they already use.

## Difference between a channel and a conversation

A conversation is a thread of messages.

A channel is a place where messages can enter or leave.

Examples:

- Telegram bot
- WhatsApp number
- Slack workspace
- Discord server
- Web chat
- CLI
- Email
- Voice

## Why channels matter

Channels make the assistant always available.

The user should be able to message an assistant from their phone, desktop, terminal, or browser.

## Channel configuration

Each channel should have:

- Provider/service name
- Connection status
- Authentication/setup state
- Routing destination assistant
- Allowed users
- Allowed groups/channels
- Mention rules
- Safety policy
- Rate limits
- Last message
- Error state

## Examples

### Telegram Assistant

```text
Channel: Telegram
Assistant: Personal Chief of Staff
Allowed users: only me
Capabilities: receive text, receive files, send replies, send reminders
Safety: external actions require approval
```

### Discord Coding Assistant

```text
Channel: Discord
Assistant: Coding Supervisor
Allowed channels: #dev, #bugs
Mention required: yes
Capabilities: answer questions, open investigation jobs, summarize logs
Safety: cannot push code without approval
```

### Web UI Assistant

```text
Channel: Web UI
Assistant: Any selected assistant
Capabilities: full configuration, chat, logs, browser view, developer mode
Safety: owner-only admin actions
```

### CLI Surface

```text
Channel: zazi CLI
Assistant: Developer Assistant
Capabilities: ask, inspect status, trigger jobs, list assistants
Safety: local permissions apply
```

## Channel routing

One message should be routable based on:

- Channel
- Sender
- Group
- Mention
- Keyword
- Assistant assignment
- Time/context
- Current job

Example routing rules:

- Telegram DMs from owner go to Personal Assistant.
- Discord messages in #bugs mentioning @zazi go to Coding Supervisor.
- WhatsApp family group messages go to Family Assistant only when mentioned.
- Web UI project chat goes to the assistant selected in that project.

## Multi-channel identity

The same assistant can have one identity across many channels.

A Personal Assistant can respond in Web UI, Telegram, and CLI while sharing memory and jobs.

## Channel safety

Messaging channels are dangerous if unrestricted.

Important safety controls:

- Allowlist users
- Require mention in groups
- Pairing flow for new users
- Read-only mode
- Approval before external actions
- Logs for every inbound and outbound message
- Easy pause/disable

## Honest states

If Telegram is not implemented, the UI should not say "Connected".

It should say:

```text
Telegram is not configured.
Connect a bot token to enable this channel.
```

If no implementation exists yet:

```text
Telegram integration is not available in this build.
```

No fake integrations.
