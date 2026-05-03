# Memory, Knowledge, and Persona

These three concepts are different and must not be merged into one vague prompt field.

## Persona

Persona is how the assistant behaves.

Examples:

- Friendly
- Direct
- Executive
- Careful
- Technical
- Supportive
- Formal
- Playful

Persona controls tone, communication style, and decision behavior.

## Knowledge

Knowledge is information the assistant can reference.

Examples:

- Uploaded files
- Project docs
- Company docs
- Notes
- Web pages
- Manuals
- Codebase summaries
- User-provided instructions

Knowledge is usually externally supplied or retrieved.

## Memory

Memory is durable learned context about the user, projects, preferences, and repeated facts.

Examples:

- User prefers concise answers.
- User is building Zsistant as an assistant platform, not a chat clone.
- The default project is called zazi.
- Browser QA should use screenshots.

Memory should be editable and inspectable.

## Why separation matters

A coding assistant may have:

- Persona: direct senior engineer
- Knowledge: repository docs and architecture files
- Memory: user preferences and past decisions

A Telegram community assistant may have:

- Persona: warm community moderator
- Knowledge: FAQ and product docs
- Memory: recurring community issues

## Assistant-specific memory

Some memories should belong only to one assistant.

Example:

```text
Browser QA Assistant remembers the app's baseline screenshot layout.
```

## Global memory

Some memories should apply everywhere.

Example:

```text
The user wants Zsistant to be a full assistant platform, not a chat clone.
```

## Memory controls

The UI should let users:

- See memories
- Add memories
- Edit memories
- Delete memories
- Disable memory for an assistant
- Make a memory global
- Make a memory assistant-specific
- Require approval before saving memory

## Knowledge controls

The UI should let users:

- Upload files
- Attach files to assistants
- Attach files to projects
- Remove files
- See indexing/parsing status
- See retrieval/source usage if available

## Persona controls

The UI should let users define:

- Tone
- Role
- Response style
- Boundaries
- Decision policy
- Formatting preferences

## Example profile

```text
Assistant: Research Assistant
Persona: rigorous, skeptical, source-focused
Knowledge: uploaded research papers, saved web pages, project notes
Memory: user prefers concise executive summaries with source trails
Tools: web search, file reader, citation manager
Channels: Web UI, Telegram
Jobs: weekly research digest
```

## Anti-pattern

Do not create one giant "system prompt" text box and pretend that is a complete assistant.

A system prompt is one part of an assistant profile, not the whole product.
