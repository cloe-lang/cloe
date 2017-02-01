package signature

import "../vm"

type Signature struct {
	positionals argumentSet
	keywords    argumentSet
}

type argumentSet struct {
	required []string
	optional []OptionalArgument
	rest     string
}

type OptionalArgument struct {
	key          string
	defaultValue *vm.Thunk
}

func NewSignature(
	pr []string, po []OptionalArgument, pp string,
	kr []string, ko []OptionalArgument, kk string) Signature {
	return Signature{
		positionals: argumentSet{
			required: pr,
			optional: po,
			rest:     pp,
		},
		keywords: argumentSet{
			required: kr,
			optional: ko,
			rest:     kk,
		},
	}
}

func (s Signature) App(args Arguments) Signature {
	return Signature{}
}
