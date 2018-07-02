package core

// Arguments represents a structured set of arguments passed to a predicate.
// It allows destructive operations to internal properties because it is
// guaranteed by Thunks that arguments objects are never reused as a function
// call creates a Thunk.
type Arguments struct {
	positionals  []Value
	expandedList Value
	keywords     []KeywordArgument
}

// NewArguments creates a new Arguments.
func NewArguments(ps []PositionalArgument, ks []KeywordArgument) Arguments {
	vs := make([]Value, 0, len(ps))
	l := Value(nil)

	for i, p := range ps {
		if p.expanded {
			l = mergePositionalArguments(ps[i:])
			break
		}

		vs = append(vs, p.value)
	}

	return Arguments{vs, l, ks}
}

// NewPositionalArguments creates an Arguments which consists of unexpanded
// positional arguments.
func NewPositionalArguments(vs ...Value) Arguments {
	return Arguments{vs, nil, nil}
}

func mergePositionalArguments(ps []PositionalArgument) Value {
	v := Value(EmptyList)

	// Optimization for a common pattern of (func a b c ... ..xs).
	// Note that Merge is O(n) but Prepend is O(1).
	if last := len(ps) - 1; ps[last].expanded {
		v = ps[last].value
		ps = ps[:last]
	}

	for i := len(ps) - 1; i >= 0; i-- {
		p := ps[i]

		if p.expanded {
			v = PApp(Merge, p.value, v)
		} else {
			v = cons(p.value, v)
		}
	}

	return v
}

func (args *Arguments) nextPositional() Value {
	if len(args.positionals) != 0 {
		v := args.positionals[0]
		args.positionals = args.positionals[1:]
		return v
	}

	if args.expandedList == nil {
		return nil
	}

	l := args.expandedList
	args.expandedList = PApp(Rest, l)
	return PApp(First, l)
}

func (args *Arguments) restPositionals() Value {
	vs := args.positionals
	l := args.expandedList
	args.positionals = nil
	args.expandedList = nil

	if l == nil {
		return NewList(vs...)
	} else if len(vs) == 0 {
		return l
	}

	return StrictPrepend(vs, l)
}

func (args *Arguments) searchKeyword(s string) Value {
	for i := len(args.keywords) - 1; i >= 0; i-- {
		k := args.keywords[i]

		if k.name == s {
			args.keywords = append(args.keywords[:i], args.keywords[i+1:]...)
			return k.value
		} else if k.name == "" {
			d, err := EvalDictionary(k.value)

			if err != nil {
				return err
			}

			k := NewString(s)

			// Using DictionaryType.{index,delete} methods is safe here
			// because the key is always StringType.
			if v, err := d.find(k); err == nil {
				args.keywords = append(
					args.keywords[:i],
					append(
						[]KeywordArgument{NewKeywordArgument("", d.delete(k))},
						args.keywords[i+1:]...)...)
				return v
			}
		}
	}

	return nil
}

func (args *Arguments) restKeywords() Value {
	ks := args.keywords
	args.keywords = nil

	d := Value(EmptyDictionary)

	for _, k := range ks {
		// Using DictionaryType.Insert method is safe here
		// because the key is always StringType.
		if k.name == "" {
			d = PApp(Merge, d, k.value)
		} else {
			d = PApp(Insert, d, NewString(k.name), k.value)
		}
	}

	return d
}

// Merge merges 2 sets of arguments into one.
func (args Arguments) Merge(old Arguments) Arguments {
	ks := append(args.keywords, old.keywords...)

	if args.expandedList == nil {
		return Arguments{append(args.positionals, old.positionals...), old.expandedList, ks}
	}

	l := Value(EmptyList)

	if old.expandedList != nil {
		l = old.expandedList
	}

	return Arguments{
		args.positionals,
		PApp(Merge, args.expandedList, StrictPrepend(old.positionals, l)),
		ks,
	}
}

func (args Arguments) checkEmptyness() Value {
	if len(args.positionals) > 0 {
		return argumentError("%d positional arguments are left", len(args.positionals))
	}

	// Testing args.expandedList is impossible because we cannot know its length
	// without evaluating it.

	// Keyword arguments are not checked in the current implementation as
	// expanded dictionaries can contain extra arguments.

	return nil
}
