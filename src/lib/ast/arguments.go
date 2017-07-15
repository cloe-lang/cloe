package ast

import (
	"fmt"
	"strings"
)

// Arguments represents arguments passed to a function.
type Arguments struct {
	positionals   []PositionalArgument
	keywords      []KeywordArgument
	expandedDicts []interface{}
}

// NewArguments creates arguments.
func NewArguments(
	ps []PositionalArgument,
	ks []KeywordArgument,
	expandedDicts []interface{}) Arguments {
	return Arguments{ps, ks, expandedDicts}
}

// Positionals returns positional arguments contained in arguments.
func (a Arguments) Positionals() []PositionalArgument {
	return a.positionals
}

// Keywords returns keyword arguments contained in arguments.
func (a Arguments) Keywords() []KeywordArgument {
	return a.keywords
}

// ExpandedDicts returns expanded dictionary arguments contained in arguments.
func (a Arguments) ExpandedDicts() []interface{} {
	return a.expandedDicts
}

func (a Arguments) String() string {
	ss := make([]string, 0, len(a.positionals)+len(a.keywords)+len(a.expandedDicts)+1)

	for _, p := range a.positionals {
		ss = append(ss, p.String())
	}

	if len(a.keywords)+len(a.expandedDicts) > 0 {
		ss = append(ss, ".")
	}

	for _, k := range a.keywords {
		ss = append(ss, k.String())
	}

	for _, d := range a.expandedDicts {
		ss = append(ss, fmt.Sprintf("%v", d))
	}

	return strings.Join(ss, " ")
}
