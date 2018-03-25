package core

type positionalParameters struct {
	parameters []string
	rest       string
}

func (ps positionalParameters) arity() int {
	n := len(ps.parameters)

	if ps.rest != "" {
		n++
	}

	return n
}

func (ps positionalParameters) bind(args *Arguments) ([]Value, Value) {
	vs := make([]Value, 0, ps.arity())

	for _, s := range ps.parameters {
		v := args.nextPositional()

		if v == nil {
			return nil, argumentError("positional argument, {} is missing", s)
		}

		vs = append(vs, v)
	}

	if ps.rest != "" {
		vs = append(vs, args.restPositionals())
	}

	return vs, nil
}
