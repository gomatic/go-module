// Package module parses repository remotes into Go module paths and derives the
// name variants a project uses: the raw name, a Go identifier, and an
// environment-variable prefix.
//
// It is pure string logic with no I/O, reusable by any domain that needs to
// reason about a project's identity. The rename domain uses it to turn the git
// origin remote into the module path that Go requires and to compute the token
// variants the rewrite engine substitutes.
package module

import (
	"slices"
	"strings"
)

type (
	// Remote is a git remote URL, e.g. "git@github.com:gomatic/template.cli.git".
	Remote string
	// Path is a Go module path, e.g. "github.com/gomatic/template.cli".
	Path string
	// Name is a project's short name — the binary and the cmd/<name> directory,
	// e.g. "template.cli".
	Name string
	// Identifier is a Name reduced to a valid Go identifier, e.g. "templatecli".
	Identifier string
	// EnvPrefix is a Name reduced to an environment-variable prefix, e.g.
	// "TEMPLATE_CLI".
	EnvPrefix string
)

// Parse turns a git remote URL into a Go module path, accepting both the scp-like
// SSH form (git@host:org/repo.git) and URL forms (https://host/org/repo.git,
// ssh://git@host/org/repo). It strips the scheme, any userinfo, and a trailing
// ".git", and reports ErrInvalidRemote when the result is not a host/org/repo
// path.
func Parse(remote Remote) (Path, error) {
	raw := strings.TrimSpace(string(remote))
	raw = strings.TrimSuffix(raw, ".git")
	raw = stripScheme(raw)
	raw = stripUserinfo(raw)
	raw = scpToPath(raw)
	if !valid(raw) {
		return "", ErrInvalidRemote.With(nil, string(remote))
	}
	return Path(raw), nil
}

// Repo returns the last segment of the module path — the repository name.
func (p Path) Repo() Name {
	s := string(p)
	if i := strings.LastIndexByte(s, '/'); i >= 0 {
		return Name(s[i+1:])
	}
	return Name(s)
}

// Identifier reduces the name to a valid lowercase Go identifier by dropping
// every character that is not a lowercase letter or digit.
func (n Name) Identifier() Identifier {
	return Identifier(strings.Map(keepIdentifier, strings.ToLower(string(n))))
}

// EnvPrefix reduces the name to an environment-variable prefix by uppercasing it
// and mapping every non-alphanumeric character to an underscore.
func (n Name) EnvPrefix() EnvPrefix {
	return EnvPrefix(strings.Map(toEnv, strings.ToUpper(string(n))))
}

// stripScheme removes a leading "scheme://" when present.
func stripScheme(raw string) string {
	if _, after, found := strings.Cut(raw, "://"); found {
		return after
	}
	return raw
}

// stripUserinfo removes a leading "user@" when present.
func stripUserinfo(raw string) string {
	if _, after, found := strings.Cut(raw, "@"); found {
		return after
	}
	return raw
}

// scpToPath converts the "host:org/repo" separator of the scp-like SSH form into
// "host/org/repo". A colon that appears before the first slash is the scp
// separator; a colon after a slash (or none) is left untouched.
func scpToPath(raw string) string {
	colon := strings.IndexByte(raw, ':')
	slash := strings.IndexByte(raw, '/')
	if colon >= 0 && (slash < 0 || colon < slash) {
		return raw[:colon] + "/" + raw[colon+1:]
	}
	return raw
}

// valid reports whether raw is a host/org/repo path: at least three non-empty,
// space-free segments.
func valid(raw string) bool {
	if strings.ContainsAny(raw, " \t") {
		return false
	}
	segments := strings.Split(raw, "/")
	return len(segments) >= 3 && !slices.Contains(segments, "")
}

// keepIdentifier keeps lowercase letters and digits, dropping everything else.
func keepIdentifier(r rune) rune {
	if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
		return r
	}
	return -1
}

// toEnv keeps uppercase letters and digits and maps everything else to '_'.
func toEnv(r rune) rune {
	if (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
		return r
	}
	return '_'
}
