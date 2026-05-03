# Anti-Patterns

These are signs the coding AI is building the wrong product.

## Anti-pattern: Assistant equals prompt preset

Wrong:

```text
Assistant = name + system prompt
```

Correct:

```text
Assistant = identity + persona + tools + channels + memory + knowledge + jobs + permissions + logs + provider strategy
```

## Anti-pattern: UI claims capability without backend reality

Wrong:

```text
Button says "Connect Telegram" but opens nothing.
```

Correct:

```text
Button opens real setup, or says integration unavailable in this build.
```

## Anti-pattern: Fake job history

Wrong:

```text
Shows "Daily briefing completed" using static demo data.
```

Correct:

```text
Shows no runs yet, or real run records.
```

## Anti-pattern: Hardcoded assistant response

Wrong:

```text
Assistant always replies from a canned local string.
```

Correct:

```text
Real provider call, or missing-provider setup state.
```

## Anti-pattern: Developer mode is only a label

Wrong:

```text
Toggle changes nothing.
```

Correct:

```text
Toggle exposes traces, diagnostics, raw structures, tool logs, provider health.
```

## Anti-pattern: Browser is decorative

Wrong:

```text
Browser tab has a fake screenshot.
```

Correct:

```text
Browser tab shows real Chrome/MCP connection or honest unavailable state.
```

## Anti-pattern: Chat is the whole app

Wrong:

```text
Home -> chat.
Everything else hidden or absent.
```

Correct:

```text
Home -> assistant manager, jobs, channels, health, recent activity.
Chat is one assistant tab.
```

## Anti-pattern: Too technical for normal users

Wrong:

```text
Raw JSON, logs, tokens, provider payloads shown everywhere.
```

Correct:

```text
Normal mode is simple. Developer Mode reveals internals.
```

## Anti-pattern: No permission model

Wrong:

```text
Assistant can use all tools by default.
```

Correct:

```text
Tools have permission levels and approval policies.
```

## Anti-pattern: No chronology

Wrong:

```text
Only conversations exist.
```

Correct:

```text
Jobs, run history, upcoming actions, logs, and timelines exist.
```

## Anti-pattern: Generic dashboard aesthetics

Wrong:

```text
Random cards, fake charts, starter template look.
```

Correct:

```text
Polished, calm, coherent, assistant-focused product UI.
```
