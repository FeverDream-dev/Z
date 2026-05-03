# External Skill Analyzer

## Goal

Analyze ClawHub/OpenClaw-style skills without executing them.

## Inputs to read first

- `AGENTS.md`
- `docs/00-master-map.md`
- Relevant docs for this phase

## Tasks

- Read SKILL.md and supporting text files.
- Detect risky patterns.
- Produce summary/risk/translation plan.
- Add fixtures for safe and malicious-looking skills.
- Never run skill commands.

## Acceptance criteria

- Safe fixture gets low risk.
- Malicious fixture gets high risk.
- No commands from skill are executed.

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
