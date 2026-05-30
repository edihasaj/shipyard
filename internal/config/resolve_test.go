package config

import "testing"

func TestInferFromURL(t *testing.T) {
	cfgs := []*Config{
		{Key: "panel"},
		{Key: "api"},
		{Key: "ent"},
	}
	cfgs[0].GitHub.Repo = "acme/panel"
	cfgs[1].GitHub.Repo = "acme/api"
	cfgs[2].Jira.ProjectKey = "DEVO"

	cases := []struct {
		url, want string
		ok        bool
	}{
		{"https://github.com/acme/panel/issues/86", "panel", true},
		{"https://github.com/acme/api/pull/12", "api", true},
		{"https://github.com/acme/api.git", "api", true},
		{"https://acme.atlassian.net/browse/DEVO-25497", "ent", true},
		{"https://acme.atlassian.net/jira/x?selectedIssue=DEVO-1", "ent", true},
		{"https://github.com/other/thing/issues/1", "", false},
		{"not a url", "", false},
	}
	for _, c := range cases {
		got, ok := InferFromURL(c.url, cfgs)
		if got != c.want || ok != c.ok {
			t.Errorf("InferFromURL(%q) = (%q,%v), want (%q,%v)", c.url, got, ok, c.want, c.ok)
		}
	}
}

func TestIsURL(t *testing.T) {
	if !IsURL("https://x.com") || IsURL("DEVO-1") {
		t.Fatal("IsURL mismatch")
	}
}
