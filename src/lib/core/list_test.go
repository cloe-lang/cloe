package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListEqual(t *testing.T) {
	for _, tss := range [][2][]*Thunk{
		{{}, {}},
		{{True}, {True}},
		{{True, False}, {True, False}},
	} {
		assert.True(t, testEqual(NewList(tss[0]...), NewList(tss[1]...)))
	}

	for _, tss := range [][2][]*Thunk{
		{{}, {True}},
		{{True}, {False}},
		{{True, True}, {True, True, True}},
	} {
		assert.True(t, !testEqual(NewList(tss[0]...), NewList(tss[1]...)))
	}
}

func TestListComparable(t *testing.T) {
	for _, tss := range [][2][]*Thunk{
		{{}, {True}},
		{{False}, {True}},
		{{True, False}, {True, True}},
		{{NewNumber(123), NewNumber(456)}, {NewNumber(123), NewNumber(2049)}},
		{{NewNumber(123), NewList()}, {NewNumber(123), Nil}},
	} {
		assert.True(t, testLess(NewList(tss[0]...), NewList(tss[1]...)))
	}
}

func TestListPrepend(t *testing.T) {
	for _, tss := range [][2][]*Thunk{
		{{}, {True}},
		{{False}, {True}},
		{{True, False}, {True, True}},
		{{NewNumber(123), NewNumber(456)}, {NewNumber(123), NewNumber(2049)}},
		{{NewNumber(123), NewList()}, {NewNumber(123), Nil}},
	} {
		l := PApp(Prepend, append(tss[0], NewList(tss[1]...))...)
		assert.True(t, testEqual(NewList(append(tss[0], tss[1]...)...), l))
	}
}

func TestListRestWithNonListValues(t *testing.T) {
	for _, th := range []*Thunk{Nil, EmptyDictionary, NewNumber(100)} {
		assert.Equal(t, "TypeError", PApp(Rest, th).Eval().(ErrorType).Name())
	}

	assert.Equal(t, "ValueError", PApp(Rest, ValueError("No way!")).Eval().(ErrorType).Name())
}

func TestListAppend(t *testing.T) {
	for _, tss := range [][2][]*Thunk{
		{{}, {True}},
		{{False}, {True}},
		{{True, False}, {True, True}},
		{{NewNumber(123), NewNumber(456)}, {NewNumber(123), NewNumber(2049)}},
		{{NewNumber(123), NewList()}, {NewNumber(123), Nil}},
	} {
		l := NewList(tss[0]...)

		for _, t := range tss[1] {
			l = PApp(Append, l, t)
		}

		assert.True(t, testEqual(NewList(append(tss[0], tss[1]...)...), l))
	}
}

func TestListMerge(t *testing.T) {
	for _, tss := range [][][]*Thunk{
		{{}, {True}},
		{{False}, {True}, {True, True}},
		{{True, False}, {True, False, False, True}},
		{{NewNumber(123), NewNumber(456)}, {NewNumber(123), NewNumber(2049)}},
		{{NewNumber(123), NewList()}, {NewNumber(123), Nil}, {True, False, True}},
	} {
		all := make([]*Thunk, 0)
		for _, ts := range tss {
			all = append(all, ts...)
		}
		l1 := NewList(all...)

		l2 := NewList(tss[0]...)
		for _, ts := range tss[1:] {
			l2 = PApp(Merge, l2, NewList(ts...))
		}

		ls := make([]*Thunk, 0)
		for _, ts := range tss {
			ls = append(ls, NewList(ts...))
		}
		l3 := PApp(Merge, ls...)

		assert.True(t, testEqual(l1, l2))
		assert.True(t, testEqual(l1, l3))
	}
}

func TestListToString(t *testing.T) {
	for _, xs := range []struct {
		expected string
		thunk    *Thunk
	}{
		{"[]", NewList()},
		{"[123]", NewList(NewNumber(123))},
		{"[123 nil]", NewList(NewNumber(123), Nil)},
		{"[[123]]", NewList(NewList(NewNumber(123)))},
		{"[nil [123]]", NewList(Nil, NewList(NewNumber(123)))},
	} {
		assert.Equal(t, StringType(xs.expected), PApp(ToString, xs.thunk).Eval())
	}
}

