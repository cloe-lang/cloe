package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func testEqual(t1, t2 *Thunk) bool {
	return compare(t1.Eval(), t2.Eval()) == 0
}

func testLess(t1, t2 *Thunk) bool {
	return t1.Eval().(comparable).compare(t2.Eval().(comparable)) < 0
}

func TestXFailLess(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fail()
		}
	}()

	compare(NewNumber(42).Eval(), NewError("you", "failed.").Eval())
}

func testCompare(t1, t2 *Thunk) NumberType {
	return PApp(Compare, t1, t2).Eval().(NumberType)
}

func TestCompareWithInvalidValues(t *testing.T) {
	for _, ts := range [][2]*Thunk{
		{True, False},
		{Nil, Nil},
		{NewNumber(0), False},
		{NewNumber(0), Nil},
		{True, Nil},
		{NewDictionary([]Value{Nil.Eval()}, []*Thunk{Nil}),
			NewDictionary([]Value{Nil.Eval()}, []*Thunk{Nil})},
		{NotNumberError(Nil), NotNumberError(Nil)},
	} {
		v := PApp(Compare, ts[0], ts[1]).Eval()

		t.Log(v)

		_, ok := v.(ErrorType)
		assert.True(t, ok)
	}
}
