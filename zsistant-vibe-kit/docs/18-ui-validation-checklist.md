# UI Validation Checklist

## Goal
Validate the Zsistant web UI manually or via Chrome DevTools MCP.

## MCP Setup (Optional)

1. Ensure Chrome DevTools MCP is configured in your OpenCode environment.
2. MCP is used **only** for local development testing.
3. No authenticated accounts are automated by default.

## Manual Checklist

### Dashboard
- [ ] Open http://localhost:8080/
- [ ] Verify page loads without console errors
- [ ] Confirm agent list displays correctly
- [ ] Check dark theme renders properly

### Chat
- [ ] Navigate to /chat?agent=<test-agent>
- [ ] Send a test message
- [ ] Verify response appears
- [ ] Check job timeline updates
- [ ] Confirm no network errors in devtools

### API
- [ ] GET /health returns `{"status":"ok"}`
- [ ] GET /api/agents returns agent list
- [ ] GET /api/providers returns provider health
- [ ] POST /api/chat creates job and returns response
- [ ] GET /api/jobs/<agent_id> returns event timeline

### Responsive
- [ ] Test at 1920x1080 (desktop)
- [ ] Test at 768x1024 (tablet)
- [ ] Test at 375x667 (mobile)

## MCP Automated Checklist

When Chrome DevTools MCP is available:

1. Open local UI at http://localhost:8080/
2. Capture screenshot of dashboard
3. Inspect console for errors
4. Inspect network panel for failed requests
5. Navigate to chat page
6. Send test message via automation
7. Capture screenshot of response
8. Record performance trace
9. Export findings to validation report

## Safety Rules

- Do NOT automate logged-in personal accounts.
- Do NOT expose cookies, tokens, or localStorage in logs.
- Do NOT test real Telegram/Discord/WhatsApp credentials.
- MCP use is strictly optional; manual testing is always valid.

## Report Template

```text
Date: [YYYY-MM-DD]
Tester: [name or MCP]
Browser: [Chrome version]
Results:
- Dashboard: [PASS/FAIL]
- Chat: [PASS/FAIL]
- API: [PASS/FAIL]
- Console Errors: [count]
- Network Errors: [count]
Notes: [any issues found]
```
