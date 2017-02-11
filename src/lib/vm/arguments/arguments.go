package signature

import ".."

// Arguments represents a structured set of arguments passed to a predicate.
type Arguments struct {
	positionals   []*vm.Thunk
	expandedList  *vm.Thunk
	keywords      []KeywordArgument
	expandedDicts []*vm.Thunk
}

func NewArguments(
	ps []PositionalArgument,
	ks []KeywordArgument,
	expandedDicts []*vm.Thunk) Arguments {
	ts := make([]*vm.Thunk, 0, len(ps))

	var i int

	for j, p := range ps {
		if p.expanded {
			i = j
			break
		}

		ts = append(ts, p.value)
	}

	l := (*vm.Thunk)(nil)

	if i < len(ps) {
		l = listPositionalArgs(ps[i:]...)
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

	// t := ps[0].value

	// for _, p := range ps[1:] {
	// 	t = vm.App(vm.Append, t, p.value)
	// }

	return nil
}
