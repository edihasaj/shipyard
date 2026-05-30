// Package assets embeds the bundled skill and config schema so the binary is
// fully self-contained (no install-time file fetching).
package assets

import _ "embed"

//go:embed skill/SKILL.md
var SkillMD []byte

//go:embed schema/_schema.yml
var SchemaYML []byte
