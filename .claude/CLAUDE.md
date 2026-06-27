# go-module

Derives Go module identity from a git remote (package `module`): `Parse(remote) → Path`, `Path.Repo() → Name`, and `Name.Identifier()`/`Name.EnvPrefix()`. Generic Go-tooling logic for scaffolding, release, and rename tools. Lives in `gomatic`.

- Owns its one sentinel, `ErrInvalidRemote`, on `gomatic/go-error` (`error.Const`).
- Gate: gofumpt, vet, staticcheck, govulncheck, gocognit ≤ 7, 100% coverage.
