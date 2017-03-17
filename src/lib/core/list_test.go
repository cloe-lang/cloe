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

func TestListOrdered(t *testing.T) {
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
		{NewList(a), 0},
		{NewList(True, a), 1},
		{NewList(a, False), 0},
		{NewList(True, False, a), 2},
		{NewList(Nil, Nil, Nil, Nil, a), 4},
		{NewList(Nil, Nil, Nil, a, Nil), 3},
	} {
		assert.True(t, testEqual(a, PApp(li.list, NewNumber(li.index))))
	}
}

func TestXFailListIndex(t *testing.T) {
	for _, li := range []struct {
		list  *Thunk
		index float64
	}{
		{NewList(), 0},
		{NewList(), 1},
		{NewList(Nil), 1},
		{NewList(Nil, Nil), 1.5},
	} {
		v := PApp(li.list, NewNumber(li.index)).Eval()
		_, ok := v.(ErrorType)
		t.Log(v)
		assert.True(t, ok)
	}
}
