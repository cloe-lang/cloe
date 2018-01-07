package core

// Arguments represents a structured set of arguments passed to a predicate.
// It allows destructive operations to internal properties because it is
// guaranteed by Thunks that arguments objects are never reused as a function
// call creates a Thunk.
type Arguments struct {
	positionals   []*Thunk
	expandedList  *Thunk
	keywords      []KeywordArgument
	expandedDicts []*Thunk
}

// NewArguments creates a new Arguments.
func NewArguments(
	ps []PositionalArgument,
	ks []KeywordArgument,
	ds []*Thunk) Arguments {
	ts := make([]*Thunk, 0, len(ps))
	l := (*Thunk)(nil)

	for i, p := range ps {
		if p.expanded {
			l = mergePositionalArguments(ps[i:])
			break
		}

		ts = append(ts, p.value)
	}

	return Arguments{ts, l, ks, ds}
}

// NewPositionalArguments creates an Arguments which consists of unexpanded
// positional arguments.
func NewPositionalArguments(ts ...*Thunk) Arguments {
	return Arguments{ts, nil, nil, nil}
}

func mergePositionalArguments(ps []PositionalArgument) *Thunk {
	t := EmptyList

	// Optimization for a common pattern of (func a b c ... ..xs).
	// Note that Merge is O(n) but Prepend is O(1).
	if last := len(ps) - 1; ps[last].expanded {
		t = ps[last].value
		ps = ps[:last]
	}

	for i := len(ps) - 1; i >= 0; i-- {
		p := ps[i]

		if p.expanded {
			t = PApp(Merge, p.value, t)
		} else {
			t = Normal(cons(p.value, t))
		}
	}

	return t
}

func (args *Arguments) nextPositional() *Thunk {
	if len(args.positionals) != 0 {
		t := args.positionals[0]
		args.positionals = args.positionals[1:]
		return t
	}

	if args.expandedList == nil {
		return nil
	}

	l := args.expandedList
	args.expandedList = PApp(Rest, l)
	return PApp(First, l)
}

func (args *Arguments) restPositionals() *Thunk {
	ts := args.positionals
	l := args.expandedList
	args.positionals = nil
	args.expandedList = nil

	if l == nil {
		return NewList(ts...)
	} else if len(ts) == 0 {
		return l
	}

	return PApp(Prepend, append(ts, l)...)
}

func (args *Arguments) searchKeyword(s string) *Thunk {
	for i, k := range args.keywords {
		if s == k.name {
			args.keywords = append(args.keywords[:i], args.keywords[i+1:]...)
			return k.value
		}
	}

	for i, t := range args.expandedDicts {
		v := t.Eval()
		d, ok := v.(DictionaryType)

		if !ok {
			return NotDictionaryError(v)
		}

		k := StringType(s)

		// Using DictionaryType.{Search,Remove} methods is safe here
		// because the key is always StringType.
		if v, ok := d.Search(k); ok {
			args.expandedDicts = append([]*Thunk{}, args.expandedDicts...)
			args.expandedDicts[i] = Normal(d.Remove(k))
			return v
		}
	}

	return nil
}

func (args *Arguments) restKeywords() *Thunk {
	ks := args.keywords
	ds := args.expandedDicts
	args.keywords = nil
	args.expandedDicts = nil

	d := emptyDictionary

	for _, k := range ks {
		v := d.insert(StringType(k.name), k.value)
		var ok bool
		d, ok = v.(DictionaryType)

		if !ok {
			return NotDictionaryError(v)
		}
	}

	return PApp(Merge, append([]*Thunk{Normal(d)}, ds...)...)
}

// Merge merges 2 sets of arguments into one.
func (args Arguments) Merge(old Arguments) Arguments {
	ks := append(args.keywords, old.keywords...)
	ds := append(args.expandedDicts, old.expandedDicts...)

	if args.expandedList == nil {
		return Arguments{append(args.positionals, old.positionals...), old.expandedList, ks, ds}
	}

	l := EmptyList

	if old.expandedList != nil {
		l = old.expandedList
	}

	for i := len(old.positionals) - 1; i >= 0; i-- {
		l = Normal(cons(old.positionals[i], l))
	}

	return Arguments{
		args.positionals,
		PApp(Merge, args.expandedList, l),
		ks,
		ds,
	}
}

func (args Arguments) empty() *Thunk {
	if len(args.positionals) > 0 {
		return argumentError("%d positional arguments are left", len(args.positionals))
	}

	// Testing args.expandedList is impossible because we cannot know its length
	// without evaluating it.

	n := 0

	for _, t := range args.expandedDicts {
		v := t.Eval()
		d, ok := v.(DictionaryType)

		if !ok {
			return NotDictionaryError(v)
		}

		n += d.Size()
	}

	if n != 0 || args.keywords != nil && len(args.keywords) > 0 {
		return argumentError("%d keyword arguments are left", len(args.keywords)+n)
	}

	return nil
}
