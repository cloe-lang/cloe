package comb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBlock(t *testing.T) {
	for _, test := range []struct {
		indent int
		source string
	}{
		{0, "foo\nfoo\nfoo\nbar"},
		{1, "foo\n foo\n foo \n bar"},
		{2, "foo\n  foo\n  foo\n  bar"},
	} {
		s := NewState(test.source)
		blank := s.Many(s.InString(" \n"))
		foo := s.And(s.String("foo"), blank)
		result, err := s.Exhaust(s.And(blank, s.Block(test.indent, foo, foo, s.String("bar"))))()

		t.Logf("%#v", result)

		assert.NotEqual(t, result, nil)
		assert.Equal(t, err, nil)
	}
}

func TestBlockFail(t *testing.T) {
	for _, test := range []struct {
		indent int
		source string
	}{
		{0, "foo\n foo\nfoo\nbar"},
		{1, "foo\nfoo\n foo\n bar"},
		{2, "foo\n  foo\n foo\n  bar"},
		{1, "foo\n foo\n foo\n  bar"},
	} {
		s := NewState(test.source)
		blank := s.Many(s.InString(" \n"))
		foo := s.And(s.String("foo"), blank)
		result, err := s.Exhaust(s.And(blank, s.Block(test.indent, foo, foo, s.String("bar"))))()

		t.Logf("%#v", err)

		assert.Equal(t, result, nil)
		assert.NotEqual(t, err, nil)
	}
}

func TestSameLine(t *testing.T) {
	for _, str := range []string{
		"foo foo",
		"   foo    foo",
	} {
		s := NewState(str)
		blank := s.Many(s.InString(" "))
		foo := s.And(s.String("foo"), blank)
		result, err := s.Exhaust(s.And(blank, s.SameLine(foo, foo)))()

		t.Logf("%#v", result)

		assert.NotEqual(t, result, nil)
		assert.Equal(t, err, nil)
	}
}

func TestSameLineFail(t *testing.T) {
	for _, str := range []string{
		"foo\nfoo",
		"   foo  \n  foo",
	} {
		s := NewState(str)
		blank := s.Many(s.InString(" \n"))
		foo := s.And(s.String("foo"), blank)
		result, err := s.Exhaust(s.And(blank, s.SameLine(foo, foo)))()

		t.Logf("%#v", result)

		assert.Equal(t, result, nil)
		assert.NotEqual(t, err, nil)
	}
}

func TestSameLineOrIndented(t *testing.T) {
	for _, str := range []string{
		"foo foo",
		"foo\n  foo",
		"foo\n   foo",
		"foo\n\n  foo",
		"  foo    foo",
		"  foo\n    foo",
		"  foo\n     foo",
		"  foo\n\n    foo",
	} {
		s := NewState(str)
		blank := s.Many(s.InString(" \n"))
		foo := s.And(s.String("foo"), blank)
		result, err := s.Exhaust(s.And(blank, s.WithPosition(s.And(foo, s.SameLineOrIndented(2, foo)))))()

		t.Logf("%#v", result)

		assert.NotEqual(t, result, nil)
		assert.Equal(t, err, nil)
	}
}

func TestSameLineOrIndentedFail(t *testing.T) {
	for _, str := range []string{
		"foo\nfoo",
		"foo\n foo",
		"  foo\nfoo",
		"  foo\n foo",
		"  foo\n  foo",
		"  foo\n   foo",
	} {
		s := NewState(str)
		blank := s.Many(s.InString(" \n"))
		foo := s.And(s.String("foo"), blank)
		result, err := s.Exhaust(s.And(blank, s.WithPosition(s.And(foo, s.SameLineOrIndented(2, foo)))))()

		t.Logf("%#v", err)

		assert.Equal(t, result, nil)
		assert.NotEqual(t, err, nil)
	}
}

func TestNoIndent(t *testing.T) {
	s := NewState("foo")
	blank := s.Many(s.InString(" \n"))
	result, err := s.Exhaust(s.And(blank, s.NoIndent(s.And(s.String("foo"), blank))))()

	t.Logf("%#v", result)

	assert.NotEqual(t, result, nil)
	assert.Equal(t, err, nil)
}

func TestNoIndentFail(t *testing.T) {
	for _, str := range []string{
		" foo",
		"  foo",
	} {
		s := NewState(str)
		blank := s.Many(s.InString(" \n"))
		result, err := s.Exhaust(s.And(blank, s.NoIndent(s.And(s.String("foo"), blank))))()

		t.Logf("%#v", err)

		assert.Equal(t, result, nil)
		assert.NotEqual(t, err, nil)
	}
}
