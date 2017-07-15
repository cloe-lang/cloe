package ast

import (
	"fmt"
	"strings"

	"github.com/tisp-lang/tisp/src/lib/debug"
)

// MutualRecursion represents a definition of mutually-recursive functions.
type MutualRecursion struct {
	letFunctions []LetFunction
	info         debug.Info
}

// NewMutualRecursion creates a mutual recursion node from mutually-recursive
// functions.
func NewMutualRecursion(fs []LetFunction, i debug.Info) MutualRecursion {
	return MutualRecursion{fs, i}
}

// LetFunctions returns let-function statements in a mutual recursion
// definition.
func (mr MutualRecursion) LetFunctions() []LetFunction {
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
