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

func (hs halfSignature) bindPositionals(args *Arguments) ([]*Thunk, *Thunk) {
	ts := make([]*Thunk, 0, hs.arity())

	for _, name := range hs.requireds {
		t := args.searchKeyword(name)

		if t == nil {
			t = args.nextPositional()
		}

		if t == nil {
			return nil, argumentError("Could not bind a required positional argument.")
		}

		ts = append(ts, t)
	}

	for _, o := range hs.optionals {
		t := args.searchKeyword(o.name)

		if t == nil {
			t = args.nextPositional()
		}

		if t == nil {
			t = o.defaultValue
		}

		ts = append(ts, t)
	}

	if hs.rest != "" {
		t := args.searchKeyword(hs.rest)

		if t == nil {
			t = args.restPositionals()
		}

		ts = append(ts, t)
	}

	if len(ts) != hs.arity() {
		return nil, argumentError("Number of arguments bound to names is different from an arity of a positional half signature.")
	}

	return ts, nil
}

func (hs halfSignature) bindKeywords(args *Arguments) ([]*Thunk, *Thunk) {
	ts := make([]*Thunk, 0, hs.arity())

	for _, name := range hs.requireds {
		t := args.searchKeyword(name)

		if t == nil {
			return nil, argumentError("Could not bind a required keyword argument.")
		}

		ts = append(ts, t)
	}

	for _, opt := range hs.optionals {
		t := args.searchKeyword(opt.name)

		if t == nil {
			t = opt.defaultValue
		}

		ts = append(ts, t)
	}

	if hs.rest != "" {
		t := args.searchKeyword(hs.rest)

		if t == nil {
			t = args.restKeywords()
		}

		ts = append(ts, t)
	}

	if len(ts) != hs.arity() {
		return nil, argumentError("Number of arguments bound to names is different from an arity of a keyword half signature.")
	}

	return ts, nil
}
