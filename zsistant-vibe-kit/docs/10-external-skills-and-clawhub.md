# 10 — External Skills and ClawHub Importer

## Goal

Zsistant should learn from external skill ecosystems without trusting them blindly.

MVP behavior: **analyze and translate, do not execute**.

## Skill importer flow

1. User points Zsistant to a skill folder, URL, or pasted Markdown.
2. Zsistant reads metadata and instructions.
3. Analyzer summarizes capabilities.
4. Analyzer flags risk indicators.
5. Analyzer maps the skill to Zsistant-native capabilities:
   - channel adapter
   - tool adapter
   - persona patch
   - job template
   - browser workflow
   - external API integration
6. Zsistant proposes an implementation plan.
7. User approves before any generated code or execution.

## Risk levels

### Low

- Pure instructions.
- No external commands.
- No credential access.
- No filesystem access.

### Medium

- Reads local files.
- Connects to an external API.
- Requires tokens.
- Suggests browser automation.

### High

- Executes commands.
- Installs dependencies.
- Downloads scripts.
- Reads browser profiles, wallets, SSH keys, password stores.
- Obfuscated shell commands.
- Hidden instruction injection.

## Outputs

The importer should produce:

- `skill-summary.md`
- `risk-report.md`
- `zsistant-translation-plan.md`
- optional generated tests after approval

## Migration from OpenClaw

The user can ask OpenClaw to export a large Markdown file with:

- installed skills
- agent/persona settings
- channels
- workflows
- cron jobs
- examples of previous tasks
- failure points

Zsistant should ingest that Markdown and propose a clean local setup.

## Never do by default

- Do not run install scripts.
- Do not execute shell blocks.
- Do not mount the user's home folder into skill execution.
- Do not import credentials.
- Do not trust marketplace ranking or popularity.
