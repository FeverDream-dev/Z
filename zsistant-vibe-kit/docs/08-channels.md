# 08 — Channel Integration Map

Mandatory MVP channels:

- CLI.
- Local web chat.
- Telegram.
- Discord.
- WhatsApp.

## Common channel adapter interface

Each channel should normalize inbound messages into:

```text
channel_type
channel_account_id
conversation_id
sender_id
agent_id
text
attachments
timestamp
raw_event_ref
```

Each outbound message should include:

```text
agent_id
conversation_id
text
attachments
reply_to
status_mode
```

## CLI

Purpose:

- Local testing.
- Admin commands.
- Quick chat.
- Agent creation.
- Status inspection.

Important commands to design:

```text
zazi doctor
zazi init
zazi serve
zazi agent create
zazi agent list
zazi chat
zazi job status
zazi channel telegram setup
zazi channel discord setup
zazi channel whatsapp setup
zazi peer allow
zazi skill scan
```

## Local web chat

Purpose:

- Non-technical user control.
- Create agents.
- Bind channels.
- Watch status.
- Approve risky actions.
- Inspect logs and workspace summaries.

## Telegram

Start with long polling for local development. Add webhooks later.

Setup fields:

- bot token
- allowed chat IDs, optional
- target agent
- display name

## Discord

Start with the simplest reliable approach for local testing.

Options:

- Gateway bot for message events.
- HTTP interactions endpoint for slash commands.

Use minimal permissions.

## WhatsApp

Use WhatsApp Cloud API concepts:

- phone number ID
- access token
- webhook verify token
- inbound webhook events
- outbound message send

WhatsApp setup has more friction, so the UI should show a checklist.

## Channel binding

A channel does not equal an agent.

The same machine can have:

- multiple Telegram bots mapped to different agents
- one Discord app with different command routes
- separate WhatsApp numbers or routing rules
- local-only agents with no external channel

## Message status

When a job takes time, channel adapters should send concise progress updates.

Examples:

```text
Working on it. I found the right inbox and am checking the latest message.
```

```text
Still working. The model provider timed out, so I switched to the fallback.
```

```text
Paused for approval: this action would send an external message.
```
