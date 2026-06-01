package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/edihasaj/shipyard/internal/agent"
	"github.com/edihasaj/shipyard/internal/config"
	"github.com/spf13/cobra"
)

var (
	flagHeadless     bool
	flagPrint        bool
	flagAgent        string
	flagAgentProfile string
)

func runShip(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return cmd.Help()
	}

	all, err := config.LoadAll()
	if err != nil {
		return err
	}

	// A pasted URL as the first token may resolve the repo by itself.
	var repoKey, task string
	first := args[0]
	if config.IsURL(first) {
		if key, ok := config.InferFromURL(first, all); ok {
			repoKey = key
			task = first
			args = args[1:]
		} else {
			return fmt.Errorf("could not infer a configured repo from URL: %s\nuse: shipyard <repo> %s", first, first)
		}
	} else {
		repoKey = first
		args = args[1:]
		if len(args) > 0 {
			task = args[0]
			args = args[1:]
		}
	}
	notes := strings.Join(args, " ")

	cfg, err := config.Load(repoKey)
	if err != nil {
		keys, _ := config.List()
		return fmt.Errorf("%w\nconfigured repos: %s", err, strings.Join(keys, ", "))
	}

	repoPath := cfg.ResolvedPath()
	if st, err := os.Stat(repoPath); err != nil || !st.IsDir() {
		return fmt.Errorf("repo path missing or not a dir: %s", repoPath)
	}

	agentBin := flagAgent
	if agentBin == "" {
		agentBin = envOr("SHIPYARD_AGENT", "claude")
	}

	// Pick the invocation profile: explicit flag/env wins, else infer from the
	// agent binary name (falling back to a generic positional-prompt profile).
	prof := agent.ProfileFor(agentBin)
	if name := flagAgentProfile; name == "" {
		name = os.Getenv("SHIPYARD_AGENT_PROFILE")
		if name != "" {
			flagAgentProfile = name
		}
	}
	if flagAgentProfile != "" {
		p, ok := agent.ProfileByName(flagAgentProfile)
		if !ok {
			fmt.Fprintf(os.Stderr, "shipyard: unknown agent profile %q; using %q (known: %s)\n",
				flagAgentProfile, p.Name, strings.Join(agent.Names(), ", "))
		}
		prof = p
	}

	taskArgs := strings.Join(strings.Fields(repoKey+" "+task+" "+notes), " ")
	prompt := prof.Prompt(taskArgs)
	print := flagHeadless || flagPrint

	fmt.Fprintf(os.Stderr, "▶ shipyard: %s  task=%s  agent=%s/%s  (%s)\n",
		repoKey, orNone(task), filepath.Base(agentBin), prof.Name, repoPath)

	cmdArgs := prof.Argv(prompt, print)

	c := exec.Command(agentBin, cmdArgs...)
	c.Dir = repoPath
	c.Stdin, c.Stdout, c.Stderr = os.Stdin, os.Stdout, os.Stderr
	return c.Run()
}

func init() {
	// flags are registered on the root via PersistentFlags in Execute's tree;
	// kept here next to the command they serve.
}

// bindShipFlags is called from root to attach run flags.
func bindShipFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&flagPrint, "print", "p", false, "headless/print mode (non-interactive)")
	cmd.Flags().BoolVar(&flagHeadless, "headless", false, "alias for --print")
	cmd.Flags().StringVar(&flagAgent, "agent", "", "agent binary to invoke (default: $SHIPYARD_AGENT or 'claude')")
	cmd.Flags().StringVar(&flagAgentProfile, "agent-profile", "", "invocation profile: claude|codex|generic (default: inferred from --agent; or $SHIPYARD_AGENT_PROFILE)")
}

func envOr(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func orNone(s string) string {
	if s == "" {
		return "<none>"
	}
	return s
}
