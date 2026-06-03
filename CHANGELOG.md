# Changelog

All notable changes to this project are documented here. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project
adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.2.0] - 2026-06-03

### Added
- Custom ticket templates via `ticket.template`: the issue-body analogue of
  `pr.template`. Accepts a filesystem path or an Obsidian URL (same resolution
  as `pr.template`), or `""` for the built-in default skeleton (Status / Ask /
  Description / Out of scope). Used when the pipeline drafts or updates a ticket
  description.
- Ticket pipeline step routes **acceptance criteria into the tracker's
  dedicated fields** (e.g. Jira "End User Acceptance Criteria" / "PR Acceptance
  Criteria") instead of the description body — discovered via issue metadata,
  written as ADF in a separate update from the markdown description — and keeps
  section headings small (`###`). Falls back to description sections when no
  such field exists.

## [0.1.1] - 2026-06-01

### Added
- Custom PR templates via `pr.template`: accepts a filesystem path or an
  Obsidian URL (`obsidian://open?vault=<v>&file=<note>`, resolved with the
  `obsidian` CLI). When set, the pipeline fills the template's prose
  placeholders and tokens while preserving headers and review checklists
  verbatim; unset keeps the built-in default skeleton.
- Per-agent invocation profiles (`internal/agent`): `claude` (default, slash
  command), `codex`, and `generic`. Non-Claude agents get the pipeline inlined
  in the prompt, so `$SHIPYARD_AGENT` works without a per-agent skill install.
  Select with `--agent-profile` / `$SHIPYARD_AGENT_PROFILE` (inferred from the
  agent binary name otherwise).
- Community health files: issue templates, PR template, `SECURITY.md`,
  `CHANGELOG.md`, README badges.

## [0.1.0] - 2026-06-01

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

[Unreleased]: https://github.com/edihasaj/shipyard/compare/v0.2.0...HEAD
[0.2.0]: https://github.com/edihasaj/shipyard/compare/v0.1.1...v0.2.0
[0.1.1]: https://github.com/edihasaj/shipyard/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/edihasaj/shipyard/releases/tag/v0.1.0
