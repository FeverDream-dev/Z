# Developer Mode and Observability

Developer Mode is a major product feature.

It is not a hidden debug panel or an afterthought.

## Purpose

Normal users need simplicity.

Developers and power users need visibility.

Developer Mode should reveal how the assistant works without cluttering the default experience.

## Developer Mode toggle

There should be a real setting:

```text
Settings -> Developer Mode -> On/Off
```

When enabled, the UI changes meaningfully.

## What Developer Mode exposes

### Assistant internals

- Active assistant profile
- Active persona/instructions
- Memory used
- Knowledge sources used
- Tool permissions
- Model/provider routing

### Conversation internals

- Raw message structure
- System/developer/user/assistant/tool separation where appropriate
- Context window usage if available
- Token counts if available
- Latency
- Streaming events
- Stop/retry metadata

### Tool traces

- Tool name
- Input summary
- Output summary
- Duration
- Error
- Approval state
- Assistant that called it
- Triggering channel/message

### MCP traces

- MCP server status
- Available tools
- Last call
- Request/response summary
- Errors
- Connection health

### Browser traces

- Browser session id
- Current URL
- Last screenshot
- Console errors
- DOM/action log if available
- Click/type/navigation history

### Provider diagnostics

- Configured providers
- Missing credentials
- Health check results
- Fallbacks
- Model selection
- Error rates

### Job diagnostics

- Job schedule
- Next run
- Last run
- Run logs
- Retry attempts
- Failure cause

### Export tools

- Export conversation as Markdown
- Export conversation as JSON
- Export assistant profile
- Export job run log
- Export browser QA report

## Developer Mode UI placements

Developer Mode can add:

- Developer tab on assistant detail
- Trace drawer in chat
- Raw event log panel
- Tool call inspector
- Provider health page
- MCP registry details
- Debug badges on UI elements

## Normal mode vs Developer Mode example

Normal mode:

```text
Browser QA Assistant checked the app and found 2 issues.
```

Developer Mode:

```text
Tool trace:
- chrome.open_url: 842 ms
- chrome.screenshot: 421 ms
- console.read: 83 ms
- visual_check: 2 issues
Provider: OpenAI GPT-X fallback disabled
Tokens: input 9123, output 804
```

## Anti-pattern

Do not expose raw debug data to everyone by default.
Do not hide all debug data so failures are impossible to understand.

Developer Mode solves both problems.
