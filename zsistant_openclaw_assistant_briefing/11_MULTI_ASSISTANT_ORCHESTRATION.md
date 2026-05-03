# Multi-Assistant Orchestration

Zsistant should support many assistants, not one generic bot.

Different assistants should have different roles, tools, memories, channels, and jobs.

## Why multiple assistants matter

A single generic assistant becomes confusing.

It cannot safely have all tools, all permissions, all personas, all channels, and all jobs at once.

Multiple assistants allow specialization.

## Examples of specialized assistants

- Personal Chief of Staff
- Browser QA Assistant
- Coding Agent Supervisor
- Research Assistant
- Telegram Community Assistant
- Calendar Assistant
- Inbox Assistant
- Finance Assistant
- Content Assistant
- Memory Curator
- Automation Supervisor
- Security Guard Assistant

## Routing

When a message or event arrives, Zsistant should route it to the right assistant.

Routing can use:

- Channel
- Sender
- Project
- Keyword
- Mention
- Tool need
- Assistant availability
- User selection
- Automation rule

## Delegation

Assistants should be able to delegate.

Example:

```text
Personal Assistant receives: "Check if the website looks broken."
Personal Assistant delegates to Browser QA Assistant.
Browser QA Assistant opens Chrome, captures screenshots, reports issues.
Personal Assistant summarizes the result to the user.
```

## Supervisor pattern

A supervisor assistant can coordinate other assistants.

Example:

```text
Automation Supervisor watches scheduled jobs.
If Browser QA fails, it asks Browser QA Assistant for details.
If code fix is needed, it asks Coding Supervisor to create a plan.
Then it reports to the user.
```

## Assistant boundaries

Each assistant should know what it owns and what it should delegate.

Examples:

- Research Assistant should not send emails.
- Browser QA Assistant should not manage calendar.
- Telegram Community Assistant should not read private files.
- Coding Supervisor should not answer family WhatsApp messages.

## Shared context

Some context may be shared:

- Global user preferences
- Project goals
- Common knowledge base
- Provider settings

Some context should be isolated:

- Private channel data
- Assistant-specific memory
- Sensitive files
- Tool permissions

## UI expectations

Zsistant should let the user:

- View all assistants
- See their relationships
- Assign default assistants to channels
- Assign assistants to projects
- See delegation history
- See which assistant handled an event
- Pause an assistant
- Duplicate an assistant
- Archive an assistant

## Orchestration examples

### Web app repair flow

1. User asks Personal Assistant: "Fix the UI."
2. Personal Assistant routes to Product/Design Assistant for critique.
3. Browser QA Assistant captures screenshots.
4. Coding Supervisor creates implementation plan.
5. Browser QA Assistant verifies the result.
6. Personal Assistant summarizes outcome.

### Research-to-content flow

1. Research Assistant collects sources.
2. Memory Curator extracts durable insights.
3. Content Assistant writes a draft.
4. Personal Assistant asks user for approval before publishing.

### Telegram support flow

1. Telegram Community Assistant receives question.
2. It searches product knowledge.
3. If uncertain, it asks Research Assistant or routes to human approval.
4. It replies or drafts a reply.

## Anti-pattern

Do not implement "assistants" as only different prompt presets in a dropdown.

An assistant is not just a prompt.

An assistant is a configured actor with tools, channels, memory, jobs, and permissions.
