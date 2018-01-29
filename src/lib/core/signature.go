package core

// Signature represents function signature.
type Signature struct {
	positionals, keywords halfSignature
}

// NewSignature defines a new Signature.
func NewSignature(
	pr []string, po []OptionalArgument, pp string,
	kr []string, ko []OptionalArgument, kk string) Signature {
	return Signature{halfSignature{pr, po, pp}, halfSignature{kr, ko, kk}}
}

// Bind binds Arguments to names defined in Signature and returns full
// arguments to be passed to a function.
func (s Signature) Bind(args Arguments) ([]Value, Value) {
	ps, err := s.positionals.bindPositionals(&args)

	if err != nil {
		return nil, err
	}

	ks, err := s.keywords.bindKeywords(&args)

	if err != nil {
		return nil, err
	}

	if err := args.empty(); err != nil {
		return nil, err
	}

	return append(ps, ks...), nil
}

func (s Signature) arity() int {
	return s.positionals.arity() + s.keywords.arity()
}
