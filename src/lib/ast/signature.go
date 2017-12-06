package ast

import (
	"strings"
)

// Signature represents a signature of a function.
type Signature struct {
	positionals halfSignature
	keywords    halfSignature
}

// halfSignature represents a half of a signature.
type halfSignature struct {
	requireds []string
	optionals []OptionalArgument
	rest      string
}

// NewSignature creates a Signature from {positional, keyword} x
// {required, optional} arguments and a positional rest argument
// and a keyword rest argument.
func NewSignature(
	pr []string, po []OptionalArgument, pp string,
	kr []string, ko []OptionalArgument, kk string) Signature {
	return Signature{
		positionals: halfSignature{pr, po, pp},
		keywords:    halfSignature{kr, ko, kk},
	}
}

// PosReqs returns positional required arguments of a signature.
func (s Signature) PosReqs() []string {
	return s.positionals.requireds
}

// PosOpts returns positional optional arguments of a signature.
func (s Signature) PosOpts() []OptionalArgument {
	return s.positionals.optionals
}

// PosRest returns a positional rest argument of a signature.
func (s Signature) PosRest() string {
	return s.positionals.rest
}

// KeyReqs returns keyword required arguments of a signature.
func (s Signature) KeyReqs() []string {
	return s.keywords.requireds
}

// KeyOpts returns keyword optional arguments of a signature.
func (s Signature) KeyOpts() []OptionalArgument {
	return s.keywords.optionals
}

// KeyRest returns a keyword rest argument of a signature.
func (s Signature) KeyRest() string {
	return s.keywords.rest
}

// NameToIndex converts an argument name into an index in arguments inside a signature.
func (s Signature) NameToIndex() map[string]int {
	m := map[string]int{}

	for i, n := range append(s.positionals.names(), s.keywords.names()...) {
		m[n] = i
	}

	return m
}

func (hs halfSignature) names() []string {
	ns := hs.requireds

	for _, o := range hs.optionals {
		ns = append(ns, o.name)
	}

	if hs.rest != "" {
		ns = append(ns, hs.rest)
	}

	return ns
}

func (hs halfSignature) String() string {
	ss := make([]string, 0, len(hs.requireds)+len(hs.optionals)+1)

	ss = append(ss, hs.requireds...)

	for _, o := range hs.optionals {
		ss = append(ss, o.String())
	}

	if hs.rest != "" {
		ss = append(ss, hs.rest)
	}

	return strings.Join(ss, " ")
}

func (s Signature) String() string {
	str := s.positionals.String()

	if ks := s.keywords.String(); ks != "" {
		str += " " + ks
	}

	return str
}
