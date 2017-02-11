package signature

import ".."

// Arguments represents a structured set of arguments passed to a predicate.
type Arguments struct {
	positionals   []*vm.Thunk
	expandedList  *vm.Thunk
	keywords      []KeywordArgument
	expandedDicts []*vm.Thunk
}

// NewArguments creates a new Arguments.
func NewArguments(
	ps []PositionalArgument,
	ks []KeywordArgument,
	expandedDicts []*vm.Thunk) Arguments {
	ts := make([]*vm.Thunk, 0, len(ps))

	l := (*vm.Thunk)(nil)

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

func mergeRestPositionalArgs(ps ...PositionalArgument) *vm.Thunk {
	if !ps[0].expanded {
		panic("First PositionalArgument must be a list.")
	}

	t := ps[0].value

	for _, p := range ps[1:] {
		if p.expanded {
			t = vm.App(vm.Merge, t, p.value)
		} else {
			t = vm.App(vm.Append, t, p.value)
		}
	}

	return t
}

func (args *Arguments) nextPositional() *vm.Thunk {
	if len(args.positionals) != 0 {
		defer func() { args.positionals = args.positionals[1:] }()
		return args.positionals[0]
	}

	if args.expandedList == nil {
		return nil
	}

	defer func() { args.expandedList = vm.App(vm.Rest, args.expandedList) }()
	return vm.App(vm.First, args.expandedList)
}

func (args Arguments) restPositionals() *vm.Thunk {
	if args.expandedList == nil {
		return vm.NewList(args.positionals...)
	}

	if len(args.positionals) == 0 {
		return args.expandedList
	}

	return vm.App(vm.Merge, vm.NewList(args.positionals...), args.expandedList)
}

func (args *Arguments) searchKeyword(s string) *vm.Thunk {
	for i, k := range args.keywords {
		if s == k.name {
			args.keywords = append(args.keywords[:i], args.keywords[i+1:]...)
			return k.value
		}
	}

	for i, t := range args.expandedDicts {
		o := t.Eval()
		d, ok := o.(vm.DictionaryType)

		if !ok {
			return vm.NotDictionaryError(o)
		}

		k := vm.StringType(s)

		if v, ok := d.Search(k); ok {
			args.expandedDicts[i] = vm.Normal(d.Remove(k))
			return v.(*vm.Thunk)
		}
	}

	return nil
}

func (args Arguments) restKeywords() *vm.Thunk {
	t := vm.EmptyDictionary

	for _, k := range args.keywords {
		t = vm.App(vm.Set, t, vm.NewString(k.name), k.value)
	}

	for _, tt := range args.expandedDicts {
		t = vm.App(vm.Merge, t, tt)
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
		new.expandedList = vm.App(
			vm.Append,
			append([]*vm.Thunk{original.expandedList}, merged.positionals...)...)

		if merged.expandedList != nil {
			new.expandedList = vm.App(vm.Merge, new.expandedList, merged.expandedList)
		}
	}

	new.keywords = append(original.keywords, merged.keywords...)
	new.expandedDicts = append(original.expandedDicts, merged.expandedDicts...)

	return new
}
