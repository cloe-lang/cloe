package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListEqual(t *testing.T) {
	for _, vss := range [][2][]Value{
		{{}, {}},
		{{True}, {True}},
		{{True, False}, {True, False}},
	} {
		assert.True(t, testEqual(NewList(vss[0]...), NewList(vss[1]...)))
	}

	for _, vss := range [][2][]Value{
		{{}, {True}},
		{{True}, {False}},
		{{True, True}, {True, True, True}},
	} {
		assert.True(t, !testEqual(NewList(vss[0]...), NewList(vss[1]...)))
	}
}

func TestListComparable(t *testing.T) {
	for _, vss := range [][2][]Value{
		{{}, {True}},
		{{False}, {True}},
		{{True, False}, {True, True}},
		{{NewNumber(123), NewNumber(456)}, {NewNumber(123), NewNumber(2049)}},
		{{NewNumber(123), EmptyList}, {NewNumber(123), Nil}},
	} {
		assert.True(t, testLess(NewList(vss[0]...), NewList(vss[1]...)))
	}
}

func TestListPrepend(t *testing.T) {
	for _, vss := range [][2][]Value{
		{{}, {}},
		{{}, {True}},
		{{False}, {True}},
		{{True, False}, {True, True}},
		{{NewNumber(123), NewNumber(456)}, {NewNumber(123), NewNumber(2049)}},
		{{NewNumber(123), EmptyList}, {NewNumber(123), Nil}},
	} {
		l := PApp(Prepend, append(vss[0], NewList(vss[1]...))...)
		assert.True(t, testEqual(NewList(append(vss[0], vss[1]...)...), l))
	}
}

func TestListPrependToMergedLists(t *testing.T) {
	l := NewList(Nil)
	assert.Equal(t,
		NewNumber(2),
		EvalPure(PApp(Size, PApp(Prepend, PApp(Merge, l, PApp(Prepend, l))))).(NumberType))
}

func TestListRestWithNonListValues(t *testing.T) {
	for _, v := range []Value{
		Nil,
		EmptyDictionary,
		NewNumber(100),
		PApp(Prepend, Nil, EmptyDictionary),
	} {
		assert.Equal(t, "TypeError", EvalPure(PApp(Rest, v)).(ErrorType).Name())
	}
}

func TestListRestErrorPropagation(t *testing.T) {
	assert.Equal(t, "ValueError", EvalPure(PApp(Rest, ValueError("No way!"))).(ErrorType).Name())
}

func TestListMerge(t *testing.T) {
	for _, vss := range [][][]Value{
		{{}, {True}},
		{{False}, {True}, {True, True}},
		{{True, False}, {True, False, False, True}},
		{{NewNumber(123), NewNumber(456)}, {NewNumber(123), NewNumber(2049)}},
		{{NewNumber(123), EmptyList}, {NewNumber(123), Nil}, {True, False, True}},
	} {
		all := make([]Value, 0)
		for _, vs := range vss {
			all = append(all, vs...)
		}
		l1 := NewList(all...)

		l2 := NewList(vss[0]...)
		for _, vs := range vss[1:] {
			l2 = PApp(Merge, l2, NewList(vs...))
		}

		ls := make([]Value, 0)
		for _, vs := range vss {
			ls = append(ls, NewList(vs...))
		}
		l3 := PApp(Merge, ls...)

		assert.True(t, testEqual(l1, l2))
		assert.True(t, testEqual(l1, l3))
	}
}

func TestListToString(t *testing.T) {
	for _, xs := range []struct {
		expected string
		thunk    Value
	}{
		{"[]", EmptyList},
		{"[123]", NewList(NewNumber(123))},
		{"[123 nil]", NewList(NewNumber(123), Nil)},
		{"[[123]]", NewList(NewList(NewNumber(123)))},
		{"[nil [123]]", NewList(Nil, NewList(NewNumber(123)))},
	} {
		assert.Equal(t, NewString(xs.expected), EvalPure(PApp(ToString, xs.thunk)))
	}
}

func TestListIndex(t *testing.T) {
	a := NewString("I'm the answer.")

	for _, c := range []struct {
		list  Value
		index float64
	}{
		{NewList(a), 1},
		{NewList(True, a), 2},
		{NewList(a, False), 1},
		{NewList(True, False, a), 3},
		{NewList(Nil, Nil, Nil, Nil, a), 5},
		{NewList(Nil, Nil, Nil, a, Nil), 4},
	} {
		assert.True(t, testEqual(a, PApp(c.list, NewNumber(c.index))))
	}
}

func TestListIndexError(t *testing.T) {
	for _, c := range []struct {
		list  Value
		index float64
	}{
		{EmptyList, 0},
		{EmptyList, 1},
		{NewList(Nil), 2},
		{NewList(Nil, Nil), 1.5},
	} {
		v := EvalPure(PApp(c.list, NewNumber(c.index)))
		_, ok := v.(ErrorType)
		t.Log(v)
		assert.True(t, ok)
	}
}

