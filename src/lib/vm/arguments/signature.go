package signature

import ".."

// Signature represents function signature.
type Signature struct {
	positionals argumentSet
	keywords    argumentSet
}

// NewSignature defines a new Signature.
func NewSignature(
	pr []string, po []OptionalArgument, pp string,
	kr []string, ko []OptionalArgument, kk string) Signature {
	return Signature{
		positionals: argumentSet{
			requireds: pr,
			optionals: po,
			rest:      pp,
		},
		keywords: argumentSet{
			requireds: kr,
			optionals: ko,
			rest:      kk,
		},
	}
}

// Bind binds Arguments to names defined in Signature and returns full
// arguments to be passed to a function.
func (s Signature) Bind(args Arguments) []*vm.Thunk {
	ts := make([]*vm.Thunk, 0, s.arity())

	// for i, p := range args.positionals {
	// 	if i != len(args.positionals)-1 &&  {
	// 	}
	// 	p.value
	// }

	if len(ts) != s.arity() {
		panic("You sucks!")
	}

	return ts
}

func (s Signature) arity() int {
	return s.positionals.size() + s.keywords.size()
}
