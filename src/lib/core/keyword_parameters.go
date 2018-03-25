package core

type keywordParameters struct {
	parameters []OptionalParameter
	rest       string
}

func (ks keywordParameters) arity() int {
	n := +len(ks.parameters)

	if ks.rest != "" {
		n++
	}

	return n
}

func (ks keywordParameters) bind(args *Arguments) []Value {
	vs := make([]Value, 0, ks.arity())

	for _, o := range ks.parameters {
		v := args.searchKeyword(o.name)

		if v == nil {
			v = o.defaultValue
		}

		vs = append(vs, v)
	}

	if ks.rest != "" {
		vs = append(vs, args.restKeywords())
	}

	return vs
}
