# Jobs, Automations, and Chronology

A real assistant works across time.

This is one of the biggest differences between an assistant platform and a chat app.

## Jobs

A job is scheduled work.

Examples:

- Every morning at 8 AM, prepare a daily briefing.
- Every Friday, summarize project progress.
- Every 15 minutes, check a website for changes.
- Every night, run browser QA on the app.
- Every Monday, clean and summarize the inbox.

## Automations

An automation is event-triggered work.

Examples:

- When a Telegram message contains "urgent", notify the user.
- When an email arrives from a VIP sender, summarize it.
- When a GitHub issue is opened, triage it.
- When a browser QA check fails, create a report.
- When a file is uploaded, extract knowledge and attach it to a project.

## Chronology

Chronology is the timeline of what happened, what is happening, and what will happen.

Zsistant should make time visible.

For every assistant, show:

- Upcoming jobs
- Recent runs
- Current tasks
- Failed tasks
- Pending approvals
- Next scheduled actions

## Job object concept

Each job should have:

- Name
- Assistant owner
- Purpose
- Schedule or trigger
- Required tools
- Required permissions
- Next run
- Last run
- Status
- Output history
- Failure history
- Retry policy
- Approval policy

## Job statuses

Use honest states:

- Draft
- Active
- Paused
- Running
- Waiting for approval
- Completed
- Failed
- Disabled due to missing setup

## Example jobs

### Daily Briefing

```text
Assistant: Personal Chief of Staff
Schedule: Every weekday at 08:00
Inputs: Calendar, email summaries, project updates, weather if available
Output: Telegram message and Web UI briefing
Approval: not required for read-only summary
```

### Nightly Browser QA

```text
Assistant: Browser QA Assistant
Schedule: Every night at 02:00
Inputs: local app URL, test checklist
Tools: Chrome MCP, screenshot capture, console log reader
Output: QA report with screenshots
Approval: not required on local dev URL
```

### Community Triage

```text
Assistant: Telegram Community Assistant
Trigger: New Telegram message in support group
Rule: Respond only when mentioned or when confidence is high
Tools: FAQ knowledge base, product docs
Output: Reply draft or direct response depending on permission
Approval: required for sensitive topics
```

### Coding Agent Supervisor

```text
Assistant: Coding Supervisor
Trigger: Sentry error or user command
Tools: repository reader, issue tracker, coding agent, browser QA
Output: investigation summary, proposed fix, PR plan
Approval: required before code-changing action
```

## UI expectations

The Web UI should include:

- Jobs list
- Calendar/timeline view
- Next run indicators
- Run history
- Logs per run
- Manual run button
- Pause/resume
- Failure details
- Approval queue

## Anti-patterns

Do not show fake scheduled jobs.
Do not show fake job history.
Do not use visual-only cron cards.
Do not claim a job ran if there is no run record.

If scheduling is not implemented, show an honest disabled state.
