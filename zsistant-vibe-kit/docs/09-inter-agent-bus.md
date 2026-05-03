# 09 — Inter-Agent Bus

## Default behavior

Agents do not talk to each other by default.

The user must explicitly allow communication.

## Communication model

Use typed JSON-like envelopes, not free-form hidden chatter.

Envelope fields:

```text
id
from_agent_id
to_agent_id
kind
request
attachments
constraints
created_at
expires_at
human_visible
```

Kinds:

- ask_summary
- ask_file
- ask_decision
- delegate_task
- return_result
- deny_request
- request_permission

## Permission model

ACL examples:

```text
agent_a may ask_summary from agent_b
agent_a may ask_file from agent_b only under shared/export/
agent_a may delegate_task to agent_b
agent_b may deny any request
```

## Audit

Every inter-agent request should be visible in both agents' audit logs.

## File transfer

Never hand over arbitrary filesystem paths.

Use:

1. exporting agent copies approved file into a shared transfer area
2. receiving agent gets a transfer reference
3. system logs hash/size/source/approval

## Human command examples

```text
Allow Sales Assistant to ask Research Assistant for summaries.
```

```text
Allow DevOps Agent and QA Agent to talk for this job only.
```

```text
Give Coder Agent the markdown migration export from OpenClaw Agent.
```
