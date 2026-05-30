package config

import (
	"regexp"
	"strings"
)

var (
	reGitHub = regexp.MustCompile(`github\.com[/:]([^/]+/[^/.]+)`)
	reJira   = regexp.MustCompile(`([A-Z][A-Z0-9]+-\d+)`)
)

// InferFromURL maps a pasted GitHub/Jira URL to a configured repo key.
// Returns ("", false) when no config matches.
//   - GitHub URL: match owner/repo against each config's github.repo.
//   - Jira URL/key: match the project-key prefix against jira.project_key.
func InferFromURL(url string, configs []*Config) (string, bool) {
	if m := reGitHub.FindStringSubmatch(url); m != nil {
		ownerRepo := strings.TrimSuffix(m[1], ".git")
		for _, c := range configs {
			if strings.EqualFold(c.GitHub.Repo, ownerRepo) {
				return c.Key, true
			}
		}
	}
	if m := reJira.FindStringSubmatch(url); m != nil {
		proj := strings.SplitN(m[1], "-", 2)[0]
		for _, c := range configs {
			if c.Jira.ProjectKey != "" && strings.EqualFold(c.Jira.ProjectKey, proj) {
				return c.Key, true
			}
		}
	}
	return "", false
}

// IsURL reports whether s looks like an http(s) URL.
func IsURL(s string) bool {
	return strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://")
}
