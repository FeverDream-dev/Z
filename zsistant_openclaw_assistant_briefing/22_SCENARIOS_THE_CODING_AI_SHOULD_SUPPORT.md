# Scenarios the Coding AI Should Support

These scenarios describe product behavior. They are not programming instructions.

## Scenario 1: Create a browser QA assistant

The user creates an assistant named Browser QA.

They set:

- Purpose: inspect the web UI and find issues
- Tools: Chrome/MCP if connected
- Jobs: nightly smoke test
- Knowledge: product design checklist
- Permissions: local dev URLs only, ask before form submit

The assistant appears in the assistant manager.

If Chrome/MCP is not connected, the assistant shows "Needs setup" instead of pretending to work.

## Scenario 2: Connect Telegram to an assistant

The user opens an assistant's Channels tab.

They choose Telegram.

The UI shows either:

- Connected bot and allowed users, or
- Setup instructions, or
- Unavailable in this build

The user can route Telegram DMs to the Personal Assistant.

Group messages require mention by default.

## Scenario 3: Schedule a daily briefing

The user creates a job for Personal Chief of Staff.

Schedule: weekdays at 08:00.

Inputs: calendar, unread important messages, project updates.

If calendar/email tools are not configured, the job shows missing dependencies.

The product does not fake a briefing.

## Scenario 4: Developer Mode debugging

A provider call fails.

Normal mode shows:

```text
The assistant could not contact the selected model. Check provider settings.
```

Developer Mode shows:

- Provider
- Model
- Error type
- Time
- Retry/fallback behavior
- Request id if available

## Scenario 5: Assistant delegation

The user asks Personal Assistant:

```text
Check if the new UI is polished enough.
```

The Personal Assistant delegates to Product Design Assistant and Browser QA Assistant.

The UI shows a timeline:

1. Request received
2. Delegated to Product Design Assistant
3. Browser QA captured screenshot
4. Findings returned
5. Summary delivered

## Scenario 6: Memory curation

After repeated conversations, Memory Curator suggests:

```text
Save memory: User wants Zsistant to be assistant-first, not a chat clone.
```

The user can approve, edit, or reject it.

## Scenario 7: Failed job

Nightly browser QA fails because Chrome MCP is disconnected.

The Jobs tab shows:

- Failed
- Cause: Chrome MCP unavailable
- Last successful run
- Retry option
- Setup link

No fake success.

## Scenario 8: Tool permission request

Assistant wants to send an email.

The UI shows:

- Draft email
- Recipient
- Reason
- Source context
- Approve
- Edit
- Deny

The assistant cannot silently send it unless permission allows.

## Scenario 9: Assistant status dashboard

The user opens Home.

They see:

- Personal Assistant active
- Browser QA needs Chrome MCP setup
- Telegram Community Assistant paused
- 3 jobs scheduled today
- 1 failed automation
- 2 approvals pending

This is an assistant command center, not just chat history.

## Scenario 10: Project assistant workspace

The user creates a project called Zsistant.

They attach:

- Product scope docs
- UI screenshots
- Coding rules
- Assistants: Coding Supervisor, Browser QA, Product Design Assistant

Project chats use project context and selected assistant identity.
