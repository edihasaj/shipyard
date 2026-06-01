# Security Policy

## Reporting a vulnerability

Please report security issues privately via GitHub's
[private vulnerability reporting](https://github.com/edihasaj/shipyard/security/advisories/new)
or by email to **edihasaj@gmail.com**. Do not open a public issue for
suspected vulnerabilities.

You can expect an initial response within a few days.

## Scope notes

shipyard is a thin launcher: it resolves a per-repo config and then **execs an
agent CLI** (default `claude`) which runs the `ship-task` skill. Keep in mind:

- **Configs are user-owned and may reference private repos, tokens, and
  internal task systems.** shipyard never commits real configs (`.shipyard/`
  is gitignored); only the schema and examples ship in this repo.
- The pipeline can push branches and open PRs only when a repo config opts in
  (`push: pr`). The default safety rail is `push: manual` / `push: ask`.
- shipyard shells out to whatever binary `$SHIPYARD_AGENT` / `--agent` names.
  Only point it at an agent you trust.
