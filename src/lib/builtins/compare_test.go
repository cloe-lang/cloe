package builtins

import (
	"testing"

	"github.com/cloe-lang/cloe/src/lib/core"
	"github.com/stretchr/testify/assert"
)

func TestLessTrue(t *testing.T) {
	for _, ts := range [][]core.Value{
		{},
		{core.NewNumber(42)},
		{core.NewString("foo")},
		{core.NewList(core.NewNumber(42))},
		{core.NewNumber(0), core.NewNumber(0.0000000000001)},
		{core.NewNumber(0), core.NewNumber(0.00000001), core.NewNumber(0.0001)},
		{core.NewString("bar"), core.NewString("baz"), core.NewString("foo")},
		{core.NewList(core.NewNumber(42)), core.NewList(core.NewNumber(2049))},
		{core.NewList(core.NewNumber(42)), core.NewList(core.NewNumber(42), core.NewNumber(42))},
		{core.NewList(core.NewList(core.NewNumber(42))),
			core.NewList(core.NewList(core.NewNumber(42), core.NewNumber(42))),
			core.NewList(core.NewList(core.NewNumber(42), core.NewNumber(42), core.NewNumber(42)))},
	} {
		assert.True(t, bool(*core.EvalPure(core.PApp(Less, ts...)).(*core.BooleanType)))
	}
}

func TestLessFalse(t *testing.T) {
	for _, ts := range [][]core.Value{
		{core.NewNumber(42), core.NewNumber(0)},
		{core.NewNumber(0), core.NewNumber(42), core.NewNumber(42)},
		{core.NewNumber(0), core.NewNumber(0), core.NewNumber(42)},
		{core.NewString("bar"), core.NewString("bar"), core.NewString("baz")},
		{core.NewString("foo"), core.NewString("bar"), core.NewString("baz")},
		{core.NewList(core.NewList(core.NewNumber(42))),
			core.NewList(core.NewList(core.NewNumber(42), core.NewNumber(42))),
			core.NewList(core.NewList(core.NewNumber(42), core.NewNumber(42)))},
		{core.NewList(core.NewList(core.NewNumber(42), core.NewNumber(42))),
			core.NewList(core.NewList(core.NewNumber(42), core.NewNumber(42))),
			core.NewList(core.NewList(core.NewNumber(42), core.NewNumber(2049)))},
	} {
		assert.True(t, !bool(*core.EvalPure(core.PApp(Less, ts...)).(*core.BooleanType)))
	}
}

func TestLessError(t *testing.T) {
	for _, ts := range [][]core.Value{
		{core.Nil},
		{core.True},
		{core.False},
		{core.EmptyDictionary},
		{core.NewNumber(42), core.Nil},
		{core.NewString("foo"), core.Nil},
	} {
		v := core.EvalPure(core.PApp(Less, ts...))
		t.Log(v)
		_, ok := v.(*core.ErrorType)
		assert.True(t, ok)
	}
}

func TestLessNoPanic(t *testing.T) {
	assert.NotPanics(t, func() {
		core.EvalPure(core.App(Less, core.NewArguments(
			[]core.PositionalArgument{
				core.NewPositionalArgument(core.NewList(core.DummyError, core.NewNumber(42)), true),
			},
			nil)))
	})
}

func TestLessEqTrue(t *testing.T) {
	for _, ts := range [][]core.Value{
		{},
		{core.NewNumber(42)},
		{core.NewString("foo")},
		{core.NewList(core.NewNumber(42))},
		{core.NewNumber(0), core.NewNumber(0.0000000000001)},
		{core.NewNumber(0), core.NewNumber(0.00000001), core.NewNumber(0.0001)},
		{core.NewString("bar"), core.NewString("baz"), core.NewString("foo")},
		{core.NewList(core.NewNumber(42)), core.NewList(core.NewNumber(2049))},
		{core.NewList(core.NewNumber(42)), core.NewList(core.NewNumber(42), core.NewNumber(42))},
		{core.NewList(core.NewList(core.NewNumber(42))),
			core.NewList(core.NewList(core.NewNumber(42), core.NewNumber(42))),
			core.NewList(core.NewList(core.NewNumber(42), core.NewNumber(42), core.NewNumber(42)))},

		{core.NewNumber(0), core.NewNumber(42), core.NewNumber(42)},
		{core.NewNumber(0), core.NewNumber(0), core.NewNumber(42)},
		{core.NewString("bar"), core.NewString("bar"), core.NewString("baz")},
		{core.NewList(core.NewList(core.NewNumber(42))),
			core.NewList(core.NewList(core.NewNumber(42), core.NewNumber(42))),
			core.NewList(core.NewList(core.NewNumber(42), core.NewNumber(42)))},
		{core.NewList(core.NewList(core.NewNumber(42), core.NewNumber(42))),
			core.NewList(core.NewList(core.NewNumber(42), core.NewNumber(42))),
			core.NewList(core.NewList(core.NewNumber(42), core.NewNumber(2049)))},
	} {
		assert.True(t, bool(*core.EvalPure(core.PApp(LessEq, ts...)).(*core.BooleanType)))
	}
}

func TestLessEqFalse(t *testing.T) {
	for _, ts := range [][]core.Value{
		{core.NewNumber(42), core.NewNumber(0)},
		{core.NewString("foo"), core.NewString("bar"), core.NewString("baz")},
	} {
		assert.True(t, !bool(*core.EvalPure(core.PApp(LessEq, ts...)).(*core.BooleanType)))
	}
}

func TestGreaterTrue(t *testing.T) {
	for _, ts := range [][]core.Value{
		{},
		{core.NewNumber(42)},
		{core.NewString("foo")},
		{core.NewList(core.NewNumber(42))},
		{core.NewNumber(0.000001), core.NewNumber(0.0000000000001)},
		{core.NewString("foo"), core.NewString("baz"), core.NewString("bar")},
	} {
		assert.True(t, bool(*core.EvalPure(core.PApp(Greater, ts...)).(*core.BooleanType)))
	}
}

func TestGreaterFalse(t *testing.T) {
	for _, ts := range [][]core.Value{
		{core.NewNumber(42), core.NewNumber(2049)},
		{core.NewString("bar"), core.NewString("baz")},
		{core.NewNumber(42), core.NewNumber(42)},
		{core.NewString("foo"), core.NewString("bar"), core.NewString("bar")},
	} {
		assert.True(t, !bool(*core.EvalPure(core.PApp(Greater, ts...)).(*core.BooleanType)))
	}
}

func TestGreaterEqTrue(t *testing.T) {
	for _, ts := range [][]core.Value{
		{},
		{core.NewNumber(42)},
		{core.NewString("foo")},
		{core.NewList(core.NewNumber(42))},
		{core.NewNumber(0.000001), core.NewNumber(0.0000000000001)},
		{core.NewString("foo"), core.NewString("baz"), core.NewString("bar")},

		{core.NewNumber(42), core.NewNumber(42)},
		{core.NewString("foo"), core.NewString("bar"), core.NewString("bar")},
	} {
		assert.True(t, bool(*core.EvalPure(core.PApp(GreaterEq, ts...)).(*core.BooleanType)))
	}
}

func TestGreaterEqFalse(t *testing.T) {
	for _, ts := range [][]core.Value{
		{core.NewNumber(42), core.NewNumber(2049)},
		{core.NewString("bar"), core.NewString("baz")},
	} {
		assert.True(t, !bool(*core.EvalPure(core.PApp(GreaterEq, ts...)).(*core.BooleanType)))
	}
}
