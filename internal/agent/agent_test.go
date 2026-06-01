package agent

import "testing"

func TestProfileForInfersFromBinaryName(t *testing.T) {
	cases := map[string]string{
		"claude":          "claude",
		"/usr/bin/claude": "claude",
		"Claude.exe":      "claude",
		"codex":           "codex",
		"my-custom-agent": "generic",
		"":                "generic",
		"/opt/llm/aider":  "generic",
	}
	for in, want := range cases {
		if got := ProfileFor(in).Name; got != want {
			t.Errorf("ProfileFor(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestProfileByNameUnknownFallsBack(t *testing.T) {
	if p, ok := ProfileByName("nope"); ok || p.Name != "generic" {
		t.Errorf("ProfileByName(nope) = (%q,%v), want (generic,false)", p.Name, ok)
	}
	if p, ok := ProfileByName("CODEX"); !ok || p.Name != "codex" {
		t.Errorf("ProfileByName(CODEX) = (%q,%v), want (codex,true)", p.Name, ok)
	}
}

func TestPromptSlashVsInlined(t *testing.T) {
	claude, _ := ProfileByName("claude")
	if got := claude.Prompt("repo TASK-1"); got != "/ship-task repo TASK-1" {
		t.Errorf("claude prompt = %q", got)
	}

	generic, _ := ProfileByName("generic")
	p := generic.Prompt("repo TASK-1")
	if len(p) < 200 { // SKILL.md is inlined, so the prompt is large
		t.Errorf("generic prompt too short to contain inlined skill: %d bytes", len(p))
	}
}

func TestArgvPrintModes(t *testing.T) {
	claude, _ := ProfileByName("claude")
	if got := claude.Argv("X", true); len(got) != 4 || got[0] != "-p" || got[3] != "acceptEdits" {
		t.Errorf("claude print argv = %v", got)
	}
	if got := claude.Argv("X", false); len(got) != 1 || got[0] != "X" {
		t.Errorf("claude interactive argv = %v", got)
	}
	codex, _ := ProfileByName("codex")
	if got := codex.Argv("X", true); len(got) != 2 || got[0] != "exec" {
		t.Errorf("codex print argv = %v", got)
	}
}
