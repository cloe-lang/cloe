package vm

// Arguments represents a structured set of arguments passed to a predicate.
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
	expandedDicts []*Thunk) Arguments {
	ts := make([]*Thunk, 0, len(ps))

	l := (*Thunk)(nil)

	for i, p := range ps {
		if p.expanded {
			l = mergeRestPositionalArgs(ps[i:]...)
			break
		}

		ts = append(ts, p.value)
	}

	return Arguments{
		positionals:   ts,
		expandedList:  l,
		keywords:      ks,
		expandedDicts: expandedDicts,
	}
}

func mergeRestPositionalArgs(ps ...PositionalArgument) *Thunk {
	if !ps[0].expanded {
		panic("First PositionalArgument must be a list.")
	}

	t := ps[0].value

	for _, p := range ps[1:] {
		if p.expanded {
			t = App(Merge, t, NewList(p.value))
		} else {
			t = App(Append, t, p.value)
		}
	}

	return t
}

func (args *Arguments) nextPositional() *Thunk {
	if len(args.positionals) != 0 {
		defer func() { args.positionals = args.positionals[1:] }()
		return args.positionals[0]
	}

	if args.expandedList == nil {
		return nil
	}

	defer func() { args.expandedList = App(Rest, args.expandedList) }()
	return App(First, args.expandedList)
}

func (args Arguments) restPositionals() *Thunk {
	if args.expandedList == nil {
		return NewList(args.positionals...)
	}

	if len(args.positionals) == 0 {
		return args.expandedList
	}

	return App(Merge, NewList(args.positionals...), NewList(args.expandedList))
}

func (args *Arguments) searchKeyword(s string) *Thunk {
	for i, k := range args.keywords {
		if s == k.name {
			args.keywords = append(args.keywords[:i], args.keywords[i+1:]...)
			return k.value
		}
	}

	for i, t := range args.expandedDicts {
		o := t.Eval()
		d, ok := o.(DictionaryType)

		if !ok {
			return NotDictionaryError(o)
		}

		k := StringType(s)

		if v, ok := d.Search(k); ok {
			args.expandedDicts[i] = Normal(d.Remove(k))
			return v.(*Thunk)
		}
	}

	return nil
}

func (args Arguments) restKeywords() *Thunk {
	t := EmptyDictionary

	for _, k := range args.keywords {
		t = App(Set, t, NewString(k.name), k.value)
	}

	for _, tt := range args.expandedDicts {
		t = App(Merge, t, NewList(tt))
	}

	return t
}

func (original Arguments) Merge(merged Arguments) Arguments {
	var new Arguments

	if new.expandedList == nil {
		new.positionals = append(original.positionals, merged.positionals...)
		new.expandedList = merged.expandedList
	} else {
		new.positionals = original.positionals
		new.expandedList = App(
			Append,
			append([]*Thunk{original.expandedList}, merged.positionals...)...)

		if merged.expandedList != nil {
			new.expandedList = App(Merge, new.expandedList, NewList(merged.expandedList))
		}
	}

	new.keywords = append(original.keywords, merged.keywords...)
	new.expandedDicts = append(original.expandedDicts, merged.expandedDicts...)

	return new
}
