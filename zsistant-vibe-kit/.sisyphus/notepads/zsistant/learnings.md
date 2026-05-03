Phase 00: Repo Foundation for Zsistant (zsistant-vibe-kit)

What I did:
- Initialized Go module: github.com/zsistant/zsistant
- Created repository skeleton under zsistant-vibe-kit with empty, compilable placeholder packages:
  - cmd/zazi (main placeholder)
  - internal/config, internal/daemon, internal/agents, internal/store, internal/jobs, internal/llm, internal/channels, internal/bus, internal/trainer, internal/skills, internal/sandbox, internal/server
  - web (placeholder)
  - tests (placeholder)
- Added a basic CI workflow at .github/workflows/ci.yml to run go test ./... and go build ./...
- Updated README.md with a Development section describing how to build/test locally

Verification performed:
- go mod tidy
- go build ./...
- go test ./...
All commands completed with success (no tests present, but packages compiled).

Open questions / next steps:
- Add real code per subsequent phases; expand tests; implement actual functionality.
