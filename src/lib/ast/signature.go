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
