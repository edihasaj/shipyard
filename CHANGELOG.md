# Changelog

All notable changes to this project are documented here. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project
adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Per-agent invocation profiles (`internal/agent`): `claude` (default, slash
  command), `codex`, and `generic`. Non-Claude agents get the pipeline inlined
  in the prompt, so `$SHIPYARD_AGENT` works without a per-agent skill install.
  Select with `--agent-profile` / `$SHIPYARD_AGENT_PROFILE` (inferred from the
  agent binary name otherwise).
- Community health files: issue templates, PR template, `SECURITY.md`,
  `CHANGELOG.md`, README badges.

## [0.1.0] - TBD

First public release.

### Added
- `shipyard <repo> <task>` — run the per-task ship pipeline via an agent CLI.
- `shipyard list` / `init` / `install-skill` / `where` commands.
- Per-repo YAML config (task source, branch/commit convention, gates, review,
  push policy) with `$SHIPYARD_HOME` → `./.shipyard` → `$XDG_CONFIG_HOME`
  resolution.
- URL → repo inference for GitHub and Jira task references.
- Embedded `ship-task` skill and documented config schema.
- CI (vet/build/test -race/smoke) and a goreleaser release pipeline
  (darwin/linux, amd64/arm64) with a Homebrew tap formula.

[Unreleased]: https://github.com/edihasaj/shipyard/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/edihasaj/shipyard/releases/tag/v0.1.0
