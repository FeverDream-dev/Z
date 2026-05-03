# Release Checklist

## Goal
Ensure every release is tested, documented, and explicitly approved.

## Pre-release

- [ ] All tests pass (`go test ./...`)
- [ ] Build succeeds (`go build ./...`)
- [ ] No secrets in source or logs
- [ ] Version bumped in `cmd/zazi/main.go`
- [ ] CHANGELOG.md updated
- [ ] UI validation passed (`zazi validate ui`)
- [ ] Skill analyzer tests pass
- [ ] Agent isolation verified
- [ ] Documentation current

## Prepare

Run:
```bash
zazi release prepare --version=v0.x.x
```

This will:
1. Run the full test suite
2. Generate a changelog draft
3. Show the release checklist
4. Output a summary report

**Does NOT publish anything.**

## Approve

A human must explicitly approve by running:
```bash
zazi release publish --version=v0.x.x --approve="I have reviewed the tests and release notes"
```

The approval string is logged for audit purposes.

## Publish (Manual)

After approval, the human may:
1. Tag the release: `git tag v0.x.x`
2. Push the tag: `git push origin v0.x.x`
3. Create GitHub release with notes
4. Attach binaries if applicable

## Safety

- `prepare` never pushes, tags, or publishes
- `publish` requires explicit approval text
- Release notes must mention test results and known risks
- No automated publishing without human gate
