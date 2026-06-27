package module

// Imported bare (the package is named error); this file declares only sentinels
// and uses no builtin error type, so the declaration reads errs.Const.
import errs "github.com/gomatic/go-error"

// ErrInvalidRemote indicates a git remote that cannot be parsed into a module path.
const ErrInvalidRemote errs.Const = "invalid remote"
