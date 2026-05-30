# Contributing to shipyard

Thanks for your interest! shipyard is a small Go CLI + an embedded skill.

## Layout

```
main.go                         entry point
internal/cli/                   cobra commands (ship, list, init, install-skill, where)
internal/config/                config loading + URL→repo inference
internal/assets/                go:embed of the skill + schema
  skill/SKILL.md                the pipeline (generalized, agent-agnostic)
  schema/_schema.yml            documented config template
examples/repos/                 illustrative configs (never real ones)
```

## Dev loop

```sh
make build      # -> ./shipyard
make test       # go test ./...
make vet
make fmt
```

## Principles

- **Configs are the user's, not ours.** shipyard ships only the schema and
  examples. Never add real repo configs to this repo.
- **The launcher stays thin.** Behavior lives in the skill; the binary resolves
  config, infers the repo, and execs the agent. Keep it that way.
- **Agent-agnostic.** Default agent is `claude`, overridable via
  `$SHIPYARD_AGENT` / `--agent`. Don't hardcode vendor specifics in the CLI.

## Commits

Conventional Commits (`feat|fix|refactor|docs|test|chore`). Small, reviewable.
