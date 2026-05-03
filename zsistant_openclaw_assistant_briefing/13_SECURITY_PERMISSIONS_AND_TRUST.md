# Security, Permissions, and Trust

An assistant that can act is powerful.

Power requires control.

## Trust model

The user must be able to see and control what each assistant can do.

Every assistant should have a permission profile.

## Permission categories

### Read permissions

- Read files
- Read calendar
- Read email
- Read web pages
- Read messages
- Read project docs

### Write/action permissions

- Send messages
- Send emails
- Create calendar events
- Edit files
- Run browser actions
- Run code or shell commands
- Trigger external APIs

### Communication permissions

- Which channels can the assistant use?
- Who can message it?
- Can it reply in groups?
- Does it need to be mentioned?

### Memory permissions

- Can it save memories automatically?
- Does it require approval?
- Can it remember sensitive facts?

## Approval policies

Use clear approval levels:

- Never allow
- Ask every time
- Allow for this session
- Allow within rules
- Always allow

Examples:

```text
Gmail send: Ask every time
Gmail draft: Allow within rules
Calendar read: Always allow
Calendar event creation: Ask every time
Browser screenshot: Allow within local dev URLs
Browser form submit: Ask every time
Shell command: Never allow unless Developer Mode and explicit opt-in
```

## Channel security

External messaging channels require extra safety.

Controls:

- User allowlist
- Group allowlist
- Pairing flow
- Mention required in groups
- Rate limiting
- Sensitive-action approval
- Easy pause
- Audit logs

## Tool safety

Every tool call should be logged.

Every dangerous tool should show:

- What will be accessed
- What will be changed
- Which assistant requested it
- Which user/channel triggered it
- Whether approval was granted

## Browser safety

Browser automation can accidentally click or submit forms.

Controls:

- Read-only browsing mode
- Local-dev-only mode
- Ask before submit
- Ask before purchase/payment
- Ask before sending forms
- Redaction for screenshots if needed

## Memory safety

Memory should be inspectable and removable.

The assistant should not silently store sensitive data unless the user explicitly enables that behavior.

## Developer Mode safety

Developer Mode can expose raw prompts, payloads, and tool traces.

It should be opt-in and clearly labeled.

## Reliability and trust

Trust also means honesty.

The product should never claim:

- A job ran if it did not run
- A tool is connected if it is not connected
- A model responded if the response is fake
- A browser screenshot exists if it was not captured
- A message was sent if it was only drafted

## Trust-building UI

Every assistant should have:

- Status badge
- Permission summary
- Recent activity
- Error log
- Pause button
- Approval queue
- Setup warnings

## Principle

The more powerful the assistant, the more visible its permissions and actions must be.
