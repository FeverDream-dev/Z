# Feature Matrix

This matrix helps prevent the product from collapsing into a chat clone.

| Area | Chat clone version | Zsistant version |
|---|---|---|
| Primary object | Conversation | Assistant |
| Home screen | Recent chats | Assistants, jobs, channels, health |
| Chat | Main product | One surface of an assistant |
| Settings | Basic theme/model | Providers, tools, MCP, channels, memory, developer mode |
| Tools | Hidden or fake | Registry with permissions and logs |
| Browser | Not present | Chrome/MCP/browser sessions and screenshots |
| Channels | Web only | Web, CLI, Telegram, WhatsApp, Discord, Slack, etc. |
| Jobs | Not present | Scheduled recurring responsibilities |
| Automations | Not present | Event-triggered workflows |
| Memory | Vague prompt | Editable global and assistant-specific memory |
| Knowledge | File upload maybe | Assistant/project knowledge sources with status |
| Personas | Prompt preset | Full assistant identity and behavior policy |
| Developers | No visibility | Developer Mode with traces and diagnostics |
| Reliability | Fake success | Honest setup, failure, and permission states |
| UI | Chat page | Assistant control center |

## Minimum product modules

| Module | Must exist? | Notes |
|---|---:|---|
| Assistant list | Yes | Real assistants, not fake cards |
| Assistant detail | Yes | Overview plus tabs |
| Chat tab | Yes | Belongs to assistant |
| Settings cog | Yes | Real route/modal |
| Developer Mode | Yes | Toggle and visible advanced UI |
| Tools tab | Yes | Real or honest setup states |
| Channels tab | Yes | Web plus external channel setup states |
| Jobs tab | Yes | Real or honest disabled state |
| Browser tab | Yes | Chrome/MCP connected or honest unavailable state |
| Knowledge tab | Yes | Files/docs/context state |
| Memory tab | Recommended | If not implemented, state clearly |
| Logs tab | Yes | Timeline of actions/errors |
| Provider config | Yes | Missing credentials must be clear |

## Scope guard

Any time the coding AI adds a feature, ask:

- Which assistant does this belong to?
- Is this real or fake?
- How does the user configure it?
- How does the user know if it worked?
- How does the user stop it?
- What happens when it fails?
- Is it visible in Developer Mode?
