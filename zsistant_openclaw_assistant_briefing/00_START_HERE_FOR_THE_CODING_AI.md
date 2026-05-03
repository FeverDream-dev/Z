# Start Here for the Coding AI

You have been building the wrong product.

You interpreted the assignment as: "make a ChatGPT-style chat interface."

That is not enough.

The intended product is Zsistant / zazi: a complete assistant platform inspired by OpenClaw-like personal AI systems, OpenAI workspaces, Claude artifacts, Perplexity spaces, browser automation, MCP tools, Telegram/WhatsApp-style surfaces, scheduled jobs, and multi-agent work.

## Core correction

Do not ask: "How do I make a nicer chat UI?"

Ask: "How do I let a user create, configure, monitor, and command real assistants that can do work across channels, tools, browser sessions, files, schedules, and other agents?"

## The product hierarchy

The correct hierarchy is:

1. Workspace
2. Assistants
3. Channels
4. Tools
5. Memory and knowledge
6. Jobs and automations
7. Conversations
8. Messages

The wrong hierarchy is:

1. Chat page
2. Messages
3. Text box
4. Some settings later

## The minimum mental model

Each assistant is a persistent digital worker.

It has:

- A name
- A purpose
- A persona
- A model/provider strategy
- A memory policy
- A knowledge base
- A set of tools
- Channel connections
- Permissions
- Scheduled jobs
- Event triggers
- Current tasks
- Logs
- Health/status
- Conversation history
- Developer diagnostics

A chat window is merely one way to talk to it.

## What the UI must communicate

When a user opens Zsistant, they should see that this is a place to manage assistants, not only send prompts.

They should be able to understand:

- Which assistants exist
- What each assistant can do
- Where each assistant is connected
- What jobs are running
- Which tools are enabled
- Which channels are live
- What each assistant remembers
- What it did recently
- What needs permission
- What failed
- What is scheduled next

## Forbidden product direction

Do not build only:

- A left sidebar with chats
- A message list
- A composer
- A fake model selector
- A fake settings cog
- Placeholder cards
- Hardcoded assistant replies
- Demo integrations

That product is a weak ChatGPT clone. Zsistant is not that.

## Required product direction

Build toward:

- Assistant manager
- Assistant profiles
- Tool registry
- MCP/Chrome/browser controls
- Telegram/WhatsApp/Discord/Slack-like channel connections
- Knowledge and memory areas
- Schedules and cron-like jobs
- Automations and triggers
- Project/workspace organization
- Real provider/model routing
- Developer mode
- Logs, traces, and observability
- Polished Apple-level UI

## If you are unsure what to build

Build the thing that makes Zsistant more like an operating system for assistants, not more like a text-message clone.
