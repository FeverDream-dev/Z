# 04 — Security and Isolation Model

Zsistant must be safe by default because it connects agents to messaging platforms, files, tools, browsers, and eventually email.

## Default deny

Agents cannot:

- Read another agent's workspace.
- Write outside their workspace.
- Talk to another agent.
- Execute third-party skills.
- Use real external credentials.
- Publish releases.
- Browse authenticated websites.

unless explicitly allowed.

## Isolation levels

### Level 1 — Logical isolation

- Per-agent workspace root.
- Path-scope checks.
- ACLs for inter-agent messages.
- Per-agent tools and channels.

This is MVP.

### Level 2 — Process isolation

- Risky tools run in worker processes.
- Worker gets a limited environment and filesystem view.
- Worker has a job timeout.

This is post-MVP.

### Level 3 — Container/microVM isolation

- Dangerous or untrusted skills run in containers or microVMs.
- No host secrets.
- No broad filesystem mounts.
- Egress network policy.

This is required before any external skill execution marketplace.

## Skill safety policy

External skills are data, not trusted code.

The analyzer should flag:

- pipe-to-shell installation instructions
- curl/wget into shell
- base64 or obfuscated commands
- credential harvesting language
- browser profile access
- wallet/SSH key access
- broad filesystem reads
- suspicious remote endpoints
- requests to disable security tools
- hidden prompt injection instructions

## Channel safety

### Telegram

- Use token from environment/config secret store only.
- Do not log token.
- Verify bot identity during setup.
- Store chat bindings per agent.

### Discord

- Prefer minimal permissions.
- Verify interaction signatures if using HTTP interactions.
- Keep bot token secret.
- Avoid broad guild permissions in MVP.

### WhatsApp

- Verify webhook token.
- Store phone number ID and access token securely.
- Keep template/message sending explicit.
- Watch rate limits and delivery status.

## Browser/MCP safety

Chrome MCP can inspect real pages and potentially authenticated sessions. It should be used for local development/testing only unless the user explicitly chooses otherwise.

## Release safety

Z Assistant may prepare a release checklist and changelog. It must not push, tag, or publish without explicit human approval.
