# WhatsApp Adapter

## Goal

Add WhatsApp Cloud API setup checklist and webhook adapter skeleton.

## Inputs to read first

- `AGENTS.md`
- `docs/00-master-map.md`
- Relevant docs for this phase

## Tasks

- Define phone number ID/access token/verify token config.
- Implement webhook verification path.
- Normalize inbound messages.
- Support dry-run send.
- Document setup friction.

## Acceptance criteria

- Webhook verify test passes.
- Dry-run inbound event creates job.
- No real token needed for tests.

## Required end-of-step report

```text
Changed:
- ...

Validated:
- command: ...
- result: ...

Manual check:
- ...

Risks / notes:
- ...

Next step:
- ...
```

## Human approval needed before

- Using real credentials.
- Installing global tools.
- Pushing to GitHub.
- Running external skill/plugin code.
