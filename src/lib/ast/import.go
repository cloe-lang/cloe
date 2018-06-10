package ast

import (
	"fmt"

	"github.com/cloe-lang/cloe/src/lib/debug"
)

// Import represents an import of a sub module.
type Import struct {
	path   string
	prefix string
	info   *debug.Info
}

// NewImport creates an Import.
func NewImport(path, prefix string, info *debug.Info) Import {
	return Import{path, prefix, info}
}

// Path returns a path to an imported sub module.
func (i Import) Path() string {
	return i.path
}

// Prefix returns a prefix apppended to members' names in the imported module.
func (i Import) Prefix() string {
	return i.prefix
}

func (i Import) String() string {
	return fmt.Sprintf("(import %v)", i.path)
}
