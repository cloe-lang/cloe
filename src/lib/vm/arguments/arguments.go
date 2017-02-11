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

func (args Arguments) search(s string) *vm.Thunk {
	for _, k := range args.keywords {
		if s == k.name {
			return k.value
		}
	}

	for _, d := range args.expandedDicts {
		o := d.Eval()
		d, ok := o.(vm.DictionaryType)

		if !ok {
			return vm.NotDictionaryError(o)
		}

		if v, ok := d.Search(vm.NewString(s).Eval()); ok {
			return v.(*vm.Thunk)
		}
	}

	return nil
}
