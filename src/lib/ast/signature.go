package ast

type Signature struct {
	positionals argumentSet
	keywords    argumentSet
}

type argumentSet struct {
	requireds []string
	optionals []OptionalArgument
	rest      string
}

func NewSignature(
	pr []string, po []OptionalArgument, pp string,
	kr []string, ko []OptionalArgument, kk string) Signature {
	return Signature{
		positionals: argumentSet{pr, po, pp},
		keywords:    argumentSet{kr, ko, kk},
	}
}

func (s Signature) PosReqs() []string {
	return s.positionals.requireds
}

func (s Signature) PosOpts() []OptionalArgument {
	return s.positionals.optionals
}

func (s Signature) PosRest() string {
	return s.positionals.rest
}

func (s Signature) KeyReqs() []string {
	return s.keywords.requireds
}

func (s Signature) KeyOpts() []OptionalArgument {
	return s.keywords.optionals
}

func (s Signature) KeyRest() string {
	return s.keywords.rest
}
