package desugar

import (
	"github.com/raviqqe/tisp/src/lib/ast"
	"github.com/raviqqe/tisp/src/lib/util"
)

type names map[string]bool

func newNames(ss ...string) names {
	ns := make(names)

	for _, s := range ss {
		ns.add(s)
	}

	return ns
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
	case []interface{}:
		ms := newNames()

		for _, s := range x {
			ms.merge(ns.find(s))
		}

		return ms
	case ast.LetVar:
		ns := ns.copy()
		delete(ns, x.Name())
		return ns.find(x.Expr())
	case ast.LetFunction:
		ns := ns.copy()

		delete(ns, x.Name())
		for _, n := range signatureToNames(x.Signature()).slice() {
			delete(ns, n)
		}

		ms := ns.find(x.Lets())
		ms.merge(ns.find(x.Body()))
		return ms
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

	util.Fail("Invalid type: %#v", x)
	panic("Unreachable")
}
