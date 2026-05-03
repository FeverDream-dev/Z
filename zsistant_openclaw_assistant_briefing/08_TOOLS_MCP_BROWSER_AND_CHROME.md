# Tools, MCP, Browser, and Chrome

Tools are the assistant's hands.

Without tools, an assistant is only a talking model.

## Tool concept

A tool is any capability that lets the assistant observe, retrieve, change, or act in the world.

Examples:

- Browser control
- Chrome MCP
- Playwright MCP
- Web search
- File reading/writing
- Email
- Calendar
- Messaging
- GitHub
- Database query
- API call
- Shell/sandbox
- Code agent
- Screenshot capture
- OCR/vision
- Knowledge retrieval

## MCP concept

MCP-like integrations should be treated as tool servers.

A tool server exposes capabilities to the assistant.

The UI should show:

- MCP server name
- Status
- Available tools
- Required permissions
- Last call
- Errors
- Configuration

## Browser and Chrome as central features

The browser is not a minor add-on.

Browser access gives the assistant eyes and hands.

The assistant can:

- Open a website
- Capture screenshots
- Inspect visible UI state
- Click buttons
- Fill forms
- Test the app
- Compare before/after UI quality
- Find console errors
- Verify flows
- Produce QA reports

This is critical to the Zsistant vision because the assistant should be able to inspect the product it is building.

## Browser Assistant example

```text
Assistant: Browser QA Assistant
Purpose: Inspect web apps and report visual/functionality issues
Tools: Chrome MCP, screenshot capture, DOM inspector, console log reader
Jobs: Nightly UI smoke test
Channels: Web UI, Discord #qa
Permissions: Can browse local dev URLs and public websites; cannot submit forms without approval
```

## Tool registry

Zsistant should have a registry of tools.

Each tool should show:

- Name
- Category
- Description
- Availability
- Setup state
- Permission level
- Assistants allowed to use it
- Last used
- Health
- Logs

## Tool categories

### Observe

- Browser screenshot
- DOM inspect
- File read
- Search
- Calendar read
- Email read

### Act

- Browser click/type
- Send message
- Create calendar event
- Draft email
- Write file
- Trigger job

### Compute

- Run code in sandbox
- Analyze data
- Call model
- Route to another agent

### Communicate

- Telegram
- WhatsApp
- Slack
- Discord
- Email

## Permission levels

Use clear labels:

- Disabled
- Read-only
- Draft-only
- Ask before acting
- Auto-act within rules
- Full access

Examples:

```text
Gmail: Draft-only, requires approval to send
Browser: Auto-act on local dev URLs, ask before submitting forms
Files: Read project folder only
Shell: Disabled
Telegram: Reply only to allowed users
```

## Tool activity timeline

Every tool use should be visible.

Example:

```text
09:01 Browser QA Assistant opened http://localhost:3000
09:01 Captured screenshot
09:02 Clicked Settings cog
09:02 Console error detected: Cannot read property provider of undefined
09:03 Created QA issue
```

## Honest unavailable state

If Chrome MCP is missing:

```text
Chrome MCP is not connected.
To enable browser control, configure a Chrome/Playwright MCP server.
Until then, browser actions are unavailable.
```

Do not fake screenshots.
Do not fake browser actions.
Do not say Chrome is connected when it is not.