func TestListToList(t *testing.T) {
	_, ok := EvalPure(PApp(ToList, EmptyList)).(ListType)
	assert.True(t, ok)
}

func TestListDelete(t *testing.T) {
	for _, c := range []struct {
		list   Value
		index  float64
		answer Value
	}{
		{NewList(Nil), 1, EmptyList},
		{NewList(Nil, True), 2, NewList(Nil)},
		{NewList(Nil, True, False), 3, NewList(Nil, True)},
	} {
		assert.True(t, testEqual(PApp(Delete, c.list, NewNumber(c.index)), c.answer))
	}
}

func TestListSize(t *testing.T) {
	for _, c := range []struct {
		list Value
		size NumberType
	}{
		{EmptyList, 0},
		{NewList(Nil), 1},
		{NewList(Nil, True), 2},
		{NewList(Nil, True, False), 3},
	} {
		assert.Equal(t, c.size, EvalPure(PApp(Size, c.list)).(NumberType))
	}
}

func TestListInclude(t *testing.T) {
	for _, c := range []struct {
		list   Value
		elem   Value
		answer BoolType
	}{
		{EmptyList, Nil, false},
		{NewList(Nil), Nil, true},
		{NewList(Nil, True), True, true},
		{NewList(Nil, False), True, false},
		{NewList(Nil, True, NewNumber(42.1), NewNumber(42), False), NewNumber(42), true},
	} {
		assert.Equal(t, c.answer, EvalPure(PApp(Include, c.list, c.elem)).(BoolType))
	}
}

func TestListInsert(t *testing.T) {
	for _, c := range []struct {
		list     Value
		index    NumberType
		elem     Value
		expected Value
	}{
		{EmptyList, 1, Nil, NewList(Nil)},
		{NewList(True), 1, False, NewList(False, True)},
		{NewList(True), 2, False, NewList(True, False)},
		{NewList(True, False), 1, Nil, NewList(Nil, True, False)},
		{NewList(True, False), 2, Nil, NewList(True, Nil, False)},
		{NewList(True, False), 3, Nil, NewList(True, False, Nil)},
	} {
		assert.True(t, testEqual(c.expected, PApp(Insert, c.list, c.index, c.elem)))
	}
}

func TestListInsertFailure(t *testing.T) {
	_, ok := EvalPure(PApp(Insert, EmptyList, NewNumber(0), Nil)).(ErrorType)
	assert.True(t, ok)
}

func TestListCompare(t *testing.T) {
	for _, c := range []struct {
		answer      int
		left, right Value
	}{
		{0, EmptyList, EmptyList},
		{0, NewList(NewNumber(42)), NewList(NewNumber(42))},
		{1, NewList(NewNumber(42)), EmptyList},
		{1, NewList(NewNumber(2049)), NewList(NewNumber(42))},
		{1, NewList(NewNumber(2049)), NewList(NewNumber(42), NewNumber(1))},
		{1, NewList(NewNumber(42), NewNumber(2049)), NewList(NewNumber(42), NewNumber(1))},
		{-1, EmptyList, NewList(NewNumber(42))},
		{-1, NewList(NewNumber(42)), NewList(NewNumber(2049))},
	} {
		assert.True(t, testCompare(c.left, c.right) == c.answer)
	}
}

func TestListFunctionsError(t *testing.T) {
	for _, v := range []Value{
		App(Prepend, NewArguments([]PositionalArgument{
			NewPositionalArgument(DummyError, true),
		}, nil, nil)),
		App(Prepend, NewArguments([]PositionalArgument{
			NewPositionalArgument(PApp(Prepend, Nil, DummyError), true),
		}, nil, nil)),
		PApp(Prepend),
		PApp(PApp(Prepend, Nil, DummyError), NewNumber(2)),
		PApp(First, EmptyList),
		PApp(Rest, EmptyList),
		PApp(Delete, EmptyList, NewNumber(1)),
		PApp(Delete, PApp(Prepend, Nil, DummyError), NewNumber(2)),
		PApp(Delete, NewList(Nil), DummyError),
		PApp(NewList(Nil), DummyError),
		PApp(NewList(Nil), NewNumber(0)),
		PApp(NewList(Nil), NewNumber(0.5)),
		PApp(Compare, NewList(Nil), NewList(Nil)),
		PApp(ToString, PApp(Prepend, Nil, DummyError)),
		PApp(Size, PApp(Prepend, Nil, DummyError)),
		PApp(Include, NewList(DummyError), Nil),
	} {
		_, ok := EvalPure(v).(ErrorType)
		assert.True(t, ok)
	}
}
