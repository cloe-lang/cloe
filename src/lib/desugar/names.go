package desugar

import (
	"github.com/raviqqe/tisp/src/lib/ast"
	"log"
)

type names map[string]bool

func newNames() names {
	return make(names)
}

func (ns names) slice() []string {
	ms := make([]string, 0, len(ns))

	for k := range ns {
		ms = append(ms, k)
	}

	return ms
}

func (ns names) copy() names {
	ms := newNames()

	for k, v := range ns {
		ms[k] = v
	}

	return ms
}

func (ns names) merge(ms names) {
	for k, v := range ms {
		ns[k] = v
	}
}

func (ns names) add(n string) {
	ns[n] = true
}

func (ns names) subtract(ms names) {
	for k := range ms {
		delete(ns, k)
	}
}

func (ns names) include(n string) bool {
	_, ok := ns[n]
	return ok
}

func (ns names) find(x interface{}) names {
	switch x := x.(type) {
	case ast.LetConst:
		ns := ns.copy()
		delete(ns, x.Name())
		return ns.find(x.Expr())
	case ast.App:
		ms := newNames()

		ms.merge(ns.find(x.Function()))
		ms.merge(ns.find(x.Arguments()))

		return ms
	case ast.Arguments:
		ms := newNames()

		for _, p := range x.Positionals() {
			ms.merge(ns.find(p.Value()))
		}

		for _, k := range x.Keywords() {
			ms.merge(ns.find(k.Value()))
		}

		for _, d := range x.ExpandedDicts() {
			ms.merge(ns.find(d))
		}

		return ms
	case string:
		ms := newNames()

		if ns.include(x) {
			ms.add(x)
		}

		return ms
	}

	log.Printf("Invalid type: %#v", x)
	panic("Unreachable")
}
