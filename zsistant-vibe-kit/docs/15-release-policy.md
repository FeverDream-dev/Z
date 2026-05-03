# 15 — Release Policy

Zsistant can prepare releases but must not publish them without human approval.

## Release stages

1. Run full tests.
2. Run security checks.
3. Generate changelog draft.
4. Generate human-readable diff summary.
5. Ask the user:

```text
We completed X, Y, and Z. Tests passed. Should I push this to GitHub and create a release?
```

6. Only after approval: commit, tag, push, or publish.

## Approval examples

Valid approval:

```text
Yes, push and release v0.1.0.
```

Not enough:

```text
Looks good.
```

## Release artifacts later

- Linux amd64/arm64 binary.
- macOS amd64/arm64 binary.
- Windows binary.
- checksums.
- install script.
- Docker image, optional.

## Public/commercial strategy

Recommended license direction:

- Open-source core under Apache-2.0 or MIT.
- Optional hosted cloud/LLM backend later.
- Optional commercial support and managed cloud.
- Keep local mode useful and not crippled.

## Public tracking

The GitHub repo should clearly link to:

- company website
- maintainer profile
- roadmap
- contribution guide
- security policy
- sponsorship/commercial contact
