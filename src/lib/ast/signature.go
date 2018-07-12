package ast

// Signature represents a signature of a function.
type Signature struct {
	positionals positionalParameters
	keywords    keywordParameters
}

// NewSignature creates a Signature from {positional, keyword} x
// {required, optional} arguments and a positional rest argument
// and a keyword rest argument.
func NewSignature(ps []string, pr string, ks []OptionalParameter, kr string) Signature {
	return Signature{
		positionals: positionalParameters{ps, pr},
		keywords:    keywordParameters{ks, kr},
	}
}

// Positionals returns positional required arguments of a signature.
func (s Signature) Positionals() []string {
	return s.positionals.parameters
}

// RestPositionals returns a positional rest argument of a signature.
func (s Signature) RestPositionals() string {
	return s.positionals.rest
}

// Keywords returns keyword optional arguments of a signature.
func (s Signature) Keywords() []OptionalParameter {
	return s.keywords.parameters
}

// RestKeywords returns a keyword rest argument of a signature.
func (s Signature) RestKeywords() string {
	return s.keywords.rest
}

// Names returns names used in a signature.
func (s Signature) Names() []string {
	return append(s.positionals.names(), s.keywords.names()...)
}

func (s Signature) String() string {
	str := s.positionals.String()

	if ks := s.keywords.String(); ks != "" {
		str += " . " + ks
	}

	return str
}
