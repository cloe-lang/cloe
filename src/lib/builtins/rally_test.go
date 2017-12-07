package builtins

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tisp-lang/tisp/src/lib/core"
	"github.com/tisp-lang/tisp/src/lib/systemt"
)

func TestRally(t *testing.T) {
	go systemt.RunDaemons()

	ts := []*core.Thunk{core.True, core.False, core.Nil}

	l1 := core.NewList(ts...)
	l2 := core.App(Rally, core.NewArguments(
		[]core.PositionalArgument{core.NewPositionalArgument(l1, true)},
		nil,
		nil))

	for i := 0; i < len(ts); i++ {
		e := core.PApp(core.First, l2)
		t.Logf("%#v", e.Eval())
		assert.True(t, bool(core.PApp(core.Include, l1, e).Eval().(core.BoolType)))

		l1 = core.PApp(core.Delete, l1, core.PApp(indexOf, l1, e))
		l2 = core.PApp(core.Rest, l2)
	}

	assert.True(t, bool(core.PApp(core.Equal, core.EmptyList, l1).Eval().(core.BoolType)))
	assert.True(t, bool(core.PApp(core.Equal, core.EmptyList, l2).Eval().(core.BoolType)))
}

func TestRallyError(t *testing.T) {
	go systemt.RunDaemons()

	ts := []*core.Thunk{
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

		if _, ok := l.Eval().(core.ErrorType); ok {
			break
		}

		l = core.PApp(core.Rest, l)
	}
}

var indexOf = core.NewLazyFunction(
	core.NewSignature([]string{"list", "elem"}, nil, "", nil, nil, ""),
	func(ts ...*core.Thunk) core.Value {
		l, e := ts[0], ts[1]

		for i := 0; ; i++ {
			if v := checkEmptyList(l, core.ValueError("A value is not in a list.")); v != nil {
				return v
			}

			v := core.PApp(core.Equal, core.PApp(core.First, l), e).Eval()
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

	if _, ok := l.Eval().(core.ErrorType); !ok {
		v := core.PApp(core.Rest, l).Eval()
		_, ok := v.(core.ErrorType)
		assert.True(t, ok)
	}
}
