package core

// Signature represents function signature.
type Signature struct {
	positionals positionalParameters
	keywords    keywordParameters
}

// NewSignature defines a new Signature.
func NewSignature(ps []string, pr string, ks []OptionalParameter, kr string) Signature {
	return Signature{positionalParameters{ps, pr}, keywordParameters{ks, kr}}
}

// Bind binds Arguments to names defined in Signature and returns full
// arguments to be passed to a function.
func (s Signature) Bind(args Arguments) ([]Value, Value) {
	ps, err := s.positionals.bind(&args)

	if err != nil {
		return nil, err
	}

	ks := s.keywords.bind(&args)

	if err := args.checkEmptyness(); err != nil {
		return nil, err
	}

	return append(ps, ks...), nil
}

func (s Signature) arity() int {
	return s.positionals.arity() + s.keywords.arity()
}
