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
		{0, "foo\nfoo\nfoo"},
		{1, " foo\n foo\n foo"},
		{2, "  foo\n  foo\n  foo"},
	} {
		s := NewState(test.source)
		blank := s.Many(s.InString(" \n"))
		result, err := s.Exhaust(s.Block(test.indent, blank, s.And(s.String("foo"), blank), s.None()))()

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
		{0, "foo\n foo\nfoo"},
		{1, " foo\nfoo\n foo"},
		{2, "  foo\n  foo\n foo"},
	} {
		s := NewState(test.source)
		blank := s.Many(s.InString(" \n"))
		result, err := s.Exhaust(s.Block(test.indent, blank, s.And(s.String("foo"), blank), s.None()))()

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
