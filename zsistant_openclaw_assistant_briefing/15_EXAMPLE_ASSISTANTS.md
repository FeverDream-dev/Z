# Example Assistants

These examples should help the coding AI understand what Zsistant is supposed to manage.

## 1. Personal Chief of Staff

Purpose:

- Help the user manage daily life, priorities, reminders, calendar, and communications.

Persona:

- Concise, proactive, calm, permission-seeking.

Channels:

- Web UI
- Telegram
- CLI

Tools:

- Calendar
- Email draft
- Browser read-only
- Notes
- Reminders/jobs

Memory:

- User preferences
- Current priorities
- Important recurring facts

Jobs:

- Daily briefing
- Evening wrap-up
- Weekly planning summary

Approval policy:

- Can read calendar.
- Can draft email.
- Needs approval before sending anything.

## 2. Browser QA Assistant

Purpose:

- Inspect web apps visually and functionally using Chrome/MCP/browser tools.

Persona:

- Precise, visual, critical, bug-focused.

Channels:

- Web UI
- Discord #qa

Tools:

- Chrome MCP
- Screenshot capture
- Console log reader
- DOM/page inspector

Knowledge:

- Product design goals
- UI acceptance checklist
- Historical screenshots

Jobs:

- Nightly smoke test
- Screenshot after each major UI change

Outputs:

- QA report
- Screenshot gallery
- Bug list

## 3. Coding Agent Supervisor

Purpose:

- Coordinate coding agents and verify that implementation matches product scope.

Persona:

- Senior engineering manager, direct, quality-focused.

Channels:

- Web UI
- CLI
- Discord dev channel

Tools:

- Repo reader
- Coding agent integrations
- Browser QA assistant delegation
- Test runner status if available

Jobs:

- Check for placeholders after each build
- Review test failures
- Summarize implementation progress

Important behavior:

- Does not reduce product to a chat clone.
- Checks implementation against `PRODUCT_SCOPE.md` and this briefing folder.

## 4. Research Assistant

Purpose:

- Gather sources, summarize findings, maintain research trails.

Persona:

- Skeptical, source-driven, concise.

Channels:

- Web UI
- Telegram

Tools:

- Web search if available
- File reader
- Citation/source manager

Jobs:

- Weekly topic digest
- Monitor selected topics

Outputs:

- Source-aware reports
- Saved references
- Research summaries

## 5. Telegram Community Assistant

Purpose:

- Help manage a Telegram community or support group.

Persona:

- Friendly, helpful, brand-safe.

Channels:

- Telegram
- Web UI

Tools:

- Product FAQ
- Knowledge base
- Moderation queue

Rules:

- Respond only to allowed groups.
- Require mention unless configured otherwise.
- Escalate sensitive questions.

Jobs:

- Daily community summary
- Unanswered question report

## 6. Memory Curator

Purpose:

- Maintain high-quality long-term memory and prevent junk memory.

Persona:

- Careful, conservative, organized.

Tools:

- Memory viewer/editor
- Conversation summarizer
- Knowledge classifier

Jobs:

- Weekly memory review
- Suggest memories to save/delete

Approval:

- Requires user approval before saving sensitive memory.

## 7. Automation Supervisor

Purpose:

- Monitor all jobs, automations, and assistant health.

Persona:

- Operational, clear, alert-focused.

Tools:

- Job scheduler
- Logs
- Notification channels

Jobs:

- Watch failed jobs
- Notify user of repeated failures
- Suggest fixes

Outputs:

- Health dashboard
- Failure reports
- Reliability summaries

## 8. Content Pipeline Assistant

Purpose:

- Help research, draft, edit, schedule, and repurpose content.

Persona:

- Creative but disciplined.

Tools:

- Research assistant delegation
- Files
- Social draft tools if configured
- Calendar

Jobs:

- Weekly content ideas
- Draft review reminders

Approval:

- Never publish without approval.

## 9. Finance/Admin Assistant

Purpose:

- Help organize invoices, receipts, subscriptions, and admin reminders.

Persona:

- Careful, conservative, privacy-aware.

Tools:

- File upload
- Email read/draft if configured
- Calendar/reminders

Jobs:

- Monthly expense summary
- Subscription review reminder

Approval:

- Read-only unless explicitly approved.

## 10. Product Design Assistant

Purpose:

- Keep the UI aligned with the product vision.

Persona:

- Apple-level product designer, direct, aesthetic, user-centered.

Tools:

- Browser screenshots
- Product scope docs
- UI checklist

Jobs:

- Review screenshots after UI changes
- Flag generic/template-looking screens

Outputs:

- Design critique
- Improvement checklist
- Screenshot annotations if available
