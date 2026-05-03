# What an Assistant Is

A chatbot answers messages.

An assistant helps the user achieve outcomes.

That difference changes everything.

## Chatbot definition

A chatbot is usually reactive:

- User sends message.
- Model replies.
- Conversation ends until the next prompt.

It has little or no durable identity. It usually does not own jobs, tools, schedules, or long-lived responsibilities.

## Assistant definition

An assistant is a persistent actor that can help across time.

It can:

- Understand the user's goals
- Remember useful context
- Use tools
- Ask for permission
- Take action
- Run jobs later
- Monitor events
- Work across channels
- Coordinate with other assistants
- Keep logs of what it did
- Report status
- Recover from failure

An assistant has identity and responsibility.

## The assistant as a digital worker

Think of each assistant as a worker with a desk.

Its desk contains:

- Instructions
- Skills
- Files
- Tools
- Calendar
- Inbox
- Browser
- Notes
- Memory
- Task board
- Logs
- Permissions

The user should be able to inspect that desk.

## The assistant lifecycle

An assistant should have a lifecycle:

1. Created
2. Configured
3. Connected to channels/tools
4. Given knowledge
5. Given permissions
6. Used in conversations
7. Assigned jobs
8. Monitored
9. Improved
10. Archived or deleted

If the product only supports "send message and receive response," it does not support an assistant lifecycle.

## The assistant contract

Every assistant should answer these questions:

- Who are you?
- What are you responsible for?
- What are you allowed to do?
- What tools can you use?
- What channels can reach you?
- What do you know?
- What do you remember?
- What are you working on now?
- What will you do later?
- What did you do recently?
- What failed?
- What needs approval?

## Assistant behavior examples

A real assistant can say:

- "I checked your calendar and found three conflicts tomorrow."
- "I drafted replies to the two urgent emails, but I need approval before sending."
- "Your Telegram automation failed because the bot token expired."
- "The website monitor job ran at 09:00 and found no changes."
- "I opened Chrome, captured the page screenshot, and found the login modal blocking the test."
- "I routed this task to the Coding Assistant because it requires repository changes."

A weak chatbot says:

- "Sure, I can help with that."

## Implication for Zsistant

Zsistant must provide assistant management, not just chat management.

The UI, data model, navigation, settings, and developer mode should all start from the assistant concept.
