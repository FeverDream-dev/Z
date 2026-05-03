# 07 — Persona Trainer

## Goal

Each Zsistant agent should adapt to the user without becoming unstable or fake.

The trainer watches how the user communicates, how the user reacts, and what outputs are accepted or corrected.

It proposes small changes to `persona.md`.

## Inputs

- User messages.
- User corrections.
- Explicit stars/ratings.
- Whether the user asks for shorter/longer answers.
- Whether the user asks for more technical or more casual tone.
- Whether the user repeats instructions.
- Whether the user accepts generated outputs.

## Style dimensions

Track these as small scores, not as permanent assumptions:

- brevity: short / balanced / detailed
- tone: formal / friendly / street / technical / executive
- explanation depth: low / medium / high
- autonomy: ask first / act with updates / act aggressively
- jargon tolerance: low / medium / high
- emoji tolerance: none / light / expressive
- structure preference: bullets / paragraphs / checklists / tables
- status preference: frequent / only important / final only

## Persona update policy

The trainer must not rewrite the whole persona after one message.

Recommended behavior:

1. Observe patterns.
2. Keep rolling notes.
3. Propose a small persona patch.
4. Apply automatically only if confidence is high and patch is low risk.
5. Ask the user for major tone/personality changes.

## Example persona notes

```markdown
## Learned communication preferences

- User prefers direct progress and dislikes vague promises.
- User likes seeing implementation maps before code.
- User uses casual language and accepts informal phrasing.
- User wants human approval before releases.
- User values real-time coding visibility through OpenCode.
```

## Trainer sub-agent

A lightweight trainer sub-agent can run cheap classification tasks:

- infer style signal
- detect frustration
- detect correction
- propose persona patch
- identify repetitive failures

The trainer should use short context windows and summaries to avoid token waste.

## Safety

Do not infer sensitive personal traits. Only adapt working style and communication preferences.
