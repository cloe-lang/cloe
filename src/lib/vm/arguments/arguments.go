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
			l = listPositionalArgs(ps[i:]...)
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

func listPositionalArgs(ps ...PositionalArgument) *vm.Thunk {
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
