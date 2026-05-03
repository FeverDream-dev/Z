# Master Prompt for the Coding AI

Paste this into the coding AI after adding this folder to the repository.

```text
You previously misunderstood the product.

You built toward a simple ChatGPT clone. That is not the intended product.

Before editing code, read the entire folder `zsistant_openclaw_assistant_briefing/`.

The product is Zsistant / zazi: an assistant-first AI workspace and control center.

It must let the user create and manage multiple assistants. Each assistant can have persona, knowledge, memory, tools, channels, scheduled jobs, automations, browser/Chrome/MCP capabilities, provider/model routing, permissions, logs, and developer diagnostics.

The central object is Assistant, not Message.
The central experience is managing digital workers, not only sending prompts.

Your mission:

1. Read all documents in `zsistant_openclaw_assistant_briefing/`.
2. Create or update `PRODUCT_SCOPE.md` in the repo using those documents as the north star.
3. Audit the current app and identify where it is only a chat clone.
4. Create or update `AUDIT_AND_REBUILD_PLAN.md`.
5. Remove fake/mock/placeholder production behavior.
6. Rebuild the product around assistants, not only conversations.

Required product areas:

- Assistant manager
- Assistant detail/control center
- Assistant chat tab
- Assistant profile/persona
- Tools registry
- MCP/browser/Chrome area
- Channels area for Web UI, Telegram, WhatsApp, Discord, Slack, CLI, etc.
- Knowledge and memory areas
- Jobs and automations
- Logs and activity timeline
- Settings cog and settings page
- Developer Mode toggle and diagnostics
- Provider/model configuration
- Honest setup states for unavailable integrations

Do not fake capabilities.

If Telegram is not connected, show a setup state.
If Chrome MCP is not connected, show a setup state.
If provider credentials are missing, show a setup state.
If jobs do not run yet, do not show fake job history.
If a button does nothing, remove it or make it real.

Acceptance criteria:

- The app clearly looks and behaves like an assistant platform.
- Assistants are first-class objects.
- Chat is only one tab/surface of an assistant.
- Settings exists and works.
- Developer Mode exists and changes the UI meaningfully.
- Tools/MCP/browser capability is honestly represented.
- Channels are represented as real connections or honest setup states.
- Jobs/automations are represented as real scheduled responsibilities or honest disabled states.
- No fake production behavior remains.
- UI is polished, simple, and Apple-like for normal users.
- Advanced details are visible in Developer Mode.

Do not continue building only chat features.

Start by reading the briefing folder and writing the audit/product scope. Then start implementing the assistant-first product structure.
```
