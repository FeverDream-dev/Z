# Acceptance Criteria

Use this checklist to judge whether the product is becoming Zsistant or falling back into a chat clone.

## Product identity

- [ ] The app clearly presents itself as an assistant platform, not only a chat app.
- [ ] Assistants are first-class objects.
- [ ] The home screen or main navigation makes assistant management obvious.
- [ ] There is a product scope document in the repository that says this explicitly.

## Assistant manager

- [ ] User can see a list/grid of assistants.
- [ ] User can create an assistant.
- [ ] User can open an assistant detail page.
- [ ] Assistant has profile, purpose, persona, tools, channels, memory/knowledge, jobs, logs, and settings.
- [ ] Assistant cards reflect real state, not fake demo stats.

## Chat

- [ ] Chat exists as one tab/surface for an assistant.
- [ ] Messages belong to a real assistant/conversation.
- [ ] Provider/model behavior is real or honestly missing.
- [ ] No hardcoded fake assistant replies in production.

## Settings

- [ ] There is a real settings cog.
- [ ] Settings page has meaningful sections.
- [ ] Developer Mode toggle exists.
- [ ] Provider configuration exists or shows honest missing state.
- [ ] Appearance/preferences are real if shown.

## Developer Mode

- [ ] Toggle changes the UI meaningfully.
- [ ] Tool traces are visible if tools exist.
- [ ] Provider/model routing details are visible if available.
- [ ] Raw conversation/request details are visible where safe.
- [ ] MCP/browser diagnostics are visible if available.

## Tools and MCP

- [ ] Tools are listed in a registry or assistant tool tab.
- [ ] Tool state is honest: disabled, needs setup, connected, failed.
- [ ] MCP server state is visible if MCP is supported.
- [ ] No fake tool calls.

## Browser/Chrome

- [ ] Browser/Chrome capability is represented as a real tool/surface.
- [ ] If connected, user can see browser status and actions.
- [ ] If unavailable, the UI clearly says setup is needed or feature is unavailable.
- [ ] No fake screenshots or browser actions.

## Channels

- [ ] Web UI is one channel.
- [ ] Telegram/WhatsApp/Discord/Slack-like channels are modeled as connections, not fake cards.
- [ ] Channel setup states are honest.
- [ ] Routing to assistants is conceptually supported.

## Jobs and automations

- [ ] Jobs are first-class.
- [ ] User can see scheduled jobs for an assistant.
- [ ] Job state is honest.
- [ ] No fake cron history.
- [ ] Failed jobs are visible.

## Memory and knowledge

- [ ] Knowledge sources are distinct from memory.
- [ ] Persona is distinct from memory and knowledge.
- [ ] Uploaded/attached knowledge has honest processing state.
- [ ] Memory is inspectable/editable if implemented.

## UX quality

- [ ] UI feels premium and cohesive.
- [ ] Empty states educate the user.
- [ ] No dead buttons.
- [ ] No placeholder panels in production paths.
- [ ] Normal users are not overwhelmed.
- [ ] Developer controls are discoverable when Developer Mode is on.

## Reliability

- [ ] Missing credentials produce setup states, not fake success.
- [ ] Failed integrations show errors.
- [ ] Logs exist for important assistant actions.
- [ ] Browser QA or manual UI verification has been performed.

## Anti-placeholder audit

Search the repository for:

- placeholder
- mock
- mocked
- fake
- dummy
- stub
- TODO
- FIXME
- demo
- sample
- lorem
- coming soon
- not implemented
- hardcoded
- no-op
- canned response

Every result must be removed, replaced, moved to test/demo-only context, or explicitly justified.

## Final acceptance sentence

The product is acceptable only when a user can say:

"I am managing assistants that can work across tools, channels, browser, jobs, memory, and models."

Not merely:

"I can chat with an AI."
