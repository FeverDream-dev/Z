# Not a ChatGPT Clone

This product can have a ChatGPT-like chat experience, but it must not stop there.

## What a simple ChatGPT clone contains

- Conversation sidebar
- Message list
- Text composer
- Model dropdown
- Basic settings
- Markdown rendering

That is the starting point, not the destination.

## Why a clone is insufficient

A clone only gives the user a place to talk.

Zsistant must give the user a place to manage work.

The user should be able to say:

- "My Telegram assistant handles community questions."
- "My browser QA assistant checks the app every night."
- "My coding assistant watches Sentry and creates repair plans."
- "My research assistant saves sources into the project."
- "My calendar assistant reminds me before meetings and prepares briefs."
- "My automation supervisor shows failed scheduled jobs."

These are not chat features. They are assistant-platform features.

## Chat is only one viewport

A conversation is only a viewport into the assistant's activity.

Other viewports include:

- Job timeline
- Tool logs
- Browser session
- Knowledge base
- Channel inbox
- Permission queue
- Memory editor
- Developer trace
- Assistant profile

## Anti-patterns to remove

### Fake capability cards

Do not show cards like "Connect Telegram" or "Browse Web" unless the product has a real setup flow or an honest disabled state.

### Dead buttons

No button should exist unless it works, opens a real configuration flow, or clearly says why it is disabled.

### Placeholder assistants

Do not show fake assistants with made-up stats unless explicitly in demo mode. Production UI should reflect actual stored assistants.

### Fake model integration

Do not simulate model responses in production. Missing API keys should show setup guidance.

### Fake browser automation

Do not claim browser or Chrome control unless actual MCP/browser integration exists.

### Fake jobs

Do not show cron jobs that do not run.

### Fake memory

Do not display memory summaries unless they are actually persisted.

## Minimum acceptable Zsistant beyond chat

At minimum, the UI should have:

- Assistant list
- Assistant detail page
- Assistant configuration
- Tools tab
- Channels tab
- Knowledge tab
- Jobs tab
- Logs tab
- Settings cog
- Developer mode toggle
- Provider configuration
- Honest unavailable states

If these do not exist, the product is still a chat clone.
