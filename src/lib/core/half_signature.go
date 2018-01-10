package core

type halfSignature struct {
	requireds []string
	optionals []OptionalArgument
	rest      string
}

func (hs halfSignature) arity() int {
	n := len(hs.requireds) + len(hs.optionals)

	if hs.rest != "" {
		n++
	}

	return n
}

func (hs halfSignature) bindPositionals(args *Arguments) ([]Value, Value) {
	vs := make([]Value, 0, hs.arity())

	for _, s := range hs.requireds {
		v := args.nextPositional()

		if v == nil {
			v = args.searchKeyword(s)
		}

		if v == nil {
			return nil, argumentError("Could not bind a required positional argument.")
		}

		vs = append(vs, v)
	}

	for _, o := range hs.optionals {
		v := args.nextPositional()

		if v == nil {
			v = args.searchKeyword(o.name)
		}

		if v == nil {
			v = o.defaultValue
		}

		vs = append(vs, v)
	}

	if hs.rest != "" {
		vs = append(vs, args.restPositionals())
	}

	return vs, nil
}

func (hs halfSignature) bindKeywords(args *Arguments) ([]Value, Value) {
	vs := make([]Value, 0, hs.arity())

	for _, s := range hs.requireds {
		v := args.searchKeyword(s)

		if v == nil {
			return nil, argumentError("Could not bind a required keyword argument.")
		}

		vs = append(vs, v)
	}

	for _, o := range hs.optionals {
		v := args.searchKeyword(o.name)

		if v == nil {
			v = o.defaultValue
		}

		vs = append(vs, v)
	}

	if hs.rest != "" {
		vs = append(vs, args.restKeywords())
	}

	return vs, nil
}
