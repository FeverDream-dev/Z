# 13 — Data Model Sketch

This is a conceptual data model. OpenCode should convert it into the chosen persistence layer.

## AgentProfile

```text
id
name
slug
role
persona_path
workspace_root
created_at
updated_at
status
```

## AgentPermissions

```text
agent_id
filesystem_scope
tools_allowed
tools_denied
channels_allowed
peers_allowed
requires_approval_for
```

## ChannelBinding

```text
id
agent_id
channel_type
external_account_id
conversation_id
display_name
enabled
created_at
```

## Job

```text
id
agent_id
source_channel
objective
status
current_step
checkpoint_summary
created_at
updated_at
completed_at
```

## JobEvent

```text
id
job_id
agent_id
event_type
message
metadata
created_at
```

## ModelProvider

```text
id
name
type
base_url
models
priority
health_status
budget_policy
```

## PersonaSignal

```text
agent_id
signal_type
score
confidence
evidence_ref
created_at
```

## PeerAcl

```text
from_agent_id
to_agent_id
allowed_kinds
scope
expires_at
created_by
```

## SkillScan

```text
id
source
summary
risk_level
risk_findings
translation_plan_path
created_at
```
