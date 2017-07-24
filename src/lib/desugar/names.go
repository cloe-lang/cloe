package desugar

import "github.com/tisp-lang/tisp/src/lib/ast"

type names map[string]bool

func newNames(ss ...string) names {
	ns := make(names, len(ss))

	for _, s := range ss {
		ns.add(s)
	}

	return ns
}

func (ns names) slice() []string {
	ms := make([]string, 0, len(ns))

	for n := range ns {
		ms = append(ms, n)
	}

	return ms
}

func (ns names) add(n string) {
	ns[n] = true
}

func (ns names) copy() names {
	ms := newNames()

	for n := range ns {
		ms.add(n)
	}

	return ms
}

func (ns names) merge(ms names) {
	for m := range ms {
		ns.add(m)
	}
}

func (ns names) delete(n string) {
	delete(ns, n)
}

func (ns names) subtract(ms names) {
	for m := range ms {
		ns.delete(m)
	}
}

func (ns names) include(n string) bool {
	_, ok := ns[n]
	return ok
}

// findInFunction finds names in a let-function node. find assumes that a given
// let-function node does not define its function recursively.
func (ns names) findInFunction(f ast.LetFunction) names {
	return ns.find(f)
}

// find finds names in a AST node. This should not be used directly.
func (ns names) find(x interface{}) names {
	switch x := x.(type) {
	case ast.LetVar:
		return ns.find(x.Expr())
	case ast.LetFunction:
		ns := ns.copy()

		for n := range signatureToNames(x.Signature()) {
			ns.delete(n)
		}

		ms := newNames()

		for _, l := range x.Lets() {
			switch l := l.(type) {
			case ast.LetVar:
				ms.merge(ns.find(l))
				ns.delete(l.Name())
			case ast.LetFunction:
				ns.delete(l.Name())
				ms.merge(ns.find(l))
			default:
				panic("Unreachable")
			}
		}

		ms.merge(ns.find(x.Body()))
		return ms
	case ast.App:
		ms := ns.find(x.Function())
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
	case ast.Switch:
		ms := ns.find(x.Value())

		for _, c := range x.Cases() {
			ms.merge(ns.find(c.Value()))
		}

		ms.merge(ns.find(x.DefaultCase()))

		return ms
	case string:
		if ns.include(x) {
			return newNames(x)
		}

		return newNames()
	}

	panic("Unreachable")
}
