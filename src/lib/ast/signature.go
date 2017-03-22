package ast

import "fmt"

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

// Arity returns a number of arguments in a signature.
func (s Signature) Arity() int {
	return s.positionals.arity() + s.keywords.arity()
}

func (as halfSignature) arity() int {
	rest := 0

	if as.rest != "" {
		rest = 1
	}

	return len(as.requireds) + len(as.optionals) + rest
}

// NameToIndex converts an argument name into an index in arguments inside a signature.
func (s Signature) NameToIndex(name string) (int, error) {
	i, ok := s.positionals.nameToIndex(name)

	if ok {
		return i, nil
	}

	j, ok := s.keywords.nameToIndex(name)

	if ok {
		return i + j, nil
	}

	return -1, fmt.Errorf("name %#v was not found in a signature", name)
}

func (as halfSignature) nameToIndex(name string) (int, bool) {
	i := 0

	for _, r := range as.requireds {
		if name == r {
			return i, true
		}
		i++
	}

	for _, o := range as.optionals {
		if name == o.name {
			return i, true
		}
		i++
	}

	if as.rest != "" {
		if name == as.rest {
			return i, true
		}
		i++
	}

	return i, false
}
