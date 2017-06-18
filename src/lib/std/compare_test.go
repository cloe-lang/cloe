package std

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tisp-lang/tisp/src/lib/core"
)

func TestLessTrue(t *testing.T) {
	for _, ts := range [][]*core.Thunk{
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
		assert.True(t, bool(core.PApp(Less, ts...).Eval().(core.BoolType)))
	}
}

func TestLessFalse(t *testing.T) {
	for _, ts := range [][]*core.Thunk{
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
		assert.True(t, !bool(core.PApp(Less, ts...).Eval().(core.BoolType)))
	}
}

func TestLessError(t *testing.T) {
	for _, ts := range [][]*core.Thunk{
		{core.Nil},
		{core.True},
		{core.False},
		{core.NewDictionary(nil, nil)},
		{core.NewNumber(42), core.Nil},
		{core.NewString("foo"), core.Nil},
	} {
		v := core.PApp(Less, ts...).Eval()
		t.Log(v)
		_, ok := v.(core.ErrorType)
		assert.True(t, ok)
	}
}

func TestLessEqTrue(t *testing.T) {
	for _, ts := range [][]*core.Thunk{
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
		assert.True(t, bool(core.PApp(LessEq, ts...).Eval().(core.BoolType)))
	}
}

func TestLessEqFalse(t *testing.T) {
	for _, ts := range [][]*core.Thunk{
		{core.NewNumber(42), core.NewNumber(0)},
		{core.NewString("foo"), core.NewString("bar"), core.NewString("baz")},
	} {
		assert.True(t, !bool(core.PApp(LessEq, ts...).Eval().(core.BoolType)))
	}
}

func TestGreaterTrue(t *testing.T) {
	for _, ts := range [][]*core.Thunk{
		{},
		{core.NewNumber(42)},
		{core.NewString("foo")},
		{core.NewList(core.NewNumber(42))},
		{core.NewNumber(0.000001), core.NewNumber(0.0000000000001)},
		{core.NewString("foo"), core.NewString("baz"), core.NewString("bar")},
	} {
		assert.True(t, bool(core.PApp(Greater, ts...).Eval().(core.BoolType)))
	}
}

func TestGreaterFalse(t *testing.T) {
	for _, ts := range [][]*core.Thunk{
		{core.NewNumber(42), core.NewNumber(2049)},
		{core.NewString("bar"), core.NewString("baz")},
		{core.NewNumber(42), core.NewNumber(42)},
		{core.NewString("foo"), core.NewString("bar"), core.NewString("bar")},
	} {
		assert.True(t, !bool(core.PApp(Greater, ts...).Eval().(core.BoolType)))
	}
}

func TestGreaterEqTrue(t *testing.T) {
	for _, ts := range [][]*core.Thunk{
		{},
		{core.NewNumber(42)},
		{core.NewString("foo")},
		{core.NewList(core.NewNumber(42))},
		{core.NewNumber(0.000001), core.NewNumber(0.0000000000001)},
		{core.NewString("foo"), core.NewString("baz"), core.NewString("bar")},

		{core.NewNumber(42), core.NewNumber(42)},
		{core.NewString("foo"), core.NewString("bar"), core.NewString("bar")},
	} {
		assert.True(t, bool(core.PApp(GreaterEq, ts...).Eval().(core.BoolType)))
	}
}

func TestGreaterEqFalse(t *testing.T) {
	for _, ts := range [][]*core.Thunk{
		{core.NewNumber(42), core.NewNumber(2049)},
		{core.NewString("bar"), core.NewString("baz")},
	} {
		assert.True(t, !bool(core.PApp(GreaterEq, ts...).Eval().(core.BoolType)))
	}
}
