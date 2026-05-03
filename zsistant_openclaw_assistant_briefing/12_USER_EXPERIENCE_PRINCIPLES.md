# User Experience Principles

Zsistant should be powerful without feeling chaotic.

The target feeling is Apple-like simplicity on top of advanced capability.

## Principle 1: Simple first, powerful on demand

Normal users should see:

- Assistants
- Chat
- Jobs
- Channels
- Settings
- Clear status

Developer details should stay hidden unless Developer Mode is enabled.

## Principle 2: Assistant-centered navigation

Navigation should make assistants the main object.

Suggested sidebar:

- Home
- Assistants
- Projects
- Jobs
- Channels
- Tools
- Browser
- Knowledge
- Logs
- Settings

Or, if simpler:

- Assistants
- Projects
- Automations
- Tools
- Settings

## Principle 3: Status is visible

Users should always know:

- What is connected
- What is disabled
- What needs setup
- What failed
- What is running
- What needs approval

## Principle 4: Honest UI

Never show fake success.

Better to show:

```text
Chrome MCP is not connected. Configure it to enable browser actions.
```

than to show a fake browser card.

## Principle 5: Polished empty states

Empty states should educate.

Example:

```text
No assistants yet.
Create your first assistant: a persistent AI worker with tools, memory, channels, and scheduled jobs.
```

Bad empty state:

```text
No data.
```

## Principle 6: Progressive disclosure

Assistant detail should start simple:

- Overview
- Chat
- Jobs
- Channels

Advanced users can open:

- Tools
- Memory
- Logs
- Developer

## Principle 7: Beautiful but functional

The UI should look premium, but every visible element should correspond to real state.

No fake graphs.
No fake activity.
No fake integrations.

## Principle 8: Approvals are clear

When an assistant wants to take an external action, the user should see:

- What it wants to do
- Why
- What data it used
- What will happen
- Approve/deny/edit controls

## Principle 9: Logs should be human-readable

Do not make logs developer-only.

A normal user should understand:

```text
The Calendar Assistant checked tomorrow's meetings and found two conflicts.
```

Developer Mode can show raw payloads separately.

## Principle 10: The user should feel in control

A powerful assistant platform can feel scary.

Control comes from:

- Permissions
- Logs
- Pause buttons
- Approval queues
- Clear setup states
- Memory editing
- Channel allowlists
- Tool restrictions

## Visual direction

Zsistant should feel:

- Calm
- Premium
- Clean
- Fast
- Spacious
- Precise
- Trustworthy
- Modern

Avoid:

- Generic dashboard templates
- Random gradients everywhere
- Cluttered cards
- Fake statistics
- Dense developer UI by default
- Tiny unreadable controls
- Dead icons

## The "eyes want to look at it" standard

The product should feel carefully designed.

This means:

- Good spacing
- Good typography
- Clear hierarchy
- Smooth motion
- Strong empty states
- Clear icons
- Responsive layout
- High-quality settings screens
- Clean assistant cards
- Beautiful but restrained surfaces

A beautiful UI with fake features is still unacceptable.
A functional UI with no polish is also not enough.
Zsistant needs both.
