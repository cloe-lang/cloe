package ir

import (
	"github.com/cloe-lang/cloe/src/lib/core"
)

// CreateFunction creates a new user-defined function.
func CreateFunction(s core.Signature, ns []string, ls []interface{}, b interface{}, f func(string) (core.Value, error)) (core.Value, error) {
	c := newCompiler(f)
	bs, cs, ss, ns, err := c.Compile(ns, ls, b)

	if err != nil {
		return nil, err
	}

	return core.NewLazyFunction(
		s,
		func(vs ...core.Value) core.Value {
			args := make([]core.Value, len(cs)+len(vs))
			copy(args, append(cs, vs...))

			i := NewInterpreter(bs, ss, ns, args)

			return i.Interpret()
		}), nil
}
