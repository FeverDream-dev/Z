# 14 — Testing and Evaluation Plan

## Test layers

### Unit tests

- config parsing
- path isolation
- provider fallback decision
- job state transitions
- persona trainer scoring
- inter-agent ACL decisions
- skill risk scanning

### Integration tests

- CLI creates agent and sends message
- web UI can list agents
- Telegram test update routes to correct agent
- Discord test event routes to correct agent
- WhatsApp test webhook routes to correct agent
- provider timeout triggers retry/fallback

### Manual tests

- create two agents
- confirm they cannot see each other
- allow one specific peer request
- confirm audit log records it
- connect a fake/test channel
- inspect job status in UI

### Chrome MCP tests

When Chrome DevTools MCP is configured, use it to:

- open the local web UI
- verify dashboard loads
- create an agent through the UI
- send a chat message
- inspect network errors
- capture console errors
- test responsive layout

## Evals

Create small eval fixtures for:

- personal trainer style adaptation
- loop detection
- skill risk detection
- model fallback status language
- agent isolation refusals

## Golden rule

A feature that cannot be observed, tested, or explained should not be considered done.
