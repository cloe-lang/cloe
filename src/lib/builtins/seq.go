package builtins

import (
	"github.com/coel-lang/coel/src/lib/core"
)

func createSeqFunction(f func(t *core.Thunk) core.Value) *core.Thunk {
	return core.NewLazyFunction(
		core.NewSignature(nil, nil, "xs", nil, nil, ""),
		func(ts ...*core.Thunk) core.Value {
			l := ts[0]

			for {
				r := core.PApp(core.Rest, l)

				v := core.PApp(core.Equal, r, core.EmptyList).Eval()
				b, ok := v.(core.BoolType)

				if !ok {
					return core.NotBoolError(v)
				}

				t := core.PApp(core.First, l)

				if b {
					return t
				}

				if err, ok := f(t).(core.ErrorType); ok {
					return err
				}

				l = r
			}
		})
}

// Seq runs arguments of pure values sequentially and returns the last one.
var Seq = createSeqFunction(func(t *core.Thunk) core.Value { return t.Eval() })

// EffectSeq runs arguments of effects sequentially and returns the last one.
var EffectSeq = createSeqFunction(func(t *core.Thunk) core.Value { return t.EvalEffect() })
