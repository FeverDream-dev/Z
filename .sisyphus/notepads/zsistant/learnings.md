Phase 01: Implemented zazi CLI skeleton with version/doctor/init/serve placeholders.
-CLI uses home-dir ~/.zazi to store config and data; init creates dirs and a sample config.yaml.
-Tests added for version, doctor, help, and init behavior; tests pass.
-Moved tests to correct module path zsistant-vibe-kit/cmd/zazi to align with go.mod and module layout.
-Verification: go test ./... on the module passes for the cmd/zazi package.

Phase 02: Config and Local State for Zsistant
- Updated tests in cmd/zazi/main_test.go to align with Phase 02 output changes:
  - Doctor output now shows "Config file:" and "Data path:" with status OK or MISSING
  - Init test now checks for new header "Zsistant configuration file" in config.yaml
  - All changes kept isolated to tests; no config or CLI logic modified
- Added internal/config with config.go, defaults.go, and loader.go implementing a simple YAML-based config system with defaults and a redacted String() for safe logging.
- Updated CLI (cmd/zazi) to use new config system: init (idempotent) and doctor (verification) commands.
- Added unit tests under internal/config/config_test.go to verify defaults, round-trip save/load, dir creation, and redaction behavior.
- Created go.mod and wired gopkg.in/yaml.v3 as a dependency.
- Verified build and tests via go build ./... and go test ./...; all tests pass.
