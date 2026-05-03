# 05 — Agent Runtime Design

## Runtime objects

### Agent

A named persona with isolated state, tool permissions, channel bindings, and job history.

Fields to design:

- id
- name
- role
- persona path
- workspace root
- enabled channels
- tool permissions
- model policy
- memory policy
- peer ACL
- status

### Job

A durable unit of work.

Fields to design:

- id
- agent id
- source channel
- user message or trigger
- objective
- status
- current step
- model/provider attempts
- tool attempts
- retry policy
- loop detector state
- created/updated timestamps
- result

### Event

Append-only audit/status record.

Examples:

- job.created
- job.started
- llm.requested
- llm.timeout
- llm.fallback
- tool.called
- tool.failed
- agent.message.sent
- agent.peer.requested
- human.approval.required
- job.completed
- job.failed

## Job lifecycle

```text
queued -> planning -> acting -> waiting -> retrying -> completed
                           └-> waiting_for_human
                           └-> loop_suspected
                           └-> failed
```

## Long-running behavior

When a provider times out, the user should not see a raw timeout unless useful. The status should become:

```text
Still working. Provider timed out, retry scheduled.
```

The system should retry with backoff, then switch provider if configured.

## Death-loop detection

Signals:

- same tool called repeatedly with same args
- same model error repeated
- no new useful event after N attempts
- token budget burned without state progress
- job objective unchanged after multiple loops
- provider fallback ping-pong

Actions:

1. Pause job.
2. Summarize loop evidence.
3. Resume from last stable checkpoint if possible.
4. Ask a critic sub-agent or human for next move.

## Checkpointing

Every job should maintain:

- objective
- plan
- completed steps
- pending steps
- files touched
- last successful tool output summary
- retry/fallback history

This allows restart after crash or provider failure.
