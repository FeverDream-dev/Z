# Reliability and Honesty Rules

These rules protect the product from becoming a fake demo.

## Rule 1: No invisible fake success

If an action did not happen, do not show it as completed.

## Rule 2: Missing setup is a real state

A missing API key, disconnected MCP server, or unconfigured Telegram bot is not a failure of design. It is a state the UI should handle well.

## Rule 3: Disabled is better than fake

A disabled button with a clear explanation is better than a button that pretends to work.

## Rule 4: Every external action needs evidence

If the assistant says it sent, created, captured, opened, or scheduled something, there should be a record.

## Rule 5: Separate demo/test from production

Demo data must never appear as production state.

If demo data exists, it must be clearly marked and isolated.

## Rule 6: Logs are part of the product

Assistants should leave activity trails.

Logs are not optional for tools, jobs, channels, and automations.

## Rule 7: State must be inspectable

Users should be able to inspect:

- Assistant profile
- Tool permissions
- Channel routing
- Job schedule
- Memory
- Knowledge
- Recent actions
- Errors

## Rule 8: Powerful tools require guardrails

Browser, shell, email, calendar, messaging, and file-write tools need permissions and approvals.

## Rule 9: Fail clearly

A failed provider call should not look like the assistant is thinking forever.

A failed job should not disappear.

A failed tool call should show a human-readable reason.

## Rule 10: Build product truth before product polish

The UI can be beautiful only if it represents real state.

A polished fake UI is worse than an honest simple UI.

But the final goal is both: real and beautiful.
