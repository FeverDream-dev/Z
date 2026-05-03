# 17 — OpenClaw Migration Playbook

## Goal

Help users migrate ideas, personas, workflows, skills, and memory from OpenClaw-like systems without copying heavy architecture or unsafe behaviors.

## Prompt to ask the old tool

```text
Export a complete Markdown summary of this agent/system for migration.
Include:
- agent names and purposes
- persona/system prompts
- installed skills and what they do
- connected channels
- cron jobs and recurring tasks
- important memory
- file/workspace layout
- tool permissions
- examples of successful workflows
- known failures or quirks
- what you need in a new system to continue working
Do not include secrets, tokens, passwords, cookies, or private keys.
```

## Zsistant migration process

1. User uploads/pastes the Markdown export.
2. Zsistant extracts entities:
   - agents
   - skills
   - channels
   - workflows
   - memory
   - permissions
3. Zsistant produces:
   - migration summary
   - risk warnings
   - proposed Zsistant agents
   - proposed persona files
   - proposed job templates
   - proposed channel bindings
4. User approves what to create.
5. Zsistant creates local agents and personas.
6. Skills remain analyze-only until reviewed.

## Migration principles

- Do not import secrets.
- Do not execute old skill code.
- Do not preserve unsafe broad permissions.
- Do not assume old memory is clean.
- Prefer clean personas and small job templates.
