package ast

import (
	"fmt"
	"strings"

	"github.com/coel-lang/coel/src/lib/debug"
)

// MutualRecursion represents a definition of mutually-recursive functions.
type MutualRecursion struct {
	letFunctions []DefFunction
	info         debug.Info
}

// NewMutualRecursion creates a mutual recursion node from mutually-recursive
// functions.
func NewMutualRecursion(fs []DefFunction, i debug.Info) MutualRecursion {
	if len(fs) < 2 {
		panic("A number of mutually recursive functions must be more than 2.")
	}

	return MutualRecursion{fs, i}
}

// DefFunctions returns let-function statements in a mutual recursion
// definition.
func (mr MutualRecursion) DefFunctions() []DefFunction {
	return mr.letFunctions
}

// DebugInfo returns debug information of mutually-recursive function definition.
func (mr MutualRecursion) DebugInfo() debug.Info {
	return mr.info
}

func (mr MutualRecursion) String() string {
	ss := make([]string, 0, len(mr.letFunctions))

	for _, l := range mr.letFunctions {
		ss = append(ss, l.String())
	}

	return fmt.Sprintf("(mr %v)", strings.Join(ss, " "))
}
