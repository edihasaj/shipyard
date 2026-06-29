// Package agent maps an agent CLI to how shipyard drives it: how the task
// prompt is rendered and how argv is built for interactive vs print mode.
//
// This is what makes shipyard agent-agnostic. Claude Code already has the
// ship-task skill installed, so it gets a "/ship-task <args>" slash command.
// Other agents don't have the skill, so their profile inlines the embedded
// pipeline (SKILL.md) ahead of the inputs — they read the whole thing in the
// prompt. Add an agent = add a profile, not new launch code.
package agent

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/edihasaj/shipyard/internal/assets"
)

// envOr returns $key, or def when unset/empty. Lets headless defaults below be
// overridden without changing the Argv signature.
func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// Profile describes how to drive one agent CLI.
type Profile struct {
	// Name is the canonical profile name (also the matched binary base name).
	Name string
	// SlashCommand is true when the agent already has the ship-task skill and
	// understands "/ship-task <args>"; when false the pipeline is inlined.
	SlashCommand bool

	interactiveArgs func(prompt string) []string
	printArgs       func(prompt string) []string
}

func positional(prompt string) []string { return []string{prompt} }

var profiles = map[string]Profile{
	// Claude Code: slash command + headless flags. acceptEdits still denies
	// bash, so a headless run can't branch/commit; pre-allow git (override via
	// SHIPYARD_HEADLESS_ALLOWED, e.g. to add gate commands).
	"claude": {
		Name:            "claude",
		SlashCommand:    true,
		interactiveArgs: positional,
		printArgs: func(p string) []string {
			allowed := envOr("SHIPYARD_HEADLESS_ALLOWED", "Bash(git:*)")
			return []string{"-p", p, "--permission-mode", "acceptEdits", "--allowedTools", allowed}
		},
	},
	// OpenAI Codex CLI: `codex "<prompt>"` interactive, `codex exec` headless.
	// `codex exec` sandboxes to read-only by default, so let it write the repo
	// (override via SHIPYARD_CODEX_SANDBOX: read-only|workspace-write|
	// danger-full-access). A `--` guards prompts that start with "---" (the
	// inlined skill's YAML front matter) from codex's arg parser.
	"codex": {
		Name:            "codex",
		SlashCommand:    false,
		interactiveArgs: func(p string) []string { return []string{"--", p} },
		printArgs: func(p string) []string {
			sandbox := envOr("SHIPYARD_CODEX_SANDBOX", "workspace-write")
			return []string{"exec", "--sandbox", sandbox, "--", p}
		},
	},
	// Fallback: pass the prompt as a single positional argument. Works for any
	// agent CLI whose first positional arg is the instruction.
	"generic": {
		Name:            "generic",
		SlashCommand:    false,
		interactiveArgs: positional,
		printArgs:       positional,
	},
}

// Names returns the built-in profile names, sorted, for help/diagnostics.
func Names() []string {
	out := make([]string, 0, len(profiles))
	for k := range profiles {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// ProfileFor returns the profile for an agent binary path, matched on its base
// name (case-insensitive, .exe stripped), falling back to the generic profile.
func ProfileFor(bin string) Profile {
	name := strings.TrimSuffix(strings.ToLower(filepath.Base(bin)), ".exe")
	if p, ok := profiles[name]; ok {
		return p
	}
	return profiles["generic"]
}

// ProfileByName returns a profile by explicit name, or the generic profile if
// the name is unknown. The returned bool reports whether the name was known.
func ProfileByName(name string) (Profile, bool) {
	if p, ok := profiles[strings.ToLower(name)]; ok {
		return p, true
	}
	return profiles["generic"], false
}

// Prompt renders the instruction sent to the agent. Skill-aware agents get a
// slash command; others get the full pipeline inlined ahead of the inputs.
func (p Profile) Prompt(args string) string {
	if p.SlashCommand {
		return strings.TrimSpace("/ship-task " + args)
	}
	return string(assets.SkillMD) + "\n\n---\nInputs (`<repo> <task-ref> [notes…]`): " + args + "\n"
}

// Argv builds the agent's command-line arguments for the rendered prompt.
func (p Profile) Argv(prompt string, print bool) []string {
	if print {
		return p.printArgs(prompt)
	}
	return p.interactiveArgs(prompt)
}
