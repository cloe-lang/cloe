package core

// Signature represents function signature.
type Signature struct {
	positionals halfSignature
	keywords    halfSignature
}

// NewSignature defines a new Signature.
func NewSignature(
	pr []string, po []OptionalArgument, pp string,
	kr []string, ko []OptionalArgument, kk string) Signature {
	return Signature{
		positionals: halfSignature{
			requireds: pr,
			optionals: po,
			rest:      pp,
		},
		keywords: halfSignature{
			requireds: kr,
			optionals: ko,
			rest:      kk,
		},
	}
}

// Bind binds Arguments to names defined in Signature and returns full
// arguments to be passed to a function.
func (s Signature) Bind(args Arguments) ([]*Thunk, *Thunk) {
	ps, err := s.positionals.bindPositionals(&args)

	if err != nil {
		return nil, err
	}

	ks, err := s.keywords.bindKeywords(&args)

	if err != nil {
		return nil, err
	}

	ts := append(ps, ks...)

	if len(ts) != s.arity() {
		return nil, argumentError("Number of arguments bound to names is different from signature's arity.")
	}

	if err := args.empty(); err != nil {
		return nil, err
	}

	return ts, nil
}

func (s Signature) arity() int {
	return s.positionals.arity() + s.keywords.arity()
}
