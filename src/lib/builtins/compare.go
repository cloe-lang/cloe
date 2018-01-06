package builtins

import "github.com/coel-lang/coel/src/lib/core"

func createCompareFunction(checkOrder func(core.NumberType) bool) *core.Thunk {
	return core.NewLazyFunction(
		core.NewSignature(nil, nil, "args", nil, nil, ""),
		func(ts ...*core.Thunk) core.Value {
			l := ts[0]

			if v := checkEmptyList(l, core.True); v != nil {
				return v
			}

			prev := core.PApp(core.First, l)

			if b, err := core.EvalBool(core.PApp(core.IsOrdered, prev)); err != nil {
				return err
			} else if !b {
				return core.NotOrderedError(prev.Eval())
			}

			for {
				l = core.PApp(core.Rest, l)

				if v := checkEmptyList(l, core.True); v != nil {
					return v
				}

				current := core.PApp(core.First, l)

				if n, err := core.EvalNumber(core.PApp(core.Compare, prev, current)); err != nil {
					return err
				} else if !checkOrder(n) {
					return core.False
				}

				prev = current
			}
		})
}

// Less checks if arguments are aligned in ascending order or not.
var Less = createCompareFunction(func(n core.NumberType) bool { return n == -1 })

// LessEq checks if arguments are aligned in ascending order or not.
var LessEq = createCompareFunction(func(n core.NumberType) bool { return n == -1 || n == 0 })

// Greater checks if arguments are aligned in ascending order or not.
var Greater = createCompareFunction(func(n core.NumberType) bool { return n == 1 })

// GreaterEq checks if arguments are aligned in ascending order or not.
var GreaterEq = createCompareFunction(func(n core.NumberType) bool { return n == 1 || n == 0 })
