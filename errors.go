package module

// Imported bare (the package is named error); this file declares only sentinels
// and uses no builtin error type, so the declaration reads error.Const.
import "github.com/gomatic/go-error"

// ErrInvalidRemote indicates a git remote that cannot be parsed into a module path.
const ErrInvalidRemote error.Const = "invalid remote"
