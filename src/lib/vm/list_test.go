package vm

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
