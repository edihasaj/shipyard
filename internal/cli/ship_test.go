package cli

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// buildBinary compiles shipyard once into a temp dir and returns its path.
func buildBinary(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	bin := filepath.Join(dir, "shipyard")
	// build from the module root (two levels up from internal/cli)
	cmd := exec.Command("go", "build", "-o", bin, "../..")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}
	return bin
}

// fakeAgent writes a small script that records argv + cwd to logPath.
func fakeAgent(t *testing.T, logPath string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "fakeagent")
	script := "#!/usr/bin/env bash\n" +
		"echo \"CWD=$(pwd)\" >> " + logPath + "\n" +
		"echo \"ARGS=$*\" >> " + logPath + "\n"
	if err := os.WriteFile(p, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	return p
}

// setupHome creates a config home with one repo config and returns (home, repoPath).
func setupHome(t *testing.T) (string, string) {
	t.Helper()
	home := t.TempDir()
	repo := t.TempDir()
	reposDir := filepath.Join(home, "repos")
	if err := os.MkdirAll(reposDir, 0o755); err != nil {
		t.Fatal(err)
	}
	cfg := "key: testrepo\npath: " + repo + "\n" +
		"task_source: github\ngithub:\n  repo: acme/testrepo\n" +
		"base_branch: main\nbranch_format: \"{type}/{slug}\"\n" +
		"push: ask\n"
	if err := os.WriteFile(filepath.Join(reposDir, "testrepo.yml"), []byte(cfg), 0o644); err != nil {
		t.Fatal(err)
	}
	return home, repo
}

func run(t *testing.T, bin, home, agent string, args ...string) (string, error) {
	t.Helper()
	cmd := exec.Command(bin, args...)
	cmd.Env = append(os.Environ(),
		"SHIPYARD_HOME="+home,
		"SHIPYARD_AGENT="+agent,
	)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func TestPlainInvoke(t *testing.T) {
	bin := buildBinary(t)
	home, repo := setupHome(t)
	log := filepath.Join(t.TempDir(), "agent.log")
	agent := fakeAgent(t, log)

	if _, err := run(t, bin, home, agent, "testrepo", "add csv export"); err != nil {
		t.Fatalf("run failed: %v", err)
	}
	b, _ := os.ReadFile(log)
	s := string(b)
	if !strings.Contains(s, "CWD="+repo) {
		t.Errorf("agent not run in repo cwd:\n%s", s)
	}
	if !strings.Contains(s, "/ship-task testrepo add csv export") {
		t.Errorf("prompt malformed:\n%s", s)
	}
}

func TestURLInference(t *testing.T) {
	bin := buildBinary(t)
	home, _ := setupHome(t)
	log := filepath.Join(t.TempDir(), "agent.log")
	agent := fakeAgent(t, log)

	if _, err := run(t, bin, home, agent, "https://github.com/acme/testrepo/issues/42"); err != nil {
		t.Fatalf("run failed: %v", err)
	}
	b, _ := os.ReadFile(log)
	if !strings.Contains(string(b), "/ship-task testrepo") {
		t.Errorf("URL did not resolve to testrepo:\n%s", b)
	}
}

func TestHeadlessFlag(t *testing.T) {
	bin := buildBinary(t)
	home, _ := setupHome(t)
	log := filepath.Join(t.TempDir(), "agent.log")
	agent := fakeAgent(t, log)

	if _, err := run(t, bin, home, agent, "-p", "testrepo", "fix bug"); err != nil {
		t.Fatalf("run failed: %v", err)
	}
	b, _ := os.ReadFile(log)
	s := string(b)
	if !strings.Contains(s, "-p") || !strings.Contains(s, "acceptEdits") {
		t.Errorf("headless flags not forwarded:\n%s", s)
	}
}

func TestUnknownRepoFails(t *testing.T) {
	bin := buildBinary(t)
	home, _ := setupHome(t)
	agent := fakeAgent(t, filepath.Join(t.TempDir(), "a.log"))

	out, err := run(t, bin, home, agent, "nonexistent", "x")
	if err == nil {
		t.Error("expected nonzero exit for unknown repo")
	}
	if !strings.Contains(out, "testrepo") {
		t.Errorf("error should list configured repos:\n%s", out)
	}
}

func TestUnresolvableURLFails(t *testing.T) {
	bin := buildBinary(t)
	home, _ := setupHome(t)
	agent := fakeAgent(t, filepath.Join(t.TempDir(), "a.log"))

	if _, err := run(t, bin, home, agent, "https://github.com/who/unknown/issues/1"); err == nil {
		t.Error("expected nonzero exit for unresolvable URL")
	}
}

func TestListAndWhere(t *testing.T) {
	bin := buildBinary(t)
	home, _ := setupHome(t)
	agent := fakeAgent(t, filepath.Join(t.TempDir(), "a.log"))

	out, err := run(t, bin, home, agent, "list")
	if err != nil || !strings.Contains(out, "testrepo") {
		t.Errorf("list failed: %v\n%s", err, out)
	}
	out, err = run(t, bin, home, agent, "where")
	if err != nil || !strings.Contains(out, home) {
		t.Errorf("where failed: %v\n%s", err, out)
	}
}
