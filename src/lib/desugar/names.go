package desugar

import "github.com/cloe-lang/cloe/src/lib/ast"

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

// findInDefFunction finds names in a let-function node. This assumes that a
// given let-function statement does not define its function recursively.
func (ns names) findInDefFunction(f ast.DefFunction) names {
	ns = ns.copy()
	ns.subtract(signatureToNames(f.Signature()))

	ms := newNames()

	for _, l := range f.Lets() {
		switch l := l.(type) {
		case ast.LetVar:
			ms.merge(ns.findInLetVar(l))
			ns.delete(l.Name())
		case ast.DefFunction:
			ns.delete(l.Name())
			ms.merge(ns.findInDefFunction(l))
		default:
			panic("unreachable")
		}
	}

	ms.merge(ns.findInExpression(f.Body()))
	return ms
}

// findInLetVar finds names in a let-variable node. This assumes that a given
// let-variable statement does not define its variable recursively.
func (ns names) findInLetVar(l ast.LetVar) names {
	return ns.findInExpression(l.Expr())
}

func (ns names) findInExpression(x interface{}) names {
	switch x := x.(type) {
	case string:
		if ns.include(x) {
			return newNames(x)
		}

		return newNames()
	case ast.AnonymousFunction:
		ns := ns.copy()
		ns.subtract(signatureToNames(x.Signature()))
		return ns.findInExpression(x.Body())
	case ast.App:
		ms := ns.findInExpression(x.Function())
		ms.merge(ns.findInExpression(x.Arguments()))
		return ms
	case ast.Arguments:
		ms := newNames()

		for _, p := range x.Positionals() {
			ms.merge(ns.findInExpression(p.Value()))
		}

		for _, k := range x.Keywords() {
			ms.merge(ns.findInExpression(k.Value()))
		}

		return ms
	case ast.Switch:
		ms := ns.findInExpression(x.Value())

		for _, c := range x.Cases() {
			ms.merge(ns.findInExpression(c.Value()))
		}

		ms.merge(ns.findInExpression(x.DefaultCase()))

		return ms
	}

	panic("unreachable")
}
