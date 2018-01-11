package builtins

import (
	"testing"

	"github.com/coel-lang/coel/src/lib/core"
	"github.com/coel-lang/coel/src/lib/systemt"
	"github.com/stretchr/testify/assert"
)

func TestRally(t *testing.T) {
	go systemt.RunDaemons()

	ts := []core.Value{core.True, core.False, core.Nil}

	l1 := core.NewList(ts...)
	l2 := core.App(Rally, core.NewArguments(
		[]core.PositionalArgument{core.NewPositionalArgument(l1, true)},
		nil,
		nil))

	for i := 0; i < len(ts); i++ {
		e := core.PApp(core.First, l2)
		t.Logf("%#v", core.EvalPure(e))
		assert.True(t, bool(core.EvalPure(core.PApp(core.Include, l1, e)).(core.BoolType)))

		l1 = core.PApp(core.Delete, l1, core.PApp(indexOf, l1, e))
		l2 = core.PApp(core.Rest, l2)
	}

	assert.True(t, bool(core.EvalPure(core.PApp(core.Equal, core.EmptyList, l1)).(core.BoolType)))
	assert.True(t, bool(core.EvalPure(core.PApp(core.Equal, core.EmptyList, l2)).(core.BoolType)))
}

func TestRallyError(t *testing.T) {
	go systemt.RunDaemons()

	ts := []core.Value{
		core.True,
		core.False,
		core.Nil,
		core.ValueError("I am the sentinel."),
	}

	l := core.App(Rally, core.NewArguments(
		[]core.PositionalArgument{core.NewPositionalArgument(core.NewList(ts...), true)},
		nil,
		nil))

	for i := 0; ; i++ {
		assert.True(t, i < len(ts))

		if _, ok := core.EvalPure(l).(core.ErrorType); ok {
			break
		}

		l = core.PApp(core.Rest, l)
	}
}

var indexOf = core.NewLazyFunction(
	core.NewSignature([]string{"list", "elem"}, nil, "", nil, nil, ""),
	func(ts ...core.Value) core.Value {
		l, e := ts[0], ts[1]

		for i := 1; ; i++ {
			if v := core.ReturnIfEmptyList(l, core.ValueError("A value is not in a list.")); v != nil {
				return v
			}

			v := core.EvalPure(core.PApp(core.Equal, core.PApp(core.First, l), e))
			if b, ok := v.(core.BoolType); !ok {
				return core.NotBoolError(v)
			} else if b {
				return core.NewNumber(float64(i))
			}

			l = core.PApp(core.Rest, l)
		}
	})

func TestRallyWithInvalidExpandedList(t *testing.T) {
	l := core.App(
		Rally,
		core.NewArguments(
			[]core.PositionalArgument{
				core.NewPositionalArgument(core.Nil, false),
				core.NewPositionalArgument(core.OutOfRangeError(), true),
			},
			nil,
			nil))

	if _, ok := core.EvalPure(l).(core.ErrorType); !ok {
		v := core.EvalPure(core.PApp(core.Rest, l))
		_, ok := v.(core.ErrorType)
		assert.True(t, ok)
	}
}
