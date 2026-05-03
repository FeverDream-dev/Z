# 16 — Chrome MCP Workflow

Chrome DevTools MCP should be used as a development/testing tool for the local web UI.

## Use cases

- Open local web UI.
- Inspect console errors.
- Inspect network calls.
- Validate API endpoints.
- Record performance traces.
- Test forms and navigation.
- Capture screenshots for review.

## Suggested OpenCode task

```text
Use Chrome DevTools MCP to open the local Zsistant web UI, create an agent, send a test message, and report any console/network errors. Do not use authenticated personal websites.
```

## Safety

- Do not automate logged-in personal accounts unless the user explicitly asks.
- Do not expose cookies, tokens, or local storage values in logs.
- Do not test real Telegram/Discord/WhatsApp credentials in browser automation.

## Manual fallback

If MCP is unavailable, OpenCode should provide:

- URL to open.
- click path.
- expected result.
- screenshots requested from user only if needed.
