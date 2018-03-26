package ast

import (
	"strings"
)

// Arguments represents arguments passed to a function.
type Arguments struct {
	positionals []PositionalArgument
	keywords    []KeywordArgument
}

// NewArguments creates arguments.
func NewArguments(ps []PositionalArgument, ks []KeywordArgument) Arguments {
	return Arguments{ps, ks}
}

// Positionals returns positional arguments contained in arguments.
func (a Arguments) Positionals() []PositionalArgument {
	return a.positionals
}

// Keywords returns keyword arguments contained in arguments.
func (a Arguments) Keywords() []KeywordArgument {
	return a.keywords
}

func (a Arguments) String() string {
	ss := make([]string, 0, len(a.positionals)+len(a.keywords)+1)

	for _, p := range a.positionals {
		ss = append(ss, p.String())
	}

	if len(a.keywords) > 0 {
		ss = append(ss, ".")
	}

	for _, k := range a.keywords {
		ss = append(ss, k.String())
	}

	return strings.Join(ss, " ")
}