func TestListIndex(t *testing.T) {
	a := NewString("I'm the answer.")

	for _, li := range []struct {
		list  *Thunk
		index float64
	}{
		{NewList(a), 1},
		{NewList(True, a), 2},
		{NewList(a, False), 1},
		{NewList(True, False, a), 3},
		{NewList(Nil, Nil, Nil, Nil, a), 5},
		{NewList(Nil, Nil, Nil, a, Nil), 4},
	} {
		assert.True(t, testEqual(a, PApp(li.list, NewNumber(li.index))))
	}
}

func TestListIndexFail(t *testing.T) {
	for _, li := range []struct {
		list  *Thunk
		index float64
	}{
		{NewList(), 0},
		{NewList(), 1},
		{NewList(Nil), 2},
		{NewList(Nil, Nil), 1.5},
	} {
		v := PApp(li.list, NewNumber(li.index)).Eval()
		_, ok := v.(ErrorType)
		t.Log(v)
		assert.True(t, ok)
	}
}

func TestListToList(t *testing.T) {
	_, ok := PApp(ToList, EmptyList).Eval().(ListType)
	assert.True(t, ok)
}

func TestListDelete(t *testing.T) {
	for _, test := range []struct {
		list   *Thunk
		index  float64
		answer *Thunk
	}{
		{NewList(Nil), 1, NewList()},
		{NewList(Nil, True), 2, NewList(Nil)},
		{NewList(Nil, True, False), 3, NewList(Nil, True)},
	} {
		assert.True(t, testEqual(PApp(Delete, test.list, NewNumber(test.index)), test.answer))
	}
}

func TestListSize(t *testing.T) {
	for _, test := range []struct {
		list *Thunk
		size NumberType
	}{
		{NewList(), 0},
		{NewList(Nil), 1},
		{NewList(Nil, True), 2},
		{NewList(Nil, True, False), 3},
	} {
		assert.Equal(t, test.size, PApp(Size, test.list).Eval().(NumberType))
	}
}

func TestListInclude(t *testing.T) {
	for _, test := range []struct {
		list   *Thunk
		elem   *Thunk
		answer BoolType
	}{
		{NewList(), Nil, false},
		{NewList(Nil), Nil, true},
		{NewList(Nil, True), True, true},
		{NewList(Nil, False), True, false},
		{NewList(Nil, True, NewNumber(42.1), NewNumber(42), False), NewNumber(42), true},
	} {
		assert.Equal(t, test.answer, PApp(Include, test.list, test.elem).Eval().(BoolType))
	}
}

func TestListInsert(t *testing.T) {
	for _, test := range []struct {
		list     *Thunk
		index    NumberType
		elem     *Thunk
		expected *Thunk
	}{
		{NewList(), 1, Nil, NewList(Nil)},
		{NewList(True), 1, False, NewList(False, True)},
		{NewList(True), 2, False, NewList(True, False)},
		{NewList(True, False), 1, Nil, NewList(Nil, True, False)},
		{NewList(True, False), 2, Nil, NewList(True, Nil, False)},
		{NewList(True, False), 3, Nil, NewList(True, False, Nil)},
	} {
		assert.True(t, testEqual(test.expected, PApp(Insert, test.list, Normal(test.index), test.elem)))
	}
}

func TestListInsertFailure(t *testing.T) {
	_, ok := PApp(Insert, EmptyList, NewNumber(0), Nil).Eval().(ErrorType)
	assert.True(t, ok)
}

func TestListCompare(t *testing.T) {
	assert.True(t, testCompare(NewList(), NewList()) == 0)
	assert.True(t, testCompare(NewList(NewNumber(42)), NewList()) == 1)
	assert.True(t, testCompare(NewList(), NewList(NewNumber(42))) == -1)
	assert.True(t, testCompare(NewList(NewNumber(42)), NewList(NewNumber(42))) == 0)
	assert.True(t, testCompare(NewList(NewNumber(42)), NewList(NewNumber(2049))) == -1)
	assert.True(t, testCompare(NewList(NewNumber(2049)), NewList(NewNumber(42))) == 1)
	assert.True(t, testCompare(NewList(NewNumber(2049)), NewList(NewNumber(42), NewNumber(1))) == 1)
	assert.True(t, testCompare(
		NewList(NewNumber(42), NewNumber(2049)),
		NewList(NewNumber(42), NewNumber(1))) == 1)
}
