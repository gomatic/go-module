package module

// Aliased to errs (its package is named errs) so the sentinel declaration below
// reads errs.Const without shadowing the builtin error type.
import errs "github.com/gomatic/go-error"

// ErrInvalidRemote indicates a git remote that cannot be parsed into a module path.
const ErrInvalidRemote errs.Const = "invalid remote"
