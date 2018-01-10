package core

// Arguments represents a structured set of arguments passed to a predicate.
// It allows destructive operations to internal properties because it is
// guaranteed by Thunks that arguments objects are never reused as a function
// call creates a Thunk.
type Arguments struct {
	positionals   []Value
	expandedList  Value
	keywords      []KeywordArgument
	expandedDicts []Value
}

// NewArguments creates a new Arguments.
func NewArguments(
	ps []PositionalArgument,
	ks []KeywordArgument,
	ds []Value) Arguments {
	vs := make([]Value, 0, len(ps))
	l := Value(nil)

	for i, p := range ps {
		if p.expanded {
			l = mergePositionalArguments(ps[i:])
			break
		}

		vs = append(vs, p.value)
	}

	return Arguments{vs, l, ks, ds}
}

// NewPositionalArguments creates an Arguments which consists of unexpanded
// positional arguments.
func NewPositionalArguments(vs ...Value) Arguments {
	return Arguments{vs, nil, nil, nil}
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
	for i, k := range args.keywords {
		if s == k.name {
			args.keywords = append(args.keywords[:i], args.keywords[i+1:]...)
			return k.value
		}
	}

	for i, v := range args.expandedDicts {
		d, ok := EvalPure(v).(DictionaryType)

		if !ok {
			return NotDictionaryError(v)
		}

		k := NewString(s)

		// Using DictionaryType.{Search,Remove} methods is safe here
		// because the key is always StringType.
		if v, ok := d.Search(k); ok {
			args.expandedDicts = append([]Value{}, args.expandedDicts...)
			args.expandedDicts[i] = d.Remove(k)
			return v
		}
	}

	return nil
}

func (args *Arguments) restKeywords() Value {
	ks := args.keywords
	ds := args.expandedDicts
	args.keywords = nil
	args.expandedDicts = nil

	d := EmptyDictionary

	for _, k := range ks {
		// Using DictionaryType.Insert method is safe here
		// because the key is always StringType.
		d = d.Insert(NewString(k.name), k.value)
	}

	return PApp(Merge, append([]Value{d}, ds...)...)
}

// Merge merges 2 sets of arguments into one.
func (args Arguments) Merge(old Arguments) Arguments {
	ks := append(args.keywords, old.keywords...)
	ds := append(args.expandedDicts, old.expandedDicts...)

	if args.expandedList == nil {
		return Arguments{append(args.positionals, old.positionals...), old.expandedList, ks, ds}
	}

	l := Value(EmptyList)

	if old.expandedList != nil {
		l = old.expandedList
	}

	return Arguments{
		args.positionals,
		PApp(Merge, args.expandedList, StrictPrepend(old.positionals, l)),
		ks,
		ds,
	}
}

func (args Arguments) empty() Value {
	if len(args.positionals) > 0 {
		return argumentError("%d positional arguments are left", len(args.positionals))
	}

	// Testing args.expandedList is impossible because we cannot know its length
	// without evaluating it.

	n := 0

	for _, t := range args.expandedDicts {
		v := EvalPure(t)
		d, ok := v.(DictionaryType)

		if !ok {
			return NotDictionaryError(v)
		}

		n += d.Size()
	}

	if n != 0 || len(args.keywords) > 0 {
		return argumentError("%d keyword arguments are left", len(args.keywords)+n)
	}

	return nil
}
