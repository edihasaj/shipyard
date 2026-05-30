// Package config loads per-repo shipyard configs and resolves the config home.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Home returns the directory holding repo configs, resolved in this order:
//   1. $SHIPYARD_HOME
//   2. ./.shipyard  (repo-local, if present)
//   3. $XDG_CONFIG_HOME/shipyard  (or ~/.config/shipyard)
func Home() string {
	if h := os.Getenv("SHIPYARD_HOME"); h != "" {
		return expand(h)
	}
	if st, err := os.Stat(".shipyard"); err == nil && st.IsDir() {
		abs, _ := filepath.Abs(".shipyard")
		return abs
	}
	base := os.Getenv("XDG_CONFIG_HOME")
	if base == "" {
		home, _ := os.UserHomeDir()
		base = filepath.Join(home, ".config")
	}
	return filepath.Join(base, "shipyard")
}

// ReposDir is where <repo>.yml config files live.
func ReposDir() string { return filepath.Join(Home(), "repos") }

// Config is one managed repo. Fields mirror repos/_schema.yml. Only what the
// launcher needs is typed strongly; the skill reads the rest of the YAML itself.
type Config struct {
	Key    string `yaml:"key"`
	Path   string `yaml:"path"`
	GitHub struct {
		Repo string `yaml:"repo"`
	} `yaml:"github"`
	Jira struct {
		ProjectKey string `yaml:"project_key"`
	} `yaml:"jira"`

	file string // source path, for diagnostics
}

// ResolvedPath returns Path with a leading ~ expanded.
func (c *Config) ResolvedPath() string { return expand(c.Path) }

// Load reads a single config by key.
func Load(key string) (*Config, error) {
	p := filepath.Join(ReposDir(), key+".yml")
	b, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("no config for %q (looked in %s)", key, ReposDir())
		}
		return nil, err
	}
	var c Config
	if err := yaml.Unmarshal(b, &c); err != nil {
		return nil, fmt.Errorf("parse %s: %w", p, err)
	}
	if c.Key == "" {
		c.Key = key
	}
	c.file = p
	return &c, nil
}

// List returns all config keys, sorted, excluding the _schema template.
func List() ([]string, error) {
	entries, err := os.ReadDir(ReposDir())
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var keys []string
	for _, e := range entries {
		n := e.Name()
		if e.IsDir() || !strings.HasSuffix(n, ".yml") {
			continue
		}
		key := strings.TrimSuffix(n, ".yml")
		if strings.HasPrefix(key, "_") {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys, nil
}

// LoadAll loads every config (used for URL inference).
func LoadAll() ([]*Config, error) {
	keys, err := List()
	if err != nil {
		return nil, err
	}
	var out []*Config
	for _, k := range keys {
		c, err := Load(k)
		if err != nil {
			continue // skip unparseable configs rather than failing the whole list
		}
		out = append(out, c)
	}
	return out, nil
}

func expand(p string) string {
	if p == "~" || strings.HasPrefix(p, "~/") {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, strings.TrimPrefix(p, "~"))
	}
	return p
}
