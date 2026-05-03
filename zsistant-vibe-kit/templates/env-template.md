# Environment Template

Do not create a real `.env` with secrets until the implementation exists.

Suggested future secret names:

```text
ZAZI_HOME=
ZAZI_OPENAI_API_KEY=
ZAZI_OPENAI_BASE_URL=
ZAZI_ZAI_API_KEY=
ZAZI_OLLAMA_BASE_URL=
ZAZI_TELEGRAM_BOT_TOKEN=
ZAZI_DISCORD_BOT_TOKEN=
ZAZI_DISCORD_PUBLIC_KEY=
ZAZI_WHATSAPP_ACCESS_TOKEN=
ZAZI_WHATSAPP_PHONE_NUMBER_ID=
ZAZI_WHATSAPP_VERIFY_TOKEN=
```

Rules:

- Never commit real secrets.
- Never print secrets in logs.
- Prefer OS keychain or encrypted local secret store later.
