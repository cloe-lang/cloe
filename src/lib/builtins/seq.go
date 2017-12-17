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
				t := core.PApp(core.First, l)

				if err, ok := f(t).(core.ErrorType); ok {
					return err
				}

				l = core.PApp(core.Rest, l)

				if v := checkEmptyList(l, t); v != nil {
					return v
				}
			}
		})
}

// Seq runs arguments of pure values sequentially and returns the last one.
var Seq = createSeqFunction(func(t *core.Thunk) core.Value { return t.Eval() })

// EffectSeq runs arguments of effects sequentially and returns the last one.
var EffectSeq = createSeqFunction(func(t *core.Thunk) core.Value { return t.EvalEffect() })
