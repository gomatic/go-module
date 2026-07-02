// Package module parses a git remote URL into a Go module path and the name
// variants derived from it: the repository name, a Go identifier, and an
// environment-variable prefix.
//
// It is pure string logic with no I/O — a reusable leaf library for any caller
// that needs to reason about a project's identity from its remote.
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
	// Name is a project's short name — the last segment of the module path,
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
	raw = stripScheme(rawRemote(raw))
	raw = stripUserinfo(rawRemote(raw))
	raw = scpToPath(rawRemote(raw))
	if !valid(rawRemote(raw)) {
		// raw has already had any "user:token@" userinfo stripped, so embedding
		// it cannot leak credentials into the error text.
		return "", ErrInvalidRemote.With(nil, raw)
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

// rawRemote is a remote URL string partway through normalization into a module path.
type rawRemote string

// stripScheme removes a leading "scheme://" when present.
func stripScheme(raw rawRemote) string {
	if _, after, found := strings.Cut(string(raw), "://"); found {
		return after
	}
	return string(raw)
}

// stripUserinfo removes a leading "user@" when present.
func stripUserinfo(raw rawRemote) string {
	if _, after, found := strings.Cut(string(raw), "@"); found {
		return after
	}
	return string(raw)
}

// scpToPath converts the "host:org/repo" separator of the scp-like SSH form into
// "host/org/repo". A colon that appears before the first slash is the scp
// separator; a colon after a slash (or none) is left untouched.
func scpToPath(raw rawRemote) string {
	colon := strings.IndexByte(string(raw), ':')
	slash := strings.IndexByte(string(raw), '/')
	if colon >= 0 && (slash < 0 || colon < slash) {
		return string(raw)[:colon] + "/" + string(raw)[colon+1:]
	}
	return string(raw)
}

// valid reports whether raw is a host/org/repo path: at least three non-empty,
// space-free segments.
func valid(raw rawRemote) bool {
	if strings.ContainsAny(string(raw), " \t") {
		return false
	}
	segments := strings.Split(string(raw), "/")
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
